package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type signals struct {
	signalThrown bool
	mu           sync.Mutex
}

func languageSwitch(sigs signals) {
	var (
		err                       error
		githubToken               = GetEnvVar("GITHUBTOKEN")
		language                  = GetEnvVar("LANGUAGE")
		repoName                  = GetEnvVar("REPOFULLNAME")
		teacherUnitTestsEnabled   = GetEnvVar("TEACHER_UNIT_TESTS")
		courseType                = GetEnvVar("COURSE_TYPE")
		assignmentName            = GetEnvVar("ASSIGNMENT_NAME")
		gradeDocsEnabled          = GetEnvVar("GRADE_DOCS")
		studentTestsEnabled       = GetEnvVar("STUDENT_TESTS_ENABLED")
		nonCodeSubmissionsEnabled = GetEnvVar("NONCODE_SUBMISSIONS")
		courseID                  = GetEnvVar("COURSE_ID")
		semesterID                = GetEnvVar("SEMESTER_ID")
		studentUserName           = GetEnvVar("SUDENT_USERNAME")
		config                    = "AutoGraderConfig"
		repoPath                  = "/opt/gradle"
		orgName                   = "IUS-CS"
		tempPath                  = "/tmp"
		testsPath                 = fmt.Sprintf("%v/%v/src/%v/%v/current/tests", tempPath, config, courseType, assignmentName)
		autoGraderPath            = fmt.Sprintf("%v/AutoGrader/%v-%v/%v/%v", tempPath, courseID, semesterID, assignmentName, studentUserName)
	)
	switch language {
	case "java":
		{
			//clone repo
			err = cloneRepo(githubToken, orgName, repoName, repoPath)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			if teacherUnitTestsEnabled == "TRUE" {
				err = cloneConfigRepo(githubToken, orgName, config, repoPath)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
				err = copyTestsToFolder(repoPath, repoName, language, testsPath)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}
			if gradeDocsEnabled == "TRUE" {
				err = createDocs(repoPath, repoName)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
				err = copyDocsToFolder(repoPath, repoName, tempPath, courseID, semesterID, assignmentName, studentUserName)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}
			if studentTestsEnabled == "TRUE" || teacherUnitTestsEnabled == "TRUE" {
				err = runUnitTests(repoPath, repoName)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
				err = copyTestResultsToFolder(repoPath, repoName, autoGraderPath)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}

			if nonCodeSubmissionsEnabled == "TRUE" {
				err = handleNonCodeSubmissions(repoPath, repoName, autoGraderPath)
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}
		}
	case "c++":
		{

		}

	case "python":
		{

		}
	default:
		{
			//log and deal with unsupported language
			fmt.Println(language + " is not supported. Nothing has been graded.")
		}
	}
	return
}

func GetEnvVar(envName string) string {
	if envName != "" {
		ret := os.Getenv(envName)
		return ret
	}
	fmt.Printf(envName + "environment variable was not found.")
	return "error"
}

func cloneRepo(githubToken string, orgName string, repoName string, repoPath string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v/%v.git", githubToken, orgName, repoName))
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func cloneConfigRepo(githubToken string, orgName string, config string, repoPath string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v/%v.git", githubToken, orgName, config))
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func copyTestsToFolder(repoPath string, repoName string, language string, autoGraderPath string) error {
	copyFromPath := fmt.Sprintf("%v/%v/src/test/%v", repoPath, repoName, language)
	cmd := exec.Command("cp", "-r", copyFromPath, autoGraderPath)
	cmd.Dir = repoPath
	err := cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	return err
}

func copyTestResultsToFolder(repoPath string, repoName string, autoGraderPath string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/test-results/test", repoPath, repoName)
	cmd := exec.Command("cp", "-r", copyFromPath, autoGraderPath)
	cmd.Dir = fmt.Sprintf("%v/%v", repoPath, repoName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	//rename
	cmd = exec.Command("mv", autoGraderPath, "tests")
	err = cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	return err
}

func copyDocsToFolder(repoPath string, repoName string, tempPath string, courseID string, semesterID string, assignmentName string, studentUserName string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/docs/javadoc", repoPath, repoName)
	copyToPath := fmt.Sprintf("%v/AutoGrader/%v-%v/%v/%v", tempPath, courseID, semesterID, assignmentName, studentUserName)
	cmd := exec.Command("cp", "-r", copyFromPath, copyToPath)
	cmd.Dir = repoPath
	err := cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	//rename
	cmd = exec.Command("mv", copyToPath, "docs")
	err = cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	return err
}

func createDocs(repoPath string, repoName string) error {
	cmd := exec.Command("gradle", "javadoc")
	cmd.Dir = fmt.Sprintf("%v/%v", repoPath, repoName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func runUnitTests(repoPath string, repoName string) error {
	cmd := exec.Command("gradle", "test")
	cmd.Dir = fmt.Sprintf("%v/%v", repoPath, repoName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func handleNonCodeSubmissions(repoPath string, repoName string, autoGraderPath string) error {
	copyFromPath := fmt.Sprintf("%v/%v/build/test-results/test", repoPath, repoName)
	cmd := exec.Command("cp", "r", copyFromPath, autoGraderPath)
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func main() {
	sc := make(chan os.Signal, 1)
	//catch sigint, sigterm & os.int
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	done := make(chan bool, 1)
	//sigStruct := signals{signalThrown, mu}
	//waitgroup languageswitch
	//go languageSwitch(sigStruct)
	time.Sleep(1 * time.Second)
	//go sigs(sigStruct)
	<-done
	fmt.Println("Exiting")

	// done: start with environment variables
	// done: switch statement
	// TODO: one language at a time
	// TODO: java>c++>duplicate
	// TODO: sigterm handling on grader level
}
