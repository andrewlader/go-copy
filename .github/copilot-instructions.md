# GitHub Copilot Instructions for go-copy

## Project Overview
This project is a Go command-line tool for copying files from a source location to one or more destination locations, with support for configuration and variable substitution.

## Coding Guidelines
- Use idiomatic Go style and naming conventions.
- Organize new features under `internal/copylib/` unless they are CLI-specific.
- Keep the main entrypoint in `cmd/go-copy/go-copy.go` minimal; delegate logic to internal packages.
- Write unit tests for new logic in `internal/copylib/`.
- Use Go modules for dependency management.

## Copilot Usage
- When asked to add a new feature, prefer extending `internal/copylib/` and update the CLI as needed.
- For configuration changes, edit files in `configs/` and update parsing logic in `internal/copylib/config.go`.
- When adding new commands or flags, update `cmd/go-copy/go-copy.go`.
- Always update or add tests in `internal/copylib/filecopier_test.go` or related test files.

## Best Practices
- Follow Go best practices for error handling and logging.
- Keep code modular and testable.
- Document public functions and exported types.

## Example Tasks
- Add a new CLI flag: update `cmd/go-copy/go-copy.go` and propagate to internal logic.
- Support new config option: update `configs/go-copy-config.yaml` and `internal/copylib/config.go`.
- Add a new file operation: implement in `internal/copylib/filecopier.go` and test in `filecopier_test.go`.

---
For more details, see the README.md and code comments.
