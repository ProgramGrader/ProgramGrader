package email

import (
	"SubmissionProducer/internal/common"
	"strconv"
	"strings"
)

func CreateTextFromData(data common.AssignmentSubmissions) string {

	var sb strings.Builder

	sb.WriteString("Assignment: " + data.AssignmentName + "\n")

	sb.WriteString("Due Date: " + data.DueDate.Format("Jan-02-06 15:04") + "\n\n")

	if data.SubmissionURL != "" {
		sb.WriteString("Attachment URL: " + data.SubmissionURL + "\n")
		sb.WriteString("Note: This url may be dead within 6 hours of receiving this email. We are currently working on a fix for this.\n\n")
	}

	sb.WriteString("Student Submissions: \n")

	sb.WriteString("Note: If a students name does not show up it was either not turned in or did not compile.\n\n")

	for _, item := range data.Data {

		var FullName string
		if item.FullName == "" {
			FullName = "No Name Provided"
		} else {
			FullName = item.FullName
		}

		sb.WriteString("\t" + FullName + " (" + item.UserName + ") - " + item.GitURL + "\n\n")

		if data.GradeTests {
			if len(item.FailedTests.TestNames) != 0 {
				sb.WriteString("\t\t" + strconv.FormatFloat((1.0-(float64(len(item.FailedTests.TestNames))/float64(item.FailedTests.Total)))*100, 'f', 2, 32) +
					"% (" + strconv.Itoa(item.FailedTests.Total-len(item.FailedTests.TestNames)) + "/" +
					strconv.Itoa(item.FailedTests.Total) + ") Completed\n")

				if !item.FailedTests.Completed {
					sb.WriteString("\t\tRequired amount of test where not completed.\n")
				}

				sb.WriteString("\t\tFailed Tests: \n")

				for _, tests := range item.FailedTests.TestNames {
					sb.WriteString("\t\t\t- " + tests + "\n")
				}
				sb.WriteString("\n")
			} else if item.FailedTests.AllPass {
				sb.WriteString("\t\t\t No Failed Tests. 100%\n")
			} else {
				sb.WriteString("\t\t\t No Tests Found\n")
			}
		}

		if data.GradeDocs {

			if len(item.MissingDocs.TestNames) != 0 {

				sb.WriteString("\t\t" + strconv.FormatFloat((1.0-(float64(len(item.MissingDocs.TestNames))/float64(item.MissingDocs.Total)))*100, 'f', 2, 32) +
					"% (" + strconv.Itoa(item.MissingDocs.Total-len(item.MissingDocs.TestNames)) + "/" +
					strconv.Itoa(item.MissingDocs.Total) + ") Completed\n")

				sb.WriteString("\t\tMissing Docs: \n")

				for _, docs := range item.MissingDocs.TestNames {
					sb.WriteString("\t\t\t- " + docs + "\n")
				}
				sb.WriteString("\n")
			} else if item.FailedTests.AllPass {
				sb.WriteString("\t\t\t No Failed Docs. 100%\n")
			} else {
				// TODO Change when docs parsing is implemented
				//sb.WriteString("No Docs Found")
			}
		}

		if data.GradeDocs || data.GradeTests {

			var total = item.MissingDocs.Total + item.FailedTests.Total
			var totalmissing = len(item.MissingDocs.TestNames) + len(item.FailedTests.TestNames)

			sb.WriteString("\t\tTotal Completed: " +
				strconv.Itoa(total-totalmissing) + "/" + strconv.Itoa(total) + " " +
				strconv.FormatFloat((1.0-float64(totalmissing)/float64(total))*100, 'f', 2, 32) +
				" Completed" + "\n\n")
		}
	}

	return sb.String()
}
