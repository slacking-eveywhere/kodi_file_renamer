package main

import (
	"flag"
	"fmt"
	"os"

	"kodi-renamer/internal/scanner"
)

func main() {
	directory := flag.String("dir", ".", "Directory to scan")
	flag.Parse()

	fmt.Println("=== Kodi Renamer - Scanner Test ===")
	fmt.Printf("Scanning: %s\n\n", *directory)

	s := scanner.NewScanner(*directory)
	files, err := s.ScanDirectory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No media files found.")
		return
	}

	movies := 0
	series := 0

	fmt.Printf("Found %d media file(s):\n\n", len(files))

	for i, file := range files {
		fmt.Printf("[%d] %s\n", i+1, file.Name)
		fmt.Printf("    Type: ")
		if file.IsMovie {
			fmt.Printf("Movie\n")
			movies++
			if file.Year > 0 {
				fmt.Printf("    Year: %d\n", file.Year)
			}
		} else if file.IsSeries {
			fmt.Printf("TV Series\n")
			series++
			fmt.Printf("    Season: %d, Episode: %d\n", file.Season, file.Episode)
		}
		fmt.Printf("    Parsed Name: '%s'\n", file.CleanName)
		fmt.Printf("    Search Query: '%s'\n", file.GetSearchQuery())
		fmt.Println()
	}

	fmt.Printf("Summary: %d movies, %d TV series episodes\n", movies, series)
}
