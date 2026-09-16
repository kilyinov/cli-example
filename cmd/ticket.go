package cmd

import (
	"fmt"
	"os"

	"github.com/kilyinov/cli-example/internal/config"
	"github.com/kilyinov/cli-example/internal/jira"
	"github.com/kilyinov/cli-example/internal/prompt"
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
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var ticketID string
		if len(args) > 0 {
			ticketID = args[0]
		} else {
			var err error
			ticketID, err = prompt.StringRequired("Ticket ID")
			if err != nil {
				return err
			}
		}
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

	var err error

	if summary == "" {
		summary, err = prompt.StringRequired("Summary")
		if err != nil {
			return err
		}
	}

	if description == "" && !cmd.Flags().Changed("description") {
		description, _ = prompt.String("Description (optional)")
	}

	if !cmd.Flags().Changed("type") {
		issueType, _ = prompt.StringWithDefault("Issue type", issueType)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	if project == "" {
		project = cfg.JIRA.Project
	}
	if project == "" {
		project, err = prompt.StringRequired("Project key")
		if err != nil {
			return err
		}
	}

	if !cmd.Flags().Changed("screenshot") {
		attachScreenshot = prompt.Confirm("Attach screenshot(s)?")
	}

	client, err := jira.NewClient(cfg.JIRA)
	if err != nil {
		return fmt.Errorf("initializing JIRA client: %w", err)
	}

	var screenshotPaths []string
	if attachScreenshot {
		for {
			fmt.Printf("Select a screen region to capture (screenshot %d)...\n", len(screenshotPaths)+1)
			path, err := screenshot.Capture()
			if err != nil {
				return fmt.Errorf("capturing screenshot: %w", err)
			}
			screenshotPaths = append(screenshotPaths, path)
			fmt.Printf("Screenshot %d captured.\n", len(screenshotPaths))

			if !prompt.Confirm("Capture another screenshot?") {
				break
			}
		}
		defer func() {
			for _, p := range screenshotPaths {
				os.Remove(p)
			}
		}()
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

	for i, path := range screenshotPaths {
		fmt.Printf("Attaching screenshot %d/%d to %s...\n", i+1, len(screenshotPaths), issue.Key)
		if err := client.AddAttachment(issue.Key, path); err != nil {
			return fmt.Errorf("attaching screenshot %d: %w", i+1, err)
		}
	}
	if len(screenshotPaths) > 0 {
		fmt.Printf("%d screenshot(s) attached.\n", len(screenshotPaths))
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
}
