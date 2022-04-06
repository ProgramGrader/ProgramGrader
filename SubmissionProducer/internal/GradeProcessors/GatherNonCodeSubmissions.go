package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	cp "github.com/otiai10/copy"
	"io/ioutil"
	"os"
	"strings"
)

func GatherNonCodeSubmissions(path string) {

	files, err := ioutil.ReadDir(path)
	common.CheckIfError(err)

	// Copy each sub-folder submission folder with the name of the user
	for _, f := range files {
		if !strings.Contains(strings.ToLower(f.Name()), "submission") && f.IsDir() {

			username := f.Name()

			if _, err := os.Stat(path + common.PathSeparator() + username + common.PathSeparator() + "submission"); !os.IsNotExist(err) {

				if _, err := os.Stat(path + common.PathSeparator() + "submission" + common.PathSeparator() + username); !os.IsNotExist(err) {
					common.MakeDir(path + common.PathSeparator() + "submission" + common.PathSeparator() + username)
				}

				err := cp.Copy(path+common.PathSeparator()+username+common.PathSeparator()+"submission",
					path+common.PathSeparator()+"submission"+common.PathSeparator()+username+common.PathSeparator()+"submission",
				)
				common.CheckIfError(err)
			}
		}
	}

}
