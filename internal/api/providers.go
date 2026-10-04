package api

import (
	"net/url"

	"github.com/sebastianneubert/tmdb/internal/models"
)

// GetWatchProviderList returns all movie streaming providers TMDb knows for a region.
func (c *Client) GetWatchProviderList(region string) (*models.ProviderListResponse, error) {
	params := url.Values{}
	params.Set("watch_region", region)

	req, err := c.createRequest("/watch/providers/movie", params)
	if err != nil {
		return nil, err
	}

	var response models.ProviderListResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
