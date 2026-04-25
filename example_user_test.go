package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// UserService
	doerUserAll    = newMockDoer(fixture.User.ListJSON)
	doerUserOne    = newMockDoer(fixture.User.SingleJSON)
	doerUserOwn    = newMockDoer(fixture.User.SingleJSON)
	doerUserAdd    = newMockDoer(fixture.User.SingleJSON)
	doerUserUpdate = newMockDoer(fixture.User.SingleJSON)
	doerUserDelete = newMockDoer(fixture.User.SingleJSON)

	// UserActivityService
	doerUserActivityList = newMockDoer(fixture.Activity.ListJSON)

	// UserRecentlyViewedService
	doerRecentlyViewedListIssues   = newMockDoer(fixture.RecentlyViewed.IssueListJSON)
	doerRecentlyViewedAddIssue     = newMockDoer(fixture.RecentlyViewed.IssueSingleJSON)
	doerRecentlyViewedListProjects = newMockDoer(fixture.RecentlyViewed.ProjectListJSON)
	doerRecentlyViewedListWikis    = newMockDoer(fixture.RecentlyViewed.WikiListJSON)
	doerRecentlyViewedAddWiki      = newMockDoer(fixture.RecentlyViewed.WikiSingleJSON)
)

func ExampleUserService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserAll),
	)

	users, _ := c.User.All(context.Background())
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserOne),
	)

	user, _ := c.User.One(context.Background(), 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Own() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserOwn),
	)

	user, _ := c.User.Own(context.Background())
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserAdd),
	)

	user, _ := c.User.Add(
		context.Background(),
		"admin",
		"password",
		"admin",
		"eguchi@nulab.example",
		backlog.RoleAdministrator,
	)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserUpdate),
	)

	user, _ := c.User.Update(
		context.Background(),
		1,
		c.User.Option.WithName("admin"),
	)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserDelete),
	)

	user, _ := c.User.Delete(context.Background(), 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserActivityService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserActivityList),
	)

	activities, _ := c.User.Activity.List(context.Background(), 1)
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleUserRecentlyViewedService_ListIssues() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedListIssues),
	)

	issues, _ := c.User.RecentlyViewed.ListIssues(context.Background())
	fmt.Printf("IssueKey: %s\n", issues[0].IssueKey)
	// Output:
	// IssueKey: TEST-1
}

func ExampleUserRecentlyViewedService_AddIssue() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedAddIssue),
	)

	issue, _ := c.User.RecentlyViewed.AddIssue(context.Background(), 1)
	fmt.Printf("IssueKey: %s\n", issue.IssueKey)
	// Output:
	// IssueKey: TEST-1
}

func ExampleUserRecentlyViewedService_ListProjects() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedListProjects),
	)

	projects, _ := c.User.RecentlyViewed.ListProjects(context.Background())
	fmt.Printf("ProjectKey: %s\n", projects[0].ProjectKey)
	// Output:
	// ProjectKey: TEST
}

func ExampleUserRecentlyViewedService_ListWikis() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedListWikis),
	)

	wikis, _ := c.User.RecentlyViewed.ListWikis(context.Background())
	fmt.Printf("Name: %s\n", wikis[0].Name)
	// Output:
	// Name: Home
}

func ExampleUserRecentlyViewedService_AddWiki() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedAddWiki),
	)

	wiki, _ := c.User.RecentlyViewed.AddWiki(context.Background(), 10)
	fmt.Printf("Name: %s\n", wiki.Name)
	// Output:
	// Name: Home
}
