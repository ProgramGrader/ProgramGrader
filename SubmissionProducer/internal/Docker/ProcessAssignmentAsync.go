package Docker

import (
	"SubmissionProducer/internal/common"
	"context"
	"github.com/google/go-github/v41/github"
	"runtime/pprof"
	"strconv"
	"sync"
)

func ProcessAssignmentAsync(wg *sync.WaitGroup, id int, repo *github.Repository, assignment common.DueAssignment) {
	defer wg.Done()
	labels := pprof.Labels("Process", *repo.Name, "id", strconv.Itoa(int(id)))
	pprof.Do(context.Background(), labels, func(_ context.Context) {
		ProcessAssignment(repo, assignment)
	})

}
