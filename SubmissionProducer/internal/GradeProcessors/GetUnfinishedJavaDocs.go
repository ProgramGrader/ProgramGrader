package GradeProcessors

import (
	"SubmissionProducer/internal/common"
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
)

func GetUnfinishedJavaDocs(path string) {

	f, err := ioutil.ReadFile(path)
	common.CheckIfError(err)

	fmt.Println(string(f))

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(f))
	common.CheckIfError(err)

	test := doc.Find("li")

	fmt.Println(test.Nodes)

}
