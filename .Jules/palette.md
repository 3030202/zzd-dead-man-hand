## 2026-02-05 - [CLI Output Formatting]
**Learning:** CLI tools outputting structured data (like JSON) should default to a pretty-printed format (indentation) for readability, with a fallback to raw output if parsing fails.
**Action:** Always wrap CLI output writers in a formatter that attempts to pretty-print JSON responses.
