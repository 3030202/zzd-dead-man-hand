## 2024-05-23 - [CLI Output Readability]
**Learning:** Users struggle to read raw JSON output in CLI tools, especially for lists of items. Relying on external tools like `jq` adds friction.
**Action:** Always pretty-print structured data (JSON/YAML) in CLI output by default, but provide a fallback or flag for raw output if needed for piping. In Go, `json.MarshalIndent` is the standard way to achieve this.
