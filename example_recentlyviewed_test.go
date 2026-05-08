package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	doerRecentlyViewedListIssues   = newMockDoer(fixture.RecentlyViewed.IssueListJSON)
	doerRecentlyViewedAddIssue     = newMockDoer(fixture.RecentlyViewed.IssueSingleJSON)
	doerRecentlyViewedListProjects = newMockDoer(fixture.RecentlyViewed.ProjectListJSON)
	doerRecentlyViewedListWikis    = newMockDoer(fixture.RecentlyViewed.WikiListJSON)
	doerRecentlyViewedAddWiki      = newMockDoer(fixture.RecentlyViewed.WikiSingleJSON)
)

func ExampleRecentlyViewedService_ListIssues() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedListIssues),
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
		backlog.WithDoer(doerRecentlyViewedAddIssue),
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
		backlog.WithDoer(doerRecentlyViewedListProjects),
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
		backlog.WithDoer(doerRecentlyViewedListWikis),
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
		backlog.WithDoer(doerRecentlyViewedAddWiki),
	)

	wiki, _ := c.RecentlyViewed.AddWiki(context.Background(), 10)
	fmt.Printf("Name: %s\n", wiki.Name)
	// Output:
	// Name: Home
}
