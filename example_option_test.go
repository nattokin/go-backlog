package backlog_test

import (
	"context"
	"fmt"
	"time"

	backlog "github.com/nattokin/go-backlog"
)

// ExampleWikiOptionService demonstrates updating a wiki page using multiple options.
func ExampleWikiOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWikiSingle),
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
		backlog.WithDoer(doerIssueList),
	)

	issues, _ := c.Issue.List(
		context.Background(),
		c.Issue.Option.WithProjectIDs([]int{10}),
		c.Issue.Option.WithKeyword("first"),
		c.Issue.Option.WithCount(50),
	)
	fmt.Printf("Count: %d, ID: %d, Summary: %s\n", len(issues), issues[0].ID, issues[0].Summary)
	// Output:
	// Count: 2, ID: 1, Summary: First issue
}

// ExampleIssueOptionService_customField demonstrates updating issues with custom fields.
func ExampleIssueOptionService_customField() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueSingle),
	)

	params := map[string]struct {
		id        int
		numValue  float64
		strValue  string
		timeValue time.Time
		itemIDs   []int
	}{
		"Num": {
			id:       1,
			numValue: 1.5,
		},
		"String": {
			id:       2,
			strValue: "test",
		},
		"Time": {
			id:        3,
			timeValue: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"Items": {
			id:      4,
			itemIDs: []int{1, 2},
		},
		"Other": {
			id:       4,
			strValue: "test",
		},
	}

	issue, _ := c.Issue.Update(
		context.Background(),
		"Issue-1",
		c.Issue.Option.WithCustomFieldNum(params["Num"].id, params["Num"].numValue),
		c.Issue.Option.WithCustomFieldString(params["String"].id, params["String"].strValue),
		c.Issue.Option.WithCustomFieldTime(params["Time"].id, params["Time"].timeValue),

		// Select items from the list custom field.
		c.Issue.Option.WithCustomFieldItems(params["Items"].id, params["Items"].itemIDs),
		// Set the "Other" text value for the list custom field.
		c.Issue.Option.WithCustomFieldOther(params["Other"].id, params["Other"].strValue),
	)
	fmt.Printf("ID: %d, Summary: %s\n", issue.ID, issue.Summary)
	// Output:
	// ID: 1, Summary: First issue
}

// ExampleIssueCommentOptionService demonstrates fetching comments with count and order options.
func ExampleIssueCommentOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerCommentList),
	)

	comments, _ := c.Issue.Comment.List(
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
		backlog.WithDoer(doerProjectList),
	)

	projects, _ := c.Project.List(
		context.Background(),
		c.Project.Option.WithAll(true),
		c.Project.Option.WithArchived(false),
	)
	fmt.Printf("Count: %d, ProjectKey: %s\n", len(projects), projects[0].ProjectKey)
	// Output:
	// Count: 3, ProjectKey: TEST
}

// ExampleUserRecentlyViewedOptionService demonstrates listing recently viewed issues with count and order.
func ExampleRecentlyViewedOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerRecentlyViewedIssueList),
	)

	issues, _ := c.RecentlyViewed.ListIssues(
		context.Background(),
		c.RecentlyViewed.Option.WithCount(20),
		c.RecentlyViewed.Option.WithOrder(backlog.OrderDesc),
	)
	fmt.Printf("IssueKey: %s\n", issues[0].IssueKey)
	// Output:
	// IssueKey: TEST-1
}

// ExampleUserOptionService demonstrates updating a user's name and role type.
func ExampleUserOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserSingle),
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

// ExampleUserStarOptionService demonstrates listing received stars with count and order.
func ExampleUserStarOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerStarList),
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
		backlog.WithDoer(doerPullRequestList),
	)

	prs, _ := c.PullRequest.List(
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
		backlog.WithDoer(doerCommentList),
	)

	comments, _ := c.PullRequest.Comment.List(
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
		backlog.WithDoer(doerActivityList),
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
		backlog.WithDoer(doerActivityList),
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
		backlog.WithDoer(doerActivityList),
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

// ExampleProjectCustomFieldOptionService demonstrates updating a custom field using multiple options.
func ExampleProjectCustomFieldOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerCustomFieldSingle),
	)

	field, _ := c.Project.CustomField.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.CustomField.Option.WithName("Sprint"),
		c.Project.CustomField.Option.WithDescription("Sprint number"),
		c.Project.CustomField.Option.WithRequired(true),
	)
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

// ExampleProjectIssueTypeOptionService demonstrates updating an issue type using multiple options.
func ExampleProjectIssueTypeOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerIssueTypeSingle),
	)

	issueType, _ := c.Project.IssueType.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.IssueType.Option.WithName("Bug"),
		c.Project.IssueType.Option.WithColor("#e30000"),
	)
	fmt.Printf("ID: %d, Name: %s, Color: %s\n", issueType.ID, issueType.Name, issueType.Color)
	// Output:
	// ID: 1, Name: Bug, Color: #e30000
}

// ExampleProjectStatusOptionService demonstrates updating a status using multiple options.
func ExampleProjectStatusOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerStatusSingle),
	)

	status, _ := c.Project.Status.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.Status.Option.WithName("Open"),
		c.Project.Status.Option.WithColor("#ed8077"),
	)
	fmt.Printf("ID: %d, Name: %s\n", status.ID, status.Name)
	// Output:
	// ID: 1, Name: Open
}

// ExampleProjectVersionOptionService demonstrates updating a version using multiple options.
func ExampleProjectVersionOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerVersionSingle),
	)

	version, _ := c.Project.Version.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.Version.Option.WithName("Version 1.0"),
		c.Project.Version.Option.WithDescription("First stable release"),
	)
	fmt.Printf("ID: %d, Name: %s\n", version.ID, version.Name)
	// Output:
	// ID: 1, Name: Version 1.0
}

// ExampleProjectWebhookOptionService demonstrates updating a webhook using multiple options.
func ExampleProjectWebhookOptionService() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerWebhookAllEvent),
	)

	wh, _ := c.Project.Webhook.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.Webhook.Option.WithName("Example Webhook"),
		c.Project.Webhook.Option.WithAllEvent(true),
	)
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}
