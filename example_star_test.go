package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
)

var (
	// StarService
	doerStarAdd    = &mockDoer{do: func(_ *http.Request) (*http.Response, error) { return &http.Response{StatusCode: 204, Body: http.NoBody}, nil }}
	doerStarRemove = &mockDoer{do: func(_ *http.Request) (*http.Response, error) { return &http.Response{StatusCode: 204, Body: http.NoBody}, nil }}
)

func ExampleStarService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerStarAdd),
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

func ExampleStarService_Remove() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerStarRemove),
	)

	err := c.Star.Remove(context.Background(), 42)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}
