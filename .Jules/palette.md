## 2024-03-20 - [CLI Output Formatting]
**Learning:** Users running CLI commands like `list` often expect human-readable output by default when in a terminal, but script-friendly (raw) output when piping.
**Action:** Detect terminal presence using `isatty` and default to pretty-printed JSON for interactive sessions, falling back to raw output for pipes.
