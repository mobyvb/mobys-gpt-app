package main

import (
	"context"

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
