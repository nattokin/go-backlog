package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// WikiOptionService
	doerWikiOption = newMockDoer(fixture.Wiki.MinimumJSON)

	// IssueOptionService
	doerIssueOption = newMockDoer(fixture.Issue.ListJSON)

	// IssueCommentOptionService
	doerIssueCommentOption = newMockDoer(fixture.Comment.ListJSON)

	// ProjectOptionService
	doerProjectOption = newMockDoer(fixture.Project.ListJSON)

	// UserOptionService
	doerUserOption = newMockDoer(fixture.User.SingleJSON)

	// UserRecentlyViewedOptionService
	doerUserRecentlyViewedOption = newMockDoer(fixture.RecentlyViewed.IssueListJSON)

	// UserStarOptionService
	doerUserStarOption = newMockDoer(fixture.Star.ListJSON)

	// PullRequestOptionService
	doerPullRequestOption = newMockDoer(fixture.PullRequest.ListJSON)

	// PullRequestCommentOptionService
	doerPullRequestCommentOption = newMockDoer(fixture.Comment.ListJSON)

	// ActivityOptionService
	doerActivityOptionSpace   = newMockDoer(fixture.Activity.ListJSON)
	doerActivityOptionProject = newMockDoer(fixture.Activity.ListJSON)
	doerActivityOptionUser    = newMockDoer(fixture.Activity.ListJSON)
)

// ExampleWikiOptionService demonstrates updating a wiki page using multiple options.
func ExampleWikiOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiOption),
	)

	wiki, _ := c.Wiki.Update(
		context.Background(),
		34,
		c.Wiki.Option.WithName("Minimum Wiki Page"),
		c.Wiki.Option.WithContent("Updated content."),
	)
	fmt.Printf("ID: %d, Name: %s\n", wiki.ID, wiki.Name)
	// Output:
	// ID: 34, Name: Minimum Wiki Page
}

// ExampleIssueOptionService demonstrates filtering issues using multiple options.
func ExampleIssueOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueOption),
	)

	issues, _ := c.Issue.All(
		context.Background(),
		c.Issue.Option.WithProjectIDs([]int{10}),
		c.Issue.Option.WithKeyword("first"),
		c.Issue.Option.WithCount(50),
	)
	fmt.Printf("Count: %d, ID: %d, Summary: %s\n", len(issues), issues[0].ID, issues[0].Summary)
	// Output:
	// Count: 2, ID: 1, Summary: First issue
}

// ExampleIssueCommentOptionService demonstrates fetching comments with count and order options.
func ExampleIssueCommentOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueCommentOption),
	)

	comments, _ := c.Issue.Comment.All(
		context.Background(),
		"PRJ-1",
		c.Issue.Comment.Option.WithCount(20),
		c.Issue.Comment.Option.WithOrder(backlog.OrderDesc),
	)
	fmt.Printf("Count: %d, ID: %d, Content: %s\n", len(comments), comments[0].ID, comments[0].Content)
	// Output:
	// Count: 2, ID: 1, Content: This is a comment.
}

// ExampleProjectOptionService demonstrates listing all projects including archived ones.
func ExampleProjectOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectOption),
	)

	projects, _ := c.Project.All(
		context.Background(),
		c.Project.Option.WithAll(true),
		c.Project.Option.WithArchived(false),
	)
	fmt.Printf("Count: %d, ProjectKey: %s\n", len(projects), projects[0].ProjectKey)
	// Output:
	// Count: 3, ProjectKey: TEST
}

// ExampleUserOptionService demonstrates updating a user's name and role type.
func ExampleUserOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserOption),
	)

	user, _ := c.User.Update(
		context.Background(),
		1,
		c.User.Option.WithName("admin"),
		c.User.Option.WithRoleType(backlog.RoleAdministrator),
	)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

// ExampleUserRecentlyViewedOptionService demonstrates listing recently viewed issues with count and order.
func ExampleUserRecentlyViewedOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserRecentlyViewedOption),
	)

	issues, _ := c.User.RecentlyViewed.ListIssues(
		context.Background(),
		c.User.RecentlyViewed.Option.WithCount(20),
		c.User.RecentlyViewed.Option.WithOrder(backlog.OrderDesc),
	)
	fmt.Printf("IssueKey: %s\n", issues[0].IssueKey)
	// Output:
	// IssueKey: TEST-1
}

// ExampleUserStarOptionService demonstrates listing received stars with count and order.
func ExampleUserStarOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserStarOption),
	)

	stars, _ := c.User.Star.List(
		context.Background(),
		1,
		c.User.Star.Option.WithCount(20),
		c.User.Star.Option.WithOrder(backlog.OrderDesc),
	)
	fmt.Printf("ID: %d, Title: %s\n", stars[0].ID, stars[0].Title)
	// Output:
	// ID: 10, Title: [TEST-1] first issue
}

// ExamplePullRequestOptionService demonstrates filtering pull requests by status and count.
func ExamplePullRequestOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestOption),
	)

	prs, _ := c.PullRequest.All(
		context.Background(),
		"TEST",
		"myrepo",
		c.PullRequest.Option.WithStatusIDs([]int{1}),
		c.PullRequest.Option.WithCount(50),
	)
	fmt.Printf("Count: %d, ID: %d, Summary: %s\n", len(prs), prs[0].ID, prs[0].Summary)
	// Output:
	// Count: 2, ID: 2, Summary: test PR
}

// ExamplePullRequestCommentOptionService demonstrates fetching pull request comments with count and order.
func ExamplePullRequestCommentOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerPullRequestCommentOption),
	)

	comments, _ := c.PullRequest.Comment.All(
		context.Background(),
		"TEST",
		"myrepo",
		1,
		c.PullRequest.Comment.Option.WithCount(20),
		c.PullRequest.Comment.Option.WithOrder(backlog.OrderDesc),
	)
	fmt.Printf("Count: %d, ID: %d, Content: %s\n", len(comments), comments[0].ID, comments[0].Content)
	// Output:
	// Count: 2, ID: 1, Content: This is a comment.
}

// ExampleStarOptionService demonstrates adding a star to an issue.
func ExampleStarOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerNoContent),
	)

	err := c.Star.Add(context.Background(), c.Star.Option.WithIssueID(1))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

// ExampleActivityOptionService_spaceWithOptions demonstrates filtering space activities by type and count.
func ExampleActivityOptionService_spaceWithOptions() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerActivityOptionSpace),
	)

	activities, _ := c.Space.Activity.List(
		context.Background(),
		c.Space.Activity.Option.WithActivityTypeIDs([]int{1, 2}),
		c.Space.Activity.Option.WithCount(20),
	)
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

// ExampleActivityOptionService_projectWithOptions demonstrates filtering project activities by type and order.
func ExampleActivityOptionService_projectWithOptions() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerActivityOptionProject),
	)

	activities, _ := c.Project.Activity.List(
		context.Background(),
		"MYPROJECT",
		c.Project.Activity.Option.WithActivityTypeIDs([]int{1, 2}),
		c.Project.Activity.Option.WithOrder(backlog.OrderDesc),
	)
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

// ExampleActivityOptionService_userWithOptions demonstrates filtering user activities by type and count.
func ExampleActivityOptionService_userWithOptions() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerActivityOptionUser),
	)

	activities, _ := c.User.Activity.List(
		context.Background(),
		1,
		c.User.Activity.Option.WithActivityTypeIDs([]int{1, 2}),
		c.User.Activity.Option.WithCount(20),
	)
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}
