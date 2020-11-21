package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nattokin/go-backlog"
)

func main() {
	// Get env and set it in variable.
	baseURL := os.Getenv("BACKLOG_BASE_URL")
	if baseURL == "" {
		log.Fatalln("You need Backlog base url.")
	}
	token := os.Getenv("BACKLOG_TOKEN")
	if token == "" {
		log.Fatalln("You need Backlog access token.")
	}

	stdin := bufio.NewScanner(os.Stdin)

	// Scan project name.
	projectKey := scanner(stdin, "project name:")

	// Get all Wikis in the project.
	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	r, err := c.Wiki.All(backlog.ProjectKey(projectKey))
	if err != nil {
		log.Fatalln(err)
	}

	// Output the count of Wikis in the project.
	fmt.Printf("%d Wikis in the project.\n", len(r))

	// Scan serch string.
	old := scanner(stdin, "serch string:")

	// TargetWiki is reflects one single target data to update Wiki.
	// It has ID of Wiki and name of wiki to update.
	type TargetWiki struct {
		wikiID int
		name   string
	}
	targets := []*TargetWiki{}
	for _, w := range r {
		if strings.Index(w.Name, old) == 0 {
			fmt.Printf("wikiID=%d, name=%s\n", w.ID, w.Name)
			targets = append(targets, &TargetWiki{w.ID, w.Name})
		}
	}
	fmt.Printf("%d name of Wikis is matched.\n", len(targets))

	// Scan replacement.
	new := scanner(stdin, "replacement:")

	// Replace name of targets.
	for _, t := range targets {
		t.name = strings.Replace(t.name, old, new, 1)
		fmt.Printf("wikiID=%d, name=%s\n", t.wikiID, t.name)
	}
	fmt.Printf("%d number of Wikis will be updated.\n", len(targets))

	// Get agreement of execution.
	agree := scanner(stdin, "Execution[y/n]:")

	// When agreement is obtained, update the name of the Wiki
	if agree == "y" || agree == "Y" || agree == "yes" || agree == "Yes" {
		for _, t := range targets {
			c.Wiki.Update(t.wikiID, c.Wiki.Option.WithFormName(t.name))
		}
	} else {
		fmt.Println("exit.")
	}
}

// Scanner returns the input string.
func scanner(stdin *bufio.Scanner, msg string) string {
	s := ""
	for s == "" {
		fmt.Print(msg)
		stdin.Scan()
		s = stdin.Text()
	}
	return s
}
