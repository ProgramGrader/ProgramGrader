package common

import (
	"context"
	"fmt"
	"github.com/google/go-github/v41/github"
	"golang.org/x/oauth2"
	"time"
)

func GetCurrentSemesterIdentifier() string {

	ctime := time.Now()

	springTime := time.Date(ctime.Year(), 5, 10, 11, 59, 00, 0, ctime.Location())
	summerTime := time.Date(ctime.Year(), 8, 10, 11, 59, 00, 0, ctime.Location())

	var semester string
	if ctime.Before(springTime) {
		semester = "SP"
	} else if ctime.Before(summerTime) {
		semester = "SU"
	} else {
		semester = "FA"
	}

	semester = semester + ctime.Format("06")

	return semester
}

func GetSemesterStartDate() time.Time {

	ctime := time.Now()
	cmonth := ctime.Month()
	cday := ctime.Day()

	var semester time.Time
	if cmonth < 5 && cday < 10 {
		semester = time.Date(ctime.Year(), 1, 1, 0, 0, 0, 0, ctime.Location())
	} else if cmonth < 8 && cday < 10 {
		semester = time.Date(ctime.Year(), 5, 10, 0, 0, 0, 0, ctime.Location())
	} else {
		semester = time.Date(ctime.Year(), 8, 10, 0, 0, 0, 0, ctime.Location())
	}

	return semester
}

func GetAndCacheCurrentRepositories() {
	Info("Caching all relevant repositories to org")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: fmt.Sprintf("%v", GithubAPIKey["token"])},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	listOption := github.ListOptions{
		Page:    0,
		PerPage: 100,
	}
	// list private repositories for org
	opt := &github.RepositoryListByOrgOptions{Type: "private",
		Sort:        "created",
		ListOptions: listOption,
	}

	currentSemesterStartDate := GetSemesterStartDate()

	// get all pages of results
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, GithubOrganization, opt)
		CheckIfError(err)
		AllRepos = append(AllRepos, repos...)

		if resp.NextPage == resp.LastPage ||
			currentSemesterStartDate.Before(AllRepos[len(AllRepos)-1].CreatedAt.Time) {
			break
		}

		opt.Page = resp.NextPage
	}
	//TODO: Add analytics call to see how many repositories are added per call
	Info("Number of repositories cached is: %v", len(AllRepos))
}
