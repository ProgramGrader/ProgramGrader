package internal

import (
	"SubmissionProducer/internal/GradeProcessors"
	"SubmissionProducer/internal/common"
	"SubmissionProducer/internal/email"
	"fmt"
	"github.com/google/go-github/v41/github"
	"github.com/jinzhu/copier"
	"github.com/shirou/gopsutil/cpu"
	"runtime"
	"strconv"
	"sync"
	"time"
)

func StartGrading(assignments []common.DueAssignment) {

	common.GetAndCacheCurrentRepositories()

	// Get a list of Repo Urls for each assignment
	// TODO parallelize with go routines each assignment to grade
	// goroutiens can handle up to 1 million concurrent goroutines on a 4gb machine
	for _, assignment := range assignments {

		var assignmentData common.AssignmentSubmissions

		assignmentData.DueDate = assignment.AssignmentDueDate
		assignmentData.Teacher = assignment.CourseTeacher
		assignmentData.AssignmentName = assignment.AssignmentName
		assignmentData.CourseNumber = assignment.CourseIDSemID
		assignmentData.CourseName = assignment.CourseType
		assignmentData.GradeDocs = assignment.AssignmentConfig.GradeDocs
		assignmentData.GradeTests = assignment.AssignmentConfig.StudentTestsEnabled || assignment.AssignmentConfig.TeacherUnitTests

		common.Info("Processing Repositories for: %v-%v-*", assignment.CourseIDSemID, assignment.AssignmentName)
		GetAssignmentUrls := GradeProcessors.GetAllDueAssignmentRepos(assignment)

		if GetAssignmentUrls == nil {
			continue
		}
		// implicit else
		common.Info("\t " + strconv.Itoa(len(GetAssignmentUrls)) + " repos found to grade.")

		var AssignmentPath string
		if runtime.GOOS == "windows" {
			AssignmentPath = fmt.Sprintf("%s\\AutoGrader\\%s\\%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName)
		} else {
			AssignmentPath = fmt.Sprintf("%s/AutoGrader/%s/%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName)
		}

		GradeProcessors.ClearExistingSubmissionFiles(AssignmentPath)

		// Parallel
		// Have to use a wait group to ensure all treads are done
		wg := sync.WaitGroup{}

		//GetAssignmentUrls = GetAssignmentUrls[:5]
		// Channel to put the returns into
		var counter = 0
		ch := make(chan common.UserAssignmentData, len(GetAssignmentUrls))
		for _, repo := range GetAssignmentUrls {
			wg.Add(1)

			temp := github.Repository{}
			err := copier.Copy(&temp, &repo)
			common.CheckIfError(err)

			// FIXME: This should only happen if the env is local
			// Keep time low to prevent I/O block from Docker Desktop
			// If over 80% of the cpu is used then we wait till it is using less then 80%
			var iterationCounter = 0
			var average = 0.0
			for {
				percentage, _ := cpu.Percent(1*time.Second, false)
				common.CheckIfError(err)

				iterationCounter = iterationCounter + 1
				average = average + percentage[0]
				//fmt.Println(average / float64(iterationCounter))

				if iterationCounter > 3 {
					if average/float64(iterationCounter) < 75 {
						//fmt.Println("break")
						iterationCounter = 0
						average = 0
						break
					}
				} else if iterationCounter > 15 {
					iterationCounter = 0
					average = 0
				}
			}

			// TODO: if on aws call fargate to create a new
			// TODO: Need to deal with subGrader being killed (TBD)
			go GradeAssignments(&wg, &temp, assignment, ch)

			if counter > 5 {
				time.Sleep(2 * time.Second)
				counter = 0
			}
			counter = counter + 1
		}

		wg.Wait()

		assignmentData.Data = make([]common.UserAssignmentData, len(GetAssignmentUrls))

		for i := range assignmentData.Data {
			temp := <-ch
			assignmentData.Data[i] = temp
		}
		close(ch)
		common.Info("\t Grading Completed!")
		// Parallel

		// Serial
		// Grade
		//for _, repo := range GetAssignmentUrls {
		//	var UserAssignmentData common.UserAssignmentData
		//
		//	UserAssignmentData.UserName = common.ParseUsername(*repo.HTMLURL, assignment)
		//	UserAssignmentData.FullName = common.GetGithubUserName(UserAssignmentData.UserName)
		//	UserAssignmentData.GitURL = *repo.URL
		//	UserAssignmentData.CommitTimes = common.GetCommitDates(*repo.HTMLURL)
		//
		//	// Remove for co-routine version
		//	Docker.ProcessAssignment(repo, assignment)
		//
		//	var OutputPath string
		//	if runtime.GOOS == "windows" {
		//		OutputPath = fmt.Sprintf("%s\\AutoGrader\\%s\\%s\\%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName, *repo.Name)
		//	} else {
		//		OutputPath = fmt.Sprintf("%s/AutoGrader/%s/%s/%s", common.TempPath(), assignment.CourseIDSemID, assignment.AssignmentName, *repo.Name)
		//	}
		//
		//	if assignment.AssignmentConfig.StudentTestsEnabled || assignment.AssignmentConfig.TeacherUnitTests {
		//		UserAssignmentData.FailedTests = GradeProcessors.GetFailedUnitTests(OutputPath, assignment.AssignmentConfig)
		//	}
		//	if assignment.AssignmentConfig.GradeDocs {
		//
		//		UserAssignmentData.MissingDocs = common.TestDocs{}
		//	}
		//
		//	assignmentData.Data = append(assignmentData.Data, UserAssignmentData)
		//
		//}
		// Serial

		if assignment.AssignmentConfig.GradeDocs {
			GradeProcessors.GatherDocs(AssignmentPath)
		}

		if assignment.AssignmentConfig.NonCodeSubmissions {
			GradeProcessors.GatherNonCodeSubmissions(AssignmentPath)
		}

		if assignment.AssignmentConfig.NonCodeSubmissions || assignment.AssignmentConfig.GradeDocs {
			GradeProcessors.ZipSubmissionFolderAndCleanDir(AssignmentPath)

			assignmentData.SubmissionURL = email.UploadFileToS3AndGetPresignedUrl(
				AssignmentPath+common.PathSeparator()+"submission.tgz",
				assignmentData.CourseNumber,
				fmt.Sprintf("%v", assignmentData.AssignmentName),
			)
			common.Info("\t Successfully Uploaded attachment to S3.")
		}

		// Process Graded assignment
		GradeProcessors.ProcessAndEmailGrades(assignmentData, assignment.AssignmentConfig)

	} //Due Assignments

	fmt.Println("All Assignments Graded list")
}
