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
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	repoName := "00000-SP22-C202-assignment-1-username2"
	repoPath := wd + "/tmp"
	err = runUnitTests(repoPath, repoName)
	if err != nil {
		fmt.Printf("Error running tests: %s", err)
	}
	autoGraderPath := repoPath + "/AutoGrader"
	err = copyTestResultsToFolder(repoPath, repoName, autoGraderPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	want := true
	got, err := exists(repoPath + "/AutoGrader/test")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CopyTestResultsToFolder failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	err = cleanup(fmt.Sprintf("%s/AutoGrader/test", repoPath))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func TestCopyDocsToFolder(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	repoName := "00000-SP22-C202-assignment-1-username2"
	repoPath := wd + "/tmp"
	err = createDocs(repoPath, repoName)
	if err != nil {
		fmt.Printf("Error creating docs: %s", err)
	}
	courseID := "00000"
	semesterID := "SP22"
	assignmentName := "C202-assignment-1"
	studentUserName := "username2"
	err = copyDocsToFolder(repoPath, repoName, courseID, semesterID, assignmentName, studentUserName)
	if err != nil {
		fmt.Printf("Error in copyDocsToFolder: %s", err)
	}
	want := true
	got, err := exists(fmt.Sprintf("%v/AutoGrader/%v-%v-%v-%v-docs/docs", repoPath, courseID, semesterID, assignmentName, studentUserName))
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CopyDocsToFolder failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	//err = cleanup(fmt.Sprintf("%v/AutoGrader/%v-%v-%v-%v-docs", repoPath, courseID, semesterID, assignmentName, studentUserName))
	err = cleanup(fmt.Sprintf("%v/AutoGrader/", repoPath))

	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
