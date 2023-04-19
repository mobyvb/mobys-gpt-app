# Moby's cool app

April 19, midnight:27

It's time for notes

So I was messing around with GPT

And I thought it would be a fun project to integrate it with Github for like personal projects and stuff.

I've already tried pair programming with it a few times (a couple times with GPT 3 and once with GPT 4) and it showed some pretty promising results.

So here's the plan.

## The Plan

I create a script/bot thing that talks with Github. It can listen for new issues, create pull requests and stuff. It'll sit on its own server with git or whatever.

So I would then create an issue, and be like "change this part of the app" or "fix this thing" or "create an alternate version of the homepage"

And then it would create a pull request, and include context in the body and update the files and all that jazz.

Sounds crazy, right? Like that's just completely ridiculous. Obviously it would probably struggle with really large apps? Would it even work with a really small app?

Would it?

How would one even send adequate information over to an LLM so that it understands the state of the repository, and generates a git commit? 

So I tried this prompt:

### START OF PROMPT

Files:
  - name: main.go
    contents:
    ```
    package main

    func main() {

    }
    ```
  - name: index.html
    contents:
    ```
    <html>
        <head></head>
        <body></body>
    </html>
    ```

Modify the files above to accomplish the following task:
Subject: Create static server
Body:
  - Update the homepage (index.html) with basic CSS and a header "moby's cool app"
  - Serve the html file from a Go server in main.go

Respond in the exact format:
Files:
  - name: main.go
    contents:
    ```
    <new main.go contents>
    ```
  - name: index.html
    contents:
    ```
    <new index.html contents>
    ```
Notes:
<additional context about your changes>

### END OF PROMPT

And GPT-4 responded in the exact format I requested. So now, this can be a template.

p.s. check ./notes/sample-prompt-chat.txt to see the first chat I tried out.

## The Plan, but more specific

* Integrate with Github to discover new issues (and pull latest `main` branch)
    - Pull author, label, subject, and body from issue
* List files and contents, and insert into template based on prompt above
    - Exclude files in a `.gptignore` file
    - Check author against a `.allowedauthors` file and ignore if not present
* Parse Github issue details into template based on prompt above.
* Send prompt to OpenAI API, response should be in the exact format requested
* Parse response from GPT-4 into a git change
    - Replace files with files based on "Files:" from the output
    - Git commit with message based on "additional context about your changes" from the output
    - Push to a new branch
    - Request to create a PR in Github
* For each comment on the PR during a code review, another prompt request should be triggered
    - Prompt would be kind of similar, but with a different prompt, below
* Update the PR, similarly to creating a new one.
* PR's must be manually merged

### START OF PROMPT

Files:
  - name: main.go
    contents:
    ```
    package main

    import (
        "io/ioutil"
        "log"
        "net/http"
    )

    func main() {
        http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
            content, err := ioutil.ReadFile("index.html")
            if err != nil {
                log.Fatal(err)
            }
            w.Write(content)
        })

        log.Println("Starting server on :8080")
        err := http.ListenAndServe(":8080", nil)
        if err != nil {
            log.Fatal(err)
        }
    }
    ```

Modify the file above to address the following comment:
Line: `err := http.ListenAndServe(":8080", nil)`
Comment: Please update this to use any available port rather than a static one.

Respond in the exact format:
Files:
  - name: main.go
    contents:
    ```
    <new main.go contents>
    ```
Notes:
<additional context about your changes>

### END OF PROMPT

see ./notes/sample-prompt-chat.txt (or git history) for results.
