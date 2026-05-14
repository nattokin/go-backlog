# issue_export

Export issues in a project to TSV format on stdout.

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_TOKEN`: Your Backlog API access token

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_TOKEN=your_token

go run . <PROJECT_KEY> [--status <id,...>]
```

## Options

| Flag | Description |
|------|-------------|
| `--status` | Comma-separated status IDs to filter (e.g. `1,2,3`) |

## Output

TSV with the following columns:

| Column | Description |
|--------|-------------|
| ID | Issue ID |
| Key | Issue key (e.g. `PROJECT-123`) |
| Summary | Issue summary |
| Status | Status name |
| Assignee | Assignee name |
| Priority | Priority name |
| Comments | Number of comments |
| Attachments | Number of attachments |
| Created | Created datetime |
| Updated | Last updated datetime |

## Example

```
go run . MYPROJECT > issues.tsv
go run . MYPROJECT --status 1,2 > open_issues.tsv
```
