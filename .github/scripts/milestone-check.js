// @ts-check
'use strict';

const https = require('https');
const fs = require('fs');

// ──────────────────────────────────────────────────────────────
//  Environment variables (injected by the workflow)
// ──────────────────────────────────────────────────────────────

const {
  GITHUB_TOKEN,
  REPO_OWNER,
  REPO_NAME,
  TAG_NAME,
  GITHUB_OUTPUT,
} = process.env;

// ──────────────────────────────────────────────────────────────
//  Entry point
// ──────────────────────────────────────────────────────────────

async function main() {
  console.log(`Looking for milestone matching tag '${TAG_NAME}'...`);

  const milestones = await githubRequest('GET', `/repos/${REPO_OWNER}/${REPO_NAME}/milestones?state=open&per_page=100`);
  const milestone = milestones.find(m => m.title === TAG_NAME);

  if (!milestone) {
    console.error(`Error: No open milestone found for tag '${TAG_NAME}'.`);
    process.exit(1);
  }

  console.log(`Found milestone #${milestone.number}: '${milestone.title}' (open issues: ${milestone.open_issues})`);

  if (milestone.open_issues !== 0) {
    console.error(`Error: Milestone '${TAG_NAME}' has ${milestone.open_issues} open issue(s). Release aborted.`);
    process.exit(1);
  }

  fs.appendFileSync(GITHUB_OUTPUT, `milestone_number=${milestone.number}\n`);
  console.log(`Milestone check passed. Number: ${milestone.number}`);
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
        'User-Agent': 'milestone-check-script',
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
