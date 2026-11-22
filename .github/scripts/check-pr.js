// @ts-check
'use strict';

// Checks the PR title and description, which become the squash merge commit
// message. All checks run and every result is reported at once.
//
// The description is checked with HTML comments removed, because
// strip-pr-comments.js removes them and its edit does not trigger a new run.

const fs = require('fs');
const { maskInlineCode, splitFooter, stripHtmlComments } = require('./pr-description');

// ──────────────────────────────────────────────────────────────
//  Environment variables (injected by the workflow)
// ──────────────────────────────────────────────────────────────

const {
  PR_TITLE,
  PR_BODY,
  PR_AUTHOR,
  GITHUB_STEP_SUMMARY,
} = process.env;

// ──────────────────────────────────────────────────────────────
//  Rules
// ──────────────────────────────────────────────────────────────

const TYPES = ['feat', 'fix', 'refactor', 'docs', 'test', 'chore', 'ci'];
const TITLE_PATTERN = new RegExp(`^(?:${TYPES.join('|')})!?: \\S`);
const TITLE_REFERENCE = /#\d+/;

// A footer line that starts with a reference keyword must be one of these.
const REFERENCE_START = /^(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?|refs?|part\s+of)\b:?\s*(?:#|[\w.-]+\/[\w.-]+#|https?:\/\/)/i;
const REFERENCE_FOOTER = /^(?:Closes|Refs) #\d+$/;

// GitHub's closing keywords, which link and close issues wherever they appear.
const CLOSING_REFERENCE = /\b(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?):?\s+(?:[\w.-]+\/[\w.-]+#\d+|#\d+|https:\/\/github\.com\/[\w.-]+\/[\w.-]+\/issues\/\d+)/i;

const BREAKING_FOOTER = /^BREAKING[ -]CHANGE: \S/;

// ──────────────────────────────────────────────────────────────
//  Checks
// ──────────────────────────────────────────────────────────────

/**
 * @typedef {{ name: string, errors: string[] }} Result
 */

/**
 * @param {string} title
 * @returns {Result[]}
 */
function checkTitle(title) {
  const format = { name: 'Title format', errors: /** @type {string[]} */ ([]) };
  if (!TITLE_PATTERN.test(title)) {
    format.errors.push(
      `The title must start with "<type>: " or "<type>!: ", where <type> is one of ${TYPES.join(', ')}. ` +
      `Got: "${title}"`,
    );
  }

  const reference = { name: 'No reference in title', errors: /** @type {string[]} */ ([]) };
  const m = TITLE_REFERENCE.exec(title);
  if (m) {
    reference.errors.push(
      `The title must not contain "${m[0]}". The PR number is appended on merge; ` +
      'put issue references in the footer block of the description.',
    );
  }

  return [format, reference];
}

/**
 * @param {string} description
 * @returns {{ results: Result[], hasBreakingFooter: boolean }}
 */
function checkDescription(description) {
  const { lines, footerStart, codeMask } = splitFooter(stripHtmlComments(description));

  const footer = { name: 'Footer references', errors: /** @type {string[]} */ ([]) };
  const outside = { name: 'No references outside the footer block', errors: /** @type {string[]} */ ([]) };
  let hasBreakingFooter = false;

  lines.forEach((line, i) => {
    if (codeMask[i]) return;
    const text = line.trim();

    if (i >= footerStart) {
      if (BREAKING_FOOTER.test(text)) hasBreakingFooter = true;
      if (REFERENCE_START.test(text) && !REFERENCE_FOOTER.test(text)) {
        footer.errors.push(
          `"${text}" is not a valid reference. Use "Closes #N" or "Refs #N", one issue per line ` +
          '(e.g. "Closes #12" and "Closes #34" on separate lines).',
        );
      }
      return;
    }

    if (BREAKING_FOOTER.test(text)) {
      outside.errors.push(
        'The "BREAKING CHANGE: ..." footer is not in the footer block. ' +
        'Put it after the last blank line of the description.',
      );
    }

    const m = CLOSING_REFERENCE.exec(maskInlineCode(line).masked);
    if (m) {
      outside.errors.push(
        `"${m[0]}" is outside the footer block, where GitHub still closes the issue. ` +
        'Move it after the last blank line of the description as "Closes #N".',
      );
    }
  });

  return { results: [footer, outside], hasBreakingFooter };
}

/**
 * The "breaking change" label is not checked: label-pr.yml keeps it in line
 * with the title "!", and its edits do not trigger a new run, so the label in
 * the event can be out of date.
 * @param {string} title
 * @param {boolean} hasBreakingFooter
 * @returns {Result}
 */
function checkBreakingChange(title, hasBreakingFooter) {
  const result = { name: 'Breaking change consistency', errors: /** @type {string[]} */ ([]) };
  const hasExclamation = /^[a-z]+!:/.test(title);

  if (hasExclamation && !hasBreakingFooter) {
    result.errors.push('The title has "!", but the description has no "BREAKING CHANGE: ..." footer.');
  }
  if (!hasExclamation && hasBreakingFooter) {
    result.errors.push('The description has a "BREAKING CHANGE: ..." footer, but the title has no "!".');
  }

  return result;
}

// ──────────────────────────────────────────────────────────────
//  Entry point
// ──────────────────────────────────────────────────────────────

function main() {
  const title = PR_TITLE ?? '';
  const description = PR_BODY ?? '';

  /** @type {Result[]} */
  const results = checkTitle(title);

  if (PR_AUTHOR === 'dependabot[bot]') {
    console.log('PR opened by Dependabot: skipping the description checks.');
  } else {
    const { results: descriptionResults, hasBreakingFooter } = checkDescription(description);
    results.push(...descriptionResults, checkBreakingChange(title, hasBreakingFooter));
  }

  report(results);
  if (results.some(r => r.errors.length > 0)) process.exit(1);
}

/**
 * @param {Result[]} results
 */
function report(results) {
  const summary = ['## PR title and description check', ''];

  for (const r of results) {
    const mark = r.errors.length === 0 ? '✅' : '❌';
    console.log(`${mark} ${r.name}`);
    summary.push(`- ${mark} ${r.name}`);

    for (const e of r.errors) {
      console.log(`::error title=${r.name}::${e}`);
      summary.push(`  - ${e}`);
    }
  }

  if (GITHUB_STEP_SUMMARY) fs.appendFileSync(GITHUB_STEP_SUMMARY, summary.join('\n') + '\n');
}

if (require.main === module) main();

module.exports = { checkTitle, checkDescription, checkBreakingChange };
