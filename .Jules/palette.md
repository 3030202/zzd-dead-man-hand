## 2025-02-23 - Pretty Printing JSON Output in CLI Tools
**Learning:** CLI tools often output raw JSON which is hard for humans to read. Defaulting to pretty-printed JSON (with indentation) significantly improves readability and user experience. It's important to provide a fallback to raw output if parsing fails.
**Action:** When implementing CLI commands that return structured data, always attempt to pretty-print the output by default. Use a mockable output writer to facilitate testing of the formatted output.
