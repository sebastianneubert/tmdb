package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func (c *Client) GetTopRatedShows(page int, language string) (*models.ShowDiscoverResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("language", language)

	req, err := c.createRequest("/tv/top_rated", params)
	if err != nil {
		return nil, err
	}

	var response models.ShowDiscoverResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetShowWatchProviders(showID int, region string) (models.RegionProviders, error) {
	apiPath := fmt.Sprintf("/tv/%d/watch/providers", showID)
	req, err := c.createRequest(apiPath, url.Values{})
	if err != nil {
		return models.RegionProviders{}, err
	}

	var response models.ShowWatchProviderResponse
	if err := c.doRequest(req, &response); err != nil {
		return models.RegionProviders{}, err
	}

	if providers, ok := response.Results[region]; ok {
		return providers, nil
	}

	return models.RegionProviders{}, fmt.Errorf("no provider data for region %s", region)
}

func (c *Client) GetShowExternalIDs(showID int) (models.ShowExternalIDs, error) {
	apiPath := fmt.Sprintf("/tv/%d/external_ids", showID)
	req, err := c.createRequest(apiPath, url.Values{})
	if err != nil {
		return models.ShowExternalIDs{}, err
	}

	var response models.ShowExternalIDs
	if err := c.doRequest(req, &response); err != nil {
		return models.ShowExternalIDs{}, err
	}

	return response, nil
}

func (c *Client) GetShowEnglishTitle(showID int) (string, error) {
	apiPath := fmt.Sprintf("/tv/%d", showID)
	params := url.Values{}
	params.Set("language", "en-US")

	req, err := c.createRequest(apiPath, params)
	if err != nil {
		return "", err
	}

	var response models.ShowDetails
	if err := c.doRequest(req, &response); err != nil {
		return "", err
	}

	return response.Name, nil
}