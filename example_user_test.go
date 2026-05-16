package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleUserService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserList),
	)

	users, _ := c.User.List(context.Background())
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_One() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserSingle),
	)

	user, _ := c.User.One(context.Background(), 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Me() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserSingle),
	)

	user, _ := c.User.Me(context.Background())
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserSingle),
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
		backlog.WithDoer(doerUserSingle),
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
		backlog.WithDoer(doerUserSingle),
	)

	user, _ := c.User.Delete(context.Background(), 1)
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Icon() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerUserIcon),
	)

	icon, _ := c.User.Icon(context.Background(), 1)
	fmt.Printf("ContentType: %s, FileName: %s\n", icon.ContentType, icon.Filename)
	// Output:
	// ContentType: image/png, FileName: icon.png
}

func ExampleUserActivityService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerActivityList),
	)

	activities, _ := c.User.Activity.List(context.Background(), 1)
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleUserStarService_Count() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerStarCount),
	)

	count, _ := c.User.Star.Count(context.Background(), 1)
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 42
}

func ExampleUserStarService_List() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerStarList),
	)

	stars, _ := c.User.Star.List(context.Background(), 1)
	fmt.Printf("ID: %d, Title: %s\n", stars[0].ID, stars[0].Title)
	// Output:
	// ID: 10, Title: [TEST-1] first issue
}
