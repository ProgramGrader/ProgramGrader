package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	os.Setenv("LANGUAGE", "java")
	want := "java"
	got := GetEnvVar("LANGUAGE")
	if want != got {
		t.Errorf("GetEnvVar failed. want: %s; got: %s", want, got)
	}
}

func TestCloneRepo(t *testing.T) {
	os.Setenv("repoName", "/ProgramGrader/00000-SP22-Assignment-1-username2")
	cloneRepo()
	cmd := exec.Command("cd", "tmp/ProgramGrader")
	err := cmd.Run()
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
