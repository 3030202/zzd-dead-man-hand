## 2024-05-22 - Improved CLI Output Readability
**Learning:** CLI tools that output raw JSON are difficult for humans to read. Pretty-printing JSON responses by default significantly improves the user experience without sacrificing machine readability (as long as it remains valid JSON).
**Action:** Always default to pretty-printed JSON for human-facing CLI commands, with a fallback to raw output if parsing fails.
