## 2025-05-23 - [CLI Output Formatting]
**Learning:** CLI tools often ignore usability by outputting raw JSON. Implementing terminal detection (`isTerminal`) to auto-format (pretty-print) JSON output significantly improves the human experience without breaking machine-readability (pipes/scripts) which continue to receive raw output.
**Action:** Always check `isatty` (or equivalent) when outputting structured data in CLI tools; default to pretty-print for TTYs and raw for non-TTYs.
