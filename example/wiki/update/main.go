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

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	r, err := c.Wiki.Update(1234, "changed name", "changed content", false)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", r)
}
