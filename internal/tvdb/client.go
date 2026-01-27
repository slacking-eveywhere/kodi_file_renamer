package tvdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"log"
)

const (
	BaseURL = "https://api4.thetvdb.com/v4"
)

type Client struct {
	apiKey     string
	token      string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Login() error {
	loginData := map[string]string{
		"apikey": c.apiKey,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login data: %w", err)
	}

	req, err := http.NewRequest("POST", BaseURL+"/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Login failed with status %d: %s", resp.StatusCode, string(body))
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	c.token = tokenResp.Data.Token
	return nil
}

func (c *Client) Search(query string) ([]Proposition, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated, call Login() first")
	}

	url := fmt.Sprintf("%s/search?query=%s", BaseURL, query)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResult SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	propositions := make([]Proposition, 0, len(searchResult.Data))
	for _, item := range searchResult.Data {
		propositions = append(propositions, Proposition{
			ID:           item.ID,
			Title:        item.Name,
			OriginalName: item.Name,
			Overview:     item.Overview,
			Year:         item.Year,
			Type:         item.Type,
			ImageURL:     item.ImageURL,
		})
	}

	return propositions, nil
}

func (c *Client) GetSeries(seriesID string) (*SeriesProposition, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated, call Login() first")
	}

	url := fmt.Sprintf("%s/series/%s", BaseURL, seriesID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create series request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute series request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get series failed with status %d: %s", resp.StatusCode, string(body))
	}

	var seriesResp SeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&seriesResp); err != nil {
		return nil, fmt.Errorf("failed to decode series response: %w", err)
	}

	genres := make([]string, 0, len(seriesResp.Data.Genres))
	for _, g := range seriesResp.Data.Genres {
		genres = append(genres, g.Name)
	}

	return &SeriesProposition{
		ID:         seriesResp.Data.ID,
		Name:       seriesResp.Data.Name,
		Overview:   seriesResp.Data.Overview,
		Year:       seriesResp.Data.Year,
		FirstAired: seriesResp.Data.FirstAired,
		Status:     seriesResp.Data.Status.Name,
		Genres:     genres,
	}, nil
}

func (c *Client) GetEpisodes(seriesID string, season int) ([]Episode, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated, call Login() first")
	}

	url := fmt.Sprintf("%s/series/%s/episodes/default?season=%d", BaseURL, seriesID, season)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create episodes request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute episodes request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get episodes failed with status %d: %s", resp.StatusCode, string(body))
	}

	var episodesResp EpisodesResponse
	if err := json.NewDecoder(resp.Body).Decode(&episodesResp); err != nil {
		return nil, fmt.Errorf("failed to decode episodes response: %w", err)
	}

	return episodesResp.Data, nil
}

func (c *Client) GetMovie(movieID string) (*MovieProposition, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated, call Login() first")
	}

	url := fmt.Sprintf("%s/movies/%s", BaseURL, movieID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute movie request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get movie failed with status %d: %s", resp.StatusCode, string(body))
	}

	var movieResp MovieResponse
	if err := json.NewDecoder(resp.Body).Decode(&movieResp); err != nil {
		return nil, fmt.Errorf("failed to decode movie response: %w", err)
	}

	genres := make([]string, 0, len(movieResp.Data.Genres))
	for _, g := range movieResp.Data.Genres {
		genres = append(genres, g.Name)
	}

	return &MovieProposition{
		ID:       movieResp.Data.ID,
		Title:    movieResp.Data.Name,
		Overview: movieResp.Data.Overview,
		Year:     movieResp.Data.Year,
		Runtime:  movieResp.Data.Runtime,
		Genres:   genres,
	}, nil
}
