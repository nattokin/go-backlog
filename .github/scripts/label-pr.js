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
} = process.env;

// ──────────────────────────────────────────────────────────────
//  Label mapping
//  Ordered from most specific (with "!") to least specific.
// ──────────────────────────────────────────────────────────────

/** @type {Array<{ pattern: RegExp, labels: string[] }>} */
const LABEL_RULES = [
  { pattern: /^feat!/,     labels: ['enhancement', 'breaking change'] },
  { pattern: /^fix!/,      labels: ['bug',         'breaking change'] },
  { pattern: /^refactor!/, labels: ['refactor',    'breaking change'] },
  { pattern: /^feat:/,     labels: ['enhancement'] },
  { pattern: /^fix:/,      labels: ['bug'] },
  { pattern: /^refactor:/, labels: ['refactor'] },
  { pattern: /^docs:/,     labels: ['documentation'] },
  { pattern: /^ci:/,       labels: ['ci'] },
  { pattern: /^chore:/,    labels: ['chore'] },
];

// ──────────────────────────────────────────────────────────────
//  Entry point
// ──────────────────────────────────────────────────────────────

async function main() {
  const prNumber = parseInt(PR_NUMBER, 10);
  const title    = PR_TITLE ?? '';

  const rule = LABEL_RULES.find(r => r.pattern.test(title));
  if (!rule) {
    console.log(`No matching label rule for title: "${title}"`);
    return;
  }

  console.log(`Matched rule for "${title}": labels = ${JSON.stringify(rule.labels)}`);
  await addLabels(prNumber, rule.labels);
  console.log(`Added labels ${JSON.stringify(rule.labels)} to PR #${prNumber}.`);
}

main().catch(err => {
  console.error(err);
  process.exit(1);
});

// ──────────────────────────────────────────────────────────────
//  GitHub API helpers
// ──────────────────────────────────────────────────────────────

/**
 * Add labels to a PR (issues endpoint accepts PR numbers).
 * @param {number} prNumber
 * @param {string[]} labels
 */
function addLabels(prNumber, labels) {
  return githubRequest(
    'POST',
    `/repos/${REPO_OWNER}/${REPO_NAME}/issues/${prNumber}/labels`,
    { labels },
  );
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
        'User-Agent': 'label-pr-script',
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
