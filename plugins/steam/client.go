package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const baseURL = "https://api.steampowered.com/IPhysicalGoodsService/CheckInventoryAvailableByPackage/v1"

type Client struct {
	httpClient *http.Client
}

type InventoryResponse struct {
	InventoryAvailable bool `json:"inventory_available"`
	HighPendingOrders  bool `json:"high_pending_orders"`
}

type apiResponse struct {
	Response InventoryResponse `json:"response"`
}

type Option func(*Client)

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}

func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) CheckInventory(packageID int, countryCode string) (*InventoryResponse, error) {
	params := url.Values{}
	params.Set("origin", "https://store.steampowered.com")
	params.Set("country_code", countryCode)
	params.Set("packageid", strconv.Itoa(packageID))

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("requesting inventory for package %d: %w", packageID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d for package %d", resp.StatusCode, packageID)
	}

	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("decoding response for package %d: %w", packageID, err)
	}

	return &apiResp.Response, nil
}
