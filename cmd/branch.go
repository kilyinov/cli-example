package cmd

import (
	"fmt"

	"github.com/kilyinov/cli-example/internal/config"
	"github.com/kilyinov/cli-example/internal/git"
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
	RunE:  runBranchCreate,
}

func runBranchCreate(cmd *cobra.Command, args []string) error {
	if !git.IsRepo() {
		return fmt.Errorf("not a git repository")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	name := args[0]
	ticketID, _ := cmd.Flags().GetString("ticket")
	prefix, _ := cmd.Flags().GetString("prefix")

	branchName := prefix + "/"
	if ticketID != "" {
		branchName += ticketID + "-"
	}
	branchName += name

	if git.BranchExists(branchName) {
		return fmt.Errorf("branch %q already exists", branchName)
	}

	dirty, err := git.HasUncommittedChanges()
	if err != nil {
		return fmt.Errorf("checking working tree: %w", err)
	}
	if dirty {
		return fmt.Errorf("you have uncommitted changes; commit or stash them first")
	}

	remote := cfg.Git.Remote
	base := cfg.Git.DefaultBranch
	startPoint := remote + "/" + base

	fmt.Printf("Fetching %s...\n", remote)
	if err := git.Fetch(remote); err != nil {
		return fmt.Errorf("fetching %s: %w", remote, err)
	}

	fmt.Printf("Creating branch %s from %s\n", branchName, startPoint)
	if err := git.CreateBranchAndSwitch(branchName, startPoint); err != nil {
		return fmt.Errorf("creating branch: %w", err)
	}

	fmt.Printf("Switched to new branch %q\n", branchName)
	return nil
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.AddCommand(branchCreateCmd)

	branchCreateCmd.Flags().StringP("ticket", "t", "", "JIRA ticket ID to include in branch name")
	branchCreateCmd.Flags().StringP("prefix", "p", "feature", "branch name prefix (feature, bugfix, hotfix)")
}
