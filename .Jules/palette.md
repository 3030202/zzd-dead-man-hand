## 2025-02-06 - JSON Pretty Printing in CLI
**Learning:** When pretty-printing JSON in Go, using `json.Unmarshal` into `interface{}` and then `json.MarshalIndent` can cause precision loss for large numbers (converted to `float64`).
**Action:** Use `json.Indent` with a `bytes.Buffer` to format JSON directly from bytes, preserving the original data integrity.
