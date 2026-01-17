# Palette's Journal

## 2024-05-22 - TTY-Aware JSON Output
**Learning:** CLI tools returning JSON are unreadable by default in terminals but useful for piping. Checking `os.Stdout.Stat()` for `os.ModeCharDevice` allows providing the best of both worlds: pretty-printed JSON for humans (TTY) and raw JSON for machines (pipes).
**Action:** Always wrap JSON output in CLI tools with a TTY check to format it for readability.
