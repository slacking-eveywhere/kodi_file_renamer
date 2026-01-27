package tmdb

type SearchResponse struct {
	Page         int            `json:"page"`
	Results      []SearchResult `json:"results"`
	TotalPages   int            `json:"total_pages"`
	TotalResults int            `json:"total_results"`
}

type SearchResult struct {
	ID           int      `json:"id"`
	MediaType    string   `json:"media_type"`
	Title        string   `json:"title"`
	Name         string   `json:"name"`
	OriginalTitle string  `json:"original_title"`
	OriginalName string   `json:"original_name"`
	Overview     string   `json:"overview"`
	ReleaseDate  string   `json:"release_date"`
	FirstAirDate string   `json:"first_air_date"`
	PosterPath   string   `json:"poster_path"`
	BackdropPath string   `json:"backdrop_path"`
	Popularity   float64  `json:"popularity"`
	VoteAverage  float64  `json:"vote_average"`
	GenreIDs     []int    `json:"genre_ids"`
}

type MovieDetails struct {
	ID               int             `json:"id"`
	Title            string          `json:"title"`
	OriginalTitle    string          `json:"original_title"`
	Overview         string          `json:"overview"`
	ReleaseDate      string          `json:"release_date"`
	Runtime          int             `json:"runtime"`
	Status           string          `json:"status"`
	Tagline          string          `json:"tagline"`
	Genres           []Genre         `json:"genres"`
	PosterPath       string          `json:"poster_path"`
	BackdropPath     string          `json:"backdrop_path"`
	Budget           int64           `json:"budget"`
	Revenue          int64           `json:"revenue"`
	Popularity       float64         `json:"popularity"`
	VoteAverage      float64         `json:"vote_average"`
	VoteCount        int             `json:"vote_count"`
	ImdbID           string          `json:"imdb_id"`
	OriginalLanguage string          `json:"original_language"`
}

type TVShowDetails struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
	Overview         string   `json:"overview"`
	FirstAirDate     string   `json:"first_air_date"`
	LastAirDate      string   `json:"last_air_date"`
	Status           string   `json:"status"`
	Type             string   `json:"type"`
	Genres           []Genre  `json:"genres"`
	PosterPath       string   `json:"poster_path"`
	BackdropPath     string   `json:"backdrop_path"`
	Popularity       float64  `json:"popularity"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	NumberOfSeasons  int      `json:"number_of_seasons"`
	NumberOfEpisodes int      `json:"number_of_episodes"`
	Seasons          []Season `json:"seasons"`
}

type Season struct {
	ID           int    `json:"id"`
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

type SeasonDetails struct {
	ID           int       `json:"id"`
	AirDate      string    `json:"air_date"`
	Name         string    `json:"name"`
	Overview     string    `json:"overview"`
	PosterPath   string    `json:"poster_path"`
	SeasonNumber int       `json:"season_number"`
	Episodes     []Episode `json:"episodes"`
}

type Episode struct {
	ID             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float64 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
	Runtime        int     `json:"runtime"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Proposition struct {
	ID           int
	Title        string
	OriginalName string
	Overview     string
	Year         string
	Type         string
	Source       string
	Popularity   float64
	VoteAverage  float64
}

type MovieProposition struct {
	ID       int
	Title    string
	Overview string
	Year     string
	Runtime  int
	Genres   []string
	Source   string
}

type SeriesProposition struct {
	ID         int
	Name       string
	Overview   string
	Year       string
	FirstAired string
	Status     string
	Genres     []string
	Source     string
}

type EpisodeInfo struct {
	SeriesName    string
	SeasonNumber  int
	EpisodeNumber int
	EpisodeName   string
	Year          string
}
