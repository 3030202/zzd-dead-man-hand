## 2025-02-18 - [CLI JSON Formatting]
**Learning:** Users interacting with CLIs benefit greatly from pretty-printed JSON output by default, rather than raw strings.
**Action:** When implementing CLI commands that output structured data, always attempt to pretty-print it, falling back to raw output only on failure.
