package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nattokin/go-backlog"
)

func main() {
	baseURL := os.Getenv("BACKLOG_BASE_URL")
	if baseURL == "" {
		log.Fatalln("You need Backlog base url.")
	}
	token := os.Getenv("BACKLOG_TOKEN")
	if token == "" {
		log.Fatalln("You need Backlog access token.")
	}

	statusFlag := flag.String("status", "", "comma-separated status IDs to filter (e.g. 1,2,3)")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		log.Fatalln("Usage: go run . <PROJECT_KEY> [--status <id,...>]")
	}
	projectKey := args[0]

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	// Resolve project key to project ID required by Issue.All.
	project, err := c.Project.One(ctx, projectKey)
	if err != nil {
		log.Fatalf("failed to get project %q: %v", projectKey, err)
	}

	opts := []backlog.RequestOption{
		c.Issue.Option.WithProjectIDs([]int{project.ID}),
		c.Issue.Option.WithCount(100),
	}

	if *statusFlag != "" {
		ids := parseIntList(*statusFlag)
		if len(ids) > 0 {
			opts = append(opts, c.Issue.Option.WithStatusIDs(ids))
		}
	}

	issues, err := c.Issue.List(ctx, opts...)
	if err != nil {
		log.Fatalf("failed to fetch issues: %v", err)
	}

	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'

	_ = w.Write([]string{"ID", "Key", "Summary", "Status", "Assignee", "Priority", "Comments", "Attachments", "Created", "Updated"})

	for _, issue := range issues {
		comments, err := c.Issue.Comment.List(ctx, issue.IssueKey)
		if err != nil {
			log.Printf("warning: failed to fetch comments for %s: %v", issue.IssueKey, err)
		}

		attachments, err := c.Issue.Attachment.List(ctx, issue.IssueKey)
		if err != nil {
			log.Printf("warning: failed to fetch attachments for %s: %v", issue.IssueKey, err)
		}

		statusName := ""
		if issue.Status != nil {
			statusName = issue.Status.Name
		}
		assigneeName := ""
		if issue.Assignee != nil {
			assigneeName = issue.Assignee.Name
		}
		priorityName := ""
		if issue.Priority != nil {
			priorityName = issue.Priority.Name
		}

		_ = w.Write([]string{
			strconv.Itoa(issue.ID),
			issue.IssueKey,
			issue.Summary,
			statusName,
			assigneeName,
			priorityName,
			strconv.Itoa(len(comments)),
			strconv.Itoa(len(attachments)),
			issue.Created.String(),
			issue.Updated.String(),
		})
	}

	w.Flush()
}

// parseIntList parses a comma-separated string of integers.
func parseIntList(s string) []int {
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		var n int
		if _, err := fmt.Sscanf(p, "%d", &n); err == nil {
			result = append(result, n)
		}
	}
	return result
}
