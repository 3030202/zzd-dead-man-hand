## 2024-05-23 - [CLI Output Formatting]
**Learning:** CLI tools returning raw JSON are hostile to users. Simple pretty-printing (indentation) significantly improves readability without changing the underlying data structure.
**Action:** Always process and format structured data (JSON/XML) for human readability in CLI outputs, while keeping a fallback for raw data if parsing fails.
