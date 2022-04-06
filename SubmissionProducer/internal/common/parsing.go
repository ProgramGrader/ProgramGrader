package common

import (
	"context"
	"errors"
	"fmt"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
	"strconv"
	"strings"
)

func SliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func RemoveDuplicates(arr []string) []string {
	var visited []string
	for _, value := range arr {
		if !SliceContains(visited, value) {
			visited = append(visited, value)
		}
	}
	return visited
}

func ParseUsername(url string, assignment DueAssignment) string {

	substrLocation1 := strings.LastIndex(strings.ToLower(url), strings.ToLower(assignment.AssignmentName))

	substrLocation2 := strings.LastIndex(strings.ToLower(url), strings.ToLower("assignment-"+strconv.Itoa(assignment.AssignmentNumber)))

	var username string
	if substrLocation1 != -1 {
		username = url[substrLocation1+len(assignment.AssignmentName)+1:]
	} else if substrLocation2 != -1 {
		username = url[substrLocation2+len("assignment-"+strconv.Itoa(assignment.AssignmentNumber))+1:]
	} else {
		CheckIfError(errors.New("could not parse username from " + url))
	}

	if strings.Contains(username, "assignment") {
		// TODO Log error
		CheckIfError(errors.New("could not parse username"))
	}

	return username
}

func GetGithubUserName(username string) string {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: fmt.Sprintf("%v", GithubAPIKey["token"])},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	var query struct {
		User struct {
			Name githubv4.String
		} `graphql:"user(login:$username)"`
	}

	variables := map[string]interface{}{
		"username": githubv4.String(username),
	}

	err := client.Query(context.Background(), &query, variables)

	CheckIfErrorWithMessage(err, fmt.Sprintf("Could not find username for username: %v ", username))

	return string(query.User.Name)

}

func Average(xs []float64) float64 {
	total := 0.0
	for _, v := range xs {
		total += v
	}
	return total / float64(len(xs))
}
