package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"os"
)

func ClearExistingSubmissionFiles(path string) {

	// See if previous zip folder exists, if so delete
	if _, err := os.Stat(path + common.PathSeparator() + "submission.tar.gz"); !os.IsNotExist(err) {
		err := os.Remove(path + common.PathSeparator() + "submission.tar.gz")
		common.CheckIfError(err)
	}

	// See if previous zip folder exists, if so delete
	if _, err := os.Stat(path + common.PathSeparator() + "submission.tar"); !os.IsNotExist(err) {
		err := os.Remove(path + common.PathSeparator() + "submission.tar")
		common.CheckIfError(err)
	}

	// See if previous zip folder exists, if so delete
	if _, err := os.Stat(path + common.PathSeparator() + "submission"); !os.IsNotExist(err) {
		err := os.RemoveAll(path + common.PathSeparator() + "submission")
		common.CheckIfError(err)
	}

	// create submissions folder at C:\tmp\AutoGrader\courseID-semesterIdentifier\AssignmentName -> this is passed in as path
	common.MakeDir(path + common.PathSeparator() + "submission")

}
