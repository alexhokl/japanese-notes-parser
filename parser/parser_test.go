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
		input                 string
		cachedLevel           string
		cachedPartOfSpeech    string
		previousKanji         string
		expectedEntry         *database.Entry
		expectedLevel         string
		expectedPartOfSpeech  string
		expectedContinuation  string
		expectedError         error
	}{
		{
			input:                "- 階 (かい) - floor",
			cachedLevel:          "N1",
			cachedPartOfSpeech:   "Nouns",
			previousKanji:        "",
			expectedEntry:        &database.Entry{Kanji: "階", Kana: "かい", English: []string{"floor"}, LinkedKanji: ""},
			expectedLevel:        "N1",
			expectedPartOfSpeech: "Nouns",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			input:                "- インフォメーション - information",
			cachedLevel:          "N1",
			cachedPartOfSpeech:   "Nouns",
			previousKanji:        "",
			expectedEntry:        &database.Entry{Kanji: "", Kana: "インフォメーション", English: []string{"information"}, LinkedKanji: ""},
			expectedLevel:        "N1",
			expectedPartOfSpeech: "Nouns",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			input:                "",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Nouns",
			previousKanji:        "",
			expectedEntry:        nil,
			expectedLevel:        "N5",
			expectedPartOfSpeech: "Nouns",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			input:                "## N4",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Nouns",
			previousKanji:        "",
			expectedEntry:        nil,
			expectedLevel:        "N4",
			expectedPartOfSpeech: "Nouns",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			input:                "### Adjectives",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Nouns",
			previousKanji:        "",
			expectedEntry:        nil,
			expectedLevel:        "N5",
			expectedPartOfSpeech: "Adjectives",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			input:                "- でも - but",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Adverbs",
			previousKanji:        "",
			expectedEntry:        &database.Entry{Kanji: "", Kana: "でも", English: []string{"but"}, LinkedKanji: ""},
			expectedLevel:        "N5",
			expectedPartOfSpeech: "Adverbs",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			input:                "  * 出す (だす) - to take out",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "verb pair",
			previousKanji:        "でる",
			expectedEntry:        &database.Entry{Kanji: "出す", Kana: "だす", English: []string{"to take out"}, LinkedKanji: "でる"},
			expectedLevel:        "N5",
			expectedPartOfSpeech: "verb pair",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			// Continuation of the previous subpoint's English text. Previously
			// silently dropped; now surfaced as a continuation string so the
			// caller can append it to the last entry's English field.
			input:                "    to take out",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "verb pair",
			previousKanji:        "",
			expectedEntry:        nil,
			expectedLevel:        "N5",
			expectedPartOfSpeech: "verb pair",
			expectedContinuation: "to take out",
			expectedError:        nil,
		},
		{
			// Continuation after a top-level point (not a subpoint).
			input:                "    floor of a building",
			cachedLevel:          "N1",
			cachedPartOfSpeech:   "Nouns",
			previousKanji:        "階",
			expectedEntry:        nil,
			expectedLevel:        "N1",
			expectedPartOfSpeech: "Nouns",
			expectedContinuation: "floor of a building",
			expectedError:        nil,
		},
		{
			// Whitespace-only lines are not continuations.
			input:                "    ",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "Nouns",
			previousKanji:        "",
			expectedEntry:        nil,
			expectedLevel:        "N5",
			expectedPartOfSpeech: "Nouns",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			// Regression guard: a subpoint line must not be mistaken for a
			// continuation even though it starts with leading whitespace.
			input:                "  * 降ろす (おろす) - to take down",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "verb pair",
			previousKanji:        "降りる",
			expectedEntry:        &database.Entry{Kanji: "降ろす", Kana: "おろす", English: []string{"to take down"}, LinkedKanji: "降りる"},
			expectedLevel:        "N5",
			expectedPartOfSpeech: "verb pair",
			expectedContinuation: "",
			expectedError:        nil,
		},
		{
			// Continuation line containing the trailing text from the bug
			// report's example.
			input:                "    unload (goods), to withdraw (money); aka 下ろす",
			cachedLevel:          "N5",
			cachedPartOfSpeech:   "verb pair",
			previousKanji:        "降ろす",
			expectedEntry:        nil,
			expectedLevel:        "N5",
			expectedPartOfSpeech: "verb pair",
			expectedContinuation: "unload (goods), to withdraw (money); aka 下ろす",
			expectedError:        nil,
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("input=%s", test.input)
		t.Run(testName, func(t *testing.T) {
			actualEntry, actualLevel, actualPartOfSpeech, actualContinuation, err := cmd.ParseLine(
				test.input,
				parser.PointRegex,
				parser.SubPointRegex,
				parser.JapaneseRegex,
				parser.LevelRegex,
				parser.PartOfSpeechRegex,
				test.cachedLevel,
				test.cachedPartOfSpeech,
				test.previousKanji,
			)
			if actualLevel != test.expectedLevel {
				t.Errorf("expected level %s, got %s", test.expectedLevel, actualLevel)
			}
			if actualPartOfSpeech != test.expectedPartOfSpeech {
				t.Errorf("expected part of speech %s, got %s", test.expectedPartOfSpeech, actualPartOfSpeech)
			}
			if actualContinuation != test.expectedContinuation {
				t.Errorf("expected continuation %q, got %q", test.expectedContinuation, actualContinuation)
			}
			if err != test.expectedError {
				t.Errorf("expected error %v, got %v", test.expectedError, err)
			}
			if test.expectedEntry != nil {
				if actualEntry == nil {
					t.Errorf("expected entry %v, got nil", test.expectedEntry)
					return
				}
				if actualEntry.Kanji != test.expectedEntry.Kanji {
					t.Errorf("expected kanji %s, got %s", test.expectedEntry.Kanji, actualEntry.Kanji)
				}
				if actualEntry.Kana != test.expectedEntry.Kana {
					t.Errorf("expected kana %s, got %s", test.expectedEntry.Kana, actualEntry.Kana)
				}
				if len(actualEntry.English) != len(test.expectedEntry.English) {
					t.Errorf("expected english entries of %v, got %v", len(test.expectedEntry.English), len(actualEntry.English))
				}
				if actualEntry.LinkedKanji != test.expectedEntry.LinkedKanji {
					t.Errorf("expected linked kanji %s, got %s", test.expectedEntry.LinkedKanji, actualEntry.LinkedKanji)
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

func TestParseLines_MultilineEnglishVerbPair(t *testing.T) {
	lines := []string{
		"## N5",
		"### verb pair",
		"- 降りる (おりる) - to get off; to descend; to go down; aka 下りる",
		"  * 降ろす (おろす) - to take down; to lower; to drop off (a passenger); to",
		"    unload (goods), to withdraw (money); aka 下ろす",
	}
	entries, err := cmd.ParseLines(lines)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	first := entries[0]
	if first.Kanji != "降りる" || first.Kana != "おりる" {
		t.Errorf("first entry: expected 降りる/おりる, got %s/%s", first.Kanji, first.Kana)
	}
	if first.LinkedKanji != "" {
		t.Errorf("first entry: expected empty LinkedKanji, got %s", first.LinkedKanji)
	}
	second := entries[1]
	if second.Kanji != "降ろす" || second.Kana != "おろす" {
		t.Errorf("second entry: expected 降ろす/おろす, got %s/%s", second.Kanji, second.Kana)
	}
	if second.LinkedKanji != "降りる" {
		t.Errorf("second entry: expected LinkedKanji 降りる, got %s", second.LinkedKanji)
	}
	expectedEnglish := []string{
		"to take down; to lower; to drop off (a passenger); to unload (goods), to withdraw (money); aka 下ろす",
	}
	if len(second.English) != len(expectedEnglish) {
		t.Fatalf("expected %d english entries, got %d (%v)", len(expectedEnglish), len(second.English), second.English)
	}
	for i, v := range expectedEnglish {
		if second.English[i] != v {
			t.Errorf("expected english[%d] = %q, got %q", i, v, second.English[i])
		}
	}
}

// - 階 (かい) - floor
// - インフォメーション - information
