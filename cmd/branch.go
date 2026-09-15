package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage feature branches",
}

var branchCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new feature branch from the default branch",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		ticketID, _ := cmd.Flags().GetString("ticket")
		prefix, _ := cmd.Flags().GetString("prefix")

		branchName := prefix + "/"
		if ticketID != "" {
			branchName += ticketID + "-"
		}
		branchName += name

		fmt.Printf("Creating branch: %s\n", branchName)
		// TODO: implement git branch creation
		return nil
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(branchCreateCmd)

	branchCreateCmd.Flags().StringP("ticket", "t", "", "JIRA ticket ID to include in branch name")
	branchCreateCmd.Flags().StringP("prefix", "p", "feature", "branch name prefix (feature, bugfix, hotfix)")
}
