package internal

import (
	"SubmissionProducer/internal/Docker"
	"SubmissionProducer/internal/GradeProcessors"
	"SubmissionProducer/internal/common"
	"fmt"
	"github.com/google/go-github/v41/github"
	"runtime"
	"sync"
)

func GradeAssignments(wg *sync.WaitGroup, repo *github.Repository, assignment common.DueAssignment, ch chan<- common.UserAssignmentData) {
	defer wg.Done()
	var UserAssignmentData common.UserAssignmentData

	UserAssignmentData.UserName = common.ParseUsername(*repo.HTMLURL, assignment)
	UserAssignmentData.FullName = common.GetGithubUserName(UserAssignmentData.UserName)
	UserAssignmentData.GitURL = *repo.HTMLURL
	UserAssignmentData.CommitTimes = common.GetCommitDates(*repo.HTMLURL)

	Docker.ProcessAssignment(repo, assignment)

	var OutputPath string
	if runtime.GOOS == "windows" {
		OutputPath = fmt.Sprintf("%s\\AutoGrader\\%s\\%s\\%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName, UserAssignmentData.UserName)
	} else {
		OutputPath = fmt.Sprintf("%s/AutoGrader/%s/%s/%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName, UserAssignmentData.UserName)
	}

	if assignment.AssignmentConfig.StudentTestsEnabled || assignment.AssignmentConfig.TeacherUnitTests {
		UserAssignmentData.FailedTests = GradeProcessors.GetFailedUnitTests(OutputPath, assignment.AssignmentConfig)
	}
	if assignment.AssignmentConfig.GradeDocs {
		// TODO LATER: Get Docs
		UserAssignmentData.MissingDocs = common.TestDocs{}
	}

	ch <- UserAssignmentData

}
