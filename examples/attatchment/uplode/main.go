package main

import (
	"fmt"
	"log"
	"os"

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

	f, err := os.Open("/path/to/test.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r, err := c.Space.Attachment.Upload(f.Name(), f)
	if err != nil {
		fmt.Println(err)
	}

	// Response
	fmt.Printf("%#v\n", r)
}
