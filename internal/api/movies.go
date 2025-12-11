package api

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/sebastianneubert/tmdb/internal/models"
)

func (c *Client) GetTopRatedMovies(page int, language string) (*models.DiscoverResponse, error) {
	params := url.Values{}
	params.Set("page", strconv.Itoa(page))
	params.Set("language", language)

	req, err := c.createRequest("/movie/top_rated", params)
	if err != nil {
		return nil, err
	}

	var response models.DiscoverResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetWatchProviders(movieID int, region string) (models.RegionProviders, error) {
	apiPath := fmt.Sprintf("/movie/%d/watch/providers", movieID)
	req, err := c.createRequest(apiPath, url.Values{})
	if err != nil {
		return models.RegionProviders{}, err
	}

	var response models.WatchProviderResponse
	if err := c.doRequest(req, &response); err != nil {
		return models.RegionProviders{}, err
	}

	if providers, ok := response.Results[region]; ok {
		return providers, nil
	}

	return models.RegionProviders{}, fmt.Errorf("no provider data for region %s", region)
}

func (c *Client) GetExternalIDs(movieID int) (models.ExternalIDs, error) {
	apiPath := fmt.Sprintf("/movie/%d/external_ids", movieID)
	req, err := c.createRequest(apiPath, url.Values{})
	if err != nil {
		return models.ExternalIDs{}, err
	}

	var response models.ExternalIDs
	if err := c.doRequest(req, &response); err != nil {
		return models.ExternalIDs{}, err
	}

	return response, nil
}

func (c *Client) GetEnglishTitle(movieID int) (string, error) {
	apiPath := fmt.Sprintf("/movie/%d", movieID)
	params := url.Values{}
	params.Set("language", "en-US")

	req, err := c.createRequest(apiPath, params)
	if err != nil {
		return "", err
	}

	var response models.MovieDetails
	if err := c.doRequest(req, &response); err != nil {
		return "", err
	}

	return response.Title, nil
}

func (c *Client) SearchMovie(query string, language string, region string) (*models.DiscoverResponse, error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("language", language)
	params.Set("include_adult", "true")
	params.Set("region", region)

  // initialize the response with the first page
  finalResponse := &models.DiscoverResponse{
    Page:         0,
    Results:      []models.Movie{},
    TotalPages:   0,
    TotalResults: 0,
  }

  req, err := c.createRequest("/search/movie", params)
  if err != nil {
    return nil, err
  }

  var response models.DiscoverResponse
  if err := c.doRequest(req, &response); err != nil {
    return nil, err
  }

  finalResponse.Results = append(finalResponse.Results, response.Results...)
  finalResponse.TotalPages = response.TotalPages
  finalResponse.TotalResults = response.TotalResults

  for page := 2; page <= finalResponse.TotalPages; page++ {
    if page > 20 {
      break
    }

    params.Set("page", strconv.Itoa(page))

    req, err := c.createRequest("/search/movie", params)
    if err != nil {
      return nil, err
    }

    var response models.DiscoverResponse
    if err := c.doRequest(req, &response); err != nil {
      return nil, err
    }

    finalResponse.Results = append(finalResponse.Results, response.Results...)
  }

  return finalResponse, nil
}

func (c *Client) GetGenres(language string) (*models.GenreListResponse, error) {
	params := url.Values{}
	params.Set("language", language)

	req, err := c.createRequest("/genre/movie/list", params)
	if err != nil {
		return nil, err
	}

	var response models.GenreListResponse
	if err := c.doRequest(req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) GetMovieDetails(movieID int, language string) (*models.Movie, error) {
	apiPath := fmt.Sprintf("/movie/%d", movieID)
	params := url.Values{}
	params.Set("language", language)

	req, err := c.createRequest(apiPath, params)
	if err != nil {
		return nil, err
	}

	var movie models.Movie
	if err := c.doRequest(req, &movie); err != nil {
		return nil, err
	}

	return &movie, nil
}