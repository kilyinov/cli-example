# devctl

A CLI tool for automating common developer workflows — feature branches, JIRA tickets, pull requests, and more.

## Install

```bash
go install github.com/kilyinov/cli-example@latest
```

## Configuration

Create `~/.devctl.yaml`:

```yaml
jira:
  base_url: https://yourorg.atlassian.net
  email: you@example.com
  token: your-jira-api-token
  project: MYPROJ

git:
  default_branch: main
  branch_prefix: feature
  remote: origin
```

You can generate a JIRA API token at https://id.atlassian.com/manage-profile/security/api-tokens.

## Usage

### Create a JIRA ticket

```bash
devctl ticket create -s "Login page broken on Safari" -T Bug
```

With a description:

```bash
devctl ticket create -s "Add dark mode" -d "Support system-level dark mode preference" -p FRONTEND
```

With screenshots (you'll be prompted to capture screen regions interactively):

```bash
devctl ticket create -s "UI misalignment on dashboard" -T Bug --screenshot
```

### View a JIRA ticket

```bash
devctl ticket view MYPROJ-123
```

### Create a feature branch

```bash
devctl branch create add-dark-mode -t MYPROJ-456
```

This creates `feature/MYPROJ-456-add-dark-mode` from the latest remote default branch.

### Create a pull request

```bash
devctl pr create -t "Add dark mode support" --draft
```
