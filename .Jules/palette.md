## 2024-05-22 - [CLI JSON Formatting]
**Learning:** CLI tools outputting structured data (like JSON) are much more usable when they default to pretty-printed output (indentation) for humans, while falling back gracefully to raw output if parsing fails.
**Action:** Always check if CLI output is intended for human consumption and format accordingly, possibly checking for TTY.
