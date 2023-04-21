package main

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetGithubClient(ctx context.Context, accessToken string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func GetLatestRepo(localPath, owner, repoName, token string) (*git.Repository, error) {
	// TODO this can be done in-memory - see https://pkg.go.dev/github.com/go-git/go-git/v5#readme-in-memory-example

	// Clone the repository
	url := fmt.Sprintf("https://github.com/%s/%s.git", owner, repoName)

	return git.PlainClone(localPath, false, &git.CloneOptions{
		URL: url,
		Auth: &http.BasicAuth{
			Username: owner,
			Password: token,
		},
	})
}

func CreatePR(ctx context.Context, client *github.Client, owner, repo, branchName, title, body string) error {
	baseBranch := "main"
	pr, _, err := client.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: &title,
		Head:  &branchName,
		Base:  &baseBranch,
		Body:  &body,
	})
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created pull request: %s\n", pr.GetHTMLURL())
	return nil
}
