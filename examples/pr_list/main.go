package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nattokin/go-backlog"
)

// prStatusIDs maps status name to Backlog pull request status ID constants.
var prStatusIDs = map[string]int{
	"open":   backlog.PullRequestStatusOpen,
	"closed": backlog.PullRequestStatusClosed,
	"merged": backlog.PullRequestStatusMerged,
	"draft":  backlog.PullRequestStatusDraft,
}

func main() {
	baseURL := os.Getenv("BACKLOG_BASE_URL")
	if baseURL == "" {
		log.Fatalln("You need Backlog base url.")
	}
	token := os.Getenv("BACKLOG_TOKEN")
	if token == "" {
		log.Fatalln("You need Backlog access token.")
	}

	statusFlag := flag.String("status", "", "filter by status: open, closed, merged, draft")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Usage: go run . [--status open|closed|merged|draft] <PROJECT_KEY> <REPO_NAME>")
	}
	projectKey := args[0]
	repoName := args[1]

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	opts := []backlog.RequestOption{}
	if *statusFlag != "" {
		id, ok := prStatusIDs[*statusFlag]
		if !ok {
			log.Fatalf("unknown status %q: must be open, closed, merged, or draft", *statusFlag)
		}
		opts = append(opts, c.PullRequest.Option.WithStatusIDs([]int{id}))
	}

	total, err := c.PullRequest.Count(ctx, projectKey, repoName, opts...)
	if err != nil {
		log.Fatalf("failed to count pull requests: %v", err)
	}
	fmt.Printf("%d pull request(s) found in %s/%s\n", total, projectKey, repoName)

	prs, err := c.PullRequest.List(ctx, projectKey, repoName, opts...)
	if err != nil {
		log.Fatalf("failed to fetch pull requests: %v", err)
	}

	for _, pr := range prs {
		statusName := ""
		if pr.Status != nil {
			statusName = pr.Status.Name
		}
		assigneeName := ""
		if pr.Assignee != nil {
			assigneeName = pr.Assignee.Name
		}

		comments, err := c.PullRequest.Comment.List(ctx, projectKey, repoName, pr.Number)
		if err != nil {
			log.Printf("warning: failed to fetch comments for PR #%d: %v", pr.Number, err)
		}

		fmt.Printf("#%d [%s] %s\n", pr.Number, statusName, pr.Summary)
		fmt.Printf("  branch   : %s -> %s\n", pr.Branch, pr.Base)
		fmt.Printf("  assignee : %s\n", assigneeName)
		fmt.Printf("  comments : %d\n", len(comments))
		fmt.Printf("  created  : %s\n", pr.Created.String())
	}
}
