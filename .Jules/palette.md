## 2026-01-21 - CLI Output Testability
**Learning:** CLI commands printed directly to os.Stdout, making output content untestable.
**Action:** Introduced a package-level 'outputWriter' variable (defaulting to os.Stdout) to allow tests to capture and assert on CLI output formatting.
