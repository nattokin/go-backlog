package backlog_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	backlog "github.com/nattokin/go-backlog"
)

// Example demonstrates basic usage: creating a client and listing wiki pages.
func Example() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wikis, err := c.Wiki.List(context.Background(), "MYPROJECT")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", wikis[0].ID, wikis[0].Name)
	// Output:
	// ID: 112, Name: test1
}

// Example_withOptions demonstrates using option methods to filter results.
func Example_withOptions() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	wikis, err := c.Wiki.List(
		context.Background(),
		"MYPROJECT",
		c.Wiki.Option.WithKeyword("test"),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Printf("ID: %d, Name: %s\n", wikis[0].ID, wikis[0].Name)
	// Output:
	// ID: 112, Name: test1
}

// ExampleNewClient demonstrates basic client initialization.
func ExampleNewClient() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(c != nil)
	// Output:
	// true
}

// ExampleNewClient_withDoer demonstrates injecting a custom HTTP client.
// This is useful for configuring timeouts, transport settings, or providing
// a mock implementation during testing. Set CheckRedirect so that the client
// does not follow redirects, which would forward the API key to another host.
func ExampleNewClient_withDoer() {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"apiKey",
		backlog.WithDoer(httpClient),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(c != nil)
	// Output:
	// true
}

// ExampleNewOAuthClient demonstrates client initialization with an OAuth 2.0
// access token. Obtaining and refreshing the token is up to the caller; create
// a new Client when the token expires.
func ExampleNewOAuthClient() {
	c, err := backlog.NewOAuthClient(
		"https://example.backlog.com",
		"accessToken",
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(c != nil)
	// Output:
	// true
}
