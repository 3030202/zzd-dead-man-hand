## 2024-05-22 - [CLI Output Formatting]
**Learning:** CLI tools returning JSON should always pretty-print the output by default. Users often use these tools manually and raw JSON is hard to parse visually.
**Action:** When implementing CLI commands that return structured data, use `json.MarshalIndent` or similar to make the output human-readable, with a fallback to raw output if parsing fails.
