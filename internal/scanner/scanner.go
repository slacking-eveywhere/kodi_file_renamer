package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	// videoExtensions lists all supported video file extensions
	videoExtensions = []string{".mkv", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v", ".iso"}

	// subtitleExtensions lists all supported subtitle file extensions
	subtitleExtensions = []string{".srt", ".sub", ".ass", ".ssa", ".vtt"}

	// seriesPatterns contains regular expressions to detect TV series episode numbering
	seriesPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)S(\d{1,2})E(\d{1,2})`),
		regexp.MustCompile(`(?i)(\d{1,2})x(\d{1,2})`),
		regexp.MustCompile(`(?i)Season\s*(\d{1,2})\s*Episode\s*(\d{1,2})`),
		regexp.MustCompile(`(?i)s(\d{1,2})\s*-\s*e(\d{1,2})`),
	}

	// yearPattern matches year formats in movie filenames
	yearPattern = regexp.MustCompile(`\((\d{4})\)|\[(\d{4})\]|(?:^|[\s\.])(\d{4})(?:[\s\.]|$)`)
)

// MediaFile represents a media file with its metadata and classification
type MediaFile struct {
	Path          string
	Name          string
	Extension     string
	IsMovie       bool
	IsSeries      bool
	Season        int
	Episode       int
	Year          int
	CleanName     string
	ParentDir     string   // Parent directory name for series files
	IsMovieFolder bool     // True if movie is a folder (contains video + subtitles/extras)
	MovieFiles    []string // All video files in movie folder
	SubtitleFiles []string // All subtitle files in movie folder
	IsBluRay      bool     // True if folder contains Blu-ray structure
	IsDVD         bool     // True if folder contains DVD structure
}

// EpisodeRenameTask represents a pending episode rename operation
type EpisodeRenameTask struct {
	File         *MediaFile
	EpisodeName  string
	NewFilename  string
	Season       int
	Episode      int
	HasError     bool
	ErrorMessage string
}

// SeriesBatchRename groups all episodes for a series folder
type SeriesBatchRename struct {
	OriginalFolderPath string
	OriginalFolderName string
	NewFolderName      string
	SeriesName         string
	SeriesYear         string
	Episodes           []EpisodeRenameTask
	NeedsFolderRename  bool
}

// Scanner scans directories for media files and extracts metadata
type Scanner struct {
	rootPath string
}

// NewScanner creates a new Scanner for the specified root directory path
func NewScanner(rootPath string) *Scanner {
	return &Scanner{
		rootPath: rootPath,
	}
}

// ScanDirectory recursively scans the root directory and returns all media files found
func (s *Scanner) ScanDirectory() ([]MediaFile, error) {
	var mediaFiles []MediaFile
	processedDirs := make(map[string]bool)

	err := filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip root directory
		if path == s.rootPath {
			return nil
		}

		// Check if this is a directory
		if info.IsDir() {
			// Skip if already processed
			if processedDirs[path] {
				return filepath.SkipDir
			}

			// Check if this is a movie folder
			movieFile, isMovieFolder := s.parseMovieFolder(path)
			if isMovieFolder {
				mediaFiles = append(mediaFiles, movieFile)
				processedDirs[path] = true
				return filepath.SkipDir // Skip processing contents
			}
			return nil
		}

		// Process individual video files (not in movie folders)
		ext := strings.ToLower(filepath.Ext(path))
		if !isVideoFile(ext) {
			return nil
		}

		// Check if parent directory is already processed as movie folder
		parentDir := filepath.Dir(path)
		if processedDirs[parentDir] {
			return nil
		}

		// This is a standalone video file or series episode
		mediaFile := s.parseFile(path, info.Name())
		mediaFiles = append(mediaFiles, mediaFile)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return mediaFiles, nil
}

// parseFile extracts metadata from a media file and classifies it as movie or TV series
func (s *Scanner) parseFile(path, filename string) MediaFile {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Get parent directory name for series organization
	parentDir := filepath.Base(filepath.Dir(path))

	mediaFile := MediaFile{
		Path:      path,
		Name:      filename,
		Extension: ext,
		ParentDir: parentDir,
	}

	// Check if it's a TV series
	season, episode, found := s.extractSeriesInfo(nameWithoutExt)
	if found {
		mediaFile.IsSeries = true
		mediaFile.Season = season
		mediaFile.Episode = episode
		mediaFile.CleanName = s.cleanSeriesName(nameWithoutExt)
	} else {
		mediaFile.IsMovie = true
		mediaFile.Year = s.extractYear(nameWithoutExt)
		mediaFile.CleanName = s.cleanMovieName(nameWithoutExt)
	}

	return mediaFile
}

// extractSeriesInfo attempts to extract season and episode numbers from a filename
func (s *Scanner) extractSeriesInfo(name string) (season, episode int, found bool) {
	for _, pattern := range seriesPatterns {
		matches := pattern.FindStringSubmatch(name)
		if len(matches) >= 3 {
			season, _ = strconv.Atoi(matches[1])
			episode, _ = strconv.Atoi(matches[2])
			return season, episode, true
		}
	}
	return 0, 0, false
}

// extractYear attempts to extract a year from a movie filename
func (s *Scanner) extractYear(name string) int {
	matches := yearPattern.FindStringSubmatch(name)
	if len(matches) > 0 {
		for i := 1; i < len(matches); i++ {
			if matches[i] != "" {
				year, err := strconv.Atoi(matches[i])
				if err == nil && year >= 1900 && year <= 2100 {
					return year
				}
			}
		}
	}
	return 0
}

// cleanSeriesName removes episode numbering and artifacts from a series filename
func (s *Scanner) cleanSeriesName(name string) string {
	// Remove series patterns
	for _, pattern := range seriesPatterns {
		name = pattern.ReplaceAllString(name, "")
	}

	// Remove common artifacts
	name = removeCommonArtifacts(name)

	return strings.TrimSpace(name)
}

// GetSeriesSearchQuery extracts clean series name from parent directory for API search
func GetSeriesSearchQuery(parentDir string) string {
	name := parentDir

	// Remove year in parentheses or brackets
	yearPattern := regexp.MustCompile(`\s*[\(\[]?\d{4}[\)\]]?\s*`)
	name = yearPattern.ReplaceAllString(name, " ")

	// Replace dots and underscores with spaces
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Remove common artifacts
	name = removeCommonArtifacts(name)

	// Clean up multiple spaces
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")

	return strings.TrimSpace(name)
}

// cleanMovieName removes year and artifacts from a movie filename
func (s *Scanner) cleanMovieName(name string) string {
	// Remove year
	name = yearPattern.ReplaceAllString(name, " ")

	// Remove common artifacts
	name = removeCommonArtifacts(name)

	return strings.TrimSpace(name)
}

// removeCommonArtifacts removes quality indicators, separators, and brackets from filenames
func removeCommonArtifacts(name string) string {
	// Replace common separators with spaces
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	// Remove quality indicators
	qualityPatterns := []string{
		`(?i)1080p`, `(?i)720p`, `(?i)480p`, `(?i)4k`,
		`(?i)bluray`, `(?i)brrip`, `(?i)webrip`, `(?i)web-dl`,
		`(?i)hdtv`, `(?i)dvdrip`, `(?i)xvid`, `(?i)x264`, `(?i)x265`,
		`(?i)hevc`, `(?i)aac`, `(?i)ac3`, `(?i)dts`,
		`(?i)\[.*?\]`, `(?i)\(.*?\)`,
	}

	for _, pattern := range qualityPatterns {
		re := regexp.MustCompile(pattern)
		name = re.ReplaceAllString(name, " ")
	}

	// Remove multiple spaces
	spacePattern := regexp.MustCompile(`\s+`)
	name = spacePattern.ReplaceAllString(name, " ")

	return name
}

// isVideoFile checks if the given extension is a supported video format
func isVideoFile(ext string) bool {
	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return true
		}
	}
	return false
}

// isSubtitleFile checks if the given extension is a supported subtitle format
func isSubtitleFile(ext string) bool {
	for _, subExt := range subtitleExtensions {
		if ext == subExt {
			return true
		}
	}
	return false
}

// parseMovieFolder checks if a directory is a movie folder and parses it
func (s *Scanner) parseMovieFolder(dirPath string) (MediaFile, bool) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return MediaFile{}, false
	}

	var videoFiles []string
	var subtitleFiles []string
	var hasBluRay bool
	var hasDVD bool

	// Check for Blu-ray structure
	for _, entry := range entries {
		name := strings.ToUpper(entry.Name())
		if entry.IsDir() && name == "BDMV" {
			hasBluRay = true
		}
		if entry.IsDir() && name == "VIDEO_TS" {
			hasDVD = true
		}
	}

	// Collect video and subtitle files
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		fullPath := filepath.Join(dirPath, entry.Name())

		if isVideoFile(ext) {
			videoFiles = append(videoFiles, fullPath)
		} else if isSubtitleFile(ext) {
			subtitleFiles = append(subtitleFiles, fullPath)
		}
	}

	// Determine if this is a movie folder:
	// 1. Has Blu-ray or DVD structure, OR
	// 2. Has video files + subtitle files, OR
	// 3. Has exactly one video file and folder name looks like a movie (has year or matches video name)
	isMovieFolder := false
	var mainVideoFile string

	if hasBluRay || hasDVD {
		isMovieFolder = true
		if len(videoFiles) > 0 {
			mainVideoFile = videoFiles[0]
		}
	} else if len(videoFiles) > 0 && len(subtitleFiles) > 0 {
		isMovieFolder = true
		mainVideoFile = videoFiles[0]
	} else if len(videoFiles) == 1 {
		// Check if folder name matches video file or contains year
		folderName := filepath.Base(dirPath)
		videoName := strings.TrimSuffix(filepath.Base(videoFiles[0]), filepath.Ext(videoFiles[0]))

		// Check for year in folder name
		hasYear := yearPattern.MatchString(folderName)

		// Check if folder name is similar to video name (allowing for year differences)
		cleanFolder := removeCommonArtifacts(folderName)
		cleanVideo := removeCommonArtifacts(videoName)
		namesMatch := strings.Contains(strings.ToLower(cleanFolder), strings.ToLower(cleanVideo)) ||
			strings.Contains(strings.ToLower(cleanVideo), strings.ToLower(cleanFolder))

		if hasYear || namesMatch {
			isMovieFolder = true
			mainVideoFile = videoFiles[0]
		}
	}

	if !isMovieFolder {
		return MediaFile{}, false
	}

	// Parse folder name as movie
	folderName := filepath.Base(dirPath)
	year := s.extractYear(folderName)
	cleanName := s.cleanMovieName(folderName)

	// If no year in folder name, try to get from main video file
	if year == 0 && mainVideoFile != "" {
		videoFileName := strings.TrimSuffix(filepath.Base(mainVideoFile), filepath.Ext(mainVideoFile))
		year = s.extractYear(videoFileName)
		if cleanName == "" || cleanName == filepath.Base(dirPath) {
			cleanName = s.cleanMovieName(videoFileName)
		}
	}

	// Determine extension from main video file
	ext := ""
	if mainVideoFile != "" {
		ext = filepath.Ext(mainVideoFile)
	} else if hasBluRay {
		ext = ".bluray"
	} else if hasDVD {
		ext = ".dvd"
	}

	return MediaFile{
		Path:          dirPath,
		Name:          folderName,
		Extension:     ext,
		IsMovie:       true,
		IsSeries:      false,
		Year:          year,
		CleanName:     cleanName,
		IsMovieFolder: true,
		MovieFiles:    videoFiles,
		SubtitleFiles: subtitleFiles,
		IsBluRay:      hasBluRay,
		IsDVD:         hasDVD,
		ParentDir:     filepath.Base(filepath.Dir(dirPath)),
	}, true
}

// GetSearchQuery returns the clean name suitable for API searches
func (m *MediaFile) GetSearchQuery() string {
	return m.CleanName
}

// GetNewFilename generates a properly formatted filename for a TV series episode
func (m *MediaFile) GetNewFilename(seriesName, episodeName string) string {
	if m.IsSeries {
		seasonStr := fmt.Sprintf("S%02d", m.Season)
		episodeStr := fmt.Sprintf("E%02d", m.Episode)
		if episodeName != "" {
			return fmt.Sprintf("%s %s%s - %s%s", seriesName, seasonStr, episodeStr, episodeName, m.Extension)
		}
		return fmt.Sprintf("%s %s%s%s", seriesName, seasonStr, episodeStr, m.Extension)
	}
	return m.Name
}

// GetEpisodeFilename returns the formatted filename for a series episode
func (m *MediaFile) GetEpisodeFilename(seriesName string, season, episode int, episodeName string) string {
	seasonStr := fmt.Sprintf("S%02d", season)
	episodeStr := fmt.Sprintf("E%02d", episode)
	if episodeName != "" {
		return fmt.Sprintf("%s %s%s - %s%s", seriesName, seasonStr, episodeStr, episodeName, m.Extension)
	}
	return fmt.Sprintf("%s %s%s%s", seriesName, seasonStr, episodeStr, m.Extension)
}

// GetMovieFilename generates a properly formatted filename for a movie
func (m *MediaFile) GetMovieFilename(title string, year string) string {
	if m.IsMovie {
		cleanTitle := strings.ReplaceAll(title, ":", " -")
		cleanTitle = strings.ReplaceAll(cleanTitle, "/", " ")
		if year != "" {
			return fmt.Sprintf("%s (%s)%s", cleanTitle, year, m.Extension)
		}
		return fmt.Sprintf("%s%s", cleanTitle, m.Extension)
	}
	return m.Name
}

// GetMovieFolderName generates a properly formatted folder name for a movie
func (m *MediaFile) GetMovieFolderName(title string, year string) string {
	if m.IsMovie {
		cleanTitle := strings.ReplaceAll(title, ":", " -")
		cleanTitle = strings.ReplaceAll(cleanTitle, "/", " ")
		if year != "" {
			return fmt.Sprintf("%s (%s)", cleanTitle, year)
		}
		return cleanTitle
	}
	return m.Name
}

// GetSeriesFolderName generates a properly formatted folder name for a TV series
func (m *MediaFile) GetSeriesFolderName(seriesName, year string) string {
	if m.IsSeries {
		cleanName := strings.ReplaceAll(seriesName, ":", " -")
		cleanName = strings.ReplaceAll(cleanName, "/", " ")
		if year != "" {
			return fmt.Sprintf("%s (%s)", cleanName, year)
		}
		return cleanName
	}
	return m.ParentDir
}

// NeedsSeriesFolderRename checks if the parent directory needs to be renamed
func (m *MediaFile) NeedsSeriesFolderRename(expectedFolderName string) bool {
	if !m.IsSeries {
		return false
	}
	return m.ParentDir != expectedFolderName
}

// GetParentDirPath returns the full path to the parent directory
func (m *MediaFile) GetParentDirPath() string {
	return filepath.Dir(m.Path)
}

// GroupSeriesByFolder organizes series episodes by their parent directory
func GroupSeriesByFolder(files []MediaFile) map[string][]*MediaFile {
	seriesMap := make(map[string][]*MediaFile)

	for i := range files {
		file := &files[i]
		if file.IsSeries {
			folderPath := file.GetParentDirPath()
			seriesMap[folderPath] = append(seriesMap[folderPath], file)
		}
	}

	return seriesMap
}
