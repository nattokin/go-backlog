# project_list

List all projects in a space with their issue types and statuses.

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_API_KEY`: Your Backlog API key

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_API_KEY=your_api_key

go run .
```

No arguments required.

## Output

```
3 project(s) found.

[MYPROJECT] My Project (ID: 12345)
  Issue Types : Bug, Task, Feature
  Statuses    : Open, In Progress, Resolved, Closed
...
```
