## 2024-03-24 - [CLI Output Testability]
**Learning:** Testing CLI output in Go is much easier when `os.Stdout` is replaced by a package-level variable (e.g., `var outputWriter io.Writer = os.Stdout`). This allows mocking with `bytes.Buffer` in tests.
**Action:** When working on CLI tools, always check if output destination is configurable. If not, refactor it to be injectable or variable-based before adding tests.
