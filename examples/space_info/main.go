package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	// Space info
	space, err := c.Space.Info(ctx)
	if err != nil {
		log.Fatalf("failed to fetch space info: %v", err)
	}
	fmt.Printf("Space     : %s (%s)\n", space.Name, space.SpaceKey)
	fmt.Printf("Timezone  : %s\n", space.Timezone)
	fmt.Printf("Lang      : %s\n", space.Lang)

	// Disk usage
	fmt.Println()
	disk, err := c.Space.DiskUsage(ctx)
	if err != nil {
		log.Fatalf("failed to fetch disk usage: %v", err)
	}
	fmt.Printf("Disk Usage (capacity: %s)\n", formatBytes(disk.Capacity))
	fmt.Printf("  Issue      : %s\n", formatBytes(disk.Issue))
	fmt.Printf("  Wiki       : %s\n", formatBytes(disk.Wiki))
	fmt.Printf("  File       : %s\n", formatBytes(disk.File))
	fmt.Printf("  Git        : %s\n", formatBytes(disk.Git))
	fmt.Printf("  Subversion : %s\n", formatBytes(disk.Subversion))

	// Recent activity
	fmt.Println()
	activities, err := c.Space.Activity.List(ctx, c.Space.Activity.Option.WithCount(10))
	if err != nil {
		log.Fatalf("failed to fetch activities: %v", err)
	}
	fmt.Printf("Recent Activity (%d entries)\n", len(activities))
	for _, a := range activities {
		user := ""
		if a.CreatedUser != nil {
			user = a.CreatedUser.Name
		}
		summary := ""
		if a.Content != nil {
			summary = a.Content.Summary
		}
		fmt.Printf("  [type:%d] %s — %s (%s)\n", a.Type, user, summary, a.Created.Format("2006-01-02 15:04:05"))
	}
}

// formatBytes formats a byte count into a human-readable string.
func formatBytes(b int) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/KB)
	default:
		return fmt.Sprintf("%d B", b)
	}
}
