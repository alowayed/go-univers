# Development

**go-univers** is a Go port of the Python [aboutcode-org/univers](https://github.com/aboutcode-org/univers) library for type-safe version comparison across software package ecosystems.

## Architecture and Philosophy

### Type Safety
Each ecosystem has separate types (`npm.Version`, `pypi.Version`) preventing cross-ecosystem mixing at compile time. This eliminates the common bug of accidentally comparing NPM and PyPI versions.

### Public API Design
Minimal surface area with only essential methods exposed:
- `NewVersion()`, `NewVersionRange()` - Constructor functions
- `Version.{String,Normalize,Compare}()` - Core version operations  
- `VersionRange.Contains()` - Range membership testing

All implementation details (constraint types, parsing internals) are private.

### Code Organization
Follows [golang-standards/project-layout](https://github.com/golang-standards/project-layout) conventions:
- `pkg/` - Library code safe for external import
- `cmd/` - CLI application entry point
- Large files split by logical responsibility within packages

Test files mirror source files. Examples organized in subdirectories to avoid package conflicts.

### Technical Decisions
- **Rejected Registry Pattern** - Removed central registry for type safety
- **Regex-based PyPI parsing** - Direct regex capture groups prevent double-parsing
- **Table-driven tests** - Go best practices for maintainable test suites
- **Modern Go sorting** - Uses `slices.SortFunc()` with existing `Compare()` methods

## Current State

```
/Users/yousef/dev/go-univers/
├── LICENSE
├── README.md
├── DEVELOPMENT.md
├── go.mod                       # github.com/alowayed/go-univers
├── cmd/
│   └── main.go                 # CLI application
└── pkg/
    ├── univers.go              # Universal interfaces (documentation)
    └── ecosystem/
        ├── npm/
        │   ├── npm.go          # Version type (public API)
        │   ├── range.go        # VersionRange type (public API) 
        │   ├── constraint.go   # constraint type (private)
        │   ├── version_test.go # Version tests
        │   ├── range_test.go   # Range tests
        │   └── constraint_test.go # Internal tests
        └── pypi/
            ├── pypi.go         # PyPI implementation
            └── pypi_test.go    # PyPI tests
```

**NPM**: Semantic versioning with full range syntax (^, ~, x-ranges, hyphen ranges, OR logic)
**PyPI**: PEP 440 compliant (epochs, prereleases, post-releases, dev releases, local versions)
**Tests**: All passing (`go test ./...`)

## Future Work

### Ecosystem Additions
- **Maven** - Java package versioning
- **RubyGems** - Ruby package versioning  
- **Debian** - Debian package versioning
- **Docker** - Container image tags
- **Go Modules** - Go module versioning

### Enhancements
- **Performance optimization** - Benchmark version parsing/comparison
- **CLI tool** - Command-line interface for version operations
- **JSON serialization** - Marshal/unmarshal support for versions/ranges
- **Fuzzing tests** - Property-based testing for edge cases

## Log

### Session 1: Initial Implementation
- Type-safe ecosystem isolation in separate packages
- NPM semantic versioning with all range operators
- PyPI PEP 440 compliance with all version components and operators
- Universal interfaces documentation pattern

### Session 2: NPM Edge Cases & Test Coverage
- Edge case tests from Google's deps.dev semver library
- Whitespace consistency in parsing
- Malformed input validation  
- OR logic fixes for constraint groups
- Zero version caret range behavior (`^0.0.1`, `^0.1.0`)

### Session 3: Go Idiomatic Cleanup
- File modularization (npm.go → npm.go, range.go, constraint.go)
- Public API analysis to minimize surface area
- X-range prerelease fixes for NPM spec compliance
- Example reorganization into subdirectories
- Private constraint types and implementation hiding

### Session 4: Native Go Sorting
- `slices.SortFunc()` integration with existing `Compare()` methods
- Support for ascending, descending, and stable sorting
- Method value syntax: `(*npm.Version).Compare`
- No wrapper types or convenience functions needed

### Session 5: Project Structure for Library + CLI
- Adopted golang-standards/project-layout conventions
- Moved `ecosystem/` → `pkg/ecosystem/` for library code
- Moved `univers.go` → `pkg/univers.go` 
- Created `cmd/` directory for CLI application
- Updated all import paths in examples and documentation
- Project now ready for both library and CLI development