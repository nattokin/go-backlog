# wiki_rename

Rename wiki pages in a project by replacing a common prefix of their names.

The program is interactive: it reads the project key, the prefix to search for and the replacement from stdin, shows the pages that will be renamed, and asks for confirmation before changing anything.

> [!WARNING]
> This example **updates wiki pages** in your space. Check the list of new names carefully before you confirm.

## Prerequisites

- `BACKLOG_BASE_URL`: Your Backlog space URL (e.g. `https://example.backlog.com`)
- `BACKLOG_API_KEY`: Your Backlog API key

## Usage

```
export BACKLOG_BASE_URL=https://example.backlog.com
export BACKLOG_API_KEY=your_api_key

go run .
```

No arguments required. The program prompts for the following values:

| Prompt | Description |
|--------|-------------|
| `project name:` | Project key (e.g. `MYPROJECT`) |
| `search string:` | Prefix to look for. Only pages whose name starts with this string are renamed |
| `replacement:` | String that replaces the prefix |
| `Execution[y/n]:` | Enter `y` or `yes` (case-insensitive) to apply the changes. Any other input exits without changes |

Errors from the individual update requests are not reported, so check the result in Backlog afterwards.

## Output

```
project name:MYPROJECT
5 Wikis in the project.
search string:Draft/
wikiID=123, name=Draft/Design
wikiID=456, name=Draft/Roadmap
2 name of Wikis is matched.
replacement:Docs/
wikiID=123, name=Docs/Design
wikiID=456, name=Docs/Roadmap
2 number of Wikis will be updated.
Execution[y/n]:y
```
