package internal

import (
	"fmt"
	"os/exec"
)

func copyTestsToFolder(repoPath string, repoName string, language string, autoGraderPath string) error {
	copyFromPath := fmt.Sprintf("%v/%v/src/test/%v", repoPath, repoName, language)
	cmd := exec.Command("cp", "-r", copyFromPath, autoGraderPath)
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func copyTestResultsToFolder(repoPath string, repoName string, autoGraderPath string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/test-results", repoPath, repoName)
	cmd := exec.Command("cp", "-r", copyFromPath, autoGraderPath)
	cmd.Dir = fmt.Sprintf("%v/%v", repoPath, repoName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func copyDocsToFolder(repoPath string, repoName string, courseID string, semesterID string, assignmentName string, studentUserName string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/docs", repoPath, repoName)
	copyToPath := fmt.Sprintf("%v-%v-%v-%v-docs", courseID, semesterID, assignmentName, studentUserName)
	cmd := exec.Command("mkdir", copyToPath)
	cmd.Dir = repoPath + "/AutoGrader"
	err := cmd.Run()
	if err != nil {
		return err
	}
	docsPath := repoPath + "/AutoGrader/" + copyToPath
	cmd = exec.Command("cp", "-r", copyFromPath, docsPath)
	cmd.Dir = repoPath + "/AutoGrader"
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error in copy")
		return err
	}
	return err
}
