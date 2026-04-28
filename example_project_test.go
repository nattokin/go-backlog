package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// ProjectService
	doerProjectAll       = newMockDoer(fixture.Project.ListJSON)
	doerProjectOne       = newMockDoer(fixture.Project.SingleJSON)
	doerProjectCreate    = newMockDoer(fixture.Project.SingleJSON)
	doerProjectUpdate    = newMockDoer(fixture.Project.SingleJSON)
	doerProjectDelete    = newMockDoer(fixture.Project.SingleJSON)
	doerProjectDiskUsage = newMockDoer(fixture.Project.DiskUsageJSON)

	// ProjectActivityService
	doerProjectActivityList = newMockDoer(fixture.Activity.ListJSON)

	// ProjectSharedFileService
	doerProjectSharedFileList = newMockDoer(fixture.SharedFile.ListJSON)

	// ProjectCategoryService
	doerProjectCategoryAll    = newMockDoer(fixture.Category.ListJSON)
	doerProjectCategoryCreate = newMockDoer(fixture.Category.SingleJSON)
	doerProjectCategoryUpdate = newMockDoer(fixture.Category.SingleJSON)
	doerProjectCategoryDelete = newMockDoer(fixture.Category.SingleJSON)

	// ProjectUserService
	doerProjectUserAll         = newMockDoer(fixture.User.ListJSON)
	doerProjectUserAdd         = newMockDoer(fixture.User.SingleJSON)
	doerProjectUserDelete      = newMockDoer(fixture.User.SingleJSON)
	doerProjectUserAddAdmin    = newMockDoer(fixture.User.SingleJSON)
	doerProjectUserAdminAll    = newMockDoer(fixture.User.ListJSON)
	doerProjectUserDeleteAdmin = newMockDoer(fixture.User.SingleJSON)

	// ProjectWebhookService
	doerProjectWebhookAll    = newMockDoer(fixture.Webhook.ListJSON)
	doerProjectWebhookCreate = newMockDoer(fixture.Webhook.AllEventJSON)
	doerProjectWebhookOne    = newMockDoer(fixture.Webhook.AllEventJSON)
	doerProjectWebhookUpdate = newMockDoer(fixture.Webhook.AllEventJSON)
	doerProjectWebhookDelete = newMockDoer(fixture.Webhook.AllEventJSON)
)

func ExampleProjectService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectAll),
	)

	projects, _ := c.Project.All(context.Background())
	fmt.Printf("ID: %d, Key: %s\n", projects[0].ID, projects[0].ProjectKey)
	// Output:
	// ID: 1, Key: TEST
}

func ExampleProjectService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectOne),
	)

	project, _ := c.Project.One(context.Background(), "TEST")
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_Create() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCreate),
	)

	project, _ := c.Project.Create(context.Background(), "TEST", "test")
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectUpdate),
	)

	project, _ := c.Project.Update(context.Background(), "TEST",
		c.Project.Option.WithName("test"),
	)
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectDelete),
	)

	project, _ := c.Project.Delete(context.Background(), "TEST")
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_DiskUsage() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectDiskUsage),
	)

	usage, _ := c.Project.DiskUsage(context.Background(), "TEST")
	fmt.Printf("ProjectID: %d, Issue: %d\n", usage.ProjectID, usage.Issue)
	// Output:
	// ProjectID: 1, Issue: 11931
}

func ExampleProjectActivityService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectActivityList),
	)

	activities, _ := c.Project.Activity.List(context.Background(), "TEST")
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleProjectCategoryService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCategoryAll),
	)

	categories, _ := c.Project.Category.All(context.Background(), "TEST")
	fmt.Printf("ID: %d, Name: %s\n", categories[0].ID, categories[0].Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectCategoryService_Create() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCategoryCreate),
	)

	category, _ := c.Project.Category.Create(context.Background(), "TEST", "Bug")
	fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectCategoryService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCategoryUpdate),
	)

	category, _ := c.Project.Category.Update(context.Background(), "TEST", 12, "Bug Fixed")
	fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectCategoryService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCategoryDelete),
	)

	category, _ := c.Project.Category.Delete(context.Background(), "TEST", 12)
	fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectSharedFileService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectSharedFileList),
	)

	files, _ := c.Project.SharedFile.List(context.Background(), "TEST")
	fmt.Printf("ID: %d, Name: %s\n", files[0].ID, files[0].Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleProjectUserService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectUserAll),
	)

	users, _ := c.Project.User.All(context.Background(), "TEST", false)
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectUserAdd),
	)

	user, _ := c.Project.User.Add(context.Background(), "TEST", 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectUserDelete),
	)

	user, _ := c.Project.User.Delete(context.Background(), "TEST", 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_AddAdmin() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectUserAddAdmin),
	)

	user, _ := c.Project.User.AddAdmin(context.Background(), "TEST", 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_AdminAll() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectUserAdminAll),
	)

	users, _ := c.Project.User.AdminAll(context.Background(), "TEST")
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_DeleteAdmin() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectUserDeleteAdmin),
	)

	user, _ := c.Project.User.DeleteAdmin(context.Background(), "TEST", 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectWebhookService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectWebhookAll),
	)

	webhooks, _ := c.Project.Webhook.All(context.Background(), "TEST")
	fmt.Printf("ID: %d, Name: %s\n", webhooks[0].ID, webhooks[0].Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_Create() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectWebhookCreate),
	)

	wh, _ := c.Project.Webhook.Create(
		context.Background(),
		"TEST",
		"notify",
		"https://example.com/webhook",
		c.Project.Webhook.Option.WithAllEvent(true),
	)
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectWebhookOne),
	)

	wh, _ := c.Project.Webhook.One(context.Background(), "TEST", 1)
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectWebhookUpdate),
	)

	wh, _ := c.Project.Webhook.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.Webhook.Option.WithName("updated"),
	)
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectWebhookDelete),
	)

	wh, _ := c.Project.Webhook.Delete(context.Background(), "TEST", 1)
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}
