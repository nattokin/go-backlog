package main

import (
	"fmt"
	"log"

	"github.com/nattokin/go-backlog"
)

func main() {
	// The base URL of Backlog API.
	baseURL := "BACKLOG_BASE_URL"
	// The token for request to Backlog API.
	token := "BACKLOG_TOKEN"

	// Create Backlog API client.
	c, err := backlog.NewClient(baseURL, token, nil)
	if err != nil {
		log.Fatalln(err)
	}

	attachmentIDs := []int{0123, 0124, 0125}

	// Attach uploded files to Wiki.
	r, err := c.Wiki.Attachment.Attach(12345, attachmentIDs)
	if err != nil {
		log.Fatalln(err)
	}

	// Response
	fmt.Printf("%#v\n", r)
}
