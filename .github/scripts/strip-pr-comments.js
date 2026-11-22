// @ts-check
'use strict';

// Removes HTML comments, such as the guidance in the PR template, from the
// PR description, because the description becomes the squash merge commit
// message.

const https = require('https');
const { stripHtmlComments } = require('./pr-description');

// ──────────────────────────────────────────────────────────────
//  Environment variables (injected by the workflow)
// ──────────────────────────────────────────────────────────────

const {
  GITHUB_TOKEN,
  REPO_OWNER,
  REPO_NAME,
  PR_NUMBER,
  PR_BODY,
} = process.env;

// ──────────────────────────────────────────────────────────────
//  Entry point
// ──────────────────────────────────────────────────────────────

async function main() {
  const prNumber = parseInt(PR_NUMBER ?? '', 10);
  const eventBody = PR_BODY ?? '';

  // Work on the current description, and leave it to the next "edited" event
  // if it changed after this event, so that a newer edit is not overwritten.
  const pr = await githubRequest('GET', `/repos/${REPO_OWNER}/${REPO_NAME}/pulls/${prNumber}`);
  const body = pr.body ?? '';
  if (body !== eventBody) {
    console.log('The description changed after this event. Skipping.');
    return;
  }

  const stripped = stripHtmlComments(body);
  if (stripped === body) {
    console.log('No HTML comments to remove.');
    return;
  }

  await githubRequest('PATCH', `/repos/${REPO_OWNER}/${REPO_NAME}/pulls/${prNumber}`, { body: stripped });
  console.log(`Removed HTML comments from the description of PR #${prNumber}.`);
}

main().catch(err => {
  console.error(err);
  process.exit(1);
});

// ──────────────────────────────────────────────────────────────
//  GitHub API helpers
// ──────────────────────────────────────────────────────────────

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
        'User-Agent': 'strip-pr-comments-script',
        'X-GitHub-Api-Version': '2022-11-28',
        ...(body
          ? { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(body) }
          : {}),
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
