package email

import (
	"SubmissionProducer/internal/common"
	"strconv"
	"strings"
)

func CreateHtmlFromData(data common.AssignmentSubmissions) string {

	var sb strings.Builder

	sb.WriteString("<div>Assignment: " + data.AssignmentName + "<br>")

	sb.WriteString("Due Date: " + data.DueDate.Format("Jan-02-06 15:04") + "<br><br>")

	if data.SubmissionURL != "" {
		sb.WriteString("Attachment URL: " + data.SubmissionURL + "<br>")
		sb.WriteString("Note: This url may be dead within 6 hours of receiving this email. We are currently working on a fix for this.<br><br>")
	}

	sb.WriteString("Student Submissions: <br>")

	sb.WriteString("Note: If a students name does not show up it was either not turned in or did not compile.<br><br>")

	for _, item := range data.Data {

		var FullName string
		if item.FullName == "" {
			FullName = "No Name Provided"
		} else {
			FullName = item.FullName
		}

		sb.WriteString("&emsp;" + FullName + " (" + item.UserName + ") - " + item.GitURL + "<br><br>")

		if len(item.FailedTests.TestNames) != 0 {
			sb.WriteString("&emsp;&emsp;<p stlye=\"color:red\"><i>" +
				strconv.FormatFloat((1.0-(float64(len(item.FailedTests.TestNames))/float64(item.FailedTests.Total)))*100, 'f', 2, 32) +
				"% (" + strconv.Itoa(item.FailedTests.Total-len(item.FailedTests.TestNames)) + "/" +
				strconv.Itoa(item.FailedTests.Total) + ") Completed</i><br>")

			if !item.FailedTests.Completed {
				sb.WriteString("&emsp;&emsp;Required amount of test where not completed.<br>")
			}

			sb.WriteString("&emsp;&emsp;Failed Tests: <br>")

			for _, tests := range item.FailedTests.TestNames {
				sb.WriteString("&emsp;&emsp;&emsp;- " + tests + "<br>")
			}
			sb.WriteString("<p><br>")
		}

		if len(item.MissingDocs.TestNames) != 0 {

			sb.WriteString("&emsp;&emsp; <p stlye=\"color:blue\"><i>" +
				strconv.FormatFloat((1.0-(float64(len(item.MissingDocs.TestNames))/float64(item.MissingDocs.Total)))*100, 'f', 2, 32) +
				"% (" + strconv.Itoa(item.MissingDocs.Total-len(item.MissingDocs.TestNames)) + "/" +
				strconv.Itoa(item.MissingDocs.Total) + ") Completed</i><br>")

			sb.WriteString("&emsp;&emsp;Missing Docs: <br>")

			for _, docs := range item.MissingDocs.TestNames {
				sb.WriteString("&emsp;&emsp;&emsp;- " + docs + "<br>")
			}
			sb.WriteString("<p><br>")
		}

		if len(item.FailedTests.TestNames) != 0 ||
			len(item.FailedTests.TestNames) != 0 {
			var total = item.MissingDocs.Total + item.FailedTests.Total
			var totalmissing = len(item.MissingDocs.TestNames) + len(item.FailedTests.TestNames)

			sb.WriteString("&emsp;&emsp;<b><i>Total Completed: " +
				strconv.Itoa(total-totalmissing) + "/" + strconv.Itoa(total) + " " +
				strconv.FormatFloat((1.0-float64(totalmissing)/float64(total))*100, 'f', 2, 32) +
				"% Completed</b></i><br><br>")
		}
	}

	sb.WriteString("</div>")
	//fmt.Println(sb.String())
	return sb.String()
}
