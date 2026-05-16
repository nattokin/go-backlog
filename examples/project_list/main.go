package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	projects, err := c.Project.List(ctx)
	if err != nil {
		log.Fatalf("failed to fetch projects: %v", err)
	}

	fmt.Printf("%d project(s) found.\n", len(projects))

	for _, p := range projects {
		fmt.Printf("\n[%s] %s (ID: %d)\n", p.ProjectKey, p.Name, p.ID)

		issueTypes, err := c.Project.IssueType.List(ctx, p.ProjectKey)
		if err != nil {
			log.Printf("warning: failed to fetch issue types for %s: %v", p.ProjectKey, err)
		}
		names := make([]string, len(issueTypes))
		for i, t := range issueTypes {
			names[i] = t.Name
		}
		fmt.Printf("  Issue Types : %s\n", strings.Join(names, ", "))

		statuses, err := c.Project.Status.List(ctx, p.ProjectKey)
		if err != nil {
			log.Printf("warning: failed to fetch statuses for %s: %v", p.ProjectKey, err)
		}
		names = make([]string, len(statuses))
		for i, s := range statuses {
			names[i] = s.Name
		}
		fmt.Printf("  Statuses    : %s\n", strings.Join(names, ", "))
	}
}
