package api

import (
	"encoding/json"
	"fmt"
	"net/url"
)

func (c *Client) SearchThreads(workspaceID int, query string, opts map[string]interface{}) ([]Thread, error) {
	endpoint := fmt.Sprintf("/threads/search?workspace_id=%d&query=%s", workspaceID, url.QueryEscape(query))

	if channelID, ok := opts["channel_id"].(int); ok {
		endpoint += fmt.Sprintf("&channel_id=%d", channelID)
	}
	if limit, ok := opts["limit"].(int); ok {
		endpoint += fmt.Sprintf("&limit=%d", limit)
	}

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var threads []Thread
	if err := json.Unmarshal(body, &threads); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return threads, nil
}

func (c *Client) SearchMessages(workspaceID int, query string, opts map[string]interface{}) ([]Comment, error) {
	endpoint := fmt.Sprintf("/comments/search?workspace_id=%d&query=%s", workspaceID, url.QueryEscape(query))

	if limit, ok := opts["limit"].(int); ok {
		endpoint += fmt.Sprintf("&limit=%d", limit)
	}

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return comments, nil
}

func (c *Client) SearchConversations(query string, opts map[string]interface{}) ([]ConversationMessage, error) {
	endpoint := fmt.Sprintf("/conversation_messages/search?query=%s", url.QueryEscape(query))

	if limit, ok := opts["limit"].(int); ok {
		endpoint += fmt.Sprintf("&limit=%d", limit)
	}

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var messages []ConversationMessage
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return messages, nil
}
