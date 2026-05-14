# space_info

Display space information, disk usage, and recent activity.

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_TOKEN`: Your Backlog API access token

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_TOKEN=your_token

go run .
```

No arguments required.

## Output

```
Space     : My Space (myspace)
Timezone  : Asia/Tokyo
Lang      : ja

Disk Usage (capacity: 10.00 GB)
  Issue      : 123.45 MB
  Wiki       : 12.34 MB
  File       : 456.78 MB
  Git        : 0 B
  Subversion : 0 B

Recent Activity (10 entries)
  [type:1] John Doe — MYPROJECT-123: Fix login bug (2026-05-14 10:00:00)
  ...
```
