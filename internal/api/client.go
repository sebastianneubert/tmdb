package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.themoviedb.org/3"

type Client struct {
	apiKey     string
	httpClient *http.Client
	timeout    time.Duration
}

func NewClient(apiKey string, timeout int) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("TMDB_API_KEY is required")
	}

	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: time.Duration(timeout) * time.Second},
		timeout:    time.Duration(timeout) * time.Second,
	}, nil
}

func (c *Client) createRequest(apiPath string, params url.Values) (*http.Request, error) {
	params.Set("api_key", c.apiKey)
	fullURL := fmt.Sprintf("%s%s?%s", baseURL, apiPath, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	return req, nil
}

func (c *Client) doRequest(req *http.Request, target interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}