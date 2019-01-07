package services

import (
	"bufio"
	"devlab/lib/errors"
	"devlab/lib/exec"
	"devlab/lib/files"
	"devlab/lib/logger"
	"os"
	"strings"
	"time"
)

// Clone pulls service repository from remote server to local machine
func Clone(contextServicesDir string, serviceName string, githabHost string, relativeGithubPath string) {
	isServiceDirExists, err := files.IsExists(contextServicesDir + "/" + serviceName)
	if errors.CheckAndReturnIfError(err) {
		return
	}

	servicesDir, err := files.AbsolutePath(contextServicesDir)

	if relativeGithubPath == "" {
		relativeGithubPath = serviceName + ".git"
	}

	if !isServiceDirExists {
		logger.Text("Cloning from github: " + githabHost + relativeGithubPath + " ...")
		exec.GitCommand(servicesDir, "git clone "+githabHost+relativeGithubPath+" "+serviceName)
	}
}

// RefreshGitRepo refreshes git repo service (refreshes service repo,
// commits or staches changes and checkout to context branch)
func RefreshGitRepo(contextServicesDir, serviceName,
	contextServiceBranch, baseServiceBranch string, context map[string]map[string]string) {
	serviceDir, _ := files.AbsolutePath(contextServicesDir + "/" + serviceName)

	if CheckRepoChanges(contextServicesDir, serviceName) {
		logger.Warn("There are some not commited changes")
		logger.Text("Please, choose action: ")
		logger.Text(" (1) commit changes ")
		logger.Text(" (2) stash changes ")
		logger.Text(" (3) nothing to do ")

		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		action := input.Text()

		switch action {
		case "1":
			logger.Text("Commiting changes ")
			CommitChanges(contextServicesDir, serviceName)
			break
		case "2":
			logger.Text("Stashing changes ")
			exec.GitCommand(serviceDir, "git stash")
			break
		case "3":
		default:
		}
	}

	CheckoutOrCreate(contextServicesDir,
		serviceName, contextServiceBranch, baseServiceBranch, context)
}

// CheckoutOrCreate creates new git branch or checkout it if branch exists
func CheckoutOrCreate(contextServicesDir string, serviceName string,
	checkoutBranch string, baseBranch string, context map[string]map[string]string) {
	serviceDir, _ := files.AbsolutePath(contextServicesDir + "/" + serviceName)
	numCheckoutBranchExistsAsRemoteStr, _ := exec.GitCommand(serviceDir, "git branch -r | grep -c "+checkoutBranch)
	currentBranch, _ := exec.GitCommand(serviceDir, "git symbolic-ref --short HEAD")

	isCheckoutBranchExistsAsRemote := strings.TrimSpace(numCheckoutBranchExistsAsRemoteStr) != "0"

	// checkoutBranch exists as remote
	if isCheckoutBranchExistsAsRemote {
		if currentBranch != checkoutBranch {
			logger.Text("Checking out to remote branch '" + checkoutBranch + "' \n")
			exec.GitCommand(serviceDir, "git checkout "+checkoutBranch)
		}

		diffLocalAndRemoteBranches, _ := exec.GitCommand(serviceDir, "git diff "+checkoutBranch+" origin/"+checkoutBranch+" --stat")

		if diffLocalAndRemoteBranches != "" {
			logger.Text("Remote branch is differnt with local branch. \n")

			BackupCurrentBranchIfNeed(contextServicesDir, serviceName, checkoutBranch)

			logger.Info("Refreshing local service folder from remote branch origin/%s.", checkoutBranch)
			exec.GitCommand(serviceDir, "git fetch origin && git reset --hard origin/"+checkoutBranch)
		}

		return
	}

	// checkoutBranch not exists as remote
	if currentBranch != checkoutBranch {
		numCheckoutBranchExistsAsLocalStr, _ := exec.GitCommand(serviceDir, "git branch | grep -c "+checkoutBranch)
		isCheckoutBranchExistsAsLocal := strings.TrimSpace(numCheckoutBranchExistsAsLocalStr) != "0"

		if isCheckoutBranchExistsAsLocal {
			logger.Text("Checking out to local branch '" + checkoutBranch + "' \n")
			exec.GitCommand(serviceDir, "git checkout "+checkoutBranch)
		} else {
			logger.Info("Creating new local branch '%s' from remote branch 'origin/%s'", checkoutBranch, baseBranch)
			exec.GitCommand(serviceDir, "git fetch origin && git checkout "+baseBranch+" && git reset --hard origin/"+baseBranch)
			exec.GitCommand(serviceDir, "git checkout -b "+checkoutBranch)
		}
	}

	if checkoutBranch != "master" && checkoutBranch != "develop" && context["git"]["auto-push"] != "false" {
		logger.Info("Pushing new local branch '%s' to remote git server", checkoutBranch)
		exec.GitCommand(serviceDir, "git push origin "+checkoutBranch)
	}
}

// BackupCurrentBranchIfNeed  backups current service branch
func BackupCurrentBranchIfNeed(contextServicesDir string, serviceName string, checkoutBranch string) {
	serviceDir, _ := files.AbsolutePath(contextServicesDir + "/" + serviceName)

	nowUnixTime := time.Now().Unix()
	logger.Info("Remote branch is differnt with local branch. \n Would you like to backup current version of local branch (create local branch '%s.backup.%s'), y|N ?", checkoutBranch, string(nowUnixTime))

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	answer := input.Text()

	if answer == "y" || answer == "Y" {
		logger.Info("Creating '%s.backup.%s' branch from current branch '%s'", checkoutBranch, string(nowUnixTime), checkoutBranch)
		exec.GitCommand(serviceDir, "git checkout -b "+checkoutBranch+".backup."+string(nowUnixTime))
	}
}

// CheckRepoChanges  checks if service repo has not commited changes
func CheckRepoChanges(contextServicesDir string, serviceName string) bool {
	serviceDir, _ := files.AbsolutePath(contextServicesDir + "/" + serviceName)
	repoChanges, _ := exec.GitCommand(serviceDir, "git status -s")
	return string(repoChanges) != ""
}

// CommitChanges commits service changes
func CommitChanges(contextServicesDir string, serviceName string) {
	serviceDir, _ := files.AbsolutePath(contextServicesDir + "/" + serviceName)

	logger.Text("Enter commit message: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	message := input.Text()

	if message != "" {
		exec.GitCommand(serviceDir, "git add --all && git commit -m '"+message+"'")
	} else {
		logger.Warn("You have not entered the commit message.")
		logger.Warn("Changes will not be commited! It will be stashed.")

		exec.GitCommand(serviceDir, "git stash")
	}
}
