// @ts-check
'use strict';

// Helpers for reading PR descriptions.
//
// The footer block is the lines after the last blank line of the description,
// one footer per line, for example:
//
//   BREAKING CHANGE: What breaks and how to migrate.
//   Closes #12
//   Refs #34

const FENCE_PATTERN = /^ {0,3}(`{3,}|~{3,})/;

// Inline code span that does not cross a blank line.
const INLINE_CODE = /(?<!`)(`+)(?!`)(?:(?!\n[ \t]*\n)[\s\S])*?(?<!`)\1(?!`)/g;

// A comment that takes whole lines, and a comment anywhere else.
// The body of a comment must not contain "-->", so that one match never
// spans two comments.
const LINE_COMMENT = /^[ \t]*<!--(?:(?!-->)[\s\S])*-->[ \t]*(?:\n|(?![\s\S]))/gm;
const INLINE_COMMENT = /<!--(?:(?!-->)[\s\S])*-->/g;

// Lines that can make up the footer block of a description written by the
// workflows themselves. Older keywords are accepted so that existing
// descriptions keep their footer.
const REFERENCE_FOOTER_LINE = /^(?:BREAKING[ -]CHANGE: .*|(?:close[sd]?|fix(?:e[sd])?|resolve[sd]?|refs?|part of):? #\d+)$/i;

/**
 * Mark which lines belong to fenced code blocks, fences included.
 * An unclosed fence runs to the end of the text.
 * @param {string[]} lines
 * @returns {boolean[]}
 */
function codeLineMask(lines) {
  const mask = [];
  /** @type {string | null} */
  let fence = null;

  for (const line of lines) {
    if (fence === null) {
      const m = FENCE_PATTERN.exec(line);
      if (m) fence = m[1];
      mask.push(fence !== null);
      continue;
    }

    mask.push(true);
    const closing = new RegExp(`^ {0,3}${fence[0] === '`' ? '`' : '~'}{${fence.length},}[ \\t]*$`);
    if (closing.test(line)) fence = null;
  }

  return mask;
}

/**
 * Replace inline code spans with placeholders.
 * @param {string} text
 * @returns {{ masked: string, restore: (s: string) => string }}
 */
function maskInlineCode(text) {
  /** @type {string[]} */
  const spans = [];
  const masked = text.replace(INLINE_CODE, m => {
    spans.push(m);
    return `\u0000${spans.length - 1}\u0000`;
  });
  return {
    masked,
    restore: s => s.replace(/\u0000(\d+)\u0000/g, (_, i) => spans[Number(i)]),
  };
}

/**
 * Remove HTML comments outside code. A comment that takes whole lines is
 * removed with its lines. An unterminated comment is left as it is.
 * Returns the description unchanged when it has no comment to remove.
 * @param {string} description
 * @returns {string}
 */
function stripHtmlComments(description) {
  const lines = description.replace(/\r\n/g, '\n').split('\n');
  const mask = codeLineMask(lines);

  // Group consecutive lines into code and text chunks.
  /** @type {Array<{ code: boolean, lines: string[] }>} */
  const chunks = [];
  lines.forEach((line, i) => {
    const last = chunks[chunks.length - 1];
    if (last && last.code === mask[i]) {
      last.lines.push(line);
    } else {
      chunks.push({ code: mask[i], lines: [line] });
    }
  });

  let changed = false;
  const result = chunks.map(chunk => {
    const text = chunk.lines.join('\n');
    if (chunk.code) return text;

    const { masked, restore } = maskInlineCode(text);
    const stripped = masked.replace(LINE_COMMENT, '').replace(INLINE_COMMENT, '');
    if (stripped === masked) return text;

    changed = true;
    return restore(stripped);
  }).join('\n');

  return changed ? result.replace(/^\n+/, '') : description;
}

/**
 * Split a description into the lines before the footer block and the
 * footer block itself. The footer block is empty when the last paragraph is
 * part of a fenced code block.
 * @param {string} description
 * @returns {{ lines: string[], footerStart: number, codeMask: boolean[] }}
 */
function splitFooter(description) {
  const lines = description.replace(/\r\n/g, '\n').trimEnd().split('\n');
  const codeMask = codeLineMask(lines);

  let footerStart = 0;
  for (let i = lines.length - 1; i >= 0; i--) {
    if (lines[i].trim() === '' && !codeMask[i]) {
      footerStart = i + 1;
      break;
    }
  }

  if (codeMask.slice(footerStart).some(Boolean)) footerStart = lines.length;

  return { lines, footerStart, codeMask };
}

/**
 * Find the footer block of a description that the workflows edit, such as a
 * parent PR. Only a last paragraph made of reference and breaking change
 * footers counts, so that a list at the end is not taken for a footer.
 * Returns lines.length when there is no such footer block.
 * @param {string[]} lines
 * @returns {number}
 */
function findReferenceFooter(lines) {
  let end = lines.length;
  while (end > 0 && lines[end - 1].trim() === '') end--;

  let start = end;
  while (start > 0 && lines[start - 1].trim() !== '') start--;

  if (start === end) return lines.length;
  for (let i = start; i < end; i++) {
    if (!REFERENCE_FOOTER_LINE.test(lines[i].trim())) return lines.length;
  }
  return start;
}

module.exports = {
  codeLineMask,
  maskInlineCode,
  stripHtmlComments,
  splitFooter,
  findReferenceFooter,
};
