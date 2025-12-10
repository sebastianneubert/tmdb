package api

import (
	"fmt"
	"net/url"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func (c *Client) SearchActor(name string, language string) (*models.ActorSearchResponse, error) {
	params := url.Values{}
	params.Set("query", name)
	params.Set("language", language)

	req, err := c.createRequest("/search/person", params)
	if err != nil {
		return nil, err
	}

	var response models.ActorSearchResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetActorCredits(actorID int, language string) (*models.ActorCreditsResponse, error) {
	apiPath := fmt.Sprintf("/person/%d/movie_credits", actorID)
	params := url.Values{}
	params.Set("language", language)

	req, err := c.createRequest(apiPath, params)
	if err != nil {
		return nil, err
	}

	var response models.ActorCreditsResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}