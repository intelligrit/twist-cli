package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type Attachment struct {
	ID         int    `json:"id"`
	Title      string `json:"title"`
	URL        string `json:"url"`
	Size       int64  `json:"size"`
	MimeType   string `json:"mime_type"`
	UploadedTS int64  `json:"uploaded_ts"`
}

func (c *Client) UploadAttachment(targetType string, targetID int, filePath string) (*Attachment, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if targetType == "thread" {
		writer.WriteField("thread_id", fmt.Sprintf("%d", targetID))
	} else if targetType == "comment" {
		writer.WriteField("comment_id", fmt.Sprintf("%d", targetID))
	} else if targetType == "conversation" {
		writer.WriteField("conversation_id", fmt.Sprintf("%d", targetID))
	} else {
		return nil, fmt.Errorf("invalid target type: must be 'thread', 'comment', or 'conversation'")
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	url := BaseURL + "/attachments/upload"
	req, err := http.NewRequest("POST", url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var attachment Attachment
	if err := json.NewDecoder(resp.Body).Decode(&attachment); err != nil {
		return nil, fmt.Errorf("failed to parse attachment response: %w", err)
	}

	return &attachment, nil
}

func (c *Client) DownloadAttachment(id int, outputPath string) error {
	endpoint := fmt.Sprintf("/attachments/getone?id=%d", id)
	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return err
	}

	var attachment Attachment
	if err := json.Unmarshal(body, &attachment); err != nil {
		return fmt.Errorf("failed to parse attachment response: %w", err)
	}

	resp, err := http.Get(attachment.URL)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *Client) GetAttachments(targetType string, targetID int) ([]Attachment, error) {
	var endpoint string
	if targetType == "thread" {
		endpoint = fmt.Sprintf("/attachments/get?thread_id=%d", targetID)
	} else if targetType == "comment" {
		endpoint = fmt.Sprintf("/attachments/get?comment_id=%d", targetID)
	} else if targetType == "conversation" {
		endpoint = fmt.Sprintf("/attachments/get?conversation_id=%d", targetID)
	} else {
		return nil, fmt.Errorf("invalid target type: must be 'thread', 'comment', or 'conversation'")
	}

	body, err := c.doRequest("GET", endpoint)
	if err != nil {
		return nil, err
	}

	var attachments []Attachment
	if err := json.Unmarshal(body, &attachments); err != nil {
		return nil, fmt.Errorf("failed to parse attachments response: %w", err)
	}

	return attachments, nil
}
