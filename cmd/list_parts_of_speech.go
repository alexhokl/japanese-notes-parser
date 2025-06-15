package cmd

import (
	"fmt"
	"sort"

	"github.com/alexhokl/helper/iohelper"
	"github.com/alexhokl/japanese-notes-parser/database"
	"github.com/alexhokl/japanese-notes-parser/parser"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type listPartsOfSpeechOptions struct {
	database string
}

var listPartsOfSpeechOpts listPartsOfSpeechOptions

// listPartsOfSpeechCmd represents the listPartsOfSpeech command
var listPartsOfSpeechCmd = &cobra.Command{
	Use:   "parts",
	Short: "List parts of speech from the specified (SQLite) database",
	RunE:  runListPartsOfSpeech,
}

func init() {
	listCmd.AddCommand(listPartsOfSpeechCmd)

	flags := listPartsOfSpeechCmd.Flags()
	flags.StringVarP(&listPartsOfSpeechOpts.database, "database", "d", "", "Database file to be read from")

	listPartsOfSpeechCmd.MarkFlagRequired("database")
}

func runListPartsOfSpeech(_ *cobra.Command, _ []string) error {
	if !iohelper.IsFileExist(listPartsOfSpeechOpts.database) {
		return fmt.Errorf("specified database does not exist")
	}

	db, err := gorm.Open(sqlite.Open(listPartsOfSpeechOpts.database), &gorm.Config{Logger: &database.CustomLogger{}})
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	partsOfSpeechList, err := database.ListPartsOfSpeech(db)
	if err != nil {
		return fmt.Errorf("failed to retrieve parts of speech: %w", err)
	}

	sort.Strings(partsOfSpeechList)

	for _, item := range partsOfSpeechList {
		if !parser.SimpleLevelRegex.MatchString(item) {
			fmt.Println(item)
		}
	}

	return nil
}
