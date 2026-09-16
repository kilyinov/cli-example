package cmd

import (
	"fmt"

	"github.com/kilyinov/cli-example/internal/prompt"
	"github.com/spf13/cobra"
)

var prCmd = &cobra.Command{
	Use:   "pr",
	Short: "Manage pull requests",
}

var prCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a pull request for the current branch",
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		draft, _ := cmd.Flags().GetBool("draft")
		base, _ := cmd.Flags().GetString("base")

		var err error

		if title == "" {
			title, err = prompt.StringRequired("PR title")
			if err != nil {
				return err
			}
		}

		if !cmd.Flags().Changed("base") {
			base, _ = prompt.StringWithDefault("Base branch", base)
		}

		if !cmd.Flags().Changed("draft") {
			draft = prompt.Confirm("Create as draft?")
		}

		mode := ""
		if draft {
			mode = " (draft)"
		}
		fmt.Printf("Creating PR%s: %s -> %s\n", mode, title, base)
		// TODO: implement GitHub/GitLab API call
		return nil
	},
}

func init() {
	rootCmd.AddCommand(prCmd)
	prCmd.AddCommand(prCreateCmd)

	prCreateCmd.Flags().StringP("title", "t", "", "pull request title")
	prCreateCmd.Flags().StringP("base", "b", "main", "base branch")
	prCreateCmd.Flags().BoolP("draft", "d", false, "create as draft PR")
}
