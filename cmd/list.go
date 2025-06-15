package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List subjects of subcommand",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify the subject to be listed")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
