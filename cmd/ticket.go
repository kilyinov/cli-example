package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ticketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "Manage JIRA tickets",
}

var ticketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new JIRA ticket",
	RunE: func(cmd *cobra.Command, args []string) error {
		summary, _ := cmd.Flags().GetString("summary")
		project, _ := cmd.Flags().GetString("project")
		issueType, _ := cmd.Flags().GetString("type")

		fmt.Printf("Creating %s ticket in %s: %s\n", issueType, project, summary)
		// TODO: implement JIRA API call
		return nil
	},
}

var ticketViewCmd = &cobra.Command{
	Use:   "view [ticket-id]",
	Short: "View details of a JIRA ticket",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ticketID := args[0]
		fmt.Printf("Fetching ticket: %s\n", ticketID)
		// TODO: implement JIRA API call
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ticketCmd)
	ticketCmd.AddCommand(ticketCreateCmd)
	ticketCmd.AddCommand(ticketViewCmd)

	ticketCreateCmd.Flags().StringP("summary", "s", "", "ticket summary")
	ticketCreateCmd.Flags().StringP("project", "p", "", "JIRA project key")
	ticketCreateCmd.Flags().StringP("type", "T", "Task", "issue type (Task, Bug, Story)")
	ticketCreateCmd.MarkFlagRequired("summary")
	ticketCreateCmd.MarkFlagRequired("project")
}
