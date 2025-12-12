package api

import (
	"encoding/json"
	"fmt"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserType string `json:"user_type"`
	Bot      bool   `json:"bot"`
	Removed  bool   `json:"removed"`
}

func (c *Client) GetWorkspaceUsers(workspaceID int) ([]User, error) {
	endpoint := fmt.Sprintf("/workspace_users/get?id=%d", workspaceID)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var users []User
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("failed to parse users response: %w", err)
	}

	return users, nil
}
