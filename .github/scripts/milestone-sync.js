// milestone-sync.js
// Keeps the milestone of a PR and the issues it closes in sync.
//
// The target issues are those referenced as "Closes #N" in the footer block
// of the PR description. "Refs #N" is ignored.
//
// - An issue has a milestone: the PR is set to it.
// - No issue has a milestone and the PR has one: the issues are set to it.
// - The PR and the issues have different milestones, or the issues have
//   different milestones from each other: the job fails and nothing is changed.
//
// Uses Node.js built-in https module (no external dependencies).

'use strict';

const https = require('https');
const { splitFooter, stripHtmlComments } = require('./pr-description');

const owner = process.env.REPO_OWNER;
const repo  = process.env.REPO_NAME;
const prNum = parseInt(process.env.PR_NUMBER ?? '', 10);
const token = process.env.GITHUB_TOKEN;

const CLOSES_FOOTER = /^Closes #(\d+)$/;

/**
 * Extract the issue numbers referenced as "Closes #N" in the footer block.
 * @param {string} body
 * @returns {number[]}
 */
function extractClosedIssues(body) {
  const { lines, footerStart, codeMask } = splitFooter(stripHtmlComments(body));
  const numbers = [];
  for (let i = footerStart; i < lines.length; i++) {
    if (codeMask[i]) continue;
    const m = CLOSES_FOOTER.exec(lines[i].trim());
    if (!m) continue;
    const n = parseInt(m[1], 10);
    if (!numbers.includes(n)) numbers.push(n);
  }
  return numbers;
}

/**
 * @typedef {{ number: number, title: string }} Milestone
 * @typedef {{ number: number, milestone: Milestone | null }} Issue
 * @typedef {{ setPR: Milestone | null, setIssues: number[], milestone: Milestone | null, error: string | null }} Plan
 */

/**
 * Decide what to change. Nothing is changed when there is an error.
 * @param {Milestone | null} prMilestone
 * @param {Issue[]} issues
 * @returns {Plan}
 */
function planSync(prMilestone, issues) {
  const plan = { setPR: null, setIssues: [], milestone: null, error: null };

  const withMilestone = issues.filter(i => i.milestone);
  const distinct = [...new Map(withMilestone.map(i => [i.milestone.number, i.milestone])).values()];

  if (distinct.length > 1) {
    const list = withMilestone.map(i => `#${i.number} "${i.milestone.title}"`).join(', ');
    plan.error = `The closed issues have different milestones: ${list}. Set the same milestone on them.`;
    return plan;
  }

  const target = distinct[0] ?? prMilestone;
  if (!target) return plan;

  if (prMilestone && prMilestone.number !== target.number) {
    const list = withMilestone.map(i => `#${i.number}`).join(', ');
    plan.error =
      `The PR milestone "${prMilestone.title}" differs from the milestone "${target.title}" of ${list}. ` +
      'Set the same milestone on the PR and the issues.';
    return plan;
  }

  plan.milestone = target;
  if (!prMilestone) plan.setPR = target;
  plan.setIssues = issues.filter(i => !i.milestone).map(i => i.number);
  return plan;
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

/**
 * @param {number} number
 * @param {Milestone} milestone
 */
function setMilestone(number, milestone) {
  return githubRequest('PATCH', `/repos/${owner}/${repo}/issues/${number}`, {
    milestone: milestone.number,
  });
}

async function run() {
  // Read the current PR, as other workflows may have changed it since the event.
  const pr = await githubRequest('GET', `/repos/${owner}/${repo}/pulls/${prNum}`);
  const issueNumbers = extractClosedIssues(pr.body ?? '');

  if (issueNumbers.length === 0) {
    console.log('No "Closes #N" footer in the PR description. Skipping.');
    return;
  }
  console.log(`Closed issues: ${issueNumbers.map(n => '#' + n).join(', ')}`);

  /** @type {Issue[]} */
  const issues = [];
  for (const n of issueNumbers) {
    const data = await githubRequest('GET', `/repos/${owner}/${repo}/issues/${n}`);
    issues.push({ number: n, milestone: data.milestone ?? null });
    console.log(`  #${n}: ${data.milestone ? `"${data.milestone.title}"` : 'no milestone'}`);
  }

  const plan = planSync(pr.milestone ?? null, issues);
  if (plan.error) {
    console.log(`::error title=Milestone mismatch::${plan.error}`);
    process.exit(1);
  }
  if (!plan.milestone) {
    console.log('Neither the PR nor the closed issues have a milestone. Nothing to do.');
    return;
  }

  if (plan.setPR) {
    await setMilestone(prNum, plan.setPR);
    console.log(`Set milestone "${plan.setPR.title}" on PR #${prNum}.`);
  }
  for (const n of plan.setIssues) {
    await setMilestone(n, plan.milestone);
    console.log(`Set milestone "${plan.milestone.title}" on issue #${n}.`);
  }
  if (!plan.setPR && plan.setIssues.length === 0) {
    console.log(`The PR and the closed issues already have milestone "${plan.milestone.title}".`);
  }
}

if (require.main === module) {
  run().catch(err => {
    console.error(err);
    process.exit(1);
  });
}

module.exports = { extractClosedIssues, planSync };
