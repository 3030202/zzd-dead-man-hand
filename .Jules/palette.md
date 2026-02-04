## 2024-10-23 - CLI JSON Readability
**Learning:** Raw JSON output in CLI tools is difficult for humans to parse, hindering quick debugging and verification. Pretty-printing JSON by default significantly improves the developer experience (DX).
**Action:** Default to indented JSON output for CLI commands that return structured data, while ensuring the output remains valid JSON for piping.
