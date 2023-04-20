package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
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

	client := GetGithubClient(ctx, token)

	// List and parse GitHub issues
	issues, _, err := client.Issues.ListByRepo(ctx, owner, repoName, nil)
	if err != nil {
		log.Fatal(err)
	}

	// TODO this code should go somewhere else
	for _, issue := range issues {
		/*
			fmt.Printf("#%d: %s\n", *issue.Number, *issue.Title)
			fmt.Printf("%s\n", *issue.Body)
			fmt.Println(*issue.User.Login)
		*/

		// TOOD allowlist of users who can do stuff
		if *issue.User.Login == owner {
			issueBody := *issue.Body

			// TODO hadnle all the different cases (e.g. "Files: " not provided)
			parts := strings.Split(issueBody, "Files:")
			body := parts[0]
			fileList := strings.Split(parts[1], ",")

			fmt.Println("we can do this one")
			issueInfo := IssueInfo{
				Subject: *issue.Title,
				Body:    body,
			}

			for _, p := range fileList {
				p = strings.TrimSpace(p)
				fmt.Println("file list item")
				fmt.Println(strings.TrimSpace(p))

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

	// Replace the contents of app/main.go
	/*
		filePath := "../app/main.go"
		newContents := []byte("test contents REPLACEME")

		err = ioutil.WriteFile(filePath, newContents, 0644)
		if err != nil {
			log.Fatal(err)
		}

		// Clone the repository
		url := fmt.Sprintf("https://github.com/%s/%s.git", owner, repoName)
		repo, err := git.PlainClone("./tmp/repo", false, &git.CloneOptions{
			URL: url,
			Auth: &http.BasicAuth{
				Username: "username", // Your GitHub username
				Password: token,
			},
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create a new git commit
		worktree, err := repo.Worktree()
		if err != nil {
			log.Fatal(err)
		}

		_, err = worktree.Add(filePath)
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*
		commit, err := worktree.Commit("test body REPLACEME", &git.CommitOptions{
			Author: &object.Signature{
				Name:  "moby-robot",          // Replace with your name
				Email: "mobyrobot@gmail.com", // Replace with your email
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	*/

	// Create a new branch
	/*
		branchName := "moby-robot-branch"
		branchRefName := plumbing.NewBranchReferenceName(branchName)

		headRef, err := repo.Head()
		if err != nil {
			log.Fatal(err)
		}

		err = repo.CreateBranch(&config.Branch{
			Name:   branchName,
			Remote: "origin",
			Merge:  branchRefName,
		})
		if err != nil {
			log.Fatal(err)
		}

		// Update the branch to point to the new commit
		err = repo.Storer.SetReference(plumbing.NewHashReference(branchRefName, headRef.Hash()))
		if err != nil {
			log.Fatal(err)
		}

		// Push the new branch to the remote repository
		remote, err := repo.Remote("origin")
		if err != nil {
			log.Fatal(err)
		}

		err = remote.Push(&git.PushOptions{
			RefSpecs: []config.RefSpec{config.RefSpec(fmt.Sprintf("%s:refs/heads/%s", branchRefName, branchName))},
			Auth: &http.BasicAuth{
				Username: "mobyvb", // Your GitHub username
				Password: token,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	*/
}