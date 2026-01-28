## 2025-02-24 - [CLI JSON Readability]
**Learning:** Raw JSON output in CLI tools is difficult to scan, forcing users to rely on external tools like `jq`. Pretty-printing JSON by default significantly improves the "out-of-the-box" experience and readability without requiring extra dependencies.
**Action:** When a CLI command outputs structured data (like JSON), default to a pretty-printed format (indented) unless raw output is explicitly required or pipe usage is detected (though in this case, we improved the default output for all cases as a safe first step).
