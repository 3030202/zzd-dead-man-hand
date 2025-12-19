## 2025-12-19 - [CLI JSON Pretty Printing]
**Learning:** CLI tools returning JSON data should format it for readability by default, but maintain raw output if parsing fails. This improves user experience significantly for data-heavy commands.
**Action:** Always use json.MarshalIndent for JSON output in CLIs unless raw output is explicitly requested or parsing fails.
