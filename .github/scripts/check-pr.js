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
const REFERENCE_START = /^(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?|refs?|references?|part\s+of|related(?:\s+to)?|see(?:\s+also)?|issues?)\b:?\s*(?:#|[\w.-]+\/[\w.-]+#|https?:\/\/)/i;
const REFERENCE_FOOTER = /^(?:Closes|Refs) #\d+$/;

// GitHub's closing keywords, which link and close issues wherever they appear.
const CLOSING_REFERENCE = /\b(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?):?\s+(?:[\w.-]+\/[\w.-]+#\d+|#\d+|https:\/\/github\.com\/[\w.-]+\/[\w.-]+\/issues\/\d+)/i;

const BREAKING_FOOTER = /^BREAKING[ -]CHANGE: \S/;

// Sections of the PR template, in the order they must appear.
const SECTIONS = ['Summary', 'Changes'];
const HEADING = /^ {0,3}(#{1,6})[ \t]+(.*?)(?:[ \t]+#+)?[ \t]*$/;
const LIST_ITEM = /^\s*(?:[-*+]|\d+[.)])\s+\S/;

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
 * Checks that the description keeps the sections of the PR template:
 * it starts with "## Summary", followed later by "## Changes". Each section
 * runs until the next heading of level 2 or above. Summary must have some
 * text, and Changes at least one list item. Other sections may follow.
 * @param {string} description
 * @returns {Result}
 */
function checkSections(description) {
  const result = { name: 'Template sections', errors: /** @type {string[]} */ ([]) };
  const { lines, codeMask } = splitFooter(stripHtmlComments(description));

  /** @type {Array<{ index: number, level: number, title: string }>} */
  const headings = [];
  lines.forEach((line, i) => {
    if (codeMask[i]) return;
    const m = HEADING.exec(line);
    if (m) headings.push({ index: i, level: m[1].length, title: m[2] });
  });

  /** @type {Map<string, { index: number, end: number }>} */
  const found = new Map();
  for (const name of SECTIONS) {
    const exact = headings.filter(h => h.level === 2 && h.title === name);
    if (exact.length === 0) {
      const similar = headings.find(h => h.title.toLowerCase() === name.toLowerCase());
      result.errors.push(
        `The description must have a "## ${name}" section, as in the PR template.` +
        (similar ? ` Found "${'#'.repeat(similar.level)} ${similar.title}" instead.` : ''),
      );
      continue;
    }
    if (exact.length > 1) {
      result.errors.push(`The description must have only one "## ${name}" section.`);
    }
    const { index } = exact[0];
    const next = headings.find(h => h.index > index && h.level <= 2);
    found.set(name, { index, end: next ? next.index : lines.length });
  }

  const summary = found.get('Summary');
  const changes = found.get('Changes');

  if (summary) {
    if (lines.slice(0, summary.index).some(line => line.trim() !== '')) {
      result.errors.push('The description must start with the "## Summary" section.');
    }
    if (!lines.slice(summary.index + 1, summary.end).some(line => line.trim() !== '')) {
      result.errors.push('The "## Summary" section is empty. Describe what the PR does and why.');
    }
  }

  if (changes) {
    if (summary && changes.index < summary.index) {
      result.errors.push('The "## Changes" section must come after the "## Summary" section.');
    }
    const hasItem = lines
      .slice(changes.index + 1, changes.end)
      .some((line, i) => !codeMask[changes.index + 1 + i] && LIST_ITEM.test(line));
    if (!hasItem) {
      result.errors.push('The "## Changes" section must list the changes as list items, e.g. "- Add ...".');
    }
  }

  return result;
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
    results.push(
      checkSections(description),
      ...descriptionResults,
      checkBreakingChange(title, hasBreakingFooter),
    );
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

module.exports = { checkTitle, checkSections, checkDescription, checkBreakingChange };
