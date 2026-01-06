## 2024-04-18 - [CLI Output Readability]
**Learning:** CLI tools often dump raw API responses (minified JSON) which are hard for humans to scan. Users shouldn't have to pipe output to `jq` just to read it.
**Action:** When a CLI command outputs structured data (like JSON), always default to a pretty-printed format (indentation) for stdout. Keep raw output as an option or fallback if needed, but optimize the default for human readability.
