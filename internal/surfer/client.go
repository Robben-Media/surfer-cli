package surfer

import (
	"context"
	"errors"
	"fmt"

	"github.com/builtbyrobben/surfer-cli/internal/api"
)

var (
	errKeywordsRequired = errors.New("keywords are required")
	errIDRequired       = errors.New("id is required")
	errURLRequired      = errors.New("url is required")
)

const defaultBaseURL = "https://app.surferseo.com/api/v1"

// Client wraps the API client with Surfer-specific methods.
type Client struct {
	*api.Client
}

// NewClient creates a new Surfer API client.
func NewClient(apiKey string) *Client {
	return &Client{
		Client: api.NewClient(apiKey,
			api.WithBaseURL(defaultBaseURL),
			api.WithUserAgent("surfer-cli/1.0"),
		),
	}
}

// ContentEditor represents a Surfer content editor.
type ContentEditor struct {
	ID       string   `json:"id"`
	Keywords []string `json:"keywords,omitempty"`
	Language string   `json:"language,omitempty"`
	Location string   `json:"location,omitempty"`
	Device   string   `json:"device,omitempty"`
	State    string   `json:"state,omitempty"`
	URL      string   `json:"url,omitempty"`
}

// ContentEditorScore represents the content score response.
type ContentEditorScore struct {
	ContentScore int `json:"content_score"`
}

// CreateContentEditorRequest is the request body for creating a content editor.
type CreateContentEditorRequest struct {
	Keywords []string `json:"keywords"`
	Location string   `json:"location,omitempty"`
	Language string   `json:"language,omitempty"`
	Device   string   `json:"device,omitempty"`
}

// ContentEditorsListResponse is the paginated response for listing content editors.
type ContentEditorsListResponse struct {
	Data       []ContentEditor `json:"data"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalItems int             `json:"total_items"`
}

// Audit represents a Surfer audit.
type Audit struct {
	ID    string `json:"id"`
	URL   string `json:"url,omitempty"`
	State string `json:"state,omitempty"`
	Score int    `json:"score,omitempty"`
}

// CreateAuditRequest is the request body for creating an audit.
type CreateAuditRequest struct {
	URL string `json:"url"`
}

// AuditsListResponse is the paginated response for listing audits.
type AuditsListResponse struct {
	Data       []Audit `json:"data"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalItems int     `json:"total_items"`
}

// ContentEditors provides methods for the Content Editor API.
func (c *Client) ContentEditors() *ContentEditorsService {
	return &ContentEditorsService{client: c}
}

// Audits provides methods for the Audit API.
func (c *Client) Audits() *AuditsService {
	return &AuditsService{client: c}
}

// ContentEditorsService handles content editor operations.
type ContentEditorsService struct {
	client *Client
}

// List returns all content editors (paginated).
func (s *ContentEditorsService) List(ctx context.Context, page, pageSize int) (*ContentEditorsListResponse, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	path := fmt.Sprintf("/content_editors?page=%d&page_size=%d", page, pageSize)

	var result ContentEditorsListResponse
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list content editors: %w", err)
	}

	return &result, nil
}

// Create creates a new content editor.
func (s *ContentEditorsService) Create(ctx context.Context, req CreateContentEditorRequest) (*ContentEditor, error) {
	if len(req.Keywords) == 0 {
		return nil, errKeywordsRequired
	}

	// Set defaults
	if req.Device == "" {
		req.Device = "desktop"
	}

	if req.Language == "" {
		req.Language = "en"
	}

	var result ContentEditor
	if err := s.client.Post(ctx, "/content_editors", req, &result); err != nil {
		return nil, fmt.Errorf("create content editor: %w", err)
	}

	return &result, nil
}

// Get returns a content editor by ID.
func (s *ContentEditorsService) Get(ctx context.Context, id string) (*ContentEditor, error) {
	if id == "" {
		return nil, errIDRequired
	}

	var result ContentEditor

	path := fmt.Sprintf("/content_editors/%s", id)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get content editor: %w", err)
	}

	return &result, nil
}

// Score returns the content score for a content editor.
func (s *ContentEditorsService) Score(ctx context.Context, id string) (*ContentEditorScore, error) {
	if id == "" {
		return nil, errIDRequired
	}

	var result ContentEditorScore

	path := fmt.Sprintf("/content_editors/%s/content_score", id)
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("get content score: %w", err)
	}

	return &result, nil
}

// AuditsService handles audit operations.
type AuditsService struct {
	client *Client
}

// List returns all audits (paginated).
func (s *AuditsService) List(ctx context.Context, page, pageSize int) (*AuditsListResponse, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	path := fmt.Sprintf("/audits?page=%d&page_size=%d", page, pageSize)

	var result AuditsListResponse
	if err := s.client.Get(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list audits: %w", err)
	}

	return &result, nil
}

// Create creates a new audit for a URL.
func (s *AuditsService) Create(ctx context.Context, url string) (*Audit, error) {
	if url == "" {
		return nil, errURLRequired
	}

	req := CreateAuditRequest{URL: url}

	var result Audit
	if err := s.client.Post(ctx, "/audits", req, &result); err != nil {
		return nil, fmt.Errorf("create audit: %w", err)
	}

	return &result, nil
}
