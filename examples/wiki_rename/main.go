package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nattokin/go-backlog"
)

func main() {
	baseURL := os.Getenv("BACKLOG_BASE_URL")
	if baseURL == "" {
		log.Fatalln("You need Backlog base url.")
	}
	token := os.Getenv("BACKLOG_TOKEN")
	if token == "" {
		log.Fatalln("You need Backlog access token.")
	}

	stdin := bufio.NewScanner(os.Stdin)

	projectKey := scanner(stdin, "project name:")

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	wikis, err := c.Wiki.All(context.Background(), projectKey)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%d Wikis in the project.\n", len(wikis))

	old := scanner(stdin, "search string:")

	// TargetWiki holds the ID and name of a wiki page to be renamed.
	type TargetWiki struct {
		wikiID int
		name   string
	}
	targets := []*TargetWiki{}
	for _, w := range wikis {
		if strings.HasPrefix(w.Name, old) {
			fmt.Printf("wikiID=%d, name=%s\n", w.ID, w.Name)
			targets = append(targets, &TargetWiki{w.ID, w.Name})
		}
	}
	fmt.Printf("%d name of Wikis is matched.\n", len(targets))

	replacement := scanner(stdin, "replacement:")

	for _, t := range targets {
		t.name = strings.Replace(t.name, old, replacement, 1)
		fmt.Printf("wikiID=%d, name=%s\n", t.wikiID, t.name)
	}
	fmt.Printf("%d number of Wikis will be updated.\n", len(targets))

	agree := scanner(stdin, "Execution[y/n]:")

	if strings.EqualFold(agree, "y") || strings.EqualFold(agree, "yes") {
		for _, t := range targets {
			c.Wiki.Update(context.Background(), t.wikiID, c.Wiki.Option.WithName(t.name))
		}
	} else {
		fmt.Println("exit.")
	}
}

// scanner returns the input string.
func scanner(stdin *bufio.Scanner, msg string) string {
	s := ""
	for s == "" {
		fmt.Print(msg)
		stdin.Scan()
		s = stdin.Text()
	}
	return s
}
