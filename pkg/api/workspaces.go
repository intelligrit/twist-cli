package api

import (
	"encoding/json"
	"fmt"
)

type Workspace struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Creator   int    `json:"creator"`
	CreatedTS int64  `json:"created_ts"`
	Plan      string `json:"plan"`
}

func (c *Client) GetWorkspaces() ([]Workspace, error) {
	body, err := c.doRequest("GET", "/workspaces/get")
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	if err := json.Unmarshal(body, &workspaces); err != nil {
		return nil, fmt.Errorf("failed to parse workspaces response: %w", err)
	}

	return workspaces, nil
}
