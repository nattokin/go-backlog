package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleRecentlyViewedService_ListIssues() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedIssueList),
	)

	issues, _ := c.RecentlyViewed.ListIssues(context.Background())
	fmt.Printf("IssueKey: %s\n", issues[0].IssueKey)
	// Output:
	// IssueKey: TEST-1
}

func ExampleRecentlyViewedService_AddIssue() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedIssueSingle),
	)

	issue, _ := c.RecentlyViewed.AddIssue(context.Background(), 1)
	fmt.Printf("IssueKey: %s\n", issue.IssueKey)
	// Output:
	// IssueKey: TEST-1
}

func ExampleRecentlyViewedService_ListProjects() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedProjectList),
	)

	projects, _ := c.RecentlyViewed.ListProjects(context.Background())
	fmt.Printf("ProjectKey: %s\n", projects[0].ProjectKey)
	// Output:
	// ProjectKey: TEST
}

func ExampleRecentlyViewedService_ListWikis() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedWikiList),
	)

	wikis, _ := c.RecentlyViewed.ListWikis(context.Background())
	fmt.Printf("Name: %s\n", wikis[0].Name)
	// Output:
	// Name: Home
}

func ExampleRecentlyViewedService_AddWiki() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedWikiSingle),
	)

	wiki, _ := c.RecentlyViewed.AddWiki(context.Background(), 10)
	fmt.Printf("Name: %s\n", wiki.Name)
	// Output:
	// Name: Home
}
