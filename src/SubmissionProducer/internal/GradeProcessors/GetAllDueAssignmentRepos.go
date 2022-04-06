package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"fmt"
	"github.com/google/go-github/v41/github"
	"regexp"
	"strconv"
)

func GetAllDueAssignmentRepos(assignment common.DueAssignment) []*github.Repository {

	var DueRepos []*github.Repository

	regReplace, err := regexp.Compile("[^0-9]")
	common.CheckIfError(err)

	AssignmentNumber := regReplace.ReplaceAll([]byte(assignment.AssignmentName), []byte(""))

	var re1 *regexp.Regexp

	var ANTemp, _ = strconv.Atoi(string(AssignmentNumber))

	if ANTemp > 9 {
		re1 = regexp.MustCompile(fmt.Sprintf("(?i)(.*)%v(-[a-zA-Z]\\d+)?-Assignment(%c|-%c)(.*)", assignment.CourseIDSemID, AssignmentNumber, AssignmentNumber))
	} else {
		re1 = regexp.MustCompile(fmt.Sprintf("(?i)(.*)%v(-[a-zA-Z]\\d+)?-Assignment(%c|-%c|\\d%c|-\\d%c)(.*)", assignment.CourseIDSemID, AssignmentNumber, AssignmentNumber, AssignmentNumber, AssignmentNumber))
	}

	common.Info("\t Checking repositories for the regex " + re1.String())

	for _, repo := range common.AllRepos {

		if re1.MatchString(*repo.Name) {
			DueRepos = append(DueRepos, repo)
		}
	}

	return DueRepos

}
