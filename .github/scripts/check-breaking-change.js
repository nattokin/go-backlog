const title = process.env.PR_TITLE || '';
const body = process.env.PR_BODY || '';
const labels = JSON.parse(process.env.PR_LABELS || '[]');

const hasExclamation = /^[a-z]+!/.test(title);
const hasBreakingLabel = labels.includes('breaking change');
const hasBreakingSection = /BREAKING CHANGE:/m.test(body);

const signals = [hasExclamation, hasBreakingLabel, hasBreakingSection];
const triggeredCount = signals.filter(Boolean).length;

if (triggeredCount === 0) {
  // Not a breaking change — all good
  console.log('Not a breaking change. No issues found.');
  process.exit(0);
}

if (triggeredCount === 3) {
  // All three present — consistent
  console.log('Breaking change is consistent: title has "!", label "breaking change" is set, and "BREAKING CHANGE:" section is present.');
  process.exit(0);
}

// Partial — fail with detailed message
const lines = [
  'Breaking change consistency check failed.',
  '',
  `  PR title has "!" suffix : ${hasExclamation  ? '✅' : '❌'}`,
  `  Label "breaking change" : ${hasBreakingLabel ? '✅' : '❌'}`,
  `  "BREAKING CHANGE:" in description : ${hasBreakingSection ? '✅' : '❌'}`,
  '',
  'All three must be present together, or none at all.',
];
console.error(lines.join('\n'));
process.exit(1);
