package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mmclsntr/lineworks-cli/auth"
)

var rootCmd = &cobra.Command{
	Use:   "lineworks",
	Short: "Command line tool for LINE WORKS API",
}

var listProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles := auth.ListConfigProfiles()
		for _, p := range profiles {
			fmt.Println(p)
		}
		return nil
	},
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(listProfilesCmd)
}
