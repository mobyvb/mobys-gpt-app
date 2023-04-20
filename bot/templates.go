package main

import (
	"bytes"
	"html/template"
)

type File struct {
	Path     string
	Contents string
}

type IssueInfo struct {
	Files   []File
	Subject string
	Body    string
}

func (i IssueInfo) String() string {
	prompt, err := PopulateIssueTemplate(i)
	if err != nil {
		panic(err)
	}
	return prompt
}

func PopulateIssueTemplate(issueInfo IssueInfo) (string, error) {
	tmpl, err := template.ParseFiles("./templates/issue.tmpl")
	if err != nil {
		return "", err
	}

	var result bytes.Buffer
	err = tmpl.Execute(&result, issueInfo)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
