## 2024-05-23 - Pretty-printing CLI JSON Output
**Learning:** Users often read CLI output directly. Raw JSON is hard to scan. Pretty-printing JSON responses when the status is OK significantly improves readability without affecting machine parsability if they redirect to a file (though typically tools like `jq` are used, default human-readable output is better).
**Action:** When a CLI command outputs structured data intended for human consumption, default to a formatted view. Use `json.MarshalIndent` or similar. Ensure fallback to raw output if parsing fails.
