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
