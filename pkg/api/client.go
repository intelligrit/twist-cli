package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BaseURL = "https://api.twist.com/api/v3"

type Client struct {
	token      string
	httpClient *http.Client
}

type APIError struct {
	Error []interface{} `json:"error"`
}

func (e *APIError) String() string {
	if len(e.Error) >= 2 {
		return fmt.Sprintf("API error %v: %v", e.Error[0], e.Error[1])
	}
	return fmt.Sprintf("API error: %v", e.Error)
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
	}
}

func (c *Client) doRequest(method, endpoint string) ([]byte, error) {
	url := BaseURL + endpoint
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(body, &apiErr); err == nil && len(apiErr.Error) > 0 {
			return nil, fmt.Errorf("%s", apiErr.String())
		}
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
