package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

var (
	mu sync.Mutex

	//secret token
	githubToken = GetEnvVar("GITHUBTOKEN")

	//env variables
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

	config        = "IUS-CS/AutoGraderConfig"
	containerName = repoName + "-Container"
	// copyToPath = "/tmp/AutoGrader/$CourseID-$SemesterID/$AssignmentName/$StudentUserName"
	tempPath = "/tmp/AutoGrader/" + courseID + "-" + semesterID + "-" + assignmentName + "-" + studentUserName + "/"
)

type signals struct {
	signalThrown bool
	mu           sync.Mutex
}

func GetEnvVar(envName string) string {
	if envName != "" {
		ret := os.Getenv(envName)
		return ret
	}
	fmt.Printf(envName + "environment variable was not found.")
	return "error"
}

func languageSwitch(sigs signals) {
	var err error
	switch language {
	case "java":
		{
			//clone repo
			err = cloneRepo()
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			if teacherUnitTestsEnabled == "TRUE" {
				err = cloneConfigRepo()
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
				//copy tests to folder
				err = copyTestsToFolder()
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}
			if gradeDocsEnabled == "TRUE" {
				fmt.Println("Grading java")
				err = createDocs()
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}
			if studentTestsEnabled == "TRUE" || teacherUnitTestsEnabled == "TRUE" {
				err = runUnitTests()
				if err != nil {
					fmt.Printf("Error: %s", err)
					return
				}
			}

			if nonCodeSubmissionsEnabled == "TRUE" {
				err = handleNonCodeSubmissions()
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

func copyToPath() string {
	var ret string
	if runtime.GOOS == "windows" {
		ret = fmt.Sprintf("%s\\AutoGrader\\%s\\%s\\%s", tempPath, courseID, assignmentName, studentUserName)
	} else {
		ret = fmt.Sprintf("%s/AutoGrader/%s/%s/%s", tempPath, courseID, assignmentName, studentUserName)
	}
	return ret
}

func copyTestsToFolder() error {
	cmd := exec.Command("cp", "--recursive")
	cmd.Dir = fmt.Sprintf("/opt/gradle/%v/src/test/java", repoName)
	err := cmd.Run()
	fmt.Printf("Error: %s", err)
	if err != nil {
		return err
	}
	return err
}

func cloneRepo() error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://github.com/%v", repoName))
	cmd.Dir = "/tmp/AutoGrader"
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func cloneConfigRepo() error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", githubToken, config))
	cmd.Dir = ""
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func createDocs() error {
	cmd := exec.Command("gradle", "javadoc")
	cmd.Dir = fmt.Sprintf("/opt/gradle/%v", repoName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	//tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/build/docs/javadoc", repoName))
	//if err != nil {
	//	return err
	//}
	return err
}

func runUnitTests() error {
	cmd := exec.Command("gradle", "test")
	cmd.Dir = fmt.Sprintf("/opt/gradle/%v", repoName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

func handleNonCodeSubmissions() error {
	//tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/submission", *repo.Name))
	cmd := exec.Command("gradle", "")
	cmd.Dir = fmt.Sprintf("/opt/gradle/%v/submission", repoName)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return err
}

//func sigs(sigs signals) {
//	sig := sigs.sc
//	fmt.Println("Signal received: ", sig)
//	sigs.done <- true
//}

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
