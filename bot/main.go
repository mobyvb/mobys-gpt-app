package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func main() {
	owner := "mobyvb"
	repoName := "mobys-gpt-app"

	// Set up the GitHub API client
	token := os.Getenv("GITHUB_GPT_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_GPT_ACCESS_TOKEN environment variable not set")
	}

	ctx := context.Background()

	//processIssues(ctx, owner, repoName, token)

	// TODO use in-memory repo so that this removal is unnecessary
	err := os.RemoveAll("/tmp/mytmprepo")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := GetLatestRepo("/tmp/mytmprepo", owner, repoName, token)
	if err != nil {
		log.Fatal(err)
	}

	err = processResponse(ctx, "/tmp/mytmprepo", repo, owner, repoName, token)
	if err != nil {
		log.Fatal(err)
	}
}

func processIssues(ctx context.Context, owner, repoName, token string) {
	client := GetGithubClient(ctx, token)

	// List and parse GitHub issues
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repoName, nil)
	if err != nil {
		log.Fatal(err)
	}

	// TODO this code should go somewhere else
	for _, issue := range issues {
		// TOOD allowlist of users who can do stuff
		if *issue.User.Login == owner {
			issueBody := *issue.Body

			// TODO hadnle all the different cases (e.g. "Files: " not provided)
			parts := strings.Split(issueBody, "Files:")
			body := parts[0]
			fileList := strings.Split(parts[1], ",")

			issueInfo := IssueInfo{
				Subject: *issue.Title,
				Body:    body,
			}

			for _, p := range fileList {
				p = strings.TrimSpace(p)

				// TODO "../" here is hardcoded but needs to go at some point
				data, err := ioutil.ReadFile("../" + p)
				if err != nil {
					panic(err)
				}

				issueInfo.Files = append(issueInfo.Files, File{
					Path:     p,
					Contents: string(data),
				})
			}

			fmt.Println(issueInfo)
		}
	}
}

func processResponse(ctx context.Context, localRepoPath string, repo *git.Repository, owner, repoName, token string) error {
	data, err := ioutil.ReadFile("response.txt")
	if err != nil {
		return err
	}
	issueResponse := ParseIssueResponse(string(data))

	// Create a worktree for a commit later
	worktree, err := repo.Worktree()
	if err != nil {
		return err
	}

	for _, f := range issueResponse.Files {

		fullPath := filepath.Join(localRepoPath, f.Path)

		// TODO standard indentation (for non-go files)
		if strings.HasSuffix(f.Path, ".go") {
			newContents, err := format.Source([]byte(f.Contents))
			if err != nil {
				return err
			}
			f.Contents = string(newContents)
		}

		err = ioutil.WriteFile(fullPath, []byte(f.Contents), 0644)
		if err != nil {
			return err
		}

		_, err = worktree.Add(f.Path)
		if err != nil {
			return err
		}
	}

	commitMessage := "TODO subject from issue\n\n" + issueResponse.Notes + "\n\nResolves #issuenumber"
	_, err = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "moby-robot-fakename",           // Replace with your name
			Email: "mobyrobot-fakeemail@gmail.com", // Replace with your email
		},
	})
	if err != nil {
		return err
	}

	// create a new branch
	//branchName := "moby-robot-branch"
	branchName := randomBranchName()
	branchRefName := plumbing.NewBranchReferenceName(branchName)

	headRef, err := repo.Head()
	if err != nil {
		return err
	}

	err = repo.CreateBranch(&config.Branch{
		Name:   branchName,
		Remote: "origin",
		Merge:  branchRefName,
	})
	if err != nil {
		return err
	}

	// Update the branch to point to the new commit
	err = repo.Storer.SetReference(plumbing.NewHashReference(branchRefName, headRef.Hash()))
	if err != nil {
		return err
	}

	// Push the new branch to the remote repository
	remote, err := repo.Remote("origin")
	if err != nil {
		return err
	}

	err = remote.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{config.RefSpec(fmt.Sprintf("%s:refs/heads/%s", branchRefName, branchName))},
		Auth: &http.BasicAuth{
			Username: owner, // Your GitHub username
			Password: token,
		},
	})
	if err != nil {
		return err
	}

	// finally, create a pull request from the newly pushed branch
	client := GetGithubClient(ctx, token)
	err = CreatePR(ctx, client, owner, repoName, branchName, "TODO replace me title", commitMessage)
	if err != nil {
		return err
	}

	return nil
}

func randomBranchName() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
