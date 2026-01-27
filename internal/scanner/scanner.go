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
	videoExtensions = []string{".mkv", ".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v"}

	// TV Show patterns
	seriesPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)S(\d{1,2})E(\d{1,2})`),
		regexp.MustCompile(`(?i)(\d{1,2})x(\d{1,2})`),
		regexp.MustCompile(`(?i)Season\s*(\d{1,2})\s*Episode\s*(\d{1,2})`),
		regexp.MustCompile(`(?i)s(\d{1,2})\s*-\s*e(\d{1,2})`),
	}

	// Year pattern for movies
	yearPattern = regexp.MustCompile(`\((\d{4})\)|\[(\d{4})\]|[\s\.](\d{4})[\s\.]`)
)

type MediaFile struct {
	Path      string
	Name      string
	Extension string
	IsMovie   bool
	IsSeries  bool
	Season    int
	Episode   int
	Year      int
	CleanName string
}

type Scanner struct {
	rootPath string
}

func NewScanner(rootPath string) *Scanner {
	return &Scanner{
		rootPath: rootPath,
	}
}

func (s *Scanner) ScanDirectory() ([]MediaFile, error) {
	var mediaFiles []MediaFile

	err := filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !isVideoFile(ext) {
			return nil
		}

		mediaFile := s.parseFile(path, info.Name())
		mediaFiles = append(mediaFiles, mediaFile)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return mediaFiles, nil
}

func (s *Scanner) parseFile(path, filename string) MediaFile {
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	mediaFile := MediaFile{
		Path:      path,
		Name:      filename,
		Extension: ext,
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

func (s *Scanner) cleanSeriesName(name string) string {
	// Remove series patterns
	for _, pattern := range seriesPatterns {
		name = pattern.ReplaceAllString(name, "")
	}

	// Remove common artifacts
	name = removeCommonArtifacts(name)

	return strings.TrimSpace(name)
}

func (s *Scanner) cleanMovieName(name string) string {
	// Remove year
	name = yearPattern.ReplaceAllString(name, " ")

	// Remove common artifacts
	name = removeCommonArtifacts(name)

	return strings.TrimSpace(name)
}

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

func isVideoFile(ext string) bool {
	for _, videoExt := range videoExtensions {
		if ext == videoExt {
			return true
		}
	}
	return false
}

func (m *MediaFile) GetSearchQuery() string {
	return m.CleanName
}

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

func (m *MediaFile) GetMovieFilename(title string, year int) string {
	if m.IsMovie {
		cleanTitle := strings.ReplaceAll(title, ":", " -")
		cleanTitle = strings.ReplaceAll(cleanTitle, "/", " ")
		if year > 0 {
			return fmt.Sprintf("%s (%d)%s", cleanTitle, year, m.Extension)
		}
		return fmt.Sprintf("%s%s", cleanTitle, m.Extension)
	}
	return m.Name
}
