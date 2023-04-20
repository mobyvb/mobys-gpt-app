package main

import (
	"bytes"
	"html/template"
	"strings"
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

type IssueResponse struct {
	Files []File
	Notes string
}

func (r IssueResponse) String() string {
	out := "Notes:\n"
	out += r.Notes + "\n\n"
	out += "Files:\n"
	for _, f := range r.Files {
		out += f.Path + ":\n```\n"
		out += f.Contents + "\n```\n"
	}

	return out
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

func ParseIssueResponse(response string) IssueResponse {
	sections := strings.Split(response, "Notes:")

	filesSection := sections[0]
	notes := strings.TrimSpace(sections[1])

	files := parseFiles(filesSection)

	return IssueResponse{
		Files: files,
		Notes: notes,
	}
}

func parseFiles(filesSection string) []File {
	fileStringList := strings.Split(filesSection, "name:")
	// first item in the list is just gonna be "Files:"
	fileStringList = fileStringList[1:]
	fileList := make([]File, len(fileStringList))
	for i, f := range fileStringList {
		fileParts := strings.Split(f, "contents:")
		path := strings.TrimSpace(fileParts[0])
		contents := strings.TrimSpace(fileParts[1])
		contents = strings.Trim(contents, "```")

		fileList[i] = File{
			Path:     path,
			Contents: contents,
		}
	}

	return fileList
}
