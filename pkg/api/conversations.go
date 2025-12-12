package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Conversation struct {
	ID           int   `json:"id"`
	UserIDs      []int `json:"user_ids"`
	MessageCount int   `json:"message_count"`
	CreatedTS    int64 `json:"created_ts"`
	IsArchived   bool  `json:"is_archived"`
	IsMuted      bool  `json:"is_muted"`
}

type ConversationMessage struct {
	ID             int    `json:"id"`
	ConversationID int    `json:"conversation_id"`
	Content        string `json:"content"`
	UserID         int    `json:"user_id"`
	CreatedTS      int64  `json:"created_ts"`
}

func (c *Client) GetConversations() ([]Conversation, error) {
	body, err := c.doRequest("GET", "/conversations/get")
	if err != nil {
		return nil, err
	}

	var conversations []Conversation
	if err := json.Unmarshal(body, &conversations); err != nil {
		return nil, fmt.Errorf("failed to parse conversations response: %w", err)
	}

	return conversations, nil
}

func (c *Client) GetConversationMessages(conversationID int) ([]ConversationMessage, error) {
	endpoint := fmt.Sprintf("/conversation_messages/get?conversation_id=%d", conversationID)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var messages []ConversationMessage
	if err := json.Unmarshal(body, &messages); err != nil {
		return nil, fmt.Errorf("failed to parse messages response: %w", err)
	}

	return messages, nil
}

func (c *Client) GetOrCreateConversation(userIDs []int) (*Conversation, error) {
	payload := map[string]interface{}{
		"user_ids": userIDs,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal conversation data: %w", err)
	}

	url := BaseURL + "/conversations/get_or_create"
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

	var conversation Conversation
	if err := json.NewDecoder(resp.Body).Decode(&conversation); err != nil {
		return nil, fmt.Errorf("failed to parse conversation response: %w", err)
	}

	return &conversation, nil
}

func (c *Client) SendConversationMessage(conversationID int, content string, recipients []int) (*ConversationMessage, error) {
	payload := map[string]interface{}{
		"conversation_id": conversationID,
		"content":         content,
	}

	if len(recipients) > 0 {
		payload["recipients"] = recipients
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal message data: %w", err)
	}

	url := BaseURL + "/conversation_messages/add"
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

	var message ConversationMessage
	if err := json.NewDecoder(resp.Body).Decode(&message); err != nil {
		return nil, fmt.Errorf("failed to parse message response: %w", err)
	}

	return &message, nil
}

func (c *Client) ArchiveConversation(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/conversations/archive"
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

func (c *Client) UnarchiveConversation(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/conversations/unarchive"
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

func (c *Client) MuteConversation(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/conversations/mute"
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

func (c *Client) UnmuteConversation(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/conversations/unmute"
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

func (c *Client) MarkConversationRead(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/conversations/mark_read"
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

func (c *Client) MarkConversationUnread(id int) error {
	payload := map[string]interface{}{
		"id": id,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	url := BaseURL + "/conversations/mark_unread"
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
