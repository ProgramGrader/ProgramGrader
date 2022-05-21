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
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	return err
}

func copyTestResultsToFolder(repoPath string, repoName string, autoGraderPath string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/test-results/test", repoPath, repoName)
	cmd := exec.Command("cp", "-r", copyFromPath, autoGraderPath)
	cmd.Dir = fmt.Sprintf("%v/%v", repoPath, repoName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	//rename
	cmd = exec.Command("mv", autoGraderPath, "tests")
	err = cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	return err
}

func copyDocsToFolder(repoPath string, repoName string, tempPath string, courseID string, semesterID string, assignmentName string, studentUserName string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/docs/javadoc", repoPath, repoName)
	copyToPath := fmt.Sprintf("%v/AutoGrader/%v-%v/%v/%v", tempPath, courseID, semesterID, assignmentName, studentUserName)
	cmd := exec.Command("cp", "-r", copyFromPath, copyToPath)
	cmd.Dir = repoPath
	err := cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	//rename
	cmd = exec.Command("mv", copyToPath, "docs")
	err = cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	return err
}
