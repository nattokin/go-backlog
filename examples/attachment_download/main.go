package main

import (
	"context"
	"flag"
	"fmt"
	"io"
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

	target := flag.String("target", "issue", "target type: issue or wiki")
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Usage: go run . [--target issue|wiki] <ISSUE_KEY_OR_WIKI_ID> <OUTPUT_DIR>")
	}
	key := args[0]
	outDir := args[1]

	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	switch *target {
	case "issue":
		downloadIssueAttachments(ctx, c, key, outDir)
	case "wiki":
		wikiID, err := strconv.Atoi(key)
		if err != nil {
			log.Fatalf("wiki ID must be an integer, got %q", key)
		}
		downloadWikiAttachments(ctx, c, wikiID, outDir)
	default:
		log.Fatalf("unknown target %q: must be issue or wiki", *target)
	}
}

func downloadIssueAttachments(ctx context.Context, c *backlog.Client, issueKey, outDir string) {
	attachments, err := c.Issue.Attachment.List(ctx, issueKey)
	if err != nil {
		log.Fatalf("failed to list attachments: %v", err)
	}
	fmt.Printf("%d attachment(s) found on issue %s\n", len(attachments), issueKey)

	for _, a := range attachments {
		fmt.Printf("  downloading %s...\n", a.Name)
		fd, err := c.Issue.Attachment.Download(ctx, issueKey, a.ID)
		if err != nil {
			log.Printf("warning: failed to download %s: %v", a.Name, err)
			continue
		}
		if err := saveFile(fd, outDir); err != nil {
			log.Printf("warning: failed to save %s: %v", a.Name, err)
		}
	}
}

func downloadWikiAttachments(ctx context.Context, c *backlog.Client, wikiID int, outDir string) {
	attachments, err := c.Wiki.Attachment.List(ctx, wikiID)
	if err != nil {
		log.Fatalf("failed to list attachments: %v", err)
	}
	fmt.Printf("%d attachment(s) found on wiki %d\n", len(attachments), wikiID)

	for _, a := range attachments {
		fmt.Printf("  downloading %s...\n", a.Name)
		fd, err := c.Wiki.Attachment.Download(ctx, wikiID, a.ID)
		if err != nil {
			log.Printf("warning: failed to download %s: %v", a.Name, err)
			continue
		}
		if err := saveFile(fd, outDir); err != nil {
			log.Printf("warning: failed to save %s: %v", a.Name, err)
		}
	}
}

// saveFile writes FileData to a file in outDir and closes the body.
func saveFile(fd *backlog.FileData, outDir string) error {
	defer fd.Body.Close()

	path := filepath.Join(outDir, fd.Filename)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, fd.Body); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}
