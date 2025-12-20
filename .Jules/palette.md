## 2024-10-27 - CLI Output Readability
**Learning:** Raw JSON output in CLI tools is difficult for humans to scan, reducing usability.
**Action:** Always pretty-print JSON responses in CLI commands using `json.MarshalIndent` or similar, while maintaining a fallback for non-JSON content.
