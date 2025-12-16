## 2024-04-18 - [CLI JSON Pretty Printing]
**Learning:** CLI tools often dump raw JSON, making it hard to read. Users expect formatted output by default.
**Action:** Always check if CLI output is JSON and pretty-print it. Use `json.MarshalIndent` or `json.NewEncoder(os.Stdout).SetIndent` for better readability.
