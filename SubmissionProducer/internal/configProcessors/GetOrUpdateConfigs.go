package configProcessors

import (
	"SubmissionProducer/internal/common"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"os"
)

func GetOrUpdateConfigs(path string) {

	_, err := os.Stat(path)

	if err != nil {

		fmt.Println("Cloning Autograder Config Repository.")

		_, err = git.PlainClone(path, false, &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: fmt.Sprintf("%v", common.GithubAPIKey["username"]),
				Password: fmt.Sprintf("%v", common.GithubAPIKey["token"]),
			},
			URL:      common.ConfigUrl,
			Progress: os.Stdout,
		})
		common.CheckIfErrorWithMessage(err, "Failed to clone repository.")

		common.Info("Successfully Cloned Config Repository")

	} else {
		r, err := git.PlainOpen(path)
		common.CheckIfError(err)

		// Get the working directory for the repository
		w, err := r.Worktree()
		common.CheckIfError(err)

		// Pull the latest changes from the origin remote and merge into the current branch
		common.Info("git pull origin")
		err = w.Pull(&git.PullOptions{
			RemoteName:   "origin",
			SingleBranch: true,
			//ReferenceName: "main",
			Auth: &http.BasicAuth{
				Username: fmt.Sprintf("%v", common.GithubAPIKey["username"]),
				Password: fmt.Sprintf("%v", common.GithubAPIKey["token"]),
			},
			//Progress: os.Stdout,
		})
		if err == nil {
			// Print the latest commit that was just pulled
			ref, err := r.Head()
			common.CheckIfError(err)

			commit, err := r.CommitObject(ref.Hash())
			common.CheckIfError(err)

			common.Info("Pulled current HEAD:%v by %v\n", commit.Hash, commit.Author)
		} else if err.Error() == "already up-to-date" {
			common.Info("Repository Already up-to-date.")
			return
		} else {
			common.CheckIfError(err)
		}

	}
}
