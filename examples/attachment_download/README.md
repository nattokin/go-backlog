# attachment_download

Download all attachments from an issue or wiki page to a local directory.

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_TOKEN`: Your Backlog API access token

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_TOKEN=your_token

# Download attachments from an issue
go run . --target issue <ISSUE_KEY> <OUTPUT_DIR>

# Download attachments from a wiki page
go run . --target wiki <WIKI_ID> <OUTPUT_DIR>
```

## Options

| Flag | Default | Description |
|------|---------|-------------|
| `--target` | `issue` | Target type: `issue` or `wiki` |

## Output

Files are saved to `<OUTPUT_DIR>/` using the original filename returned by the API.

```
3 attachment(s) found on issue MYPROJECT-42
  downloading report.pdf...
  downloading screenshot.png...
  downloading data.csv...
```

## Note

This example demonstrates the `io.Reader`-based `FileData` response and the caller's responsibility to close `FileData.Body` after use.
