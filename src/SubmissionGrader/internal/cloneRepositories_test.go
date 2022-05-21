package internal

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCloneRepo(t *testing.T) {
	err := os.Setenv("repoName", "/ProgramGrader/00000-SP22-Assignment-1-username2")
	if err != nil {
		return
	}
	repoName := "test"
	// FIXME: err := cloneRepo()
	if err != nil {

		fmt.Printf(repoName)
	}
	cmd := exec.Command("cd", "tmp/ProgramGrader")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	cmd = exec.Command("ls")
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
}
