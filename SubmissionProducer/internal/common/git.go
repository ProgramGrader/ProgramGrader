package common

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"time"
)

func GetCommitDates(url string) []string {

	var r *git.Repository
	tries := 0
	for {
		temp, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
			URL: url,
			Auth: &http.BasicAuth{
				Username: fmt.Sprintf("%v", GithubAPIKey["username"]),
				Password: fmt.Sprintf("%v", GithubAPIKey["token"]),
			},
			//Progress: os.Stdout,
		})

		if err == nil {
			r = temp
			break
		} else if err.Error() == "authentication required" {
			//TODO Log error
			CheckIfErrorWithMessage(err, "Token Has Expired. Please Refresh.")
		}

		tries++
		if tries > 3 {
			CheckIfErrorWithMessage(err, "Could not get repository after 3 tries.")
		}

		// Exponential backoff
		<-time.After(time.Duration(1*(2.0^tries)) * time.Second)
	}
	// retrieves the branch pointed by HEAD
	ref, err := r.Head()

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckIfError(err)

	var CommitTimeSlice []string

	err = cIter.ForEach(
		func(c *object.Commit) error {
			commitTime := c.Committer.When.UTC().String()
			CommitTimeSlice = append(CommitTimeSlice, commitTime)
			return nil
		},
	)
	CheckIfError(err)

	CommitTimeSlice = RemoveDuplicates(CommitTimeSlice)

	return CommitTimeSlice

}
