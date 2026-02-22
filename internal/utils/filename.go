package utils

import (
	"regexp"
	"strings"
)

// SanitizeFilename removes or replaces characters that are invalid in filenames
// for both Unix and Windows filesystems
func SanitizeFilename(name string) string {
	// Replace characters that are invalid on Windows: < > : " / \ | ? *
	// Also handle characters that can cause issues on Unix: /
	replacements := map[string]string{
		":":  " -", // Colon to dash (common in titles)
		"/":  " ",  // Forward slash to space
		"\\": " ",  // Backslash to space
		"<":  "",   // Less than - remove
		">":  "",   // Greater than - remove
		"|":  "-",  // Pipe to dash
		"?":  "",   // Question mark - remove
		"*":  "",   // Asterisk - remove
		"\"": "'",  // Double quote to single quote
	}

	result := name
	for old, new := range replacements {
		result = strings.ReplaceAll(result, old, new)
	}

	// Remove control characters (ASCII 0-31)
	result = regexp.MustCompile(`[\x00-\x1F]`).ReplaceAllString(result, "")

	// Trim leading/trailing spaces and dots (Windows doesn't allow these)
	result = strings.Trim(result, " .")

	// Replace multiple spaces with single space
	result = regexp.MustCompile(`\s+`).ReplaceAllString(result, " ")

	return result
}
