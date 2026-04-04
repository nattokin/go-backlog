// milestone-sync.js
// Reads the PR description and sets the milestone on the PR
// based on the milestone of any linked issue.
//
// Recognized link patterns (case-insensitive), matched on the LAST LINE only:
//   Closes #N  /  Fixes #N  /  Resolves #N  /  Part of #N
//
// Uses Node.js built-in https module (no external dependencies).

'use strict';

const https = require('https');

const owner = process.env.REPO_OWNER;
const repo  = process.env.REPO_NAME;
const prNum = parseInt(process.env.PR_NUMBER, 10);
const body  = process.env.PR_BODY ?? '';
const token = process.env.GITHUB_TOKEN;

// Extract issue numbers from the last line of the description only.
// Matches: Closes #N, Fixes #N, Resolves #N, Part of #N (case-insensitive)
function extractIssueNumbers(text) {
  const lastLine = text.trimEnd().split('\n').pop() ?? '';
  const pattern = /(?:closes|fixes|resolves|part\s+of)\s+#(\d+)/gi;
  const numbers = [];
  let match;
  while ((match = pattern.exec(lastLine)) !== null) {
    const n = parseInt(match[1], 10);
    if (!numbers.includes(n)) numbers.push(n);
  }
  return numbers;
}

function githubRequest(method, path, data) {
  return new Promise((resolve, reject) => {
    const reqBody = data ? JSON.stringify(data) : undefined;
    const options = {
      hostname: 'api.github.com',
      path,
      method,
      headers: {
        'Authorization': `Bearer ${token}`,
        'Accept': 'application/vnd.github+json',
        'User-Agent': 'milestone-sync-script',
        'X-GitHub-Api-Version': '2022-11-28',
        ...(reqBody
          ? { 'Content-Type': 'application/json', 'Content-Length': Buffer.byteLength(reqBody) }
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
    if (reqBody) req.write(reqBody);
    req.end();
  });
}

async function getMilestone(issueNumber) {
  try {
    const data = await githubRequest('GET', `/repos/${owner}/${repo}/issues/${issueNumber}`);
    return data.milestone ?? null;
  } catch (err) {
    console.log(`  Failed to fetch issue #${issueNumber}: ${err.message}`);
    return null;
  }
}

async function run() {
  const issueNumbers = extractIssueNumbers(body);

  if (issueNumbers.length === 0) {
    console.log('No linked issues found in the last line of PR description. Skipping.');
    return;
  }

  console.log(`Linked issues (from last line): ${issueNumbers.map(n => '#' + n).join(', ')}`);

  // Use the milestone of the first linked issue that has one.
  let targetMilestone = null;
  for (const n of issueNumbers) {
    console.log(`Checking issue #${n}...`);
    const milestone = await getMilestone(n);
    if (milestone) {
      console.log(`  Found milestone: "${milestone.title}" (id=${milestone.number})`);
      targetMilestone = milestone;
      break;
    } else {
      console.log('  No milestone set.');
    }
  }

  if (!targetMilestone) {
    console.log('No milestone found on any linked issue. Skipping.');
    return;
  }

  // Fetch current PR milestone to avoid unnecessary API calls.
  const pr = await githubRequest('GET', `/repos/${owner}/${repo}/pulls/${prNum}`);
  if (pr.milestone && pr.milestone.number === targetMilestone.number) {
    console.log(`PR already has milestone "${targetMilestone.title}". Nothing to do.`);
    return;
  }

  await githubRequest('PATCH', `/repos/${owner}/${repo}/issues/${prNum}`, {
    milestone: targetMilestone.number,
  });

  console.log(`Set milestone "${targetMilestone.title}" on PR #${prNum}.`);
}

run().catch(err => {
  console.error(err);
  process.exit(1);
});
