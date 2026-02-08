package renamer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Renamer handles file renaming operations with optional dry-run mode
type Renamer struct {
	dryRun bool
}

// NewRenamer creates a new Renamer instance with the specified dry-run mode
func NewRenamer(dryRun bool) *Renamer {
	return &Renamer{
		dryRun: dryRun,
	}
}

// RenameFile renames a file from oldPath to newFilename in the same directory
func (r *Renamer) RenameFile(oldPath, newFilename string) error {
	return r.renameFileWithOutput(oldPath, newFilename, false)
}

// RenameFileSilent renames a file silently without output (for batch operations)
func (r *Renamer) RenameFileSilent(oldPath, newFilename string) error {
	return r.renameFileWithOutput(oldPath, newFilename, true)
}

// renameFileWithOutput is the internal implementation with optional silent mode
func (r *Renamer) renameFileWithOutput(oldPath, newFilename string, silent bool) error {
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newFilename)

	if oldPath == newPath {
		return nil
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("file already exists: %s", newPath)
	}

	if r.dryRun {
		if !silent {
			fmt.Printf("[DRY RUN] Would rename:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
		}
		return nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	if !silent {
		fmt.Printf("Renamed:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
	}
	return nil
}

// SetDryRun enables or disables dry-run mode
func (r *Renamer) SetDryRun(dryRun bool) {
	r.dryRun = dryRun
}

// IsDryRun returns whether dry-run mode is currently enabled
func (r *Renamer) IsDryRun() bool {
	return r.dryRun
}

// RenameSeriesFolder renames a series parent directory (silent in non-dry-run mode)
func (r *Renamer) RenameSeriesFolder(oldDirPath, newFolderName string) (string, error) {
	parentDir := filepath.Dir(oldDirPath)
	newDirPath := filepath.Join(parentDir, newFolderName)

	if oldDirPath == newDirPath {
		return oldDirPath, nil
	}

	if _, err := os.Stat(newDirPath); err == nil {
		return "", fmt.Errorf("directory already exists: %s", newDirPath)
	}

	if r.dryRun {
		return newDirPath, nil
	}

	if err := os.Rename(oldDirPath, newDirPath); err != nil {
		return "", fmt.Errorf("failed to rename folder: %w", err)
	}

	return newDirPath, nil
}

// RenameFileInSeriesFolder renames both the series folder and the file inside it
func (r *Renamer) RenameFileInSeriesFolder(oldFilePath, newFolderName, newFileName string) error {
	oldDirPath := filepath.Dir(oldFilePath)
	oldFileName := filepath.Base(oldFilePath)

	// First, rename the folder
	newDirPath, err := r.RenameSeriesFolder(oldDirPath, newFolderName)
	if err != nil {
		return err
	}

	// Then, rename the file in the new folder location
	oldFileInNewDir := filepath.Join(newDirPath, oldFileName)
	newFilePath := filepath.Join(newDirPath, newFileName)

	if oldFileInNewDir == newFilePath {
		return nil
	}

	if _, err := os.Stat(newFilePath); err == nil {
		return fmt.Errorf("file already exists: %s", newFilePath)
	}

	if r.dryRun {
		fmt.Printf("[DRY RUN] Would rename file:\n  FROM: %s\n  TO:   %s\n\n", oldFileInNewDir, newFilePath)
		return nil
	}

	if err := os.Rename(oldFileInNewDir, newFilePath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	fmt.Printf("Renamed file:\n  FROM: %s\n  TO:   %s\n\n", oldFileInNewDir, newFilePath)
	return nil
}

// MoveSeriesFolder moves an entire series folder to a new location
func (r *Renamer) MoveSeriesFolder(oldDirPath, newDirPath string) error {
	if oldDirPath == newDirPath {
		return nil
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(newDirPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if _, err := os.Stat(newDirPath); err == nil {
		return fmt.Errorf("directory already exists: %s", newDirPath)
	}

	if r.dryRun {
		fmt.Printf("[DRY RUN] Would move folder:\n  FROM: %s\n  TO:   %s\n\n", oldDirPath, newDirPath)
		return nil
	}

	if err := os.Rename(oldDirPath, newDirPath); err != nil {
		return fmt.Errorf("failed to move folder: %w", err)
	}

	fmt.Printf("Moved folder:\n  FROM: %s\n  TO:   %s\n\n", oldDirPath, newDirPath)
	return nil
}

// MoveRenameMovieFile moves and renames a standalone movie file to a target directory with a new folder and filename
// This creates: outputDir/folderName/newFilename (and moves accompanying subtitles)
func (r *Renamer) MoveRenameMovieFile(oldPath, outputDir, folderName, newFilename string) error {
	oldDir := filepath.Dir(oldPath)
	oldFilename := filepath.Base(oldPath)
	oldNameWithoutExt := strings.TrimSuffix(oldFilename, filepath.Ext(oldFilename))
	newNameWithoutExt := strings.TrimSuffix(newFilename, filepath.Ext(newFilename))

	// Create movie folder path
	movieFolderPath := filepath.Join(outputDir, folderName)
	newPath := filepath.Join(movieFolderPath, newFilename)

	if oldPath == newPath {
		return nil
	}

	if r.dryRun {
		fmt.Printf("[DRY RUN] Would create folder: %s\n", movieFolderPath)
		fmt.Printf("[DRY RUN] Would move movie:\n  FROM: %s\n  TO:   %s\n", oldPath, newPath)

		// Check for subtitle files
		subtitles, _ := findSubtitleFiles(oldDir, oldNameWithoutExt)
		for _, subPath := range subtitles {
			subExt := filepath.Ext(subPath)
			newSubPath := filepath.Join(movieFolderPath, newNameWithoutExt+subExt)
			fmt.Printf("[DRY RUN] Would move subtitle:\n  FROM: %s\n  TO:   %s\n", subPath, newSubPath)
		}
		fmt.Println()
		return nil
	}

	// Create the movie folder
	if err := os.MkdirAll(movieFolderPath, 0755); err != nil {
		return fmt.Errorf("failed to create movie folder: %w", err)
	}

	// Check if destination file already exists
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("file already exists: %s", newPath)
	}

	// Move the video file
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to move movie file: %w", err)
	}
	fmt.Printf("Moved movie:\n  FROM: %s\n  TO:   %s\n", oldPath, newPath)

	// Move accompanying subtitle files
	subtitles, _ := findSubtitleFiles(oldDir, oldNameWithoutExt)
	for _, subPath := range subtitles {
		subExt := filepath.Ext(subPath)
		newSubPath := filepath.Join(movieFolderPath, newNameWithoutExt+subExt)
		if err := os.Rename(subPath, newSubPath); err != nil {
			fmt.Printf("Warning: failed to move subtitle %s: %v\n", subPath, err)
		} else {
			fmt.Printf("Moved subtitle:\n  FROM: %s\n  TO:   %s\n", subPath, newSubPath)
		}
	}
	fmt.Println()
	return nil
}

// MoveRenameMovieFolder moves and renames a movie folder to a target directory
// This moves the entire folder with all contents (video files, subtitles, disc structures)
func (r *Renamer) MoveRenameMovieFolder(oldDirPath, outputDir, newFolderName string) error {
	newDirPath := filepath.Join(outputDir, newFolderName)

	if oldDirPath == newDirPath {
		return nil
	}

	// Ensure parent directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if _, err := os.Stat(newDirPath); err == nil {
		return fmt.Errorf("directory already exists: %s", newDirPath)
	}

	if r.dryRun {
		fmt.Printf("[DRY RUN] Would move movie folder:\n  FROM: %s\n  TO:   %s\n\n", oldDirPath, newDirPath)
		return nil
	}

	if err := os.Rename(oldDirPath, newDirPath); err != nil {
		return fmt.Errorf("failed to move movie folder: %w", err)
	}

	fmt.Printf("Moved movie folder:\n  FROM: %s\n  TO:   %s\n\n", oldDirPath, newDirPath)
	return nil
}

// RenameMovieFileInFolder renames the main video file inside a movie folder (and accompanying subtitles)
func (r *Renamer) RenameMovieFileInFolder(folderPath, oldFilename, newFilename string) error {
	oldPath := filepath.Join(folderPath, oldFilename)
	newPath := filepath.Join(folderPath, newFilename)

	oldNameWithoutExt := strings.TrimSuffix(oldFilename, filepath.Ext(oldFilename))
	newNameWithoutExt := strings.TrimSuffix(newFilename, filepath.Ext(newFilename))

	if oldPath == newPath {
		return nil
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("file already exists: %s", newPath)
	}

	if r.dryRun {
		fmt.Printf("[DRY RUN] Would rename movie file:\n  FROM: %s\n  TO:   %s\n", oldPath, newPath)

		// Check for subtitle files
		subtitles, _ := findSubtitleFiles(folderPath, oldNameWithoutExt)
		for _, subPath := range subtitles {
			subExt := filepath.Ext(subPath)
			newSubPath := filepath.Join(folderPath, newNameWithoutExt+subExt)
			fmt.Printf("[DRY RUN] Would rename subtitle:\n  FROM: %s\n  TO:   %s\n", subPath, newSubPath)
		}
		fmt.Println()
		return nil
	}

	// Rename the video file
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename movie file: %w", err)
	}
	fmt.Printf("Renamed movie file:\n  FROM: %s\n  TO:   %s\n", oldPath, newPath)

	// Rename accompanying subtitle files
	subtitles, _ := findSubtitleFiles(folderPath, oldNameWithoutExt)
	for _, subPath := range subtitles {
		subExt := filepath.Ext(subPath)
		newSubPath := filepath.Join(folderPath, newNameWithoutExt+subExt)
		if err := os.Rename(subPath, newSubPath); err != nil {
			fmt.Printf("Warning: failed to rename subtitle %s: %v\n", subPath, err)
		} else {
			fmt.Printf("Renamed subtitle:\n  FROM: %s\n  TO:   %s\n", subPath, newSubPath)
		}
	}
	fmt.Println()
	return nil
}

// findSubtitleFiles finds all subtitle files matching the video filename
func findSubtitleFiles(dir, videoNameWithoutExt string) ([]string, error) {
	var subtitles []string
	subtitleExts := []string{".srt", ".sub", ".ass", ".ssa", ".vtt"}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		nameWithoutExt := strings.TrimSuffix(name, filepath.Ext(name))
		ext := strings.ToLower(filepath.Ext(name))

		// Check if it's a subtitle file with matching name
		for _, subExt := range subtitleExts {
			if ext == subExt && strings.HasPrefix(nameWithoutExt, videoNameWithoutExt) {
				subtitles = append(subtitles, filepath.Join(dir, name))
				break
			}
		}
	}

	return subtitles, nil
}
