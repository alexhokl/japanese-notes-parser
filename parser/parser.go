package parser

import (
	"regexp"
	"strings"

	"github.com/alexhokl/helper/regexhelper"
)

func ParseHeaderLine(line string, levelRegex *regexp.Regexp, partOfSpeechRegex *regexp.Regexp) (level string, partOfSpeech string, err error) {
	levelCaptureGroups := regexhelper.FindNamedGroupMatchedStrings(levelRegex, line)
	if len(levelCaptureGroups) > 0 {
		level = levelCaptureGroups["level"]
		return
	}
	partOfSpeechCaptureGroups := regexhelper.FindNamedGroupMatchedStrings(partOfSpeechRegex, line)
	if len(partOfSpeechCaptureGroups) > 0 {
		partOfSpeech = partOfSpeechCaptureGroups["partOfSpeech"]
		return
	}
	return
}

func ParseEnglish(english string) []string {
	untrimmedList := strings.Split(english, "/")
	list := make([]string, len(untrimmedList))
	for i, item := range untrimmedList {
		list[i] = strings.TrimSpace(item)
	}
	return list
}
