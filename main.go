package main

import (
	"fmt"
	"os"
	"strings"
)

// TODO: change binary name to git-flow

func main() {
	if len(os.Args) < 2 {
		exit(1, "Usage: git-flow command")
	}
	args := os.Args[2:]

	switch os.Args[1] {
	case "change":
		cmdChange(args...)
	case "sync":
		cmdSync(args...)
	default:
		exit(1, "Unrecognized command: %q", os.Args[1])
	}

	// TODO:
	// check --fast-forward / merge.ff=only
	// check pull.rebase = true
	// check rebase.autoStash
	// other

	// git config --global merge.ff only
	// git config --global pull.rebase true
	// git config --global rebase.autoStash true

	// git config --global push.default simple
	// git config --global init.defaultbranch main
	// git config --global core.editor vim -Nu NONE
	// git config --global diff.algorithm histogram
	// git config --global core.logallrefupdates true

	// diff.colorMoved = default
	// diff.colorMovedWS = ?
}

const prefixChangeBranch = "change/0/"

func cmdChange(args ...string) {
	if len(args) > 1 {
		exit(1, "invalid number of arguments")
	}
	if len(args) == 1 {
		targetBranch := prefixChangeBranch + args[0]
		checkoutOrCreate(targetBranch)
		return
	}
	ensureChangeBranch()
	shell("git", "commit", "--amend")
}

func ensureChangeBranch() {
	// ensure we're on a change branch
	if !strings.HasPrefix(currentBranch(), prefixChangeBranch) {
		exit(1, "current branch is not a change branch: %q", currentBranch())
	}
}

func checkoutOrCreate(branch string) {
	if !branchExists(branch) {
		createBranchAndCommit(branch)
		return
	}
	checkoutBranch(branch)
}

func branchExists(branch string) bool {
	branches, err := capture("git", "branch")
	if err != nil {
		exit(1, "could not get the list of branches: %v", err)
	}
	exist := strings.Contains(branches, fmt.Sprintf(" %s\n", branch))
	return exist
}

func createBranchAndCommit(name string) {
	shell("git", "checkout", "-t", "-b", name)
	shell("git", "commit", "--allow-empty", "-m", "created branch: "+name)
}

func checkoutBranch(name string) {
	shell("git", "checkout", name)
}

func currentBranch() string {
	branch, err := capture("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		exit(1, "could not get the current branch name: %v", err)
	}
	return strings.TrimSpace(branch)
}

func originBranch() string {
	branch, err := capture("git", "rev-parse", "--abbrev-ref", currentBranch()+"@{u}")
	if err != nil {
		exit(1, "could not get the origin branch name: %v", err)
	}
	return strings.TrimSpace(branch)
}

func cmdSync(args ...string) {
	if len(args) > 0 {
		exit(1, "invalid number of arguments")
	}
	ensureChangeBranch()
	// pull new commits from remote origin branch
	shell("git", "pull", "--rebase", "origin", originBranch())
}
