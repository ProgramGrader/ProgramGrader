package internal

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TestCloneRepo(t *testing.T) {
	githubToken := os.Getenv("GITHUBTOKEN")
	orgName := "ProgramGrader"
	repoName := "00000-SP22-C202-assignment-1-username2"
	repoPath := "/Users/josephlyell/Documents/SubmissionGraderTests"
	err := cloneRepo(githubToken, orgName, repoName, repoPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	want := true
	got, err := exists(repoPath + "/" + repoName)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CloneRepo failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
}

func TestCloneConfigRepo(t *testing.T) {
	githubToken := os.Getenv("GITHUBTOKEN")
	orgName := "ProgramGrader"
	config := "GraderConfigTest"
	repoPath := "/Users/josephlyell/Documents/SubmissionGraderTests"
	err := cloneConfigRepo(githubToken, orgName, config, repoPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	want := true
	got, err := exists(repoPath + "/" + config)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	if want != got {
		t.Errorf("CloneConfigRepo failed. want: %s; got: %s", strconv.FormatBool(want), strconv.FormatBool(got))
	}
}
