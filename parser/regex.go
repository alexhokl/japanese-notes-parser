package parser

import "regexp"

var (
	LevelRegex        = regexp.MustCompile(`^## (?P<level>\w+)`)
	PartOfSpeechRegex = regexp.MustCompile(`^### (?P<partOfSpeech>(.*))`)
	PointRegex        = regexp.MustCompile(`- (?P<japanese>\W+) - (?P<english>(.*))`)
	JapaneseRegex     = regexp.MustCompile(`(?P<kanji>\W+) \((?P<kana>\W+)\)`)
	SimpleLevelRegex  = regexp.MustCompile(`^N\d`)
)
