package main

import (
	"fmt"
	"log"

	"github.com/nattokin/go-backlog"
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

	fpath := "/path/to/test.txt"
	fname := "name.txt"
	r, err := c.Wiki.Attachment.Upload(fpath, fname)
	if err != nil {
		fmt.Println(err)
	}

	// Response
	fmt.Printf("%#v\n", r)
}
