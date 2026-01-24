## 2026-01-24 - CLI JSON Output Readability
**Learning:** Users find raw JSON output in CLI tools difficult to read. Pretty-printing (indenting) JSON by default significantly improves usability without requiring external tools like `jq`.
**Action:** When implementing CLI commands that output structured data (JSON), always attempt to pretty-print it first. Fallback to raw output if parsing fails to ensure robustness.
