package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleProjectService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	projects, err := c.Project.List(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Key: %s\n", projects[0].ID, projects[0].ProjectKey)
	// Output:
	// ID: 1, Key: TEST
}

func ExampleProjectService_One() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	project, err := c.Project.One(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	project, err := c.Project.Create(context.Background(), "TEST", "test")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	project, err := c.Project.Update(context.Background(), "TEST",
		c.Project.Option.WithName("test"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	project, err := c.Project.Delete(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Key: %s\n", project.ID, project.ProjectKey)
	// Output:
	// ID: 6, Key: TEST
}

func ExampleProjectService_DiskUsage() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	usage, err := c.Project.DiskUsage(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ProjectID: %d, Issue: %d\n", usage.ProjectID, usage.Issue)
	// Output:
	// ProjectID: 1, Issue: 11931
}

func ExampleProjectService_Icon() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	icon, err := c.Project.Icon(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ContentType: %s, FileName: %s\n", icon.ContentType, icon.Filename)
	// Output:
	// ContentType: image/png, FileName: test.png
}

func ExampleProjectActivityService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	activities, err := c.Project.Activity.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleProjectCategoryService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	categories, err := c.Project.Category.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", categories[0].ID, categories[0].Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectCategoryService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	category, err := c.Project.Category.Create(context.Background(), "TEST", "Bug")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectCategoryService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	category, err := c.Project.Category.Update(context.Background(), "TEST", 12, "Bug Fixed")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectCategoryService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	category, err := c.Project.Category.Delete(context.Background(), "TEST", 12)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", category.ID, category.Name)
	// Output:
	// ID: 12, Name: Bug
}

func ExampleProjectCustomFieldService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fields, err := c.Project.CustomField.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", fields[0].ID, fields[0].Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	field, err := c.Project.CustomField.Create(context.Background(), "TEST", backlog.CustomFieldTypeText, "Sprint")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	field, err := c.Project.CustomField.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.CustomField.Option.WithName("Sprint Updated"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	field, err := c.Project.CustomField.Delete(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_AddListItem() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	field, err := c.Project.CustomField.AddListItem(context.Background(), "TEST", 1, "Item1")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_UpdateListItem() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	field, err := c.Project.CustomField.UpdateListItem(context.Background(), "TEST", 1, 10, "Item1 Updated")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_DeleteListItem() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	field, err := c.Project.CustomField.DeleteListItem(context.Background(), "TEST", 1, 10)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectIssueTypeService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	issueTypes, err := c.Project.IssueType.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s, Color: %s\n", issueTypes[0].ID, issueTypes[0].Name, issueTypes[0].Color)
	// Output:
	// ID: 1, Name: Bug, Color: #e30000
}

func ExampleProjectIssueTypeService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	issueType, err := c.Project.IssueType.Create(context.Background(), "TEST", "Bug", "#e30000")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s, Color: %s\n", issueType.ID, issueType.Name, issueType.Color)
	// Output:
	// ID: 1, Name: Bug, Color: #e30000
}

func ExampleProjectIssueTypeService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	issueType, err := c.Project.IssueType.Update(context.Background(), "TEST", 1,
		c.Project.IssueType.Option.WithName("Bug Updated"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", issueType.ID, issueType.Name)
	// Output:
	// ID: 1, Name: Bug
}

func ExampleProjectIssueTypeService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	issueType, err := c.Project.IssueType.Delete(context.Background(), "TEST", 1, 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", issueType.ID, issueType.Name)
	// Output:
	// ID: 1, Name: Bug
}

func ExampleProjectSharedFileService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	files, err := c.Project.SharedFile.List(context.Background(), "TEST", "/design/")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", files[0].ID, files[0].Name)
	// Output:
	// ID: 454403, Name: 01_buz.png
}

func ExampleProjectSharedFileService_Download() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	file, err := c.Project.SharedFile.Download(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ContentType: %s, FileName: %s\n", file.ContentType, file.Filename)
	// Output:
	// ContentType: image/png, FileName: shared.png
}

func ExampleProjectUserService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	users, err := c.Project.User.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_List_excludeGroupMembers() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	users, err := c.Project.User.List(
		context.Background(),
		"TEST",
		c.Project.User.Option.WithExcludeGroupMembers(true),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_Add() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.Project.User.Add(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.Project.User.Delete(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_AddAdmin() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.Project.User.AddAdmin(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_AdminList() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	users, err := c.Project.User.AdminList(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectUserService_DeleteAdmin() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.Project.User.DeleteAdmin(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleProjectWebhookService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	webhooks, err := c.Project.Webhook.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", webhooks[0].ID, webhooks[0].Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wh, err := c.Project.Webhook.Create(
		context.Background(),
		"TEST",
		"notify",
		"https://example.com/webhook",
		c.Project.Webhook.Option.WithAllEvent(true),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_One() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wh, err := c.Project.Webhook.One(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wh, err := c.Project.Webhook.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.Webhook.Option.WithName("updated"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectWebhookService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wh, err := c.Project.Webhook.Delete(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", wh.ID, wh.Name)
	// Output:
	// ID: 1, Name: Example Webhook
}

func ExampleProjectStatusService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	statuses, err := c.Project.Status.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", statuses[0].ID, statuses[0].Name)
	// Output:
	// ID: 1, Name: Open
}

func ExampleProjectStatusService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	status, err := c.Project.Status.Create(context.Background(), "TEST", "Open", "#ed8077")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", status.ID, status.Name)
	// Output:
	// ID: 1, Name: Open
}

func ExampleProjectStatusService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	status, err := c.Project.Status.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.Status.Option.WithName("Updated"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", status.ID, status.Name)
	// Output:
	// ID: 1, Name: Open
}

func ExampleProjectStatusService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	status, err := c.Project.Status.Delete(context.Background(), "TEST", 1, 2)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", status.ID, status.Name)
	// Output:
	// ID: 1, Name: Open
}

func ExampleProjectStatusService_UpdateOrder() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	statuses, err := c.Project.Status.UpdateOrder(context.Background(), "TEST", []int{1, 2})
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", statuses[0].ID, statuses[0].Name)
	// Output:
	// ID: 1, Name: Open
}

func ExampleProjectVersionService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	versions, err := c.Project.Version.List(context.Background(), "TEST")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", versions[0].ID, versions[0].Name)
	// Output:
	// ID: 1, Name: Version 1.0
}

func ExampleProjectVersionService_Create() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	version, err := c.Project.Version.Create(
		context.Background(),
		"TEST",
		"ver1",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", version.ID, version.Name)
	// Output:
	// ID: 1, Name: Version 1.0
}

func ExampleProjectVersionService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	version, err := c.Project.Version.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.Version.Option.WithName("updated"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", version.ID, version.Name)
	// Output:
	// ID: 1, Name: Version 1.0
}

func ExampleProjectVersionService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	version, err := c.Project.Version.Delete(context.Background(), "TEST", 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", version.ID, version.Name)
	// Output:
	// ID: 1, Name: Version 1.0
}
