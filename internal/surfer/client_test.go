package surfer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/builtbyrobben/surfer-cli/internal/api"
)

func newTestClient(server *httptest.Server) *Client {
	return &Client{
		Client: api.NewClient("test-key",
			api.WithBaseURL(server.URL),
			api.WithUserAgent("surfer-cli/test"),
		),
	}
}

func TestContentEditors_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/content_editors" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected page=1, got %s", r.URL.Query().Get("page"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentEditorsListResponse{
			Data: []ContentEditor{
				{ID: "ce-1", Keywords: []string{"seo"}, State: "ready"},
				{ID: "ce-2", Keywords: []string{"content"}, State: "processing"},
			},
			Page:       1,
			PageSize:   10,
			TotalItems: 2,
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.ContentEditors().List(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 2 {
		t.Fatalf("expected 2 editors, got %d", len(result.Data))
	}

	if result.Data[0].ID != "ce-1" {
		t.Errorf("expected first editor ID 'ce-1', got %q", result.Data[0].ID)
	}

	if result.TotalItems != 2 {
		t.Errorf("expected TotalItems 2, got %d", result.TotalItems)
	}
}

func TestContentEditors_List_DefaultPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify defaults were applied for invalid values
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected page=1 default, got %s", r.URL.Query().Get("page"))
		}

		if r.URL.Query().Get("page_size") != "10" {
			t.Errorf("expected page_size=10 default, got %s", r.URL.Query().Get("page_size"))
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentEditorsListResponse{Data: []ContentEditor{}})
	}))
	defer server.Close()

	client := newTestClient(server)

	_, err := client.ContentEditors().List(context.Background(), 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContentEditors_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body CreateContentEditorRequest
		json.NewDecoder(r.Body).Decode(&body)

		if len(body.Keywords) != 2 {
			t.Errorf("expected 2 keywords, got %d", len(body.Keywords))
		}

		if body.Device != "desktop" {
			t.Errorf("expected device 'desktop', got %q", body.Device)
		}

		if body.Language != "en" {
			t.Errorf("expected language 'en', got %q", body.Language)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentEditor{
			ID:       "new-ce",
			Keywords: body.Keywords,
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.ContentEditors().Create(context.Background(), CreateContentEditorRequest{
		Keywords: []string{"seo", "content"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "new-ce" {
		t.Errorf("expected ID 'new-ce', got %q", result.ID)
	}
}

func TestContentEditors_Create_NoKeywords(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.ContentEditors().Create(context.Background(), CreateContentEditorRequest{})
	if err == nil {
		t.Fatal("expected error for empty keywords")
	}
}

func TestContentEditors_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/content_editors/ce-123" {
			t.Errorf("expected path /content_editors/ce-123, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentEditor{
			ID:       "ce-123",
			Keywords: []string{"test"},
			State:    "ready",
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.ContentEditors().Get(context.Background(), "ce-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "ce-123" {
		t.Errorf("expected ID 'ce-123', got %q", result.ID)
	}
}

func TestContentEditors_Get_EmptyID(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.ContentEditors().Get(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestContentEditors_Score(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/content_editors/ce-123/content_score" {
			t.Errorf("expected score path, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ContentEditorScore{ContentScore: 85})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.ContentEditors().Score(context.Background(), "ce-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ContentScore != 85 {
		t.Errorf("expected score 85, got %d", result.ContentScore)
	}
}

func TestContentEditors_Score_EmptyID(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.ContentEditors().Score(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty ID")
	}
}

func TestAudits_List(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/audits" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuditsListResponse{
			Data: []Audit{
				{ID: "a-1", URL: "https://example.com", Score: 72},
			},
			Page:       1,
			PageSize:   10,
			TotalItems: 1,
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Audits().List(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Data) != 1 {
		t.Fatalf("expected 1 audit, got %d", len(result.Data))
	}

	if result.Data[0].Score != 72 {
		t.Errorf("expected score 72, got %d", result.Data[0].Score)
	}
}

func TestAudits_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body CreateAuditRequest
		json.NewDecoder(r.Body).Decode(&body)

		if body.URL != "https://example.com" {
			t.Errorf("expected URL 'https://example.com', got %q", body.URL)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Audit{
			ID:    "new-audit",
			URL:   body.URL,
			State: "pending",
		})
	}))
	defer server.Close()

	client := newTestClient(server)

	result, err := client.Audits().Create(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != "new-audit" {
		t.Errorf("expected ID 'new-audit', got %q", result.ID)
	}
}

func TestAudits_Create_EmptyURL(t *testing.T) {
	client := &Client{Client: api.NewClient("key")}

	_, err := client.Audits().Create(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty URL")
	}
}
