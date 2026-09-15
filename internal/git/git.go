package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(out))
	if err != nil {
		return output, fmt.Errorf("git %s: %s", args[0], output)
	}
	return output, nil
}

func IsRepo() bool {
	_, err := run("rev-parse", "--git-dir")
	return err == nil
}

func Fetch(remote string) error {
	_, err := run("fetch", remote)
	return err
}

func CurrentBranch() (string, error) {
	return run("rev-parse", "--abbrev-ref", "HEAD")
}

func BranchExists(name string) bool {
	_, err := run("rev-parse", "--verify", name)
	return err == nil
}

func HasUncommittedChanges() (bool, error) {
	out, err := run("status", "--porcelain")
	if err != nil {
		return false, err
	}
	return out != "", nil
}

func CreateBranchAndSwitch(branchName, startPoint string) error {
	_, err := run("checkout", "-b", branchName, startPoint)
	return err
}
