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
// ──────────────────────────────────────────────────────────────

/** @type {Record<string, string>} */
const TYPE_LABELS = {
  feat:     'enhancement',
  fix:      'bug',
  refactor: 'refactor',
  docs:     'documentation',
  test:     'test',
  chore:    'chore',
  ci:       'ci',
};

const BREAKING_LABEL = 'breaking change';

// Labels this workflow manages. Other labels are left untouched.
const MANAGED_LABELS = [...Object.values(TYPE_LABELS), BREAKING_LABEL];

/**
 * Return the managed labels a PR title implies, or null when the title has
 * no known type prefix.
 * @param {string} title
 * @returns {string[] | null}
 */
function expectedLabels(title) {
  const m = /^([a-z]+)(!)?:/.exec(title);
  if (!m || !Object.hasOwn(TYPE_LABELS, m[1])) return null;

  const labels = [TYPE_LABELS[m[1]]];
  if (m[2]) labels.push(BREAKING_LABEL);
  return labels;
}

/**
 * Compute the managed labels to add and to remove.
 * @param {string[]} expected
 * @param {string[]} current
 * @returns {{ add: string[], remove: string[] }}
 */
function diffLabels(expected, current) {
  return {
    add:    expected.filter(l => !current.includes(l)),
    remove: current.filter(l => MANAGED_LABELS.includes(l) && !expected.includes(l)),
  };
}

// ──────────────────────────────────────────────────────────────
//  Entry point
// ──────────────────────────────────────────────────────────────

async function main() {
  const prNumber = parseInt(PR_NUMBER ?? '', 10);
  const title    = PR_TITLE ?? '';

  const expected = expectedLabels(title);
  if (!expected) {
    console.log(`No matching label rule for title: "${title}"`);
    return;
  }

  const current = (await listLabels(prNumber)).map(l => l.name);
  const { add, remove } = diffLabels(expected, current);

  if (add.length === 0 && remove.length === 0) {
    console.log(`Labels already match "${title}": ${JSON.stringify(expected)}`);
    return;
  }

  if (add.length > 0) {
    await addLabels(prNumber, add);
    console.log(`Added labels ${JSON.stringify(add)} to PR #${prNumber}.`);
  }
  for (const label of remove) {
    await removeLabel(prNumber, label);
    console.log(`Removed label "${label}" from PR #${prNumber}.`);
  }
}

if (require.main === module) {
  main().catch(err => {
    console.error(err);
    process.exit(1);
  });
}

module.exports = { expectedLabels, diffLabels };

// ──────────────────────────────────────────────────────────────
//  GitHub API helpers
// ──────────────────────────────────────────────────────────────

/**
 * List the labels of a PR (issues endpoint accepts PR numbers).
 * @param {number} prNumber
 * @returns {Promise<Array<{ name: string }>>}
 */
function listLabels(prNumber) {
  return githubRequest(
    'GET',
    `/repos/${REPO_OWNER}/${REPO_NAME}/issues/${prNumber}/labels?per_page=100`,
  );
}

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
 * Remove a label from a PR. A label that is already gone is ignored.
 * @param {number} prNumber
 * @param {string} label
 */
async function removeLabel(prNumber, label) {
  try {
    await githubRequest(
      'DELETE',
      `/repos/${REPO_OWNER}/${REPO_NAME}/issues/${prNumber}/labels/${encodeURIComponent(label)}`,
    );
  } catch (err) {
    if (/** @type {{ statusCode?: number }} */ (err).statusCode !== 404) throw err;
  }
}

/**
 * Minimal GitHub API request using Node.js built-in https.
 * @param {string} method
 * @param {string} path
 * @param {object} [data]
 * @returns {Promise<any>}
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
          const err = new Error(`GitHub API ${method} ${path} failed: ${res.statusCode} ${raw}`);
          Object.assign(err, { statusCode: res.statusCode });
          reject(err);
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
