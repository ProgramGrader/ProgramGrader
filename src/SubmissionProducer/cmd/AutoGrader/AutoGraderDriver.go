package main

import (
	"SubmissionProducer/internal"
	"SubmissionProducer/internal/GradeProcessors"
	"SubmissionProducer/internal/common"
)

// https://docs.docker.com/engine/api/sdk/examples/#run-a-container
func main() {

	common.GithubAPIKey = common.GetAWSSecret("Autograder/GithubToken")

	// Find any assignment that was due in the past and was not graded, return assignments to grade assignmentsAndDueDates :=
	DueAssignments := GradeProcessors.GetCurrentlyDueAssignments()

	// Update image before grading
	common.UpdateGradingImage()

	if DueAssignments == nil {
		common.Info("No work to do!")
	} else {
		// Get Assignment Configs and spawn grading processes
		internal.StartGrading(DueAssignments)
	}
}
