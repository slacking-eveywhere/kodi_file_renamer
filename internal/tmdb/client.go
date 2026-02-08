package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// BaseURL is the base URL for The Movie Database (TMDb) API v3
	BaseURL = "https://api.themoviedb.org/3"
)

// Client represents a TMDb API client
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new TMDb API client with the provided API key
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured checks if the client has a valid API key configured
func (c *Client) IsConfigured() bool {
	return c.apiKey != ""
}

// Search performs a multi-search across movies and TV shows on TMDb
func (c *Client) Search(query string) ([]Proposition, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	encodedQuery := url.QueryEscape(query)
	searchURL := fmt.Sprintf("%s/search/multi?api_key=%s&query=%s&include_adult=false", BaseURL, c.apiKey, encodedQuery)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode search response: %w", err)
	}

	propositions := make([]Proposition, 0)
	for _, result := range searchResp.Results {
		prop := c.resultToProposition(result)
		if prop.Type == "movie" || prop.Type == "tv" {
			propositions = append(propositions, prop)
		}
	}

	return propositions, nil
}

// resultToProposition converts a TMDb search result to a unified Proposition
func (c *Client) resultToProposition(result SearchResult) Proposition {
	prop := Proposition{
		ID:          result.ID,
		Overview:    result.Overview,
		Source:      "tmdb",
		Popularity:  result.Popularity,
		VoteAverage: result.VoteAverage,
	}

	if result.MediaType == "movie" {
		prop.Title = result.Title
		prop.OriginalName = result.OriginalTitle
		prop.Type = "movie"
		if result.ReleaseDate != "" && len(result.ReleaseDate) >= 4 {
			prop.Year = result.ReleaseDate[:4]
		}
	} else if result.MediaType == "tv" {
		prop.Title = result.Name
		prop.OriginalName = result.OriginalName
		prop.Type = "tv"
		if result.FirstAirDate != "" && len(result.FirstAirDate) >= 4 {
			prop.Year = result.FirstAirDate[:4]
		}
	} else {
		prop.Type = result.MediaType
	}

	return prop
}

// GetMovie retrieves detailed information about a movie by TMDb ID
func (c *Client) GetMovie(movieID int) (*MovieProposition, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	movieURL := fmt.Sprintf("%s/movie/%d?api_key=%s", BaseURL, movieID, c.apiKey)

	req, err := http.NewRequest("GET", movieURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute movie request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get movie failed with status %d: %s", resp.StatusCode, string(body))
	}

	var movieDetails MovieDetails
	if err := json.NewDecoder(resp.Body).Decode(&movieDetails); err != nil {
		return nil, fmt.Errorf("failed to decode movie response: %w", err)
	}

	genres := make([]string, 0, len(movieDetails.Genres))
	for _, g := range movieDetails.Genres {
		genres = append(genres, g.Name)
	}

	year := ""
	if movieDetails.ReleaseDate != "" && len(movieDetails.ReleaseDate) >= 4 {
		year = movieDetails.ReleaseDate[:4]
	}

	return &MovieProposition{
		ID:       movieDetails.ID,
		Title:    movieDetails.Title,
		Overview: movieDetails.Overview,
		Year:     year,
		Runtime:  movieDetails.Runtime,
		Genres:   genres,
		Source:   "tmdb",
	}, nil
}

// GetTVShow retrieves detailed information about a TV show by TMDb ID
func (c *Client) GetTVShow(tvID int) (*SeriesProposition, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	tvURL := fmt.Sprintf("%s/tv/%d?api_key=%s", BaseURL, tvID, c.apiKey)

	req, err := http.NewRequest("GET", tvURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create tv request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tv request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get tv show failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tvDetails TVShowDetails
	if err := json.NewDecoder(resp.Body).Decode(&tvDetails); err != nil {
		return nil, fmt.Errorf("failed to decode tv response: %w", err)
	}

	genres := make([]string, 0, len(tvDetails.Genres))
	for _, g := range tvDetails.Genres {
		genres = append(genres, g.Name)
	}

	year := ""
	if tvDetails.FirstAirDate != "" && len(tvDetails.FirstAirDate) >= 4 {
		year = tvDetails.FirstAirDate[:4]
	}

	return &SeriesProposition{
		ID:         tvDetails.ID,
		Name:       tvDetails.Name,
		Overview:   tvDetails.Overview,
		Year:       year,
		FirstAired: tvDetails.FirstAirDate,
		Status:     tvDetails.Status,
		Genres:     genres,
		Source:     "tmdb",
	}, nil
}

// GetEpisode retrieves information about a specific episode by TV show ID, season, and episode number
func (c *Client) GetEpisode(tvID, seasonNumber, episodeNumber int) (*EpisodeInfo, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	// First get the TV show details for the name
	tvShow, err := c.GetTVShow(tvID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tv show: %w", err)
	}

	// Get the season details
	seasonURL := fmt.Sprintf("%s/tv/%d/season/%d?api_key=%s", BaseURL, tvID, seasonNumber, c.apiKey)

	req, err := http.NewRequest("GET", seasonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create season request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute season request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get season failed with status %d: %s", resp.StatusCode, string(body))
	}

	var seasonDetails SeasonDetails
	if err := json.NewDecoder(resp.Body).Decode(&seasonDetails); err != nil {
		return nil, fmt.Errorf("failed to decode season response: %w", err)
	}

	// Find the specific episode
	for _, ep := range seasonDetails.Episodes {
		if ep.EpisodeNumber == episodeNumber {
			year := tvShow.Year
			if ep.AirDate != "" && len(ep.AirDate) >= 4 {
				year = ep.AirDate[:4]
			}

			return &EpisodeInfo{
				SeriesName:    tvShow.Name,
				SeasonNumber:  seasonNumber,
				EpisodeNumber: episodeNumber,
				EpisodeName:   ep.Name,
				Year:          year,
			}, nil
		}
	}

	return nil, fmt.Errorf("episode S%02dE%02d not found", seasonNumber, episodeNumber)
}

// SearchMovie performs a search specifically for movies on TMDb
func (c *Client) SearchMovie(query string) ([]Proposition, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	encodedQuery := url.QueryEscape(query)
	searchURL := fmt.Sprintf("%s/search/movie?api_key=%s&query=%s&include_adult=false", BaseURL, c.apiKey, encodedQuery)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create movie search request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute movie search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("movie search failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode movie search response: %w", err)
	}

	propositions := make([]Proposition, 0, len(searchResp.Results))
	for _, result := range searchResp.Results {
		year := ""
		if result.ReleaseDate != "" && len(result.ReleaseDate) >= 4 {
			year = result.ReleaseDate[:4]
		}

		propositions = append(propositions, Proposition{
			ID:           result.ID,
			Title:        result.Title,
			OriginalName: result.OriginalTitle,
			Overview:     result.Overview,
			Year:         year,
			Type:         "movie",
			Source:       "tmdb",
			Popularity:   result.Popularity,
			VoteAverage:  result.VoteAverage,
		})
	}

	return propositions, nil
}

// SearchTV performs a search specifically for TV shows on TMDb
func (c *Client) SearchTV(query string) ([]Proposition, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("TMDB API key not configured")
	}

	encodedQuery := url.QueryEscape(query)
	searchURL := fmt.Sprintf("%s/search/tv?api_key=%s&query=%s&include_adult=false", BaseURL, c.apiKey, encodedQuery)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create tv search request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute tv search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("tv search failed with status %d: %s", resp.StatusCode, string(body))
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode tv search response: %w", err)
	}

	propositions := make([]Proposition, 0, len(searchResp.Results))
	for _, result := range searchResp.Results {
		year := ""
		if result.FirstAirDate != "" && len(result.FirstAirDate) >= 4 {
			year = result.FirstAirDate[:4]
		}

		propositions = append(propositions, Proposition{
			ID:           result.ID,
			Title:        result.Name,
			OriginalName: result.OriginalName,
			Overview:     result.Overview,
			Year:         year,
			Type:         "tv",
			Source:       "tmdb",
			Popularity:   result.Popularity,
			VoteAverage:  result.VoteAverage,
		})
	}

	return propositions, nil
}

// FormatID converts an integer ID to a string
func FormatID(id int) string {
	return strconv.Itoa(id)
}

// ParseID converts a string ID to an integer
func ParseID(idStr string) (int, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid ID format: %w", err)
	}
	return id, nil
}
