# wiki_backup

Back up all wiki pages in a project to local Markdown files, including attachments and version history.

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_TOKEN`: Your Backlog API access token

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_TOKEN=your_token

go run . <PROJECT_KEY> <OUTPUT_DIR>
```

## Output

```
OUTPUT_DIR/
  Home.md
  Getting_Started.md
  attachments/
    123/
      diagram.png
    456/
      spec.pdf
```

Each wiki page is saved as `<OUTPUT_DIR>/<wiki-name>.md`. Attachments are saved under `<OUTPUT_DIR>/attachments/<wiki-id>/`. Characters unsafe for filenames (`/ \ : * ? " < > |`) in wiki names are replaced with `_`.

```
3 wiki page(s) found in MYPROJECT
  [123] Home — 5 version(s), 1 attachment(s)
  [456] Getting Started — 2 version(s), 0 attachment(s)
  [789] API Reference — 8 version(s), 3 attachment(s)
```
