package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	// "runtime"
)

var (
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
	tempPath = "/tmp/AutoGrader/" + courseID + "-" + semesterID + "/" + assignmentName + "/" + studentUserName
)

func GetEnvVar(envName string) string {
	if envName != "" {
		ret := os.Getenv("envName")
		return ret
	}
	fmt.Printf(envName + "environment variable was not found.")
	return "error"
}

func languageSwitch(ctx context.Context) error {
	var err error
	switch language {
	case "java":
		{
			//clone repo
			err = cloneRepo(githubToken, repoName)
			if err != nil {
				return err
			}
			if teacherUnitTestsEnabled == "TRUE" {
				err = cloneConfigRepo(githubToken)
				if err != nil {
					return err
				}
				//copy tests to folder
				err = copyTestsToFolder(repoName)
				if err != nil {
					return err
				}
			}
			if gradeDocsEnabled == "TRUE" {
				fmt.Println("Grading java")
				err = gradeDocs(repoName)
				if err != nil {
					return err
				}
			}
			if studentTestsEnabled == "TRUE" || teacherUnitTestsEnabled == "TRUE" {
				//err = runUnitTests()
				//if err != nil {
				//	return err
				//}
			}

			if nonCodeSubmissionsEnabled == "TRUE" {
				//err = handleNonCodeSubmissions()
				//if err != nil {
				//	return err
				//}
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
	return err
}

func startContainer() error {
	ctx := context.Background()
	err := languageSwitch(ctx)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	return err
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

func copyTestsToFolder(repoName string) error {
	cmd := exec.Command("cp", "--recursive", fmt.Sprintf("/opt/gradle/%v/src/test/java", repoName))
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func cloneRepo(githubToken string, repoName string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", githubToken, repoName))
	cmd.Dir = "/opt/gradle"
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func cloneConfigRepo(githubToken string) error {
	cmd := exec.Command("git", "clone", fmt.Sprintf("https://%v@github.com/%v.git", githubToken, config))
	cmd.Dir = "/tmp"
	err := cmd.Run()
	if err != nil {
		return err
	}
	return err
}

func gradeDocs(repoName string) error {
	cmd := exec.Command("gradle", "javadoc", fmt.Sprintf("/opt/gradle/%v", repoName))
	cmd.Dir = fmt.Sprintf("/opt/gradle/%v", repoName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	//tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/build/docs/javadoc", repoName))
	//if err != nil {
	//	return err
	//}

	return err
}

//func runUnitTests() error {
//	cmd := exec.Command("gradle", "test", fmt.Sprintf("/opt/gradle/%v", repoName)
//	cmd.Dir = fmt.Sprintf("/opt/gradle/%v", repoName)
//	err := cmd.Run()
//	if err != nil {
//		return err
//	}
//
//	tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/build/test-results/test", *repo.Name))
//	//common.CheckIfError(err)
//
//	if err == nil {
//		common.HandleTarStream(tarStream, CopyToPath)
//
//		err = os.Rename(CopyToPath+common.PathSeparator()+"test", CopyToPath+common.PathSeparator()+"tests")
//		common.CheckIfError(err)
//	}
//	return err
//}

//func handleNonCodeSubmissions() error {
//	tarStream, _, err := cli.CopyFromContainer(ctx, containerName, fmt.Sprintf("/opt/gradle/%v/submission", *repo.Name))
//
//	if err == nil {
//		common.HandleTarStream(tarStream, CopyToPath)
//	}
//}

func main() {
	err := startContainer()
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	// done: start with environment variables
	// done: switch statement
	// TODO: one language at a time
	// TODO: java>c++>duplicate
	// TODO: sigterm handling on grader level
}
