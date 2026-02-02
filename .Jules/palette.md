## 2024-05-21 - CLI Output Testability
**Learning:** Hardcoded `os.Stdout` in CLI commands prevents verifying output formatting and user feedback in tests.
**Action:** Always inject an `io.Writer` (defaulting to `os.Stdout`) in CLI commands to allow mocking with `bytes.Buffer` for robust UX testing.
