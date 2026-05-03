package backlog_test

import (
	"context"
	"fmt"
	"net/http"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	doerExample            = newMockDoer(fixture.Wiki.ListJSON)
	doerExampleWithOptions = newMockDoer(fixture.Wiki.ListJSON)
)

// Example demonstrates basic usage: creating a client and listing wiki pages.
func Example() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerExample),
	)

	wikis, _ := c.Wiki.All(context.Background(), "MYPROJECT")
	fmt.Printf("ID: %d, Name: %s\n", wikis[0].ID, wikis[0].Name)
	// Output:
	// ID: 112, Name: test1
}

// Example_withOptions demonstrates using option methods to filter results.
func Example_withOptions() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerExampleWithOptions),
	)

	wikis, _ := c.Wiki.All(
		context.Background(),
		"MYPROJECT",
		c.Wiki.Option.WithKeyword("test"),
	)
	fmt.Printf("ID: %d, Name: %s\n", wikis[0].ID, wikis[0].Name)
	// Output:
	// ID: 112, Name: test1
}

// ExampleNewClient demonstrates basic client initialization.
func ExampleNewClient() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"token",
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
// a mock implementation during testing.
func ExampleNewClient_withDoer() {
	c, err := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(http.DefaultClient),
	)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(c != nil)
	// Output:
	// true
}
