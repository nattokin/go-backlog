package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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

	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Usage: go run . <ISSUE_ID> <FILE_PATH>")
	}
	issueID, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf("ISSUE_ID must be an integer, got %q", args[0])
	}
	issueIDStr := strconv.Itoa(issueID)
	filePath := args[1]

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	// Step 1: Upload the file to the space.
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	fileName := filepath.Base(filePath)
	fmt.Printf("Uploading %s to space...\n", fileName)
	attachment, err := c.Space.Attachment.Upload(ctx, fileName, f)
	if err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}
	fmt.Printf("Uploaded: ID=%d, name=%s\n", attachment.ID, attachment.Name)

	// Step 2: Attach the uploaded file to the issue via Issue.Update.
	fmt.Printf("Attaching to issue %d...\n", issueID)
	_, err = c.Issue.Update(ctx, issueIDStr, c.Issue.Option.WithAttachmentIDs([]int{attachment.ID}))
	if err != nil {
		log.Fatalf("failed to attach file to issue: %v", err)
	}
	fmt.Println("Done.")

	// Confirm: list attachments on the issue.
	attachments, err := c.Issue.Attachment.List(ctx, issueIDStr)
	if err != nil {
		log.Printf("warning: failed to list attachments: %v", err)
		return
	}
	fmt.Printf("Attachments on issue %d (%d total):\n", issueID, len(attachments))
	for _, a := range attachments {
		fmt.Printf("  ID=%d, name=%s\n", a.ID, a.Name)
	}
}
