package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL = "https://api.hubapi.com"
)

// Client represents a HubSpot API client
type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewClient creates a new HubSpot API client
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey:  apiKey,
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Contact represents a HubSpot contact
type Contact struct {
	ID         string                 `json:"id"`
	Properties map[string]interface{} `json:"properties"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  string                 `json:"updatedAt"`
}

// ContactResponse represents the response from HubSpot API
type ContactResponse struct {
	Results []Contact `json:"results"`
	Paging  *Paging   `json:"paging,omitempty"`
}

// Paging represents pagination information
type Paging struct {
	Next *NextPage `json:"next,omitempty"`
}

// NextPage represents next page information
type NextPage struct {
	After string `json:"after"`
}

// Property represents a contact property
type Property struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	FieldType   string `json:"fieldType"`
	Description string `json:"description"`
	Options     []struct {
		Label string `json:"label"`
		Value string `json:"value"`
	} `json:"options,omitempty"`
}

// PropertiesResponse represents the response for properties
type PropertiesResponse struct {
	Results []Property `json:"results"`
}

// doRequest performs an HTTP request to the HubSpot API
func (c *Client) doRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ListContacts retrieves all contacts with pagination
func (c *Client) ListContacts(limit int, after string) (*ContactResponse, error) {
	endpoint := "/crm/v3/objects/contacts"
	params := url.Values{}
	params.Add("limit", fmt.Sprintf("%d", limit))
	if after != "" {
		params.Add("after", after)
	}
	params.Add("properties", "email,firstname,lastname,hs_lead_status,lifecyclestage")

	if len(params) > 0 {
		endpoint += "?" + params.Encode()
	}

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var contactResp ContactResponse
	if err := json.Unmarshal(respBody, &contactResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &contactResp, nil
}

// GetContact retrieves a specific contact by ID
func (c *Client) GetContact(contactID string) (*Contact, error) {
	endpoint := fmt.Sprintf("/crm/v3/objects/contacts/%s", contactID)
	params := url.Values{}
	params.Add("properties", "email,firstname,lastname,hs_lead_status,lifecyclestage")

	endpoint += "?" + params.Encode()

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var contact Contact
	if err := json.Unmarshal(respBody, &contact); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &contact, nil
}

// CreateContact creates a new contact
func (c *Client) CreateContact(properties map[string]interface{}) (*Contact, error) {
	endpoint := "/crm/v3/objects/contacts"

	requestBody := map[string]interface{}{
		"properties": properties,
	}

	respBody, err := c.doRequest("POST", endpoint, requestBody)
	if err != nil {
		return nil, err
	}

	var contact Contact
	if err := json.Unmarshal(respBody, &contact); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &contact, nil
}

// UpdateContact updates an existing contact
func (c *Client) UpdateContact(contactID string, properties map[string]interface{}) (*Contact, error) {
	endpoint := fmt.Sprintf("/crm/v3/objects/contacts/%s", contactID)

	requestBody := map[string]interface{}{
		"properties": properties,
	}

	respBody, err := c.doRequest("PATCH", endpoint, requestBody)
	if err != nil {
		return nil, err
	}

	var contact Contact
	if err := json.Unmarshal(respBody, &contact); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &contact, nil
}

// DeleteContact deletes a contact by ID
func (c *Client) DeleteContact(contactID string) error {
	endpoint := fmt.Sprintf("/crm/v3/objects/contacts/%s", contactID)
	_, err := c.doRequest("DELETE", endpoint, nil)
	return err
}

// SearchContacts searches for contacts using HubSpot's search API
// The query parameter can be a property name and value in format "property=value"
// or just a value to search in email, firstname, and lastname fields
func (c *Client) SearchContacts(query string, limit int) (*ContactResponse, error) {
	endpoint := "/crm/v3/objects/contacts/search"

	// Parse query - if it contains "=", treat as property=value, otherwise search in common fields
	var filterGroups []map[string]interface{}

	if strings.Contains(query, "=") {
		// Property-based search
		parts := strings.SplitN(query, "=", 2)
		if len(parts) == 2 {
			property := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			filterGroups = []map[string]interface{}{
				{
					"filters": []map[string]interface{}{
						{
							"propertyName": property,
							"operator":     "EQ",
							"value":        value,
						},
					},
				},
			}
		}
	} else {
		// Search in email field (most common use case)
		value := strings.TrimSpace(query)
		filterGroups = []map[string]interface{}{
			{
				"filters": []map[string]interface{}{
					{
						"propertyName": "email",
						"operator":     "CONTAINS_TOKEN",
						"value":        value,
					},
				},
			},
		}
	}

	requestBody := map[string]interface{}{
		"filterGroups": filterGroups,
		"limit":        limit,
		"properties":   []string{"email", "firstname", "lastname", "hs_lead_status", "lifecyclestage"},
	}

	respBody, err := c.doRequest("POST", endpoint, requestBody)
	if err != nil {
		return nil, err
	}

	var contactResp ContactResponse
	if err := json.Unmarshal(respBody, &contactResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &contactResp, nil
}

// ListProperties retrieves all contact properties
func (c *Client) ListProperties() ([]Property, error) {
	endpoint := "/crm/v3/properties/contacts"

	respBody, err := c.doRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var propsResp PropertiesResponse
	if err := json.Unmarshal(respBody, &propsResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return propsResp.Results, nil
}
