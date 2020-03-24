package main

import (
	"fmt"
	"log"

	backlog "github.com/nattokin/go-backlog"
)

func main() {
	// The base URL of Backlog API.
	baseURL := "BACKLOG_BASE_URL"
	// The tokun for request to Backlog API.
	token := "BACKLOG_TOKEN"

	// Create Backlog API client.
	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	// Get a Wiki by ID of the Wiki.
	// You get struct where represented the Wiki.
	r, err := c.Wiki.One(12345)
	if err != nil {
		log.Fatalln(err)
	}

	// Response
	fmt.Printf("%#v\n", r)
}
