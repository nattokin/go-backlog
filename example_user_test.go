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
