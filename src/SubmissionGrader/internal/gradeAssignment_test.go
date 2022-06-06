package internal

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	err := os.Setenv("LANGUAGE", "java")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	want := "java"
	got := GetEnvVar("LANGUAGE")
	if want != got {
		t.Errorf("GetEnvVar failed. want: %s; got: %s", want, got)
	}
}

func TestGradeAssignment(t *testing.T) {
	err := os.Setenv("REPOFULLNAME", "00000-SP22-C202-assignment-1-username2")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = os.Setenv("LANGUAGE", "java")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = os.Setenv("ORGNAME", "ProgramGrader")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = os.Setenv("TEACHERUNITTESTSENABLED", "FALSE")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = os.Setenv("GRADEDOCSENABLED", "FALSE")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = os.Setenv("STUDENTTESTSENABLED", "FALSE")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	err = os.Setenv("NONCODESUBMISSIONSENABLED", "FALSE")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}
	repoName := os.Getenv("REPOFULLNAME")
	err = os.Setenv("TEMPPATH", wd+"/"+repoName)

	var signalVar *GoSafeVar[bool]
	var wg *sync.WaitGroup
	fmt.Printf("Setting up os signal watcher.")
	WatchForSignal(signalVar)
	fmt.Printf("Starting Grading assignment.")
	go GradeAssignment(signalVar, wg)

}
