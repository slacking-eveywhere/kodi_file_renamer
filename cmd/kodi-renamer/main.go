package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"kodi-renamer/internal/api"
	"kodi-renamer/internal/renamer"
	"kodi-renamer/internal/scanner"
	"kodi-renamer/internal/ui"
)

var (
	tvdbAPIKey string
	tmdbAPIKey string
	directory  string
	dryRun     bool
	autoMode   bool
)

func init() {
	flag.StringVar(&tvdbAPIKey, "tvdb-key", "", "TVDB API Key")
	flag.StringVar(&tmdbAPIKey, "tmdb-key", "", "TMDB API Key")
	flag.StringVar(&directory, "dir", ".", "Directory to scan for media files")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run mode - don't actually rename files")
	flag.BoolVar(&autoMode, "auto", false, "Automatic mode - select first match")
}

func main() {
	flag.Parse()

	// Check environment variables if flags not provided
	if tvdbAPIKey == "" {
		tvdbAPIKey = os.Getenv("TVDB_API_KEY")
	}
	if tmdbAPIKey == "" {
		tmdbAPIKey = os.Getenv("TMDB_API_KEY")
	}

	// Validate at least one API key is provided
	if tvdbAPIKey == "" && tmdbAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: At least one API key is required\n\n")
		fmt.Fprintf(os.Stderr, "Provide via flags:\n")
		fmt.Fprintf(os.Stderr, "  -tvdb-key 'your-tvdb-key'\n")
		fmt.Fprintf(os.Stderr, "  -tmdb-key 'your-tmdb-key'\n\n")
		fmt.Fprintf(os.Stderr, "Or via environment variables:\n")
		fmt.Fprintf(os.Stderr, "  export TVDB_API_KEY='your-tvdb-key'\n")
		fmt.Fprintf(os.Stderr, "  export TMDB_API_KEY='your-tmdb-key'\n\n")
		fmt.Fprintf(os.Stderr, "You can set both for combined results\n")
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	interactive := ui.NewInteractive()
	fileScanner := scanner.NewScanner(directory)
	fileRenamer := renamer.NewRenamer(dryRun)

	interactive.PrintInfo(fmt.Sprintf("Scanning directory: %s", directory))

	// Initialize API manager
	apiManager, err := api.NewManager(tvdbAPIKey, tmdbAPIKey)
	if err != nil {
		return fmt.Errorf("failed to initialize API manager: %w", err)
	}

	configuredAPIs := apiManager.GetConfiguredAPIs()
	interactive.PrintSuccess(fmt.Sprintf("Authenticated with: %v", configuredAPIs))

	mediaFiles, err := fileScanner.ScanDirectory()
	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	if len(mediaFiles) == 0 {
		interactive.PrintWarning("No media files found")
		return nil
	}

	interactive.PrintInfo(fmt.Sprintf("Found %d media file(s)", len(mediaFiles)))

	movies := 0
	series := 0
	for _, file := range mediaFiles {
		if file.IsMovie {
			movies++
		} else if file.IsSeries {
			series++
		}
	}
	interactive.PrintInfo(fmt.Sprintf("  Movies: %d, TV Series: %d\n", movies, series))

	for idx, file := range mediaFiles {
		fmt.Printf("\n[%d/%d] Processing: %s\n", idx+1, len(mediaFiles), file.Name)

		if file.IsSeries {
			if err := processSeries(&file, apiManager, interactive, fileRenamer); err != nil {
				interactive.PrintError(fmt.Sprintf("Failed to process series: %v", err))
				continue
			}
		} else if file.IsMovie {
			if err := processMovie(&file, apiManager, interactive, fileRenamer); err != nil {
				interactive.PrintError(fmt.Sprintf("Failed to process movie: %v", err))
				continue
			}
		}
	}

	interactive.PrintSuccess("Processing complete!")
	return nil
}

func processSeries(file *scanner.MediaFile, apiManager *api.Manager, interactive *ui.Interactive, fileRenamer *renamer.Renamer) error {
	searchQuery := file.GetSearchQuery()
	fmt.Printf("Searching for: '%s'\n", searchQuery)

	propositions, err := apiManager.SearchSeries(searchQuery)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(propositions) == 0 {
		interactive.PrintWarning("No results found")
		return nil
	}

	var selectedIndex int
	if autoMode {
		selectedIndex = 0
		interactive.PrintInfo(fmt.Sprintf("Auto-selecting first result from %s", propositions[0].Source))
	} else {
		options := make([]string, len(propositions))
		for i, prop := range propositions {
			options[i] = fmt.Sprintf("[%s] %s (%s)", prop.Source, prop.Title, prop.Year)
		}

		selectedIndex, err = interactive.SelectFromList(
			fmt.Sprintf("Select TV series for '%s'", file.CleanName),
			options,
		)
		if err != nil {
			return err
		}

		if selectedIndex == -1 {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	selectedProp := propositions[selectedIndex]
	seriesDetails, err := apiManager.GetSeries(selectedProp.ID, selectedProp.Source)
	if err != nil {
		return fmt.Errorf("failed to get series details: %w", err)
	}

	interactive.DisplaySeriesInfo(seriesDetails.Name, seriesDetails.Year, seriesDetails.Status, seriesDetails.Genres)

	episodeInfo, err := apiManager.GetEpisode(selectedProp.ID, selectedProp.Source, file.Season, file.Episode)
	if err != nil {
		return fmt.Errorf("failed to get episode: %w", err)
	}

	interactive.DisplayEpisodeInfo(file.Season, file.Episode, episodeInfo.EpisodeName)

	newFilename := file.GetNewFilename(seriesDetails.Name, episodeInfo.EpisodeName)

	if !autoMode && !dryRun {
		if !interactive.Confirm(fmt.Sprintf("Rename to '%s'?", newFilename)) {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	return fileRenamer.RenameFile(file.Path, newFilename)
}

func processMovie(file *scanner.MediaFile, apiManager *api.Manager, interactive *ui.Interactive, fileRenamer *renamer.Renamer) error {
	searchQuery := file.GetSearchQuery()
	fmt.Printf("Searching for: '%s'\n", searchQuery)

	propositions, err := apiManager.SearchMovies(searchQuery)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(propositions) == 0 {
		interactive.PrintWarning("No results found")
		return nil
	}

	var selectedIndex int
	if autoMode {
		selectedIndex = 0
		interactive.PrintInfo(fmt.Sprintf("Auto-selecting first result from %s", propositions[0].Source))
	} else {
		options := make([]string, len(propositions))
		for i, prop := range propositions {
			options[i] = fmt.Sprintf("[%s] %s (%s)", prop.Source, prop.Title, prop.Year)
		}

		selectedIndex, err = interactive.SelectFromList(
			fmt.Sprintf("Select movie for '%s'", file.CleanName),
			options,
		)
		if err != nil {
			return err
		}

		if selectedIndex == -1 {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	selectedProp := propositions[selectedIndex]
	movieDetails, err := apiManager.GetMovie(selectedProp.ID, selectedProp.Source)
	if err != nil {
		return fmt.Errorf("failed to get movie details: %w", err)
	}

	interactive.DisplayMovieInfo(movieDetails.Title, movieDetails.Year, movieDetails.Runtime, movieDetails.Genres)

	year, _ := strconv.Atoi(movieDetails.Year)
	newFilename := file.GetMovieFilename(movieDetails.Title, year)

	if !autoMode && !dryRun {
		if !interactive.Confirm(fmt.Sprintf("Rename to '%s'?", newFilename)) {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	return fileRenamer.RenameFile(file.Path, newFilename)
}
