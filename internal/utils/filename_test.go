package utils

import "testing"

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal filename",
			input:    "Normal Movie Title",
			expected: "Normal Movie Title",
		},
		{
			name:     "Colon in title",
			input:    "Movie: The Sequel",
			expected: "Movie - The Sequel",
		},
		{
			name:     "Forward slash",
			input:    "AC/DC Live",
			expected: "AC DC Live",
		},
		{
			name:     "Backward slash",
			input:    "Path\\to\\something",
			expected: "Path to something",
		},
		{
			name:     "Multiple invalid characters",
			input:    "Title: Part 1/2 <Extended>",
			expected: "Title - Part 1 2 Extended",
		},
		{
			name:     "Question mark",
			input:    "Who? What? Where?",
			expected: "Who What Where",
		},
		{
			name:     "Asterisk",
			input:    "File*Name*Test",
			expected: "FileNameTest",
		},
		{
			name:     "Pipe character",
			input:    "Option A | Option B",
			expected: "Option A - Option B",
		},
		{
			name:     "Double quotes",
			input:    `Movie "The Best"`,
			expected: "Movie 'The Best'",
		},
		{
			name:     "Less than and greater than",
			input:    "A<B>C",
			expected: "ABC",
		},
		{
			name:     "Multiple spaces",
			input:    "Title   with    spaces",
			expected: "Title with spaces",
		},
		{
			name:     "Leading and trailing spaces",
			input:    "  Title  ",
			expected: "Title",
		},
		{
			name:     "Leading and trailing dots",
			input:    "..Title..",
			expected: "Title",
		},
		{
			name:     "Complex real-world example",
			input:    "Marvel's What If...?: Season 2",
			expected: "Marvel's What If... - Season 2",
		},
		{
			name:     "Episode with slash",
			input:    "The One Where They're Up All Night / The One With The Routine",
			expected: "The One Where They're Up All Night The One With The Routine",
		},
		{
			name:     "Movie with multiple colons",
			input:    "Star Wars: Episode V: The Empire Strikes Back",
			expected: "Star Wars - Episode V - The Empire Strikes Back",
		},
		{
			name:     "All invalid characters",
			input:    `<>:"/\|?*`,
			expected: "-' -",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only spaces",
			input:    "     ",
			expected: "",
		},
		{
			name:     "Only invalid characters that get removed",
			input:    "<>?*",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilename(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeFilename_RealWorldExamples(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "IT Crowd episode",
			input:    "The IT Crowd S01E01 - Yesterday's Jam",
			expected: "The IT Crowd S01E01 - Yesterday's Jam",
		},
		{
			name:     "Movie with year",
			input:    "Alien: Covenant (2017)",
			expected: "Alien - Covenant (2017)",
		},
		{
			name:     "Series with special chars",
			input:    "Mr. Robot S01E01 - eps1.0_hellofriend.mov",
			expected: "Mr. Robot S01E01 - eps1.0_hellofriend.mov",
		},
		{
			name:     "Movie with subtitle separator",
			input:    "Borat: Cultural Learnings of America for Make Benefit Glorious Nation of Kazakhstan",
			expected: "Borat - Cultural Learnings of America for Make Benefit Glorious Nation of Kazakhstan",
		},
		{
			name:     "Episode with question mark",
			input:    "S01E01 - Where Are You?",
			expected: "S01E01 - Where Are You",
		},
		{
			name:     "Movie with slash in title",
			input:    "24/7 (2010)",
			expected: "24 7 (2010)",
		},
		{
			name:     "Episode name with multiple invalid chars",
			input:    "Episode Title: Part 1/2 <Director's Cut>",
			expected: "Episode Title - Part 1 2 Director's Cut",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeFilename(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
