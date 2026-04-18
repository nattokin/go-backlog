package backlog_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

type projectFixedDoer struct {
	body string
}

func (d *projectFixedDoer) Do(_ *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(d.body)),
	}, nil
}

var (
	// ProjectService
	doerProjectAll    = &projectFixedDoer{body: fixture.Project.ListJSON}
	doerProjectOne    = &projectFixedDoer{body: fixture.Project.SingleJSON}
	doerProjectCreate = &projectFixedDoer{body: fixture.Project.SingleJSON}
	doerProjectUpdate = &projectFixedDoer{body: fixture.Project.SingleJSON}
	doerProjectDelete = &projectFixedDoer{body: fixture.Project.SingleJSON}

	// ProjectActivityService
	doerProjectActivityList = &projectFixedDoer{body: fixture.Activity.ListJSON}

	// ProjectUserService
	doerProjectUserAll         = &projectFixedDoer{body: fixture.User.ListJSON}
	doerProjectUserAdd         = &projectFixedDoer{body: fixture.User.SingleJSON}
	doerProjectUserDelete      = &projectFixedDoer{body: fixture.User.SingleJSON}
	doerProjectUserAddAdmin    = &projectFixedDoer{body: fixture.User.SingleJSON}
	doerProjectUserAdminAll    = &projectFixedDoer{body: fixture.User.ListJSON}
	doerProjectUserDeleteAdmin = &projectFixedDoer{body: fixture.User.SingleJSON}
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
