## 2024-05-23 - [CLI Output Readability]
**Learning:** Raw JSON output in CLI tools is difficult to read and verify for users, especially for list operations.
**Action:** Always default to pretty-printed JSON output for CLI commands that return structured data, while maintaining a fallback to raw output for non-JSON responses.
