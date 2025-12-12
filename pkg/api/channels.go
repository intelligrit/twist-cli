package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Channel struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	WorkspaceID int    `json:"workspace_id"`
	Public      bool   `json:"public"`
	Archived    bool   `json:"archived"`
	Color       int    `json:"color"`
	Icon        int    `json:"icon"`
	CreatedTS   int64  `json:"created_ts"`
}

func (c *Client) GetChannels(workspaceID int, archived bool) ([]Channel, error) {
	endpoint := fmt.Sprintf("/channels/get?workspace_id=%d", workspaceID)
	if archived {
		endpoint += "&archived=true"
	}
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var channels []Channel
	if err := json.Unmarshal(body, &channels); err != nil {
		return nil, fmt.Errorf("failed to parse channels response: %w", err)
	}

	return channels, nil
}

func (c *Client) GetChannel(id int) (*Channel, error) {
	endpoint := fmt.Sprintf("/channels/getone?id=%d", id)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var channel Channel
	if err := json.Unmarshal(body, &channel); err != nil {
		return nil, fmt.Errorf("failed to parse channel response: %w", err)
	}

	return &channel, nil
}

func (c *Client) CreateChannel(workspaceID int, name string, opts map[string]interface{}) (*Channel, error) {
	payload := map[string]interface{}{
		"workspace_id": workspaceID,
		"name":         name,
	}

	for k, v := range opts {
		payload[k] = v
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal channel data: %w", err)
	}

	url := BaseURL + "/channels/add"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var channel Channel
	if err := json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, fmt.Errorf("failed to parse channel response: %w", err)
	}

	return &channel, nil
}

func (c *Client) UpdateChannel(id int, updates map[string]interface{}) (*Channel, error) {
	payload := map[string]interface{}{
		"id": id,
	}

	for k, v := range updates {
		payload[k] = v
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update data: %w", err)
	}

	url := BaseURL + "/channels/update"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var channel Channel
	if err := json.NewDecoder(resp.Body).Decode(&channel); err != nil {
		return nil, fmt.Errorf("failed to parse channel response: %w", err)
	}

	return &channel, nil
}

func (c *Client) ArchiveChannel(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/channels/archive"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) UnarchiveChannel(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/channels/unarchive"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) DeleteChannel(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/channels/remove"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) AddChannelUser(channelID, userID int) error {
	payload := map[string]interface{}{
		"id":      channelID,
		"user_id": userID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/channels/add_user"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) RemoveChannelUser(channelID, userID int) error {
	payload := map[string]interface{}{
		"id":      channelID,
		"user_id": userID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/channels/remove_user"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return nil
}
