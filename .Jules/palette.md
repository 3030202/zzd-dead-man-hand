# Palette's Journal

## 2025-02-14 - [CLI Output Readability]
**Learning:** For backend-only tools, "UX" often means "readable CLI output". Users expect structured data (like JSON) to be formatted by default, rather than raw blobs.
**Action:** When working on CLI tools, always check if commands outputting JSON/XML are pretty-printed. Add a fallback to raw output if parsing fails.
