package internal

import (
	"fmt"
	"strconv"
	"testing"
)

func TestCopyTestsToFolder(t *testing.T) {
	repoPath := "/Users/josephlyell/Documents/SubmissionGraderTests"
	repoName := "00000-SP22-C202-assignment-1-username2"
	language := "java"
	autoGraderPath := "/tmp/AutoGrader/00000-SP22/assignment-1/username2"
	err := copyTestsToFolder(repoPath, repoName, language, autoGraderPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	want := true
	got, err := exists(repoPath + "/" + autoGraderPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CopyTestsToFolder failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
