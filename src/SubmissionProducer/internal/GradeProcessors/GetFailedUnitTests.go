package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"github.com/joshdk/go-junit"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetFailedUnitTests(testDir string, config common.AssignmentConfig) common.TestDocs {

	libRegEx, err := regexp.Compile("^.+\\.(xml)$")
	common.CheckIfError(err)

	var FilesToProcess []string

	testDirToSearch := testDir + common.PathSeparator() + "tests"

	err = filepath.Walk(testDirToSearch, func(path string, info os.FileInfo, err error) error {
		if err == nil && libRegEx.MatchString(info.Name()) {
			FilesToProcess = append(FilesToProcess, path)
		}
		return nil
	})
	common.CheckIfError(err)

	var tests common.TestDocs

	tests.Completed = true
	tests.AllPass = true

	studentTestNumberMap := make(map[string]int)

	for _, file := range FilesToProcess {
		suites, err := junit.IngestFile(file)
		common.CheckIfErrorWithMessage(err, "failed to ingest JUnit xml")

		for _, suite := range suites {

			for _, test := range suite.Tests {

				if !strings.Contains(strings.ToLower(test.Name), "teacher") {

					suiteName := strings.Split(test.Name, ".")[0]

					studentTestNumberMap[suiteName] = studentTestNumberMap[suiteName] + 1
				}
				if test.Error != nil || test.Status == "failed" {
					tests.TestNames = append(tests.TestNames, test.Name)
					tests.AllPass = false
				}
				tests.Total = tests.Total + 1

			} //tests
		} //suites
	} // files

	for _, v := range studentTestNumberMap {
		if v < config.NumberStudentTestsRequired {
			tests.Completed = false
			break
		}
	}

	return tests

}
