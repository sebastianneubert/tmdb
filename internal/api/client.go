package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
	"bytes"

	"github.com/sebastianneubert/tmdb/internal/config"
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
    cfg := config.Get()

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    // --- DEBUGGING STEP START ---

    // 1. Read the entire response body into a byte slice.
    bodyBytes, readErr := io.ReadAll(resp.Body)
    if readErr != nil {
        return fmt.Errorf("failed to read response body: %w", readErr)
    }

    // 2. Print the raw JSON response for debugging.
    // Use json.MarshalIndent for a pretty-printed output.
    // NOTE: This uses the "encoding/json" package.

    if cfg.DEBUG {
      var raw map[string]interface{}
      if err := json.Unmarshal(bodyBytes, &raw); err == nil {
          // Pretty-print if it's valid JSON
          prettyJSON, _ := json.MarshalIndent(raw, "", "  ")
          fmt.Println("--- API Response Body (Pretty) ---")
          fmt.Println(string(prettyJSON))
          fmt.Println("----------------------------------")
      } else {
          // If it's not JSON (e.g., HTML error), print it as raw text
          fmt.Println("--- API Response Body (Raw Text) ---")
          fmt.Println(string(bodyBytes))
          fmt.Println("------------------------------------")
      }
    }

    // 3. Create a new io.Reader from the byte slice.
    // This allows the JSON decoder to read the data.
    newBodyReader := bytes.NewReader(bodyBytes)

    // --- DEBUGGING STEP END ---

    if resp.StatusCode != http.StatusOK {
        // Use the read body for the error message
        return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
    }

    // Use the new reader for decoding.
    if err := json.NewDecoder(newBodyReader).Decode(target); err != nil {
        return fmt.Errorf("failed to decode response: %w", err)
    }

    return nil
}