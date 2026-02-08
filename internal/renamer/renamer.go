package renamer

import (
	"fmt"
	"os"
	"path/filepath"
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

// MoveFile moves a file to a new location (different directory)
func (r *Renamer) MoveFile(oldPath, newPath string) error {
	if oldPath == newPath {
		return nil
	}

	// Ensure parent directory exists
	newDir := filepath.Dir(newPath)
	if err := os.MkdirAll(newDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("file already exists: %s", newPath)
	}

	if r.dryRun {
		fmt.Printf("[DRY RUN] Would move:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
		return nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to move file: %w", err)
	}

	fmt.Printf("Moved:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
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
