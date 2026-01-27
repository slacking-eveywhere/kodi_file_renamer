package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Interactive struct {
	reader *bufio.Reader
}

func NewInteractive() *Interactive {
	return &Interactive{
		reader: bufio.NewReader(os.Stdin),
	}
}

func (i *Interactive) SelectFromList(title string, options []string) (int, error) {
	if len(options) == 0 {
		return -1, fmt.Errorf("no options available")
	}

	fmt.Printf("\n%s\n", title)
	fmt.Println(strings.Repeat("=", len(title)))

	for idx, option := range options {
		fmt.Printf("%d. %s\n", idx+1, option)
	}
	fmt.Printf("%d. Skip / None\n\n", len(options)+1)

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

func (i *Interactive) Confirm(message string) bool {
	fmt.Printf("%s (y/n): ", message)
	input, err := i.reader.ReadString('\n')
	if err != nil {
		return false
	}

	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes"
}

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

func (i *Interactive) DisplaySeriesInfo(name, year, status string, genres []string) {
	fmt.Printf("\nSeries: %s (%s)\n", name, year)
	fmt.Printf("Status: %s\n", status)
	if len(genres) > 0 {
		fmt.Printf("Genres: %s\n", strings.Join(genres, ", "))
	}
}

func (i *Interactive) DisplayEpisodeInfo(seasonNum, episodeNum int, episodeName string) {
	fmt.Printf("Episode: S%02dE%02d - %s\n", seasonNum, episodeNum, episodeName)
}

func (i *Interactive) DisplayMovieInfo(title, year string, runtime int, genres []string) {
	fmt.Printf("\nMovie: %s (%s)\n", title, year)
	if runtime > 0 {
		fmt.Printf("Runtime: %d minutes\n", runtime)
	}
	if len(genres) > 0 {
		fmt.Printf("Genres: %s\n", strings.Join(genres, ", "))
	}
}

func (i *Interactive) PrintError(message string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
}

func (i *Interactive) PrintWarning(message string) {
	fmt.Fprintf(os.Stderr, "WARNING: %s\n", message)
}

func (i *Interactive) PrintInfo(message string) {
	fmt.Println(message)
}

func (i *Interactive) PrintSuccess(message string) {
	fmt.Printf("✓ %s\n", message)
}
