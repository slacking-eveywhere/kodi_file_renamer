package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// EpisodeDisplay represents episode data for table display
type EpisodeDisplay struct {
	Season       int
	Episode      int
	CurrentName  string
	NewName      string
	EpisodeName  string
	HasError     bool
	ErrorMessage string
}

const (
	TITLE_LIST_COMPILATION_TAB_SMALL = 4
	TITLE_LIST_COMPILATION_TAB_LARGE = 5
)

// MovieOption represents a movie option for selection with detailed information
type MovieOption struct {
	Title   string
	Year    string
	Runtime int
	Genres  []string
	Source  string
}

// SeriesOption represents a TV series option for selection with detailed information
type SeriesOption struct {
	Name   string
	Year   string
	Status string
	Genres []string
	Source string
}

// Interactive provides interactive user interface functionality for user prompts and selections
type Interactive struct {
	reader *bufio.Reader
}

// NewInteractive creates a new Interactive instance for user interaction
func NewInteractive() *Interactive {
	return &Interactive{
		reader: bufio.NewReader(os.Stdin),
	}
}

// SelectMovieFromList displays a table of movie options and prompts the user to select one, returning -1 if skipped
func (i *Interactive) SelectMovieFromList(title string, movies []MovieOption) (int, error) {
	if len(movies) == 0 {
		return -1, fmt.Errorf("no options available")
	}

	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))
	fmt.Println()

	// Calculate column widths
	maxTitle := 20
	maxYear := 4
	maxRuntime := 7
	maxGenres := 30
	maxSource := 6

	for _, movie := range movies {
		if len(movie.Title) > maxTitle {
			maxTitle = len(movie.Title)
		}
		if len(strings.Join(movie.Genres, ", ")) > maxGenres {
			maxGenres = len(strings.Join(movie.Genres, ", "))
		}
	}

	// Print header
	header := fmt.Sprintf("%-3s  %-*s  %-*s  %-*s  %-*s  %-*s",
		"#", maxTitle, "Title", maxYear, "Year", maxRuntime, "Runtime", maxGenres, "Genres", maxSource, "Source")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))

	// Print movie rows
	for idx, movie := range movies {
		runtimeStr := "-"
		if movie.Runtime > 0 {
			runtimeStr = fmt.Sprintf("%d min", movie.Runtime)
		}
		genresStr := "-"
		if len(movie.Genres) > 0 {
			genresStr = strings.Join(movie.Genres, ", ")
			if len(genresStr) > maxGenres {
				genresStr = genresStr[:maxGenres-3] + "..."
			}
		}
		titleStr := movie.Title
		if len(titleStr) > maxTitle {
			titleStr = titleStr[:maxTitle-3] + "..."
		}

		fmt.Printf("%-3d  %-*s  %-*s  %-*s  %-*s  %-*s\n",
			idx+1, maxTitle, titleStr, maxYear, movie.Year, maxRuntime, runtimeStr, maxGenres, genresStr, maxSource, movie.Source)
	}

	// Print skip option
	fmt.Printf("%-3d  Skip / None\n\n", len(movies)+1)

	for {
		fmt.Print("Select an option (number): ")
		input, err := i.reader.ReadString('\n')
		if err != nil {
			return -1, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if choice < 1 || choice > len(movies)+1 {
			fmt.Printf("Invalid choice. Please select between 1 and %d.\n", len(movies)+1)
			continue
		}

		if choice == len(movies)+1 {
			return -1, nil
		}

		return choice - 1, nil
	}
}

// SelectSeriesFromList displays a table of TV series options and prompts the user to select one, returning -1 if skipped
func (i *Interactive) SelectSeriesFromList(title string, series []SeriesOption) (int, error) {
	if len(series) == 0 {
		return -1, fmt.Errorf("no options available")
	}

	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))
	fmt.Println()

	// Calculate column widths
	maxName := 20
	maxYear := 4
	maxStatus := 10
	maxGenres := 30
	maxSource := 6

	for _, s := range series {
		if len(s.Name) > maxName {
			maxName = len(s.Name)
		}
		if len(s.Status) > maxStatus {
			maxStatus = len(s.Status)
		}
		if len(strings.Join(s.Genres, ", ")) > maxGenres {
			maxGenres = len(strings.Join(s.Genres, ", "))
		}
	}

	// Print header
	header := fmt.Sprintf("%-3s  %-*s  %-*s  %-*s  %-*s  %-*s",
		"#", maxName, "Name", maxYear, "Year", maxStatus, "Status", maxGenres, "Genres", maxSource, "Source")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)))

	// Print series rows
	for idx, s := range series {
		statusStr := "-"
		if s.Status != "" {
			statusStr = s.Status
		}
		genresStr := "-"
		if len(s.Genres) > 0 {
			genresStr = strings.Join(s.Genres, ", ")
			if len(genresStr) > maxGenres {
				genresStr = genresStr[:maxGenres-3] + "..."
			}
		}
		nameStr := s.Name
		if len(nameStr) > maxName {
			nameStr = nameStr[:maxName-3] + "..."
		}

		fmt.Printf("%-3d  %-*s  %-*s  %-*s  %-*s  %-*s\n",
			idx+1, maxName, nameStr, maxYear, s.Year, maxStatus, statusStr, maxGenres, genresStr, maxSource, s.Source)
	}

	// Print skip option
	fmt.Printf("%-3d  Skip / None\n\n", len(series)+1)

	for {
		fmt.Print("Select an option (number): ")
		input, err := i.reader.ReadString('\n')
		if err != nil {
			return -1, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if choice < 1 || choice > len(series)+1 {
			fmt.Printf("Invalid choice. Please select between 1 and %d.\n", len(series)+1)
			continue
		}

		if choice == len(series)+1 {
			return -1, nil
		}

		return choice - 1, nil
	}
}

// SelectFromList displays a list of options and prompts the user to select one, returning -1 if skipped
func (i *Interactive) SelectFromList(title string, options []string) (int, error) {
	if len(options) == 0 {
		return -1, fmt.Errorf("no options available")
	}

	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))

	for idx, option := range options {
		var tab string
		if idx+1 < 10 {
			tab = strings.Repeat(" ", TITLE_LIST_COMPILATION_TAB_LARGE)
		} else {
			tab = strings.Repeat(" ", TITLE_LIST_COMPILATION_TAB_SMALL)
		}
		fmt.Printf("%d.%s%s\n", idx+1, tab, option)
	}
	tab := strings.Repeat(" ", TITLE_LIST_COMPILATION_TAB_LARGE)
	fmt.Printf("%d.%sSkip / None\n\n", len(options)+1, tab)

	for {
		fmt.Print("Select an option (number): ")
		input, err := i.reader.ReadString('\n')
		if err != nil {
			return -1, fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if choice < 1 || choice > len(options)+1 {
			fmt.Printf("Invalid choice. Please select between 1 and %d.\n", len(options)+1)
			continue
		}

		if choice == len(options)+1 {
			return -1, nil
		}

		return choice - 1, nil
	}
}

// Confirm prompts the user with a yes/no question and returns true if confirmed
func (i *Interactive) Confirm(message string) bool {
	fmt.Printf("%s (y/n): ", message)
	input, err := i.reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

// DisplayProposition displays a formatted search result proposition with metadata
func (i *Interactive) DisplayProposition(index int, title, overview, year, propType string, genres []string) {
	fmt.Printf("\n--- Option %d ---\n", index+1)
	fmt.Printf("Title:    %s\n", title)
	fmt.Printf("Year:     %s\n", year)
	fmt.Printf("Type:     %s\n", propType)
	if len(genres) > 0 {
		fmt.Printf("Genres:   %s\n", strings.Join(genres, ", "))
	}
	if overview != "" {
		maxLen := 150
		if len(overview) > maxLen {
			overview = overview[:maxLen] + "..."
		}
		fmt.Printf("Overview: %s\n", overview)
	}
}

// DisplaySeriesInfo displays detailed information about a TV series
func (i *Interactive) DisplaySeriesInfo(name, year, status string, genres []string) {
	fmt.Printf("\nSeries: %s (%s)\n", name, year)
	fmt.Printf("Status: %s\n", status)
	if len(genres) > 0 {
		fmt.Printf("Genres: %s\n", strings.Join(genres, ", "))
	}
}

// DisplayEpisodeInfo displays information about a specific episode
func (i *Interactive) DisplayEpisodeInfo(seasonNum, episodeNum int, episodeName string) {
	fmt.Printf("Episode: S%02dE%02d - %s\n", seasonNum, episodeNum, episodeName)
}

// DisplayMovieInfo displays detailed information about a movie
func (i *Interactive) DisplayMovieInfo(title, year string, runtime int, genres []string) {
	fmt.Printf("\nMovie: %s (%s)\n", title, year)
	if runtime > 0 {
		fmt.Printf("Runtime: %d minutes\n", runtime)
	}
	if len(genres) > 0 {
		fmt.Printf("Genres: %s\n", strings.Join(genres, ", "))
	}
}

// PrintError prints an error message to stderr
func (i *Interactive) PrintError(message string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
}

// PrintWarning prints a warning message to stderr
func (i *Interactive) PrintWarning(message string) {
	fmt.Fprintf(os.Stderr, "WARNING: %s\n", message)
}

// PrintInfo prints an informational message to stdout
func (i *Interactive) PrintInfo(message string) {
	fmt.Println(message)
}

// PrintSuccess prints a success message with a checkmark to stdout
func (i *Interactive) PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}

// DisplayEpisodeRenameTable displays a table of pending episode renames with warnings for unknown episodes
func (i *Interactive) DisplayEpisodeRenameTable(seriesName string, episodes []EpisodeDisplay) {
	fmt.Printf("\n=== Pending Renames for: %s ===\n\n", seriesName)

	// Calculate column widths
	maxCurrent := 30
	maxNew := 30
	maxEpisodeName := 40

	for _, ep := range episodes {
		if len(ep.CurrentName) > maxCurrent {
			maxCurrent = len(ep.CurrentName)
		}
		if len(ep.NewName) > maxNew {
			maxNew = len(ep.NewName)
		}
		if len(ep.EpisodeName) > maxEpisodeName {
			maxEpisodeName = len(ep.EpisodeName)
		}
	}

	// Print header
	header := fmt.Sprintf("%-8s  %-*s  %-*s  %-*s  %s",
		"Episode", maxCurrent, "Current Name", maxNew, "New Name", maxEpisodeName, "Episode Title", "Status")
	fmt.Println(header)
	fmt.Println(strings.Repeat("-", len(header)+10))

	// Print episode rows
	hasErrors := false
	for _, ep := range episodes {
		currentStr := ep.CurrentName
		if len(currentStr) > maxCurrent {
			currentStr = currentStr[:maxCurrent-3] + "..."
		}
		newStr := ep.NewName
		if len(newStr) > maxNew {
			newStr = newStr[:maxNew-3] + "..."
		}
		episodeStr := ep.EpisodeName
		if len(episodeStr) > maxEpisodeName {
			episodeStr = episodeStr[:maxEpisodeName-3] + "..."
		}

		status := "✓ OK"
		if ep.HasError {
			status = "✗ " + ep.ErrorMessage
			hasErrors = true
		}

		fmt.Printf("S%02dE%02d     %-*s  %-*s  %-*s  %s\n",
			ep.Season, ep.Episode, maxCurrent, currentStr, maxNew, newStr, maxEpisodeName, episodeStr, status)
	}

	fmt.Println()

	if hasErrors {
		fmt.Println("⚠ WARNING: Some episodes have errors and cannot be renamed!")
		fmt.Println("   Please resolve these issues before proceeding.")
		fmt.Println()
	}
}

// HasEpisodeErrors checks if any episodes in the list have errors
func (i *Interactive) HasEpisodeErrors(episodes []EpisodeDisplay) bool {
	for _, ep := range episodes {
		if ep.HasError {
			return true
		}
	}
	return false
}
