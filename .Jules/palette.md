## 2024-03-24 - [CLI JSON Pretty Printing]
**Learning:** Terminal users often struggle to read raw JSON output. Detect TTY and pretty-print JSON for immediate readability, but keep raw output for pipes/scripts.
**Action:** When implementing CLI commands that output structured data, always check `isTerminal` to decide between pretty-printed and raw format.
