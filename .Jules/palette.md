## 2024-05-23 - JSON Output Readability
**Learning:** CLI tools often dump raw JSON responses which are hard to read for humans.
**Action:** Always pretty-print JSON output in CLI tools using `json.MarshalIndent` or similar, unless the user explicitly requests raw output (e.g. for piping).
