// Example: Get the number of Wiki pages in a project.
// This example demonstrates how to count Wiki pages
// in a specified project using the Backlog API client.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/nattokin/go-backlog"
)

func main() {
	baseURL := os.Getenv("BACKLOG_BASE_URL")
	token := os.Getenv("BACKLOG_TOKEN")
	if baseURL == "" || token == "" {
		log.Fatalln("BACKLOG_BASE_URL and BACKLOG_TOKEN must be set")
	}

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	// ID or key of the project.
	projectKey := "MYPROJECT"

	count, err := c.Wiki.Count(context.Background(), projectKey)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("Count: %d\n", count)
}
