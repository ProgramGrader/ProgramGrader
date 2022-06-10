package internal

import (
	"fmt"
	"os"
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

func cleanup(path string) error {
	err := os.RemoveAll(path)
	return err
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	tempPath := wd + "/tmp"
	exists, err := exists(tempPath)
	if exists {
		err = os.RemoveAll(tempPath)
		if err != nil {
			fmt.Printf("Error: %s", err)
		}
	}
	err = os.Mkdir("tmp", 0700)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	githubToken := os.Getenv("GITHUBTOKEN")
	orgName := "ProgramGrader"
	repoName := "00000-SP22-C202-assignment-1-username2"

	err = cloneRepo(githubToken, orgName, repoName, tempPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
}
