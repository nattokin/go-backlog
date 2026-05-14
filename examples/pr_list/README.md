# pr_list

List pull requests in a repository with comment counts.

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_TOKEN`: Your Backlog API access token

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_TOKEN=your_token

go run . <PROJECT_KEY> <REPO_NAME> [--status open|closed|merged|draft]
```

## Options

| Flag | Description |
|------|-------------|
| `--status` | Filter by status: `open`, `closed`, `merged`, or `draft` |

## Output

```
3 pull request(s) found in MYPROJECT/myrepo
#42 [Open] Fix login bug
  branch   : feature/fix-login -> main
  assignee : John Doe
  comments : 5
  created  : 2026-01-15 10:00:00 +0000 UTC
...
```
