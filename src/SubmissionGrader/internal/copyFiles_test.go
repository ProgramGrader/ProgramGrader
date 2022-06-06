package internal

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestCopyTestsToFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	repoPath := wd + "/tmp"
	repoName := "00000-SP22-C202-assignment-1-username2"
	language := "java"
	err = copyTestsToFolder(repoPath, repoName, language, repoPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	want := true
	got, err := exists(repoPath + "/" + language)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CopyTestsToFolder failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	err = cleanup(repoPath + "/" + language)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func TestCopyTestResultsToFolder(t *testing.T) {
	//test after unit test results created
}

func TestCopyDocsToFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	repoName := "00000-SP22-C202-assignment-1-username2"
	repoPath := wd + "/tmp/" + repoName
	tempPath := "tmp"
	courseID := "00000"
	semesterID := "SP22"
	assignmentName := "C202-assignment-1"
	studentUserName := "username2"

	err = copyDocsToFolder(repoPath, repoName, tempPath, courseID, semesterID, assignmentName, studentUserName)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	want := true
	got, err := exists(repoPath + "/docs")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CopyDocsToFolder failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	err = cleanup(repoPath + "/docs")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
