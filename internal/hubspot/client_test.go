package hubspot

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiKey := "test-api-key"
	client := NewClient(apiKey)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.apiKey != apiKey {
		t.Errorf("Expected API key %s, got %s", apiKey, client.apiKey)
	}

	if client.baseURL != baseURL {
		t.Errorf("Expected base URL %s, got %s", baseURL, client.baseURL)
	}

	if client.client == nil {
		t.Fatal("HTTP client is nil")
	}
}

func TestClient_doRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Error("Missing or incorrect Authorization header")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "data"}`))
	}))
	defer server.Close()

	client := NewClient("test-api-key")
	client.baseURL = server.URL

	_, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Errorf("doRequest failed: %v", err)
	}
}

func TestContact_String(t *testing.T) {
	contact := Contact{
		ID: "123",
		Properties: map[string]interface{}{
			"email": "test@example.com",
		},
	}

	if contact.ID != "123" {
		t.Errorf("Expected ID 123, got %s", contact.ID)
	}

	email, ok := contact.Properties["email"].(string)
	if !ok || email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %v", contact.Properties["email"])
	}
}

