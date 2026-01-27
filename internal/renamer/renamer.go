package renamer

import (
	"fmt"
	"os"
	"path/filepath"
)

type Renamer struct {
	dryRun bool
}

func NewRenamer(dryRun bool) *Renamer {
	return &Renamer{
		dryRun: dryRun,
	}
}

func (r *Renamer) RenameFile(oldPath, newFilename string) error {
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newFilename)

	if oldPath == newPath {
		return nil
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("file already exists: %s", newPath)
	}

	if r.dryRun {
		fmt.Printf("[DRY RUN] Would rename:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
		return nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename file: %w", err)
	}

	fmt.Printf("Renamed:\n  FROM: %s\n  TO:   %s\n\n", oldPath, newPath)
	return nil
}

func (r *Renamer) SetDryRun(dryRun bool) {
	r.dryRun = dryRun
}

func (r *Renamer) IsDryRun() bool {
	return r.dryRun
}
