package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"kodi-renamer/internal/api"
	"kodi-renamer/internal/renamer"
	"kodi-renamer/internal/scanner"
	"kodi-renamer/internal/ui"
)

var (
	tvdbAPIKey       string
	tmdbAPIKey       string
	movieToRenameDir string
	movieRenamedDir  string
	serieToRenameDir string
	serieRenamedDir  string
	dryRun           bool
	autoMode         bool
)

// init initializes command-line flags for the application
func init() {
	flag.StringVar(&tvdbAPIKey, "tvdb-key", "", "TVDB API Key")
	flag.StringVar(&tmdbAPIKey, "tmdb-key", "", "TMDB API Key")
	flag.StringVar(&movieToRenameDir, "movie-to-rename", "", "Directory containing movies to rename")
	flag.StringVar(&movieRenamedDir, "movie-renamed", "", "Directory for renamed movies")
	flag.StringVar(&serieToRenameDir, "serie-to-rename", "", "Directory containing series to rename")
	flag.StringVar(&serieRenamedDir, "serie-renamed", "", "Directory for renamed series")
	flag.BoolVar(&dryRun, "dry-run", false, "Dry run mode - don't actually rename files")
	flag.BoolVar(&autoMode, "auto", false, "Automatic mode - select first match")
}

// main is the entry point of the application
func main() {
	flag.Parse()

	// Check environment variables if flags not provided
	if tvdbAPIKey == "" {
		tvdbAPIKey = os.Getenv("TVDB_API_KEY")
	}
	if tmdbAPIKey == "" {
		tmdbAPIKey = os.Getenv("TMDB_API_KEY")
	}
	if movieToRenameDir == "" {
		movieToRenameDir = os.Getenv("MOVIE_TO_RENAME_DIR")
	}
	if movieRenamedDir == "" {
		movieRenamedDir = os.Getenv("MOVIE_RENAMED_DIR")
	}
	if serieToRenameDir == "" {
		serieToRenameDir = os.Getenv("SERIE_TO_RENAME_DIR")
	}
	if serieRenamedDir == "" {
		serieRenamedDir = os.Getenv("SERIE_RENAMED_DIR")
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

	// Validate directory configuration
	if movieToRenameDir == "" && serieToRenameDir == "" {
		fmt.Fprintf(os.Stderr, "Error: At least one input directory is required\n\n")
		fmt.Fprintf(os.Stderr, "Provide via flags:\n")
		fmt.Fprintf(os.Stderr, "  -movie-to-rename 'path/to/movies'\n")
		fmt.Fprintf(os.Stderr, "  -serie-to-rename 'path/to/series'\n\n")
		fmt.Fprintf(os.Stderr, "Or via environment variables:\n")
		fmt.Fprintf(os.Stderr, "  export MOVIE_TO_RENAME_DIR='path/to/movies'\n")
		fmt.Fprintf(os.Stderr, "  export SERIE_TO_RENAME_DIR='path/to/series'\n\n")
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run executes the main application logic for scanning and renaming media files
func run() error {
	interactive := ui.NewInteractive()
	fileRenamer := renamer.NewRenamer(dryRun)

	// Initialize API manager
	apiManager, err := api.NewManager(tvdbAPIKey, tmdbAPIKey)
	if err != nil {
		return fmt.Errorf("failed to initialize API manager: %w", err)
	}

	configuredAPIs := apiManager.GetConfiguredAPIs()
	interactive.PrintSuccess(fmt.Sprintf("Authenticated with: %v", configuredAPIs))

	// Process series if directory is configured
	if serieToRenameDir != "" {
		interactive.PrintInfo(fmt.Sprintf("Scanning series directory: %s", serieToRenameDir))

		fileScanner := scanner.NewScanner(serieToRenameDir)
		mediaFiles, err := fileScanner.ScanDirectory()
		if err != nil {
			return fmt.Errorf("failed to scan series directory: %w", err)
		}

		// Filter for series only
		var seriesFiles []scanner.MediaFile
		for _, file := range mediaFiles {
			if file.IsSeries {
				seriesFiles = append(seriesFiles, file)
			}
		}

		if len(seriesFiles) > 0 {
			interactive.PrintInfo(fmt.Sprintf("Found %d series episode(s)", len(seriesFiles)))

			// Group series episodes by folder for batch processing
			seriesByFolder := scanner.GroupSeriesByFolder(seriesFiles)

			// Process series in batches by folder
			processedFolders := make(map[string]bool)
			for folderPath := range seriesByFolder {
				if processedFolders[folderPath] {
					continue
				}
				if !autoMode || !dryRun {
					fmt.Print("\033[H\033[2J")
				}
				interactive.PrintInfo(fmt.Sprintf("Processing series folder: %s", folderPath))

				if err := processSeriesBatch(seriesByFolder[folderPath], apiManager, interactive, fileRenamer, serieRenamedDir); err != nil {
					interactive.PrintError(fmt.Sprintf("Failed to process series folder: %v", err))
				}
				processedFolders[folderPath] = true
			}
		} else {
			interactive.PrintWarning("No series files found")
		}
	}

	// Process movies if directory is configured
	if movieToRenameDir != "" {
		interactive.PrintInfo(fmt.Sprintf("\nScanning movie directory: %s", movieToRenameDir))

		fileScanner := scanner.NewScanner(movieToRenameDir)
		mediaFiles, err := fileScanner.ScanDirectory()
		if err != nil {
			return fmt.Errorf("failed to scan movie directory: %w", err)
		}

		// Filter for movies only
		var movieFiles []scanner.MediaFile
		for _, file := range mediaFiles {
			if file.IsMovie {
				movieFiles = append(movieFiles, file)
			}
		}

		if len(movieFiles) > 0 {
			interactive.PrintInfo(fmt.Sprintf("Found %d movie(s)", len(movieFiles)))

			// Process movies individually
			movieCount := 0
			for _, file := range movieFiles {
				movieCount++
				if !autoMode || !dryRun {
					fmt.Print("\033[H\033[2J")
				}
				fmt.Printf("\n[Movie %d/%d] Processing: %s\n", movieCount, len(movieFiles), file.Name)

				if err := processMovie(&file, apiManager, interactive, fileRenamer, movieRenamedDir); err != nil {
					interactive.PrintError(fmt.Sprintf("Failed to process movie: %v", err))
					continue
				}
			}
		} else {
			interactive.PrintWarning("No movie files found")
		}
	}

	interactive.PrintSuccess("Processing complete!")
	return nil
}

// processSeriesBatch handles batch processing of all episodes in a series folder
func processSeriesBatch(episodes []*scanner.MediaFile, apiManager *api.Manager, interactive *ui.Interactive, fileRenamer *renamer.Renamer, outputDir string) error {
	if len(episodes) == 0 {
		return nil
	}

	// Use first episode to search for series
	firstEpisode := episodes[0]
	searchQuery := firstEpisode.GetSearchQuery()
	fmt.Printf("Searching for series: '%s'\n", searchQuery)

	propositions, err := apiManager.SearchSeries(searchQuery)
	if err != nil {
		return fmt.Errorf("search failed: %w", err)
	}

	if len(propositions) == 0 {
		interactive.PrintWarning("No results found")
		return nil
	}

	var selectedIndex int
	var seriesDetails *api.UnifiedSeriesProposition

	if autoMode {
		selectedIndex = 0
		interactive.PrintInfo(fmt.Sprintf("Auto-selecting first result from %s", propositions[0].Source))
		seriesDetails, err = apiManager.GetSeries(propositions[0].ID, propositions[0].Source)
		if err != nil {
			return fmt.Errorf("failed to get series details: %w", err)
		}
	} else {
		// Fetch details for all series to display in table
		fmt.Println("Fetching series details...")
		seriesOptions := make([]ui.SeriesOption, 0, len(propositions))

		for _, prop := range propositions {
			details, err := apiManager.GetSeries(prop.ID, prop.Source)
			if err != nil {
				seriesOptions = append(seriesOptions, ui.SeriesOption{
					Name:   prop.Title,
					Year:   prop.Year,
					Status: "",
					Genres: []string{},
					Source: prop.Source,
				})
				continue
			}

			seriesOptions = append(seriesOptions, ui.SeriesOption{
				Name:   details.Name,
				Year:   details.Year,
				Status: details.Status,
				Genres: details.Genres,
				Source: details.Source,
			})
		}

		selectedIndex, err = interactive.SelectSeriesFromList(
			fmt.Sprintf("Select TV series for '%s'", firstEpisode.CleanName),
			seriesOptions,
		)
		if err != nil {
			return err
		}

		if selectedIndex == -1 {
			interactive.PrintInfo("Skipped")
			return nil
		}

		seriesDetails, err = apiManager.GetSeries(propositions[selectedIndex].ID, propositions[selectedIndex].Source)
		if err != nil {
			return fmt.Errorf("failed to get series details: %w", err)
		}
	}

	interactive.DisplaySeriesInfo(seriesDetails.Name, seriesDetails.Year, seriesDetails.Status, seriesDetails.Genres)

	// Fetch episode details for all episodes
	fmt.Println("\nFetching episode details...")

	episodeList := make([]ui.EpisodeDisplay, 0, len(episodes))
	hasErrors := false

	for _, ep := range episodes {
		episodeInfo, err := apiManager.GetEpisode(propositions[selectedIndex].ID, propositions[selectedIndex].Source, ep.Season, ep.Episode)

		display := ui.EpisodeDisplay{
			Season:      ep.Season,
			Episode:     ep.Episode,
			CurrentName: ep.Name,
		}

		if err != nil {
			display.HasError = true
			display.ErrorMessage = "NOT FOUND"
			display.NewName = ""
			display.EpisodeName = "Unknown"
			hasErrors = true
		} else {
			display.EpisodeName = episodeInfo.EpisodeName
			display.NewName = ep.GetNewFilename(seriesDetails.Name, episodeInfo.EpisodeName)
			display.HasError = false
		}

		episodeList = append(episodeList, display)
	}

	// Display the table
	interactive.DisplayEpisodeRenameTable(seriesDetails.Name, episodeList)

	// Check for errors
	if hasErrors {
		interactive.PrintError("Cannot proceed with renaming due to unknown episodes")
		return fmt.Errorf("some episodes could not be found in the database")
	}

	// Generate folder name
	newFolderName := firstEpisode.GetSeriesFolderName(seriesDetails.Name, seriesDetails.Year)
	needsFolderRename := firstEpisode.NeedsSeriesFolderRename(newFolderName)

	if needsFolderRename {
		interactive.PrintInfo(fmt.Sprintf("\nSeries folder will be renamed to: %s", newFolderName))
	}

	// Ask for confirmation
	if !autoMode && !dryRun {
		if !interactive.Confirm(fmt.Sprintf("\nProceed with renaming %d episode(s)?", len(episodeList))) {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	// Perform renames
	if dryRun {
		if needsFolderRename {
			interactive.PrintSuccess(fmt.Sprintf("Would rename folder to: %s", newFolderName))
		}
		for _, ep := range episodeList {
			if !ep.HasError {
				interactive.PrintSuccess(fmt.Sprintf("Would rename: %s -> %s", ep.CurrentName, ep.NewName))
			}
		}
		return nil
	}

	// Rename or move folder
	var newFolderPath string
	if outputDir != "" {
		// Move to output directory
		newFolderPath = filepath.Join(outputDir, newFolderName)
		if err := fileRenamer.MoveSeriesFolder(firstEpisode.GetParentDirPath(), newFolderPath); err != nil {
			return fmt.Errorf("failed to move folder: %w", err)
		}
	} else if needsFolderRename {
		// Just rename in place
		newFolderPath, err = fileRenamer.RenameSeriesFolder(firstEpisode.GetParentDirPath(), newFolderName)
		if err != nil {
			return fmt.Errorf("failed to rename folder: %w", err)
		}
	} else {
		newFolderPath = firstEpisode.GetParentDirPath()
	}

	// Rename episodes silently
	successCount := 0
	for i, ep := range episodeList {
		if ep.HasError {
			continue
		}

		episodeInfo, _ := apiManager.GetEpisode(propositions[selectedIndex].ID, propositions[selectedIndex].Source, ep.Season, ep.Episode)
		newFilename := episodes[i].GetNewFilename(seriesDetails.Name, episodeInfo.EpisodeName)

		// Update path if folder was renamed
		filePath := episodes[i].Path
		if needsFolderRename {
			filePath = filepath.Join(newFolderPath, episodes[i].Name)
		}

		if err := fileRenamer.RenameFileSilent(filePath, newFilename); err != nil {
			interactive.PrintError(fmt.Sprintf("Failed to rename %s: %v", ep.CurrentName, err))
			continue
		}
		successCount++
	}

	interactive.PrintSuccess(fmt.Sprintf("Successfully renamed %d/%d episodes", successCount, len(episodeList)))
	return nil
}

// processMovie handles the search, selection, and renaming workflow for movies
func processMovie(file *scanner.MediaFile, apiManager *api.Manager, interactive *ui.Interactive, fileRenamer *renamer.Renamer, outputDir string) error {
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
	var movieDetails *api.UnifiedMovieProposition

	if autoMode {
		selectedIndex = 0
		interactive.PrintInfo(fmt.Sprintf("Auto-selecting first result from %s", propositions[0].Source))
		movieDetails, err = apiManager.GetMovie(propositions[0].ID, propositions[0].Source)
		if err != nil {
			return fmt.Errorf("failed to get movie details: %w", err)
		}
	} else {
		// Fetch details for all movies to display in table
		fmt.Println("Fetching movie details...")
		movieOptions := make([]ui.MovieOption, 0, len(propositions))

		for _, prop := range propositions {
			details, err := apiManager.GetMovie(prop.ID, prop.Source)
			if err != nil {
				// If we can't get details, create a basic option
				movieOptions = append(movieOptions, ui.MovieOption{
					Title:   prop.Title,
					Year:    prop.Year,
					Runtime: 0,
					Genres:  []string{},
					Source:  prop.Source,
				})
				continue
			}

			movieOptions = append(movieOptions, ui.MovieOption{
				Title:   details.Title,
				Year:    details.Year,
				Runtime: details.Runtime,
				Genres:  details.Genres,
				Source:  details.Source,
			})
		}

		selectedIndex, err = interactive.SelectMovieFromList(
			fmt.Sprintf("Select movie for '%s'", file.CleanName),
			movieOptions,
		)
		if err != nil {
			return err
		}

		if selectedIndex == -1 {
			interactive.PrintInfo("Skipped")
			return nil
		}

		// Get the full details for the selected movie (already fetched above)
		movieDetails, err = apiManager.GetMovie(propositions[selectedIndex].ID, propositions[selectedIndex].Source)
		if err != nil {
			return fmt.Errorf("failed to get movie details: %w", err)
		}
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

	// Move to output directory if specified, otherwise rename in place
	if outputDir != "" {
		newPath := filepath.Join(outputDir, newFilename)
		return fileRenamer.MoveFile(file.Path, newPath)
	}
	return fileRenamer.RenameFile(file.Path, newFilename)
}
