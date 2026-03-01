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
}

type Notification struct {
	Title    string
	Body     string
	Priority int
	Tags     []string
	ClickURL string
}

func NewClient(baseURL, topic string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
		topic:   topic,
	}
}

func (c *Client) Send(n Notification) error {
	url := fmt.Sprintf("%s/%s", c.baseURL, c.topic)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(n.Body))
	if err != nil {
		return fmt.Errorf("creating ntfy request: %w", err)
	}

	req.Header.Set("Title", n.Title)
	req.Header.Set("Priority", strconv.Itoa(n.Priority))

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
