package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"os"
)

func ZipSubmissionFolderAndCleanDir(sourcePath string) {

	err := common.Zipit(sourcePath+common.PathSeparator()+"submission", sourcePath+common.PathSeparator()+"submission.tgz")
	common.CheckIfError(err)

	if common.GetEnv("env", "") != "" {
		err = os.RemoveAll(sourcePath + common.PathSeparator() + "submission")
		common.CheckIfError(err)
		//err = os.Remove(sourcePath + common.PathSeparator() + "submission.zip")
		//common.CheckIfError(err)
	}
}
