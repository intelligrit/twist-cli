package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Thread struct {
	ID              int      `json:"id"`
	Title           string   `json:"title"`
	Content         string   `json:"content"`
	ChannelID       int      `json:"channel_id"`
	WorkspaceID     int      `json:"workspace_id"`
	Creator         int      `json:"creator"`
	PostedTS        int64    `json:"posted_ts"`
	LastUpdatedTS   int64    `json:"last_updated_ts"`
	CommentCount    int      `json:"comment_count"`
	Starred         bool     `json:"starred"`
	Pinned          bool     `json:"pinned"`
	Archived        bool     `json:"archived"`
	Participants    []int    `json:"participants"`
}

type Comment struct {
	ID            int    `json:"id"`
	Content       string `json:"content"`
	ThreadID      int    `json:"thread_id"`
	Creator       int    `json:"creator"`
	PostedTS      int64  `json:"posted_ts"`
	LastUpdatedTS int64  `json:"last_updated_ts"`
}

func (c *Client) GetThreads(channelID int) ([]Thread, error) {
	endpoint := fmt.Sprintf("/threads/get?channel_id=%d", channelID)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var threads []Thread
	if err := json.Unmarshal(body, &threads); err != nil {
		return nil, fmt.Errorf("failed to parse threads response: %w", err)
	}

	return threads, nil
}

func (c *Client) GetThread(id int) (*Thread, error) {
	endpoint := fmt.Sprintf("/threads/getone?id=%d", id)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var thread Thread
	if err := json.Unmarshal(body, &thread); err != nil {
		return nil, fmt.Errorf("failed to parse thread response: %w", err)
	}

	return &thread, nil
}

func (c *Client) GetComments(threadID int) ([]Comment, error) {
	endpoint := fmt.Sprintf("/comments/get?thread_id=%d", threadID)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("failed to parse comments response: %w", err)
	}

	return comments, nil
}

func (c *Client) CreateThread(channelID int, title, content string, recipients []int) (*Thread, error) {
	payload := map[string]interface{}{
		"channel_id": channelID,
		"title":      title,
		"content":    content,
	}

	if len(recipients) > 0 {
		payload["recipients"] = recipients
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal thread data: %w", err)
	}

	url := BaseURL + "/threads/add"
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

	var thread Thread
	if err := json.NewDecoder(resp.Body).Decode(&thread); err != nil {
		return nil, fmt.Errorf("failed to parse thread response: %w", err)
	}

	return &thread, nil
}

func (c *Client) PostComment(threadID int, content string, recipients []int) (*Comment, error) {
	payload := map[string]interface{}{
		"thread_id": threadID,
		"content":   content,
	}

	if len(recipients) > 0 {
		payload["recipients"] = recipients
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal comment data: %w", err)
	}

	url := BaseURL + "/comments/add"
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

	var comment Comment
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return nil, fmt.Errorf("failed to parse comment response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return &comment, nil
}

func (c *Client) UpdateThread(id int, updates map[string]interface{}) (*Thread, error) {
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

	url := BaseURL + "/threads/update"
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

	var thread Thread
	if err := json.NewDecoder(resp.Body).Decode(&thread); err != nil {
		return nil, fmt.Errorf("failed to parse thread response: %w", err)
	}

	return &thread, nil
}

func (c *Client) DeleteThread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/threads/remove"
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

func (c *Client) PinThread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/threads/pin"
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

func (c *Client) UnpinThread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/threads/unpin"
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

func (c *Client) StarThread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/threads/star"
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

func (c *Client) UnstarThread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/threads/unstar"
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

func (c *Client) ArchiveThread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/threads/archive"
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

func (c *Client) UnarchiveThread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/threads/unarchive"
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

func (c *Client) UpdateComment(id int, content string) (*Comment, error) {
	payload := map[string]interface{}{
		"id":      id,
		"content": content,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/comments/update"
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

	var comment Comment
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return nil, fmt.Errorf("failed to parse comment response: %w", err)
	}

	return &comment, nil
}

func (c *Client) DeleteComment(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/comments/remove"
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
