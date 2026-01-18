## 2025-02-23 - CLI JSON Output Formatting
**Learning:** Raw JSON output from CLI commands is difficult for humans to read in a terminal, but tools like `jq` expect raw JSON.
**Action:** Implement terminal detection. When outputting to a TTY, pretty-print JSON with indentation. When outputting to a pipe or file, preserve raw formatting.
