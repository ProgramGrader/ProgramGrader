package internal

import (
	"fmt"
	"os/exec"
)

func cloneRepo(githubToken string, orgName string, repoName string, repoPath string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v/%v.git", githubToken, orgName, repoName))
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func cloneConfigRepo(githubToken string, orgName string, config string, repoPath string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v/%v.git", githubToken, orgName, config))
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}
