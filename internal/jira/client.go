package jira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kilyinov/cli-example/internal/config"
)

type Client struct {
	baseURL    string
	email      string
	token      string
	httpClient *http.Client
}

func NewClient(cfg config.JIRAConfig) (*Client, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("jira.base_url is not configured")
	}
	if cfg.Email == "" {
		return nil, fmt.Errorf("jira.email is not configured")
	}
	if cfg.Token == "" {
		return nil, fmt.Errorf("jira.token is not configured")
	}
	return &Client{
		baseURL:    cfg.BaseURL,
		email:      cfg.Email,
		token:      cfg.Token,
		httpClient: &http.Client{},
	}, nil
}

type CreateIssueRequest struct {
	Project     string
	Summary     string
	Description string
	IssueType   string
}

type createIssuePayload struct {
	Fields issueFields `json:"fields"`
}

type issueFields struct {
	Project   idField `json:"project"`
	Summary   string  `json:"summary"`
	IssueType idField `json:"issuetype"`

	Description *atlasDoc `json:"description,omitempty"`
}

type idField struct {
	Key  string `json:"key,omitempty"`
	Name string `json:"name,omitempty"`
}

type atlasDoc struct {
	Type    string      `json:"type"`
	Version int         `json:"version"`
	Content []atlasNode `json:"content"`
}

type atlasNode struct {
	Type    string      `json:"type"`
	Content []atlasText `json:"content,omitempty"`
}

type atlasText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type CreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

func (c *Client) CreateIssue(req CreateIssueRequest) (*CreateIssueResponse, error) {
	payload := createIssuePayload{
		Fields: issueFields{
			Project:   idField{Key: req.Project},
			Summary:   req.Summary,
			IssueType: idField{Name: req.IssueType},
		},
	}

	if req.Description != "" {
		payload.Fields.Description = &atlasDoc{
			Type:    "doc",
			Version: 1,
			Content: []atlasNode{
				{
					Type: "paragraph",
					Content: []atlasText{
						{Type: "text", Text: req.Description},
					},
				},
			},
		}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/rest/api/3/issue", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	c.setAuth(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("JIRA API returned %d: %s", resp.StatusCode, string(respBody))
	}

	var result CreateIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}
	return &result, nil
}

func (c *Client) AddAttachment(issueKey string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("creating form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return fmt.Errorf("copying file data: %w", err)
	}
	writer.Close()

	httpReq, err := http.NewRequest("POST", c.baseURL+"/rest/api/3/issue/"+issueKey+"/attachments", &buf)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("X-Atlassian-Token", "no-check")
	c.setAuth(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("JIRA API returned %d: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (c *Client) setAuth(req *http.Request) {
	req.SetBasicAuth(c.email, c.token)
}
