package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatalln("Usage: go run . <PROJECT_KEY> <OUTPUT_DIR>")
	}
	projectKey := args[0]
	outDir := args[1]

	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatalf("failed to create output directory: %v", err)
	}

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	wikis, err := c.Wiki.All(ctx, projectKey)
	if err != nil {
		log.Fatalf("failed to fetch wiki list: %v", err)
	}
	fmt.Printf("%d wiki page(s) found in %s\n", len(wikis), projectKey)

	for _, w := range wikis {
		// Fetch full content.
		page, err := c.Wiki.One(ctx, w.ID)
		if err != nil {
			log.Printf("warning: failed to fetch wiki %d (%s): %v", w.ID, w.Name, err)
			continue
		}

		// Save page content as Markdown.
		mdPath := filepath.Join(outDir, safeFilename(page.Name)+".md")
		if err := os.WriteFile(mdPath, []byte(page.Content), 0644); err != nil {
			log.Printf("warning: failed to write %s: %v", mdPath, err)
		}

		// Fetch and save version history.
		histories, err := c.Wiki.History.List(ctx, w.ID)
		if err != nil {
			log.Printf("warning: failed to fetch history for wiki %d: %v", w.ID, err)
		}

		// Download attachments.
		attachments, err := c.Wiki.Attachment.List(ctx, w.ID)
		if err != nil {
			log.Printf("warning: failed to fetch attachments for wiki %d: %v", w.ID, err)
		}
		attachDir := filepath.Join(outDir, "attachments", fmt.Sprintf("%d", w.ID))
		if len(attachments) > 0 {
			if err := os.MkdirAll(attachDir, 0755); err != nil {
				log.Printf("warning: failed to create attachment dir: %v", err)
			}
			for _, a := range attachments {
				fd, err := c.Wiki.Attachment.Download(ctx, w.ID, a.ID)
				if err != nil {
					log.Printf("warning: failed to download attachment %s: %v", a.Name, err)
					continue
				}
				if err := saveFile(fd, attachDir); err != nil {
					log.Printf("warning: failed to save attachment %s: %v", a.Name, err)
				}
			}
		}

		fmt.Printf("  [%d] %s — %d version(s), %d attachment(s)\n",
			w.ID, page.Name, len(histories), len(attachments))
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

// safeFilename replaces characters that are unsafe for filenames.
func safeFilename(name string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	return replacer.Replace(name)
}
