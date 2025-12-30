## 2024-10-24 - CLI Output Formatting
**Learning:** Users running CLI commands like `action list` expect readable output (JSON), but raw `io.Copy` dumps unreadable strings.
**Action:** When printing API responses in CLI tools, always attempt to decode as JSON and pretty-print. Fall back to raw output only if decoding fails. Using `json.NewDecoder().UseNumber()` preserves data fidelity.
