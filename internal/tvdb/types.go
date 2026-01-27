package tvdb

import "time"

type SearchResult struct {
	Data []SearchItem `json:"data"`
}

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

type SeriesResponse struct {
	Data SeriesData `json:"data"`
}

type SeriesData struct {
	ID              int64    `json:"id"`
	Name            string   `json:"name"`
	Slug            string   `json:"slug"`
	Overview        string   `json:"overview"`
	FirstAired      string   `json:"firstAired"`
	LastAired       string   `json:"lastAired"`
	Year            string   `json:"year"`
	Status          Status   `json:"status"`
	ArtworkURL      string   `json:"image"`
	OriginalNetwork string   `json:"originalNetwork"`
	Genres          []Genre  `json:"genres"`
}

type Status struct {
	Name string `json:"name"`
}

type Genre struct {
	Name string `json:"name"`
}

type EpisodesResponse struct {
	Data []Episode `json:"data"`
}

type Episode struct {
	ID             int64  `json:"id"`
	SeriesID       int64  `json:"seriesId"`
	Name           string `json:"name"`
	Aired          string `json:"aired"`
	Runtime        int    `json:"runtime"`
	SeasonNumber   int    `json:"seasonNumber"`
	EpisodeNumber  int    `json:"number"`
	Overview       string `json:"overview"`
	Image          string `json:"image"`
	IsMovie        int    `json:"isMovie"`
	Year           string `json:"year"`
}

type MovieResponse struct {
	Data MovieData `json:"data"`
}

type MovieData struct {
	ID              int64       `json:"id"`
	Name            string      `json:"name"`
	Slug            string      `json:"slug"`
	Overview        string      `json:"overview"`
	Year            string      `json:"year"`
	Runtime         int         `json:"runtime"`
	Status          Status      `json:"status"`
	Genres          []Genre     `json:"genres"`
	Translations    Translation `json:"nameTranslations"`
	Image           string      `json:"image"`
}

type Translation struct {
	Eng string `json:"eng"`
	Fra string `json:"fra"`
}

type TokenResponse struct {
	Data Token `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

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

type Proposition struct {
	ID           string
	Title        string
	OriginalName string
	Overview     string
	Year         string
	Type         string
	ImageURL     string
}

type SeriesProposition struct {
	ID           int64
	Name         string
	Overview     string
	Year         string
	FirstAired   string
	Status       string
	Genres       []string
}

type MovieProposition struct {
	ID       int64
	Title    string
	Overview string
	Year     string
	Runtime  int
	Genres   []string
}

type EpisodeInfo struct {
	SeriesName    string
	SeasonNumber  int
	EpisodeNumber int
	EpisodeName   string
	Year          string
}
