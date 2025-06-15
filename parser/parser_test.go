package parser_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/alexhokl/japanese-notes-parser/cmd"
	"github.com/alexhokl/japanese-notes-parser/database"
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
			input:                "### Nouns (position)",
			expectedLevel:        "",
			expectedPartOfSpeech: "Nouns (position)",
			expectedError:        nil,
		},
		{
			input:                "### Adjectives (い)",
			expectedLevel:        "",
			expectedPartOfSpeech: "Adjectives (い)",
			expectedError:        nil,
		},
		{
			input:                "### Auxiliary verbs",
			expectedLevel:        "",
			expectedPartOfSpeech: "Auxiliary verbs",
			expectedError:        nil,
		},
		{
			input:                "### する verbs",
			expectedLevel:        "",
			expectedPartOfSpeech: "する verbs",
			expectedError:        nil,
		},
		{
			input:                "### 連体詞",
			expectedLevel:        "",
			expectedPartOfSpeech: "連体詞",
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
			actualLevel, actualPartOfSpeech, err := parser.ParseHeaderLine(test.input, parser.LevelRegex, parser.PartOfSpeechRegex)
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

func TestParseLine(t *testing.T) {
	tests := []struct {
		input                string
		cachedLevel          string
		cachedPartOfSpeech   string
		expectedEntry        *database.Entry
		expectedLevel        string
		expectedPartOfSpeech string
		expectedError        error
	}{
		{
			input:                "- 階 (かい) - floor",
			cachedLevel:          "N1",
			cachedPartOfSpeech:   "Nouns",
			expectedEntry:        &database.Entry{Kanji: "階", Kana: "かい", English: []string{"floor"}},
			expectedLevel:        "N1",
			expectedPartOfSpeech: "Nouns",
			expectedError:        nil,
		},
		{
			input:                "- インフォメーション - information",
			cachedLevel:          "N1",
			cachedPartOfSpeech:   "Nouns",
			expectedEntry:        &database.Entry{Kanji: "", Kana: "インフォメーション", English: []string{"information"}},
			expectedLevel:        "N1",
			expectedPartOfSpeech: "Nouns",
			expectedError:        nil,
		},
		{
			input:                "",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Nouns",
			expectedEntry:        nil,
			expectedLevel:        "N5",
			expectedPartOfSpeech: "Nouns",
			expectedError:        nil,
		},
		{
			input:                "## N4",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Nouns",
			expectedEntry:        nil,
			expectedLevel:        "N4",
			expectedPartOfSpeech: "Nouns",
			expectedError:        nil,
		},
		{
			input:                "### Adjectives",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Nouns",
			expectedEntry:        nil,
			expectedLevel:        "N5",
			expectedPartOfSpeech: "Adjectives",
			expectedError:        nil,
		},
		{
			input:                "- でも - but",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Adverbs",
			expectedEntry:        &database.Entry{Kanji: "", Kana: "でも", English: []string{"but"}},
			expectedLevel:        "N5",
			expectedPartOfSpeech: "Adverbs",
			expectedError:        nil,
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("input=%s", test.input)
		t.Run(testName, func(t *testing.T) {
			actualEntry, actualLevel, actualPartOfSpeech, err := cmd.ParseLine(
				test.input,
				parser.PointRegex,
				parser.JapaneseRegex,
				parser.LevelRegex,
				parser.PartOfSpeechRegex,
				test.cachedLevel,
				test.cachedPartOfSpeech,
			)
			if actualLevel != test.expectedLevel {
				t.Errorf("expected level %s, got %s", test.expectedLevel, actualLevel)
			}
			if actualPartOfSpeech != test.expectedPartOfSpeech {
				t.Errorf("expected part of speech %s, got %s", test.expectedPartOfSpeech, actualPartOfSpeech)
			}
			if err != test.expectedError {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			}
			if test.expectedEntry != nil {
				if actualEntry.Kanji != test.expectedEntry.Kanji {
					t.Errorf("expected kanji %s, got %s", test.expectedEntry.Kanji, actualEntry.Kanji)
				}
				if actualEntry.Kana != test.expectedEntry.Kana {
					t.Errorf("expected kana %s, got %s", test.expectedEntry.Kana, actualEntry.Kana)
				}
				if len(actualEntry.English) != len(test.expectedEntry.English) {
					t.Errorf("expected english entries of %v, got %v", len(test.expectedEntry.English), len(actualEntry.English))
				}
				for i, v := range test.expectedEntry.English {
					if actualEntry.English[i] != v {
						t.Errorf("expected english %s, got %s", v, actualEntry.English[i])
					}
				}
			}
		})
	}
}

// - 階 (かい) - floor
// - インフォメーション - information
