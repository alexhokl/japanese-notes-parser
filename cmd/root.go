package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "japanese-notes-parser",
	Short:        "A parser for creating database entries from Japanese notes",
	SilenceUsage: true,
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
