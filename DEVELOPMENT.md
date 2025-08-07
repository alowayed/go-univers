# Development

## For AI Agents

Do the following:

1. Read the README.md, CONTRIBUTING.md, and CLAUDE.md files in this directory.
2. Read https://go.dev/doc/effective_go and https://google.github.io/styleguide/go/best-practices.html.
3. Read https://ossf.github.io/osv-schema/, focusing on the ecosystems and any information about versioning.
4. Skim this directory to understand the code. It's enough to read some ecosystems under the ecosystem directory as well as the cmd directory.
5. Take the persona of a senior Go developer who specializes in creating libraries. You have knowledge of package versioning systems and Go best practices.

Once the steps above are complete, continue reading this file.

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

## Future Work

See [GitHub Issues](https://github.com/alowayed/go-univers/issues) for planned features and improvements. Look for issues labeled "good first issue" for contribution opportunities.

## Recent Major Changes

### Added Ecosystems (2025)
- **Alpine**: Alpine Linux package versioning with suffix and build components
- **Cargo**: Rust/Cargo SemVer 2.0 with caret/tilde constraints
- **Composer**: PHP/Composer versioning with stability flags and branch names  
- **NuGet**: .NET/NuGet SemVer 2.0 with revision components and bracket notation
- **RubyGems**: Ruby Gem versioning with pessimistic constraints

### Enhanced CI/CD Pipeline
- Multi-OS testing (Linux, macOS, Windows) for cross-platform compatibility
- Comprehensive linting with golangci-lint (~60+ linters)
- Automated code formatting verification (gofmt, go mod tidy)
- Security scanning and Go best practices enforcement

### Architecture Improvements
- All ecosystems follow consistent patterns (table-driven tests, type safety)
- Performance optimizations (pre-parsed version objects in constraints)
- Comprehensive edge case testing and error handling
- Modern Go idioms (slices.Contains, proper error wrapping)