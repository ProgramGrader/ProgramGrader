package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	cp "github.com/otiai10/copy"
	"io/ioutil"
	"os"
	"strings"
)

func GatherDocs(path string) {

	files, err := ioutil.ReadDir(path)
	common.CheckIfError(err)

	// Remove the submission folder from the list of student folders
	files = files[:len(files)-1]

	// Copy each sub folder submission folder with the name of the user
	for _, f := range files {
		if !strings.Contains(f.Name(), "docs") && f.IsDir() {

			username := f.Name()

			if _, err := os.Stat(path + common.PathSeparator() + username + common.PathSeparator() + "docs"); !os.IsNotExist(err) {

				if _, err := os.Stat(path + common.PathSeparator() + "submission" + common.PathSeparator() + username); !os.IsNotExist(err) {
					common.MakeDir(path + common.PathSeparator() + "submission" + common.PathSeparator() + username)
				}
				common.CheckIfError(err)

				err := cp.Copy(path+common.PathSeparator()+username+common.PathSeparator()+"docs",
					path+common.PathSeparator()+"submission"+common.PathSeparator()+username+common.PathSeparator()+"docs",
				)
				common.CheckIfError(err)
			}
		}
	}

	err = common.DownloadFile(common.DocOpenerURL,
		path+common.PathSeparator()+"submission"+common.PathSeparator()+"DocsOpener.jar",
	)
	common.CheckIfError(err)

}
