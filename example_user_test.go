package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

func ExampleUserService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	users, err := c.User.List(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", users[0].ID, users[0].UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_One() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.User.One(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Me() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.User.Me(context.Background())
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Add() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.User.Add(
		context.Background(),
		"admin",
		"password",
		"admin",
		"eguchi@nulab.example",
		backlog.RoleAdministrator,
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Update() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.User.Update(
		context.Background(),
		1,
		c.User.Option.WithName("admin"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Delete() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	user, err := c.User.Delete(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, UserID: %s\n", user.ID, user.UserID)
	// Output:
	// ID: 1, UserID: admin
}

func ExampleUserService_Icon() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	icon, err := c.User.Icon(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ContentType: %s, FileName: %s\n", icon.ContentType, icon.Filename)
	// Output:
	// ContentType: image/png, FileName: icon.png
}

func ExampleUserActivityService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	activities, err := c.User.Activity.List(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Type: %d\n", activities[0].ID, activities[0].Type)
	// Output:
	// ID: 3153, Type: 2
}

func ExampleUserStarService_Count() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	count, err := c.User.Star.Count(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("Count: %d\n", count)
	// Output:
	// Count: 42
}

func ExampleUserStarService_List() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	stars, err := c.User.Star.List(context.Background(), 1)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Title: %s\n", stars[0].ID, stars[0].Title)
	// Output:
	// ID: 10, Title: [TEST-1] first issue
}
