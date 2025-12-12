package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Reaction struct {
	ID       int    `json:"id"`
	Emoji    string `json:"emoji"`
	UserID   int    `json:"user_id"`
	ObjectID int    `json:"object_id"`
}

func (c *Client) AddReaction(objectType string, objectID int, emoji string) (*Reaction, error) {
	if objectType != "thread" && objectType != "comment" {
		return nil, fmt.Errorf("invalid object type: must be 'thread' or 'comment'")
	}

	payload := map[string]interface{}{
		"object_type": objectType,
		"object_id":   objectID,
		"emoji":       emoji,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/reactions/add"
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

	var reaction Reaction
	if err := json.NewDecoder(resp.Body).Decode(&reaction); err != nil {
		return nil, fmt.Errorf("failed to parse reaction response: %w", err)
	}

	return &reaction, nil
}

func (c *Client) RemoveReaction(objectType string, objectID int, emoji string) error {
	if objectType != "thread" && objectType != "comment" {
		return fmt.Errorf("invalid object type: must be 'thread' or 'comment'")
	}

	payload := map[string]interface{}{
		"object_type": objectType,
		"object_id":   objectID,
		"emoji":       emoji,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/reactions/remove"
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

func (c *Client) GetReactions(objectType string, objectID int) ([]Reaction, error) {
	var endpoint string
	if objectType == "thread" {
		endpoint = fmt.Sprintf("/reactions/get?thread=%d", objectID)
	} else if objectType == "comment" {
		endpoint = fmt.Sprintf("/reactions/get?comment=%d", objectID)
	} else {
		return nil, fmt.Errorf("invalid object type: must be 'thread' or 'comment'")
	}

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var reactions []Reaction
	if err := json.Unmarshal(body, &reactions); err != nil {
		return nil, fmt.Errorf("failed to parse reactions response: %w", err)
	}

	return reactions, nil
}
