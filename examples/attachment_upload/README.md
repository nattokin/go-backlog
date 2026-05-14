# attachment_upload

Upload a local file to the space and attach it to a specified issue.

This example demonstrates the standard Backlog API two-step flow:
1. Upload the file to the space via `Space.Attachment.Upload` to obtain an attachment ID
2. Attach it to the issue via `Issue.Update` with `WithAttachmentIDs`

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_TOKEN`: Your Backlog API access token

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_TOKEN=your_token

go run . <ISSUE_ID> <FILE_PATH>
```

## Arguments

| Argument | Description |
|----------|-------------|
| `ISSUE_ID` | Numeric ID of the issue |
| `FILE_PATH` | Path to the local file to upload |

## Output

```
Uploading report.pdf to space...
Uploaded: ID=123, name=report.pdf
Attaching to issue 456...
Done.
Attachments on issue 456 (2 total):
  ID=122, name=previous.png
  ID=123, name=report.pdf
```
