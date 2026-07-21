package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alexhokl/helper/iohelper"
	"github.com/alexhokl/helper/regexhelper"
	"github.com/alexhokl/japanese-notes-parser/database"
	"github.com/alexhokl/japanese-notes-parser/parser"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type createOptions struct {
	file      string
	database  string
	overwrite bool
}

var createOpts createOptions

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a SQLite database by parsing a note file",
	RunE:  runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)

	flags := createCmd.Flags()
	flags.StringVarP(&createOpts.file, "file", "f", "", "Input note file")
	flags.StringVarP(&createOpts.database, "database", "d", "", "Output database file")
	flags.BoolVar(&createOpts.overwrite, "overwrite", false, "Overwrite the database file if it exists")

	createCmd.MarkFlagRequired("file")
	createCmd.MarkFlagRequired("database")
}

func runCreate(_ *cobra.Command, _ []string) error {
	// check if note file exists
	if !iohelper.IsFileExist(createOpts.file) {
		return fmt.Errorf("input file %s does not exist", createOpts.file)
	}

	// checks if database file exists
	if !createOpts.overwrite && iohelper.IsFileExist(createOpts.database) {
		return fmt.Errorf("output database file %s already exists", createOpts.database)
	}
	if createOpts.overwrite && iohelper.IsFileExist(createOpts.database) {
		if err := os.Remove(createOpts.database); err != nil {
			return fmt.Errorf("failed to delete output database file: %w", err)
		}
		fmt.Printf("output database file %s has been removed\n", createOpts.database)
	}

	db, err := gorm.Open(sqlite.Open(createOpts.database), &gorm.Config{Logger: database.NewCustomLogger()})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Migrate the schema
	if err := database.AutoMigrate(db); err != nil {
		return fmt.Errorf("failed to migrate database schema: %w", err)
	}

	lines, err := iohelper.ReadLinesFromFile(createOpts.file)
	if err != nil {
		return fmt.Errorf("failed to read lines from file: %w", err)
	}

	entries, err := ParseLines(lines)
	if err != nil {
		return err
	}

	// TODO: bulk upsert with transaction
	for _, entry := range entries {
		if err := database.Upsert(db, entry); err != nil {
			return fmt.Errorf("failed to upsert entry: %w", err)
		}
	}

	fmt.Println("database created")
	return nil
}

func ParseLines(lines []string) ([]*database.Entry, error) {
	cachedLevel := ""
	cachedPartOfSpeech := ""
	previousKanji := ""
	var entries []*database.Entry
	var lastEntry *database.Entry
	for _, line := range lines {
		entry, level, partOfSpeech, continuation, err := ParseLine(
			line,
			parser.PointRegex,
			parser.SubPointRegex,
			parser.JapaneseRegex,
			parser.LevelRegex,
			parser.PartOfSpeechRegex,
			cachedLevel,
			cachedPartOfSpeech,
			previousKanji,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line: %w", err)
		}
		if continuation != "" && lastEntry != nil {
			if len(lastEntry.English) == 0 {
				lastEntry.English = []string{continuation}
			} else {
				lastEntry.English[len(lastEntry.English)-1] = lastEntry.English[len(lastEntry.English)-1] + " " + continuation
			}
			continue
		}
		if entry == nil {
			cachedLevel = level
			cachedPartOfSpeech = partOfSpeech
			previousKanji = ""
			lastEntry = nil
			continue
		}
		previousKanji = entry.Kanji
		entries = append(entries, entry)
		lastEntry = entry
	}
	return entries, nil
}

func ParseLine(line string, pointRegex, subPointRegex, japaneseRegex, levelRegex, partOfSpeechRegex *regexp.Regexp, cachedLevel string, cachedPartOfSpeech string, previousKanji string) (*database.Entry, string, string, string, error) {
	isSubPoint := false
	captureGroups := regexhelper.FindNamedGroupMatchedStrings(pointRegex, line)
	if len(captureGroups) == 0 {
		// checks if it is a sub point before checking for level or part of speech
		subPointCaptureGroups := regexhelper.FindNamedGroupMatchedStrings(subPointRegex, line)
		if len(subPointCaptureGroups) == 0 {
			// checks for level or part of speech
			level, partOfSpeech, err := parser.ParseHeaderLine(line, levelRegex, partOfSpeechRegex)
			if err != nil {
				// ignore
				return nil, cachedLevel, cachedPartOfSpeech, "", nil
			}
			if level != "" {
				return nil, level, cachedPartOfSpeech, "", nil
			}
			if partOfSpeech != "" {
				return nil, cachedLevel, partOfSpeech, "", nil
			}
			// no regex matched: treat non-empty line as a continuation of
			// the previous entry's English text. Empty/whitespace-only
			// lines are ignored.
			trimmed := strings.TrimSpace(line)
			if trimmed == "" {
				return nil, cachedLevel, cachedPartOfSpeech, "", nil
			}
			return nil, cachedLevel, cachedPartOfSpeech, trimmed, nil
		}
		isSubPoint = true
		captureGroups = subPointCaptureGroups
	}
	japanese := captureGroups["japanese"]
	if japanese == "" {
		return nil, cachedLevel, cachedPartOfSpeech, "", nil
	}
	english := captureGroups["english"]
	japaneseCaptureGroups := regexhelper.FindNamedGroupMatchedStrings(japaneseRegex, japanese)
	kanji := strings.TrimSpace(japaneseCaptureGroups["kanji"])
	kana := strings.TrimSpace(japaneseCaptureGroups["kana"])
	if kanji == "" && kana == "" {
		// assuming japanese is katakana
		kana = japanese
	}

	linkedKanji := ""
	if isSubPoint && previousKanji != "" {
		linkedKanji = previousKanji
	}

	entry := &database.Entry{
		Kanji:       kanji,
		Kana:        kana,
		English:     parser.ParseEnglish(english),
		Labels:      []string{cachedLevel, cachedPartOfSpeech},
		LinkedKanji: linkedKanji,
	}

	return entry, cachedLevel, cachedPartOfSpeech, "", nil
}
