package internal

import (
	"fmt"
	"os"
	"sync"
)

func GetEnvVar(envName string) string {
	if envName != "" {
		ret := os.Getenv(envName)
		return ret
	}
	fmt.Printf(envName + "environment variable was not found.")
	return "error"
}

func GradeAssignment(signalVar *GoSafeVar[bool], wg *sync.WaitGroup) {
	defer wg.Done()

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
		orgName                   = "IUS-CS"
		tempPath                  = "/tmp"
		testsPath                 = fmt.Sprintf("%v/%v/src/%v/%v/current/tests", tempPath, config, courseType, assignmentName)
		autoGraderPath            = fmt.Sprintf("%v/AutoGrader/%v-%v/%v/%v", tempPath, courseID, semesterID, assignmentName, studentUserName)
	)

	var repoPath string

	switch language {
	case "java":
		{
			repoPath = "/opt/gradle"
		}
	case "c++":
		{
			repoPath = "/opt/gradle"
		}

	case "python":
		{
			repoPath = "/opt/gradle"
		}
	default:
		{
			//log and deal with unsupported language
			fmt.Println(language + " is not supported. Nothing has been graded.")
		}
	}

	if GetValue(signalVar) {
		// TODO: put submission back in the queue
		return
	}

	//clone repo
	err = cloneRepo(githubToken, orgName, repoName, repoPath)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}

	if GetValue(signalVar) {
		// TODO: put submission back in the queue
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

	if GetValue(signalVar) {
		// TODO: put submission back in the queue
		return
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

	if GetValue(signalVar) {
		// TODO: put submission back in the queue
		return
	}

	if nonCodeSubmissionsEnabled == "TRUE" {
		err = handleNonCodeSubmissions(repoPath, repoName, autoGraderPath)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return
		}
	}

	return
}
