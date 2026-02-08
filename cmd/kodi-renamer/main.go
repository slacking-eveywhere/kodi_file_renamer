package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
	interactive      *ui.Interactive
)

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

func main() {
	flag.Parse()

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

	if movieToRenameDir == "" && serieToRenameDir == "" {
		fmt.Fprintf(os.Stderr, "Error: At least one input directory is required\n\n")
		fmt.Fprintf(os.Stderr, "Provide via flags:\n")
		fmt.Fprintf(os.Stderr, "  -movie-to-rename 'path/to/movies'\n")
		fmt.Fprintf(os.Stderr, "  -serie-to-rename 'path/to/series'\n\n")
		os.Exit(1)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	interactive = ui.NewInteractive()
	apiManager := api.NewManager(tvdbAPIKey, tmdbAPIKey)
	fileRenamer := renamer.NewRenamer(dryRun)

	if dryRun {
		interactive.PrintInfo("Running in DRY RUN mode - no changes will be made")
	}

	configuredAPIs := apiManager.GetConfiguredAPIs()
	interactive.PrintSuccess(fmt.Sprintf("Authenticated with: %v", configuredAPIs))

	if movieToRenameDir != "" {
		interactive.PrintHeader("Processing Movies")
		movieScanner := scanner.NewScanner(movieToRenameDir)
		mediaFiles, err := movieScanner.ScanDirectory()
		if err != nil {
			return fmt.Errorf("failed to scan movie directory: %w", err)
		}

		movies := make([]*scanner.MediaFile, 0)
		for i := range mediaFiles {
			if mediaFiles[i].IsMovie {
				movies = append(movies, &mediaFiles[i])
			}
		}

		if len(movies) == 0 {
			interactive.PrintInfo("No movies found in directory")
		} else {
			interactive.PrintInfo(fmt.Sprintf("Found %d movie(s)", len(movies)))

			for _, movie := range movies {
				if err := processMovie(movie, apiManager, interactive, fileRenamer, movieRenamedDir); err != nil {
					interactive.PrintError(fmt.Sprintf("Failed to process %s: %v", movie.Name, err))
					if !autoMode {
						if !interactive.Confirm("Continue with next movie?") {
							break
						}
					}
				}
			}
		}
	}

	if serieToRenameDir != "" {
		interactive.PrintHeader("Processing Series")
		seriesScanner := scanner.NewScanner(serieToRenameDir)
		mediaFiles, err := seriesScanner.ScanDirectory()
		if err != nil {
			return fmt.Errorf("failed to scan series directory: %w", err)
		}

		series := make([]*scanner.MediaFile, 0)
		for i := range mediaFiles {
			if mediaFiles[i].IsSeries {
				series = append(series, &mediaFiles[i])
			}
		}

		if len(series) == 0 {
			interactive.PrintInfo("No series found in directory")
		} else {
			interactive.PrintInfo(fmt.Sprintf("Found %d episode(s)", len(series)))

			seriesMap := make(map[string][]*scanner.MediaFile)
			for _, ep := range series {
				parentDir := ep.ParentDir
				seriesMap[parentDir] = append(seriesMap[parentDir], ep)
			}

			for parentDir, episodes := range seriesMap {
				if err := processSeriesBatch(parentDir, episodes, apiManager, interactive, fileRenamer, serieRenamedDir); err != nil {
					interactive.PrintError(fmt.Sprintf("Failed to process series %s: %v", parentDir, err))
					if !autoMode {
						if !interactive.Confirm("Continue with next series?") {
							break
						}
					}
				}
			}
		}
	}

	interactive.PrintSuccess("Processing complete!")
	return nil
}

func processSeriesBatch(parentDir string, episodes []*scanner.MediaFile, apiManager *api.Manager, interactive *ui.Interactive, fileRenamer *renamer.Renamer, outputDir string) error {
	if len(episodes) == 0 {
		return nil
	}

	interactive.PrintHeader(fmt.Sprintf("Processing Series: %s", parentDir))

	firstEpisode := episodes[0]
	searchQuery := scanner.GetSeriesSearchQuery(parentDir)

	fmt.Printf("Searching for series: '%s' (from folder: %s)\n", searchQuery, parentDir)

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
		seriesDetails, err = apiManager.GetSeries(propositions[selectedIndex].ID, propositions[selectedIndex].Source)
		if err != nil {
			return fmt.Errorf("failed to get series details: %w", err)
		}
	} else {
		fmt.Println("Fetching series details...")
		seriesOptions := make([]ui.SeriesOption, 0, len(propositions))

		for _, prop := range propositions {
			details, err := apiManager.GetSeries(prop.ID, prop.Source)
			if err != nil {
				seriesOptions = append(seriesOptions, ui.SeriesOption{
					Name:   prop.Name,
					Year:   prop.Year,
					Status: "",
					Source: prop.Source,
				})
				continue
			}

			seriesOptions = append(seriesOptions, ui.SeriesOption{
				Name:   details.Name,
				Year:   details.Year,
				Status: details.Status,
				Source: details.Source,
			})
		}

		selectedIndex, err = interactive.SelectSeriesFromList(
			fmt.Sprintf("Select series for '%s'", firstEpisode.CleanName),
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

	interactive.DisplaySeriesInfo(seriesDetails.Name, seriesDetails.Year, seriesDetails.Status)

	batch := &scanner.SeriesBatchRename{
		OriginalFolderPath: filepath.Dir(firstEpisode.Path),
		OriginalFolderName: parentDir,
		NewFolderName:      seriesDetails.GetFolderName(),
		SeriesName:         seriesDetails.Name,
		SeriesYear:         seriesDetails.Year,
		Episodes:           make([]scanner.EpisodeRenameTask, 0),
		NeedsFolderRename:  parentDir != seriesDetails.GetFolderName(),
	}

	fmt.Println("\nFetching episode details...")
	for _, ep := range episodes {
		episodeDetails, err := apiManager.GetEpisode(seriesDetails.ID, seriesDetails.Source, ep.Season, ep.Episode)
		if err != nil {
			batch.Episodes = append(batch.Episodes, scanner.EpisodeRenameTask{
				File:         ep,
				EpisodeName:  "Unknown Episode",
				NewFilename:  ep.GetEpisodeFilename(seriesDetails.Name, ep.Season, ep.Episode, "Unknown Episode"),
				Season:       ep.Season,
				Episode:      ep.Episode,
				HasError:     true,
				ErrorMessage: err.Error(),
			})
			continue
		}

		newFilename := ep.GetEpisodeFilename(seriesDetails.Name, ep.Season, ep.Episode, episodeDetails.Name)
		batch.Episodes = append(batch.Episodes, scanner.EpisodeRenameTask{
			File:        ep,
			EpisodeName: episodeDetails.Name,
			NewFilename: newFilename,
			Season:      ep.Season,
			Episode:     ep.Episode,
			HasError:    false,
		})
	}

	batchInfo := &ui.BatchRenameInfo{
		SeriesName: batch.SeriesName,
		Episodes:   make([]ui.BatchEpisodeTask, 0, len(batch.Episodes)),
	}

	for _, task := range batch.Episodes {
		batchInfo.Episodes = append(batchInfo.Episodes, ui.BatchEpisodeTask{
			CurrentName:  task.File.Name,
			NewFilename:  task.NewFilename,
			EpisodeName:  task.EpisodeName,
			Season:       task.Season,
			Episode:      task.Episode,
			HasError:     task.HasError,
			ErrorMessage: task.ErrorMessage,
		})
	}

	interactive.DisplayEpisodeBatch(batchInfo)

	if !autoMode && !dryRun {
		if !interactive.Confirm("Proceed with renaming all episodes?") {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	if outputDir != "" {
		newFolderPath := filepath.Join(outputDir, batch.NewFolderName)

		if !dryRun {
			if err := os.MkdirAll(newFolderPath, 0755); err != nil {
				return fmt.Errorf("failed to create output folder: %w", err)
			}
		}

		for _, task := range batch.Episodes {
			if task.HasError {
				interactive.PrintWarning(fmt.Sprintf("Skipping S%02dE%02d: %s", task.Season, task.Episode, task.ErrorMessage))
				continue
			}

			oldPath := task.File.Path
			newPath := filepath.Join(newFolderPath, task.NewFilename)

			if dryRun {
				fmt.Printf("[DRY RUN] Would move:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
			} else {
				if err := os.Rename(oldPath, newPath); err != nil {
					interactive.PrintError(fmt.Sprintf("Failed to move %s: %v", task.File.Name, err))
				} else {
					fmt.Printf("Moved:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
				}
			}
		}
	} else {
		if batch.NeedsFolderRename {
			newFolderPath, err := fileRenamer.RenameSeriesFolder(batch.OriginalFolderPath, batch.NewFolderName)
			if err != nil {
				return fmt.Errorf("failed to rename series folder: %w", err)
			}
			batch.OriginalFolderPath = newFolderPath
		}

		for _, task := range batch.Episodes {
			if task.HasError {
				interactive.PrintWarning(fmt.Sprintf("Skipping S%02dE%02d: %s", task.Season, task.Episode, task.ErrorMessage))
				continue
			}

			oldPath := task.File.Path
			if batch.NeedsFolderRename && !dryRun {
				oldPath = filepath.Join(batch.OriginalFolderPath, task.File.Name)
			}

			if err := fileRenamer.RenameFileSilent(oldPath, task.NewFilename); err != nil {
				interactive.PrintError(fmt.Sprintf("Failed to rename %s: %v", task.File.Name, err))
			}
		}

		if !dryRun {
			fmt.Printf("\nSuccessfully renamed %d episode(s) in series '%s'\n\n", len(batch.Episodes), batch.SeriesName)
		}
	}

	return nil
}

func processMovie(file *scanner.MediaFile, apiManager *api.Manager, interactive *ui.Interactive, fileRenamer *renamer.Renamer, outputDir string) error {
	searchQuery := file.GetSearchQuery()
	year := file.Year

	if file.IsMovieFolder {
		if file.IsBluRay {
			fmt.Printf("Processing Blu-ray folder: '%s'\n", file.Name)
		} else if file.IsDVD {
			fmt.Printf("Processing DVD folder: '%s'\n", file.Name)
		} else {
			fmt.Printf("Processing movie folder: '%s' (%d video files, %d subtitles)\n",
				file.Name, len(file.MovieFiles), len(file.SubtitleFiles))
		}
	}

	fmt.Printf("Searching for: '%s (%d)'\n", searchQuery, year)

	propositions, err := apiManager.SearchMovies(searchQuery, year)
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
		fmt.Println("Fetching movie details...")
		movieOptions := make([]ui.MovieOption, 0, len(propositions))

		for _, prop := range propositions {
			details, err := apiManager.GetMovie(prop.ID, prop.Source)
			if err != nil {
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

		movieDetails, err = apiManager.GetMovie(propositions[selectedIndex].ID, propositions[selectedIndex].Source)
		if err != nil {
			return fmt.Errorf("failed to get movie details: %w", err)
		}
	}

	interactive.DisplayMovieInfo(movieDetails.Title, movieDetails.Year, movieDetails.Runtime, movieDetails.Genres)

	if file.IsMovieFolder {
		return processMovieFolder(file, movieDetails, fileRenamer, outputDir)
	}

	return processStandaloneMovie(file, movieDetails, fileRenamer, outputDir)
}

func processMovieFolder(file *scanner.MediaFile, movieDetails *api.UnifiedMovieProposition, fileRenamer *renamer.Renamer, outputDir string) error {
	title := movieDetails.Title
	year := movieDetails.Year

	newFolderName := file.GetMovieFolderName(title, year)

	targetDir := outputDir
	if targetDir == "" {
		targetDir = filepath.Dir(filepath.Dir(file.Path))
	}

	var mainVideoFile string
	if len(file.MovieFiles) > 0 {
		mainVideoFile = filepath.Base(file.MovieFiles[0])
	}

	if !autoMode && !dryRun {
		if !interactive.Confirm(fmt.Sprintf("Move/rename folder to '%s'?", newFolderName)) {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	if err := fileRenamer.MoveRenameMovieFolder(file.Path, targetDir, newFolderName); err != nil {
		return err
	}

	if mainVideoFile != "" {
		newFileName := file.GetMovieFilename(title, year)
		newFolderPath := filepath.Join(targetDir, newFolderName)

		if err := fileRenamer.RenameMovieFileInFolder(newFolderPath, mainVideoFile, newFileName); err != nil {
			interactive.PrintWarning(fmt.Sprintf("Failed to rename video file: %v", err))
		}
	}

	return nil
}

func processStandaloneMovie(file *scanner.MediaFile, movieDetails *api.UnifiedMovieProposition, fileRenamer *renamer.Renamer, outputDir string) error {
	title := movieDetails.Title
	year := movieDetails.Year

	folderName := file.GetMovieFolderName(title, year)
	newFilename := file.GetMovieFilename(title, year)

	targetDir := outputDir
	if targetDir == "" {
		targetDir = filepath.Dir(file.Path)
	}

	if !autoMode && !dryRun {
		if !interactive.Confirm(fmt.Sprintf("Create folder '%s' and move/rename to '%s'?", folderName, newFilename)) {
			interactive.PrintInfo("Skipped")
			return nil
		}
	}

	return fileRenamer.MoveRenameMovieFile(file.Path, targetDir, folderName, newFilename)
}
