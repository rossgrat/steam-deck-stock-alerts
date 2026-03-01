package ntfy

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	topic      string
	token      string
}

type Notification struct {
	Title    string
	Body     string
	Priority int
	Tags     []string
	ClickURL string
}

type Option func(*Client)

func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

func NewClient(baseURL, topic string, opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		topic:   topic,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Send(n Notification) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, c.topic)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(n.Body))
	if err != nil {
		return fmt.Errorf("creating ntfy request: %w", err)
	}

	req.Header.Set("Title", n.Title)
	req.Header.Set("Priority", strconv.Itoa(n.Priority))

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	if len(n.Tags) > 0 {
		req.Header.Set("Tags", strings.Join(n.Tags, ","))
	}

	if n.ClickURL != "" {
		req.Header.Set("Click", n.ClickURL)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending ntfy notification: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ntfy returned status %d", resp.StatusCode)
	}

	return nil
}
