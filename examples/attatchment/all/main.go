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
	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	// You get all attachments of the Wiki.
	r, err := c.Wiki.Attachment.List(12345)
	if err != nil {
		log.Fatalln(err)
	}

	// Response
	for _, a := range r {
		fmt.Printf("%#v\n", a)
	}
}
