## 2025-01-25 - CLI Output Testability
**Learning:** Hardcoding `os.Stdout` in CLI tools makes testing output difficult and noisy.
**Action:** Always use a configurable `io.Writer` (e.g., `var outputWriter io.Writer = os.Stdout`) for all CLI output to enable easy mocking in tests.
