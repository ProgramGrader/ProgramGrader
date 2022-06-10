package internal

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	githubToken := os.Getenv("GITHUBTOKEN")
	orgName := "ProgramGrader"
	repoName := "00000-SP22-C202-assignment-1-username2"
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	err = os.Mkdir("tmp/cloneTest", 0700)
	repoPath := wd + "/tmp/cloneTest"
	err = cloneRepo(githubToken, orgName, repoName, repoPath)
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
	err = cleanup(repoPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}

func TestCloneConfigRepo(t *testing.T) {
	githubToken := os.Getenv("GITHUBTOKEN")
	orgName := "ProgramGrader"
	config := "GraderConfigTest"
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	repoPath := wd + "/tmp"
	err = cloneConfigRepo(githubToken, orgName, config, repoPath)
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
	err = cleanup(repoPath + "/" + config)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
