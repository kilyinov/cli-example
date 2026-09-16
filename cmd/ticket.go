package cmd

import (
	"fmt"
	"os"

	"github.com/kilyinov/cli-example/internal/config"
	"github.com/kilyinov/cli-example/internal/jira"
	"github.com/kilyinov/cli-example/internal/screenshot"
	"github.com/spf13/cobra"
)

var ticketCmd = &cobra.Command{
	Use:   "ticket",
	Short: "Manage JIRA tickets",
}

var ticketCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new JIRA ticket",
	RunE:  runTicketCreate,
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

func runTicketCreate(cmd *cobra.Command, args []string) error {
	summary, _ := cmd.Flags().GetString("summary")
	project, _ := cmd.Flags().GetString("project")
	issueType, _ := cmd.Flags().GetString("type")
	description, _ := cmd.Flags().GetString("description")
	attachScreenshot, _ := cmd.Flags().GetBool("screenshot")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if project == "" {
		project = cfg.JIRA.Project
	}
	if project == "" {
		return fmt.Errorf("project is required: use --project or set jira.project in config")
	}

	client, err := jira.NewClient(cfg.JIRA)
	if err != nil {
		return fmt.Errorf("initializing JIRA client: %w", err)
	}

	var screenshotPath string
	if attachScreenshot {
		fmt.Println("Select a screen region to capture...")
		path, err := screenshot.Capture()
		if err != nil {
			return fmt.Errorf("capturing screenshot: %w", err)
		}
		screenshotPath = path
		defer os.Remove(screenshotPath)
		fmt.Println("Screenshot captured.")
	}

	fmt.Printf("Creating %s in %s: %s\n", issueType, project, summary)

	issue, err := client.CreateIssue(jira.CreateIssueRequest{
		Project:     project,
		Summary:     summary,
		Description: description,
		IssueType:   issueType,
	})
	if err != nil {
		return fmt.Errorf("creating ticket: %w", err)
	}

	fmt.Printf("Created %s (%s/browse/%s)\n", issue.Key, cfg.JIRA.BaseURL, issue.Key)

	if screenshotPath != "" {
		fmt.Printf("Attaching screenshot to %s...\n", issue.Key)
		if err := client.AddAttachment(issue.Key, screenshotPath); err != nil {
			return fmt.Errorf("attaching screenshot: %w", err)
		}
		fmt.Println("Screenshot attached.")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(ticketCmd)
	ticketCmd.AddCommand(ticketCreateCmd)
	ticketCmd.AddCommand(ticketViewCmd)

	ticketCreateCmd.Flags().StringP("summary", "s", "", "ticket summary")
	ticketCreateCmd.Flags().StringP("description", "d", "", "ticket description")
	ticketCreateCmd.Flags().StringP("project", "p", "", "JIRA project key (overrides config)")
	ticketCreateCmd.Flags().StringP("type", "T", "Task", "issue type (Task, Bug, Story)")
	ticketCreateCmd.Flags().Bool("screenshot", false, "capture and attach a screenshot via screencapture")
	ticketCreateCmd.MarkFlagRequired("summary")
}
