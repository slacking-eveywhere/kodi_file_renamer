package api

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"kodi-renamer/internal/tmdb"
	"kodi-renamer/internal/tvdb"
)

// Manager orchestrates multiple API clients (TVDB and TMDB) for media search
type Manager struct {
	tvdbClient *tvdb.Client
	tmdbClient *tmdb.Client
	hasTVDB    bool
	hasTMDB    bool
}

// UnifiedProposition represents a search result from any API source
type UnifiedProposition struct {
	ID           string
	Name         string
	Title        string
	OriginalName string
	Overview     string
	Year         string
	Type         string
	Source       string
	Score        float64
}

func (p *UnifiedProposition) GetYearAsInt() int {
	year, err := strconv.Atoi(p.Year)
	if err != nil {
		fmt.Printf("Invalid year to format %s", p.Year)
	}
	return year
}

// UnifiedMovieProposition represents detailed movie information from any API source
type UnifiedMovieProposition struct {
	ID       string
	Title    string
	Overview string
	Year     string
	Runtime  int
	Genres   []string
	Source   string
}

// UnifiedSeriesProposition represents detailed TV series information from any API source
type UnifiedSeriesProposition struct {
	ID         string
	Name       string
	Overview   string
	Year       string
	FirstAired string
	Status     string
	Genres     []string
	Source     string
}

// GetFolderName returns the properly formatted folder name for the series
func (s *UnifiedSeriesProposition) GetFolderName() string {
	cleanName := strings.ReplaceAll(s.Name, ":", " -")
	cleanName = strings.ReplaceAll(cleanName, "/", " ")
	if s.Year != "" {
		return fmt.Sprintf("%s (%s)", cleanName, s.Year)
	}
	return cleanName
}

// UnifiedEpisodeInfo represents episode details from any API source
type UnifiedEpisodeInfo struct {
	SeriesName    string
	SeasonNumber  int
	EpisodeNumber int
	EpisodeName   string
	Name          string // Episode name (alias for EpisodeName for consistency)
	Year          string
	Source        string
}

// NewManager creates a new API manager with the provided API keys
func NewManager(tvdbAPIKey, tmdbAPIKey string) *Manager {
	m := &Manager{}

	if tvdbAPIKey != "" {
		m.tvdbClient = tvdb.NewClient(tvdbAPIKey)
		m.hasTVDB = true
		// Login to TVDB if configured
		if err := m.tvdbClient.Login(); err != nil {
			fmt.Printf("Warning: failed to authenticate with TVDB: %v\n", err)
			m.hasTVDB = false
		}
	}

	if tmdbAPIKey != "" {
		m.tmdbClient = tmdb.NewClient(tmdbAPIKey)
		m.hasTMDB = true
	}

	return m
}

// GetConfiguredAPIs returns a list of API names that are currently configured
func (m *Manager) GetConfiguredAPIs() []string {
	apis := []string{}
	if m.hasTVDB {
		apis = append(apis, "TVDB")
	}
	if m.hasTMDB {
		apis = append(apis, "TMDB")
	}
	return apis
}

// Search performs a general search across all configured APIs
func (m *Manager) Search(query string) ([]UnifiedProposition, error) {
	var allProps []UnifiedProposition

	// Search TVDB
	if m.hasTVDB {
		tvdbResults, err := m.tvdbClient.Search(query)
		if err != nil {
			// Don't fail completely, just log and continue
			fmt.Printf("TVDB search warning: %v\n", err)
		} else {
			for _, prop := range tvdbResults {
				allProps = append(allProps, UnifiedProposition{
					ID:           prop.ID,
					Title:        prop.Title,
					OriginalName: prop.OriginalName,
					Overview:     prop.Overview,
					Year:         prop.Year,
					Type:         prop.Type,
					Source:       "tvdb",
					Score:        0.0,
				})
			}
		}
	}

	// Search TMDB
	if m.hasTMDB {
		tmdbResults, err := m.tmdbClient.Search(query)
		if err != nil {
			fmt.Printf("TMDB search warning: %v\n", err)
		} else {
			for _, prop := range tmdbResults {
				allProps = append(allProps, UnifiedProposition{
					ID:           strconv.Itoa(prop.ID),
					Title:        prop.Title,
					OriginalName: prop.OriginalName,
					Overview:     prop.Overview,
					Year:         prop.Year,
					Type:         mapTMDBType(prop.Type),
					Source:       "tmdb",
					Score:        prop.Popularity + (prop.VoteAverage * 10),
				})
			}
		}
	}

	// Sort by score (TMDB has popularity, TVDB results go first)
	sort.SliceStable(allProps, func(i, j int) bool {
		if allProps[i].Source == "tvdb" && allProps[j].Source == "tmdb" {
			return true
		}
		if allProps[i].Source == "tmdb" && allProps[j].Source == "tvdb" {
			return false
		}
		return allProps[i].Score > allProps[j].Score
	})

	return allProps, nil
}

// SearchMovies searches specifically for movies across all configured APIs
func (m *Manager) SearchMovies(query string, year int) ([]UnifiedProposition, error) {
	var allProps []UnifiedProposition

	// Search TVDB
	if m.hasTVDB {
		tvdbResults, err := m.tvdbClient.Search(query)
		if err != nil {
			fmt.Printf("TVDB search warning: %v\n", err)
		} else {
			for _, prop := range tvdbResults {
				if prop.Type == "movie" {
					allProps = append(allProps, UnifiedProposition{
						ID:           prop.ID,
						Title:        prop.Title,
						OriginalName: prop.OriginalName,
						Overview:     prop.Overview,
						Year:         prop.Year,
						Type:         "movie",
						Source:       "tvdb",
						Score:        0.0,
					})
				}
			}
		}
	}

	// Search TMDB
	if m.hasTMDB {
		tmdbResults, err := m.tmdbClient.SearchMovie(query, year)
		if err != nil {
			fmt.Printf("TMDB search warning: %v\n", err)
		} else {
			for _, prop := range tmdbResults {
				allProps = append(allProps, UnifiedProposition{
					ID:           strconv.Itoa(prop.ID),
					Title:        prop.Title,
					OriginalName: prop.OriginalName,
					Overview:     prop.Overview,
					Year:         prop.Year,
					Type:         "movie",
					Source:       "tmdb",
					Score:        prop.Popularity + (prop.VoteAverage * 10),
				})
			}
		}
	}

	sort.SliceStable(allProps, func(i, j int) bool {
		if allProps[i].Source == "tvdb" && allProps[j].Source == "tmdb" {
			return true
		}
		if allProps[i].Source == "tmdb" && allProps[j].Source == "tvdb" {
			return false
		}
		firstYear := allProps[i].GetYearAsInt() - year
		nextYear := allProps[j].GetYearAsInt() - year
		return math.Abs(float64(firstYear)) < math.Abs(float64(nextYear))
		// return allProps[i].Score > allProps[j].Score
	})

	return allProps, nil
}

// SearchSeries searches specifically for TV series across all configured APIs
func (m *Manager) SearchSeries(query string) ([]UnifiedProposition, error) {
	var allProps []UnifiedProposition

	// Search TVDB
	if m.hasTVDB {
		tvdbResults, err := m.tvdbClient.Search(query)
		if err != nil {
			fmt.Printf("TVDB search warning: %v\n", err)
		} else {
			for _, prop := range tvdbResults {
				if prop.Type == "series" {
					allProps = append(allProps, UnifiedProposition{
						ID:           prop.ID,
						Name:         prop.Title,
						Title:        prop.Title,
						OriginalName: prop.OriginalName,
						Overview:     prop.Overview,
						Year:         prop.Year,
						Type:         "series",
						Source:       "tvdb",
						Score:        0.0,
					})
				}
			}
		}
	}

	// Search TMDB
	if m.hasTMDB {
		tmdbResults, err := m.tmdbClient.SearchTV(query, 0)
		if err != nil {
			fmt.Printf("TMDB search warning: %v\n", err)
		} else {
			for _, prop := range tmdbResults {
				allProps = append(allProps, UnifiedProposition{
					ID:           strconv.Itoa(prop.ID),
					Name:         prop.Title,
					Title:        prop.Title,
					OriginalName: prop.OriginalName,
					Overview:     prop.Overview,
					Year:         prop.Year,
					Type:         "series",
					Source:       "tmdb",
					Score:        prop.Popularity + (prop.VoteAverage * 10),
				})
			}
		}
	}

	sort.SliceStable(allProps, func(i, j int) bool {
		if allProps[i].Source == "tvdb" && allProps[j].Source == "tmdb" {
			return true
		}
		if allProps[i].Source == "tmdb" && allProps[j].Source == "tvdb" {
			return false
		}
		return allProps[i].Score > allProps[j].Score
	})

	return allProps, nil
}

// GetMovie retrieves detailed movie information by ID from the specified API source
func (m *Manager) GetMovie(id, source string) (*UnifiedMovieProposition, error) {
	switch source {
	case "tvdb":
		if !m.hasTVDB {
			return nil, fmt.Errorf("TVDB not configured")
		}
		movie, err := m.tvdbClient.GetMovie(id)
		if err != nil {
			return nil, err
		}
		return &UnifiedMovieProposition{
			ID:       strconv.FormatInt(movie.ID, 10),
			Title:    movie.Title,
			Overview: movie.Overview,
			Year:     movie.Year,
			Runtime:  movie.Runtime,
			Genres:   movie.Genres,
			Source:   "tvdb",
		}, nil

	case "tmdb":
		if !m.hasTMDB {
			return nil, fmt.Errorf("TMDB not configured")
		}
		movieID, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("invalid TMDB ID: %w", err)
		}
		movie, err := m.tmdbClient.GetMovie(movieID)
		if err != nil {
			return nil, err
		}
		return &UnifiedMovieProposition{
			ID:       strconv.Itoa(movie.ID),
			Title:    movie.Title,
			Overview: movie.Overview,
			Year:     movie.Year,
			Runtime:  movie.Runtime,
			Genres:   movie.Genres,
			Source:   "tmdb",
		}, nil

	default:
		return nil, fmt.Errorf("unknown source: %s", source)
	}
}

// GetSeries retrieves detailed TV series information by ID from the specified API source
func (m *Manager) GetSeries(id, source string) (*UnifiedSeriesProposition, error) {
	switch source {
	case "tvdb":
		if !m.hasTVDB {
			return nil, fmt.Errorf("TVDB not configured")
		}
		series, err := m.tvdbClient.GetSeries(id)
		if err != nil {
			return nil, err
		}
		return &UnifiedSeriesProposition{
			ID:         strconv.FormatInt(series.ID, 10),
			Name:       series.Name,
			Overview:   series.Overview,
			Year:       series.Year,
			FirstAired: series.FirstAired,
			Status:     series.Status,
			Genres:     series.Genres,
			Source:     "tvdb",
		}, nil

	case "tmdb":
		if !m.hasTMDB {
			return nil, fmt.Errorf("TMDB not configured")
		}
		seriesID, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("invalid TMDB ID: %w", err)
		}
		series, err := m.tmdbClient.GetTVShow(seriesID)
		if err != nil {
			return nil, err
		}
		return &UnifiedSeriesProposition{
			ID:         strconv.Itoa(series.ID),
			Name:       series.Name,
			Overview:   series.Overview,
			Year:       series.Year,
			FirstAired: series.FirstAired,
			Status:     series.Status,
			Genres:     series.Genres,
			Source:     "tmdb",
		}, nil

	default:
		return nil, fmt.Errorf("unknown source: %s", source)
	}
}

// GetEpisode retrieves specific episode information by series ID, season, and episode number from the specified API source
func (m *Manager) GetEpisode(id, source string, season, episode int) (*UnifiedEpisodeInfo, error) {
	switch source {
	case "tvdb":
		if !m.hasTVDB {
			return nil, fmt.Errorf("TVDB not configured")
		}
		episodes, err := m.tvdbClient.GetEpisodes(id, season)
		if err != nil {
			return nil, err
		}

		// Get series name
		series, err := m.tvdbClient.GetSeries(id)
		if err != nil {
			return nil, err
		}

		for _, ep := range episodes {
			if ep.SeasonNumber == season && ep.EpisodeNumber == episode {
				return &UnifiedEpisodeInfo{
					SeriesName:    series.Name,
					SeasonNumber:  season,
					EpisodeNumber: episode,
					EpisodeName:   ep.Name,
					Name:          ep.Name,
					Year:          ep.Year,
					Source:        "tvdb",
				}, nil
			}
		}
		return nil, fmt.Errorf("episode S%02dE%02d not found", season, episode)

	case "tmdb":
		if !m.hasTMDB {
			return nil, fmt.Errorf("TMDB not configured")
		}
		seriesID, err := strconv.Atoi(id)
		if err != nil {
			return nil, fmt.Errorf("invalid TMDB ID: %w", err)
		}
		episodeInfo, err := m.tmdbClient.GetEpisode(seriesID, season, episode)
		if err != nil {
			return nil, err
		}
		return &UnifiedEpisodeInfo{
			SeriesName:    episodeInfo.SeriesName,
			SeasonNumber:  episodeInfo.SeasonNumber,
			EpisodeNumber: episodeInfo.EpisodeNumber,
			EpisodeName:   episodeInfo.EpisodeName,
			Name:          episodeInfo.EpisodeName,
			Year:          episodeInfo.Year,
			Source:        "tmdb",
		}, nil

	default:
		return nil, fmt.Errorf("unknown source: %s", source)
	}
}

// mapTMDBType converts TMDB media types to unified type names
func mapTMDBType(tmdbType string) string {
	switch tmdbType {
	case "tv":
		return "series"
	case "movie":
		return "movie"
	default:
		return tmdbType
	}
}
