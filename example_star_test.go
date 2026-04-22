package backlog_test

import (
	"context"
	"fmt"
	"net/http"

	backlog "github.com/nattokin/go-backlog"
)

var (
	// StarService
	doerStarAdd    = newMockDoer("")
	doerStarRemove = newMockDoer("")
)

func ExampleStarService_Add() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(&mockDoer{do: func(_ *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
		}}),
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
		backlog.WithDoer(&mockDoer{do: func(_ *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: http.StatusNoContent, Body: http.NoBody}, nil
		}}),
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
