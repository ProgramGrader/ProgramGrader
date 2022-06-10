package internal

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

//TODO: internal/tmp to gitignore file

func TestCreateDocs(t *testing.T) {
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
	want := true
	got, err := exists(repoPath + "/" + repoName + "/build/docs/javadoc")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CreateDocs failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	err = cleanup(repoPath + "/" + repoName + "/build/docs/javadoc")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func TestRunUnitTests(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	repoName := "00000-SP22-C202-assignment-1-username2"
	repoPath := wd + "/tmp"
	err = runUnitTests(repoPath, repoName)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	want := true
	got, err := exists(repoPath + "/" + repoName + "/build/test-results")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("RunUnitTests failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	err = cleanup(repoPath + "/" + repoName + "/build/test-results")
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

//func TestHandleNonCodeSubmissions(t *testing.T) {
//	wd, err := os.Getwd()
//	if err != nil {
//		fmt.Printf("Error: %s", err)
//	}
//	repoName := "00000-SP22-C202-assignment-1-username2"
//	repoPath := wd + "/tmp"
//	autoGraderPath := repoPath + "AutoGrader" + repoName
//	err = handleNonCodeSubmissions(repoPath, repoName, autoGraderPath)
//	if err != nil {
//		fmt.Printf("Error: %s", err)
//	}
//	want := true
//	got, err := exists("TODO add path to noncode submissions")
//	if err != nil {
//		fmt.Printf("Error: %s", err)
//	}
//	if want != got {
//		t.Errorf("HandleNonCodeSubmissions failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
//	}
//	//add cleanup
//}