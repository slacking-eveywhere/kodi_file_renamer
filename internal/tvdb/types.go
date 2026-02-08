package tvdb

import "time"

// SearchResult represents the response from a TheTVDB search API call
type SearchResult struct {
	Data []SearchItem `json:"data"`
}

// SearchItem represents a single search result item from TheTVDB
type SearchItem struct {
	ID           string   `json:"tvdb_id"`
	Name         string   `json:"name"`
	FirstAired   string   `json:"first_air_time"`
	Overview     string   `json:"overview"`
	Type         string   `json:"type"`
	Year         string   `json:"year"`
	ImageURL     string   `json:"image_url"`
	Translations []string `json:"translations"`
}

// SeriesResponse represents the response from a TheTVDB series API call
type SeriesResponse struct {
	Data SeriesData `json:"data"`
}

// SeriesData contains detailed information about a TV series
type SeriesData struct {
	ID              int64   `json:"id"`
	Name            string  `json:"name"`
	Slug            string  `json:"slug"`
	Overview        string  `json:"overview"`
	FirstAired      string  `json:"firstAired"`
	LastAired       string  `json:"lastAired"`
	Year            string  `json:"year"`
	Status          Status  `json:"status"`
	ArtworkURL      string  `json:"image"`
	OriginalNetwork string  `json:"originalNetwork"`
	Genres          []Genre `json:"genres"`
}

// Status represents the current status of a series or movie
type Status struct {
	Name string `json:"name"`
}

// Genre represents a content genre
type Genre struct {
	Name string `json:"name"`
}

// EpisodesResponse represents the response from a TheTVDB episodes API call
type EpisodesResponse struct {
	Data []Episode `json:"data"`
}

// Episode represents a single episode of a TV series
type Episode struct {
	ID            int64  `json:"id"`
	SeriesID      int64  `json:"seriesId"`
	Name          string `json:"name"`
	Aired         string `json:"aired"`
	Runtime       int    `json:"runtime"`
	SeasonNumber  int    `json:"seasonNumber"`
	EpisodeNumber int    `json:"number"`
	Overview      string `json:"overview"`
	Image         string `json:"image"`
	IsMovie       int    `json:"isMovie"`
	Year          string `json:"year"`
}

// MovieResponse represents the response from a TheTVDB movie API call
type MovieResponse struct {
	Data MovieData `json:"data"`
}

// MovieData contains detailed information about a movie
type MovieData struct {
	ID           int64       `json:"id"`
	Name         string      `json:"name"`
	Slug         string      `json:"slug"`
	Overview     string      `json:"overview"`
	Year         string      `json:"year"`
	Runtime      int         `json:"runtime"`
	Status       Status      `json:"status"`
	Genres       []Genre     `json:"genres"`
	Translations Translation `json:"nameTranslations"`
	Image        string      `json:"image"`
}

// Translation contains translated names in different languages
type Translation struct {
	Eng string `json:"eng"`
	Fra string `json:"fra"`
}

// TokenResponse represents the authentication token response from TheTVDB
type TokenResponse struct {
	Data Token `json:"data"`
}

// Token contains the authentication token
type Token struct {
	Token string `json:"token"`
}

// MediaFile represents a media file with its metadata
type MediaFile struct {
	Path      string
	Name      string
	Extension string
	IsMovie   bool
	IsSeries  bool
	Season    int
	Episode   int
	Year      int
	Duration  time.Duration
}

// Proposition represents a simplified search result for user selection
type Proposition struct {
	ID           string
	Title        string
	OriginalName string
	Overview     string
	Year         string
	Type         string
	ImageURL     string
}

// SeriesProposition represents detailed TV series information for user display
type SeriesProposition struct {
	ID         int64
	Name       string
	Overview   string
	Year       string
	FirstAired string
	Status     string
	Genres     []string
}

// MovieProposition represents detailed movie information for user display
type MovieProposition struct {
	ID       int64
	Title    string
	Overview string
	Year     string
	Runtime  int
	Genres   []string
}

// EpisodeInfo contains specific episode details for renaming purposes
type EpisodeInfo struct {
	SeriesName    string
	SeasonNumber  int
	EpisodeNumber int
	EpisodeName   string
	Year          string
}
