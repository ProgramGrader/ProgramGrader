package internal

import (
	"fmt"
	"os/exec"
)

func createDocs(repoPath string, repoName string) error {
	cmd := exec.Command("gradle", "javadoc")
	cmd.Dir = fmt.Sprintf("%v/%v", repoPath, repoName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func runUnitTests(repoPath string, repoName string) error {
	cmd := exec.Command("gradle", "test")
	cmd.Dir = fmt.Sprintf("%v/%v", repoPath, repoName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func handleNonCodeSubmissions(repoPath string, repoName string, autoGraderPath string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/test-results/test", repoPath, repoName)
	cmd := exec.Command("cp", "r", copyFromPath, autoGraderPath)
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}
