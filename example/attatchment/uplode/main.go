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

	fPath := "/path/to/test.txt"
	fName := "name.txt"
	r, err := c.Wiki.Attachment.Uploade(fPath, fName)
	if err != nil {
		fmt.Println(err)
	}

	// Response
	fmt.Printf("%#v\n", r)
}
