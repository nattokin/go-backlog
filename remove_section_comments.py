#!/usr/bin/env python3
"""
Remove redundant "// --- ... ---" section-header comments from test files.

These comments (e.g. "// --- validation errors ---", "// --- fail-fast: nil option ---")
duplicate information already carried by the subtest names, which are what actually
appears in `go test -v` output, CI failure logs, and `-run` filters. Blank-line
grouping between blocks is preserved; only the comment line itself is removed.

Run from the repository root:
    python3 remove_section_comments.py
"""
import re
import glob

COMMENT_RE = re.compile(r'^\s*// --- .+ ---\s*$')

def main():
    total_files = 0
    total_removed = 0
    for path in sorted(glob.glob("internal/**/*_test.go", recursive=True)):
        with open(path, encoding="utf-8") as f:
            lines = f.readlines()
        new_lines = [l for l in lines if not COMMENT_RE.match(l)]
        removed = len(lines) - len(new_lines)
        if removed:
            with open(path, "w", encoding="utf-8") as f:
                f.writelines(new_lines)
            total_files += 1
            total_removed += removed
            print(f"{path}: removed {removed}")
    print(f"\nTotal: {total_removed} comment lines removed across {total_files} files")

if __name__ == "__main__":
    main()