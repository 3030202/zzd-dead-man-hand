## 2024-10-24 - [CLI Output Readability]
**Learning:** CLI tools returning raw JSON dumps are difficult for users to scan, significantly reducing usability for list commands.
**Action:** Always attempt to pretty-print JSON responses in CLI clients, using `json.MarshalIndent` for readability, while maintaining a robust fallback to raw output for non-JSON data.
