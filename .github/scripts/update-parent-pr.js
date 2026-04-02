// @ts-check
'use strict';

const https = require('https');

// ──────────────────────────────────────────────────────────────
//  Environment variables (injected by the workflow)
// ──────────────────────────────────────────────────────────────

const {
  GITHUB_TOKEN,
  REPO_OWNER,
  REPO_NAME,
  PR_NUMBER,
  PR_TITLE,
  PR_BODY,
  PR_BASE_REF,
  PR_MERGED,
  EVENT_ACTION,
  CHANGES_TITLE,
  CHANGES_BODY,
} = process.env;

// ──────────────────────────────────────────────────────────────
//  Entry point
// ──────────────────────────────────────────────────────────────

async function main() {
  // Skip edited events that don't change title or body
  if (EVENT_ACTION === 'edited') {
    const changesTitle = CHANGES_TITLE && CHANGES_TITLE !== 'null';
    const changesBody  = CHANGES_BODY  && CHANGES_BODY  !== 'null';
    if (!changesTitle && !changesBody) {
      console.log('No title or body change detected, skipping.');
      return;
    }
  }

  const childNumber = parseInt(PR_NUMBER, 10);
  const childTitle  = PR_TITLE ?? '';
  const childBody   = PR_BODY  ?? '';
  const baseBranch  = PR_BASE_REF ?? '';

  // Find parent PR: open PR whose HEAD branch equals the child's base branch.
  // i.e. the PR that is trying to merge `baseBranch` into something else.
  // Note: do NOT encode the slash in branch names for GitHub API head param.
  const openPRs = await listOpenPRsByHead(baseBranch);
  const parentPR = openPRs.find(pr => pr.number !== childNumber);
  if (!parentPR) {
    console.log(`No parent PR found with head branch "${baseBranch}".`);
    return;
  }
  console.log(`Found parent PR #${parentPR.number}: ${parentPR.title}`);

  // Extract child PR's Changes section content
  const childChanges = extractChangesLines(childBody);
  console.log(`Child PR Changes lines: ${JSON.stringify(childChanges)}`);

  // Normalize and update parent PR description
  let parentBody = parentPR.body ?? '';
  parentBody = normalizeDescription(parentBody);
  parentBody = upsertChildEntry(parentBody, childNumber, childTitle, childChanges);

  await updatePR(parentPR.number, parentBody);
  console.log(`Updated parent PR #${parentPR.number} description.`);
}

main().catch(err => {
  console.error(err);
  process.exit(1);
});

// ──────────────────────────────────────────────────────────────
//  GitHub API helpers
// ──────────────────────────────────────────────────────────────

/**
 * Find open PRs whose HEAD branch matches the given branch name.
 * Note: GitHub API head param requires `owner:branch` format without encoding slashes.
 * @param {string} headBranch
 * @returns {Promise<Array<{number: number, title: string, body: string|null}>>}
 */
function listOpenPRsByHead(headBranch) {
  return githubRequest('GET', `/repos/${REPO_OWNER}/${REPO_NAME}/pulls?state=open&head=${REPO_OWNER}:${headBranch}&per_page=100`);
}

/**
 * @param {number} prNumber
 * @param {string} body
 */
function updatePR(prNumber, body) {
  return githubRequest('PATCH', `/repos/${REPO_OWNER}/${REPO_NAME}/pulls/${prNumber}`, { body });
}

/**
 * Minimal GitHub API request using Node.js built-in https.
 * @param {string} method
 * @param {string} path
 * @param {object} [data]
 */
function githubRequest(method, path, data) {
  return new Promise((resolve, reject) => {
    const body = data ? JSON.stringify(data) : undefined;
    const options = {
      hostname: 'api.github.com',
      path,
      method,
      headers: {
        'Authorization': `Bearer ${GITHUB_TOKEN}`,
        'Accept': 'application/vnd.github+json',
        'User-Agent': 'update-parent-pr-script',
        'X-GitHub-Api-Version': '2022-11-28',
        ...(body ? { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(body) } : {}),
      },
    };

    const req = https.request(options, res => {
      let raw = '';
      res.on('data', chunk => { raw += chunk; });
      res.on('end', () => {
        if (res.statusCode >= 400) {
          reject(new Error(`GitHub API ${method} ${path} failed: ${res.statusCode} ${raw}`));
          return;
        }
        resolve(raw ? JSON.parse(raw) : null);
      });
    });

    req.on('error', reject);
    if (body) req.write(body);
    req.end();
  });
}

// ──────────────────────────────────────────────────────────────
//  Description manipulation helpers
// ──────────────────────────────────────────────────────────────

/**
 * Extract the top-level list items under `## Changes` from a PR body.
 * Returns an array of strings (each item text without leading "- ").
 * @param {string} body
 * @returns {string[]}
 */
function extractChangesLines(body) {
  const lines = body.split('\n');
  const start = lines.findIndex(l => /^##\s+Changes\s*$/.test(l));
  if (start === -1) return [];

  const result = [];
  for (let i = start + 1; i < lines.length; i++) {
    const line = lines[i];
    if (/^##/.test(line)) break;
    if (/^- /.test(line)) result.push(line.slice(2).trim());
  }
  return result;
}

/**
 * Ensure the body has `## Summary` at the top and `## Changes`
 * before any closing keywords.
 * @param {string} body
 * @returns {string}
 */
function normalizeDescription(body) {
  let lines = body.split('\n');

  // Ensure ## Summary exists at the top
  if (!lines.some(l => /^##\s+Summary\s*$/.test(l))) {
    lines = ['## Summary', '', ...lines];
  }

  // Ensure ## Changes exists and is placed before closing keywords
  if (!lines.some(l => /^##\s+Changes\s*$/.test(l))) {
    const closingPattern = /^(Closes|Fixes|Resolves)\s+#\d+/i;
    const closingIdx = lines.findIndex(l => closingPattern.test(l));

    if (closingIdx !== -1) {
      // Insert before the closing keyword line (with surrounding blank lines)
      const insert = ['## Changes', ''];
      const prev = lines[closingIdx - 1];
      if (prev !== undefined && prev.trim() !== '') insert.unshift('');
      lines.splice(closingIdx, 0, ...insert);
    } else {
      // Insert after ## Summary section
      const summaryIdx = lines.findIndex(l => /^##\s+Summary\s*$/.test(l));
      let insertIdx = summaryIdx + 1;
      while (insertIdx < lines.length && !/^##/.test(lines[insertIdx])) {
        insertIdx++;
      }
      lines.splice(insertIdx, 0, '', '## Changes', '');
    }
  }

  return lines.join('\n');
}

/**
 * Upsert the child PR entry (parent row + indented lines) in ## Changes.
 * - If a row containing `(#childNumber)` exists: replace it and its indented lines.
 * - Otherwise: append to the end of the ## Changes section.
 * @param {string} body
 * @param {number} childNumber
 * @param {string} childTitle
 * @param {string[]} childChanges
 * @returns {string}
 */
function upsertChildEntry(body, childNumber, childTitle, childChanges) {
  const lines = body.split('\n');

  // Build the new entry lines
  const newEntry = [`- ${childTitle} (#${childNumber})`];
  for (const change of childChanges) {
    newEntry.push(`  - ${change}`);
  }

  const changesIdx = lines.findIndex(l => /^##\s+Changes\s*$/.test(l));
  if (changesIdx === -1) return body; // Should not happen after normalize

  // Find the boundary of ## Changes section (next heading or end of file)
  let sectionEnd = lines.length;
  for (let i = changesIdx + 1; i < lines.length; i++) {
    if (/^##/.test(lines[i])) { sectionEnd = i; break; }
  }

  // Search for existing parent row within ## Changes section
  const markerPattern = new RegExp(`\\(#${childNumber}\\)`);
  let existingIdx = -1;
  for (let i = changesIdx + 1; i < sectionEnd; i++) {
    if (/^- /.test(lines[i]) && markerPattern.test(lines[i])) {
      existingIdx = i;
      break;
    }
  }

  if (existingIdx !== -1) {
    // Find the extent of the existing entry (parent row + indented rows)
    let endIdx = existingIdx + 1;
    while (endIdx < sectionEnd && /^  /.test(lines[endIdx])) endIdx++;
    lines.splice(existingIdx, endIdx - existingIdx, ...newEntry);
  } else {
    // Append after last non-empty line within the section
    let insertIdx = sectionEnd;
    for (let i = sectionEnd - 1; i > changesIdx; i--) {
      if (lines[i].trim() !== '') { insertIdx = i + 1; break; }
    }
    lines.splice(insertIdx, 0, ...newEntry);
  }

  return lines.join('\n');
}
