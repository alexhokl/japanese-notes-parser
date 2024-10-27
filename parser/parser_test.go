package parser_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/alexhokl/japanese-notes-parser/parser"
)

func TestParseEnglish(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			input:    "mother",
			expected: []string{"mother"},
		},
		{
			input:    "older sister",
			expected: []string{"older sister"},
		},
		{
			input:    "teacher / master",
			expected: []string{"teacher", "master"},
		},
		{
			input:    "-ish / -wise / -like",
			expected: []string{"-ish", "-wise", "-like"},
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("input=%s", test.input)
		t.Run(testName, func(t *testing.T) {
			actual := parser.ParseEnglish(test.input)
			for i, v := range test.expected {
				if actual[i] != v {
					t.Errorf("expected %s, got %s", v, actual[i])
				}
			}
		})
	}
}

func TestParseHeaderLine(t *testing.T) {
	levelRegex, _ := regexp.Compile(`^## (?P<level>\w+)`)
	partOfSpeechRegex, _ := regexp.Compile(`^### (?P<partOfSpeech>\w+)`)

	tests := []struct {
		input                string
		levelRegex           *regexp.Regexp
		partOfSpeechRegex    *regexp.Regexp
		expectedLevel        string
		expectedPartOfSpeech string
		expectedError        error
	}{
		{
			input:                "## N5",
			expectedLevel:        "N5",
			expectedPartOfSpeech: "",
			expectedError:        nil,
		},
		{
			input:                "### Nouns",
			expectedLevel:        "",
			expectedPartOfSpeech: "Nouns",
			expectedError:        nil,
		},
		{
			input:                "# Vocabulary",
			expectedLevel:        "",
			expectedPartOfSpeech: "",
			expectedError:        nil,
		},
		{
			input:                "- [something](link)",
			expectedLevel:        "",
			expectedPartOfSpeech: "",
			expectedError:        nil,
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("input=%s", test.input)
		t.Run(testName, func(t *testing.T) {
			actualLevel, actualPartOfSpeech, err := parser.ParseHeaderLine(test.input, levelRegex, partOfSpeechRegex)
			if actualLevel != test.expectedLevel {
				t.Errorf("expected %s, got %s", test.expectedLevel, actualLevel)
			}
			if actualPartOfSpeech != test.expectedPartOfSpeech {
				t.Errorf("expected %s, got %s", test.expectedPartOfSpeech, actualPartOfSpeech)
			}
			if err != test.expectedError {
				t.Errorf("expected %v, got %v", test.expectedError, err)
			}
		})
	}
}
