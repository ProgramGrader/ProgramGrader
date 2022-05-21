package internal

import (
	"os"
	"testing"
)

func TestGetEnvVar(t *testing.T) {
	err := os.Setenv("LANGUAGE", "java")
	if err != nil {
		return
	}
	want := "java"
	got := GetEnvVar("LANGUAGE")
	if want != got {
		t.Errorf("GetEnvVar failed. want: %s; got: %s", want, got)
	}
}
