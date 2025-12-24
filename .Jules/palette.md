## 2025-12-24 - CLI Output Formatting
**Learning:** The CLI tool `dmh-cli` was outputting raw JSON responses from the server directly to stdout. This makes it hard for users to read, especially when dealing with lists of objects.
**Action:** Implemented pretty-printing for JSON output in the `listActions` command. If the response is not valid JSON, it falls back to raw output. Future CLI commands returning data should follow this pattern or use a table printer for even better readability.
