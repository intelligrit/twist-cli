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
	endpoint := fmt.Sprintf("/comments/get_all?thread_id=%d", threadID)
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

func (c *Client) PostComment(threadID int, content string) (*Comment, error) {
	payload := map[string]interface{}{
		"thread_id": threadID,
		"content":   content,
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
