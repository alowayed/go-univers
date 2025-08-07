# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Required Reading

**IMPORTANT**: Always read these files first to understand the project before working on any tasks:
- README.md - Project overview, supported ecosystems, usage examples, and current capabilities
- CONTRIBUTING.md - Contribution guidelines and development workflow
- DEVELOPMENT.md - Extended development documentation and architecture details

## Development Commands

### Testing
```bash
# Run all tests
go test ./...

# Run tests for specific ecosystem (see pkg/ecosystem/ for all available ecosystems)
go test ./pkg/ecosystem/npm/...
go test ./pkg/ecosystem/pypi/...

# Run CLI tests
go test ./cmd/cli/...
```

### Building
```bash
# Build the CLI binary
go build -o univers ./cmd

# Build and test in one command
go build -o univers ./cmd && go test ./...
```

### Code Quality
```bash
# Format code
go fmt ./...

# Vet code for potential issues
go vet ./...

# Run linting (configured with golangci-lint)
golangci-lint run
```

## Architecture Overview

go-univers is a type-safe library for version comparison across different software package ecosystems. The key architectural principle is **ecosystem isolation** - each ecosystem has its own types to prevent accidental cross-ecosystem version mixing at compile time.

### Core Design Patterns

1. **Type Safety**: Each ecosystem (see `pkg/ecosystem/` directory) defines its own `Version` and `VersionRange` types
2. **Generic Interfaces**: Universal interfaces in `pkg/univers/univers.go` define contracts
3. **Factory Pattern**: Each ecosystem provides `NewVersion()` and `NewVersionRange()` constructors
4. **Interface Compliance**: `pkg/ecosystem/ecosystem.go` contains compile-time interface checks

### Directory Structure

```
├── README.md                   # Project overview, usage examples, and documentation
├── CONTRIBUTING.md             # Contribution guidelines and development workflow
├── CLAUDE.md                   # Development guidelines for Claude Code
├── DEVELOPMENT.md              # Extended development documentation
├── LICENSE                     # Project license
├── go.mod                      # Go module dependencies
pkg/
├── univers/
│   └── univers.go              # Universal interfaces (Version, VersionRange, Ecosystem)
└── ecosystem/
    ├── ecosystem.go            # Interface compliance verification
    ├── npm/                    # NPM semantic versioning
    │   ├── npm.go             # Public API (Version, VersionRange types)
    │   ├── version.go         # Version implementation
    │   ├── range.go           # Range implementation
    │   └── *_test.go          # Comprehensive test suite
    ├── pypi/                   # PyPI PEP 440 versioning
    │   ├── pypi.go            # Public API
    │   ├── version.go         # PEP 440 version parsing
    │   ├── range.go           # PEP 440 range operators
    │   └── *_test.go          # Test suite
    ├── gomod/                  # Go module versioning
    │   ├── gomod.go           # Public API
    │   ├── version.go         # Semantic + pseudo-version support
    │   ├── range.go           # Go module constraints
    │   └── *_test.go          # Test suite
    └── maven/                  # Maven versioning
        ├── maven.go           # Public API
        ├── version.go         # Maven version parsing with qualifiers
        ├── range.go           # Maven range operators (brackets)
        └── *_test.go          # Test suite

cmd/
├── README.md                   # CLI usage documentation
├── main.go                     # CLI entry point
└── cli/
    ├── cli.go                 # CLI runner and argument parsing
    ├── commands.go            # Command implementations (compare, sort, contains)
    └── *_test.go              # CLI test suite
```

### Key Implementation Details

- **Alpine**: Alpine package versioning with suffix and build component support
- **Cargo**: SemVer 2.0 with Rust-specific caret/tilde constraints and wildcard matching
- **Composer**: PHP package versioning with stability flags and branch name support
- **Go**: Go module versioning with pseudo-version pattern support
- **Maven**: Maven versioning with qualifier precedence and bracket range notation
- **NPM**: Semantic versioning with range operators and OR logic
- **NuGet**: SemVer 2.0 with .NET extensions (revision component, bracket notation)
- **PyPI**: Complete PEP 440 support (epochs, prereleases, post-releases, local versions)
- **RubyGems**: Ruby gem versioning with pessimistic constraint (~>) operator

### Testing Strategy

- **Table-driven tests**: All ecosystems use Go's idiomatic table-driven test pattern
- **Edge case coverage**: Comprehensive test suites include malformed input validation
- **CLI testing**: Command-line interface has full test coverage for all operations
- **Interface compliance**: Compile-time verification ensures all types implement required interfaces

### Public API Minimalism

Each ecosystem exposes only essential functions:
- `NewVersion(string) (Version, error)` - Parse version strings
- `NewVersionRange(string) (VersionRange, error)` - Parse range strings  
- `Version.Compare(other) int` - Compare versions (-1, 0, 1)
- `VersionRange.Contains(version) bool` - Test range membership
- `Version.String() string` - Original string representation

All parsing internals, constraint types, and implementation details are private.

### CLI Architecture

The CLI follows the pattern: `univers <ecosystem> <command> [args]`

Commands:
- `compare <v1> <v2>` - Compare two versions (outputs -1, 0, 1)
- `sort <v1> <v2> ...` - Sort versions in ascending order
- `contains <range> <version>` - Check if version satisfies range (outputs true/false)

See `pkg/ecosystem/` directory for all supported ecosystems.

### Development Guidelines

1. **Type Safety First**: Never allow cross-ecosystem version operations
2. **Test Coverage**: All new functionality requires comprehensive table-driven tests
   - Test function names must follow the pattern `TestStructName_MethodName` (e.g., `TestEcosystem_NewVersion`, `TestVersion_Compare`)
   - Only test PUBLIC methods and functions - never test private/internal functions
   - Follow existing test patterns in other ecosystems for consistency
3. **API Stability**: Keep public APIs minimal and stable
4. **Go Idioms**: Follow golang-standards/project-layout and effective Go practices
5. **Error Handling**: Provide clear, actionable error messages for invalid input
6. **Documentation**: Update README.md for any new ecosystem or major feature additions
7. **Contributing**: Follow guidelines in CONTRIBUTING.md for code submissions and development workflow

### Issue Completion Process

When asked to complete a GitHub issue, ALWAYS follow this standardized process:

1. **Branch Management**:
   ```bash
   git checkout -b feat/descriptive-feature-name
   ```

2. **Issue Analysis**:
   - Fetch issue details using `gh issue view <issue-number>`
   - Create todo list to track all required tasks
   - Research requirements from issue description and any linked resources

3. **Research Phase**:
   - Study linked documentation, specifications, or reference implementations
   - Examine existing ecosystem patterns in the codebase for consistency
   - Use WebFetch tool for external documentation when needed

4. **Implementation**:
   - Follow existing architectural patterns (see directory structure above)
   - Create new ecosystem under `pkg/ecosystem/<ecosystem>/` with:
     - `<ecosystem>.go` - Public API (Ecosystem struct with Name constant)
     - `version.go` - Version implementation 
     - `range.go` - VersionRange implementation
     - `<ecosystem>_test.go` - Ecosystem.Name() test
     - `version_test.go` - Version parsing and comparison tests
     - `range_test.go` - Range parsing and Contains() tests

5. **Integration**:
   - Add ecosystem to CLI in `cmd/cli/cli.go` (import and ecosystemToRun map)
   - Add interface compliance checks in `pkg/ecosystem/ecosystem.go`
   - Update README.md supported ecosystems table and add usage examples

6. **Quality Assurance** (ALWAYS run in this order):
   ```bash
   go fmt ./...           # Format code
   go vet ./...           # Check for issues
   go test ./...          # Run all tests
   golangci-lint run      # Comprehensive linting
   ```

7. **Documentation Updates**:
   - Add ecosystem to README.md supported ecosystems table
   - Add CLI usage examples (compare, sort, contains commands)
   - Add code example in ecosystem examples section
   - Keep examples concise but demonstrative of key features

8. **Commit and PR**:
   ```bash
   git add .
   git commit -s -m "feat: add <ecosystem> ecosystem support

   - Implement <ecosystem> version parsing following <specification>
   - Add comprehensive test coverage with table-driven tests
   - Support <key features> with proper <behavior> handling
   - Add CLI integration (compare/sort/contains commands)
   - Update documentation with usage examples
   
   Fixes #<issue-number>
   
   🤖 Generated with [Claude Code](https://claude.ai/code)
   
   Co-Authored-By: Claude <noreply@anthropic.com>"
   
   git push -u origin feat/descriptive-feature-name
   gh pr create --title "feat: add <ecosystem> ecosystem support" --body "Implements <ecosystem> ecosystem support as requested in #<issue-number>"
   ```

9. **Verification**:
   - Test CLI commands manually to ensure they work correctly
   - Verify all quality checks pass
   - Ensure documentation examples are accurate

This process ensures consistency, quality, and completeness for all ecosystem additions.

### Adding New Ecosystems

1. Create new package under `pkg/ecosystem/<ecosystem>/`
2. Implement `Version` and `VersionRange` types with required methods
3. Implement `Ecosystem` interface with `NewVersion()` and `NewVersionRange()`
4. Add comprehensive table-driven tests
5. Add interface compliance check in `pkg/ecosystem/ecosystem.go`
6. Extend CLI to support new ecosystem in `cmd/cli/commands.go`
7. Update README.md with ecosystem documentation
8. Follow contribution process outlined in CONTRIBUTING.md

### Common Patterns

**Version Parsing**: Use regex with named capture groups for complex formats (see PyPI implementation)
**Range Operations**: Implement as slice of constraints with AND/OR logic
**Pseudo-versions**: Handle special version formats (Go module pseudo-versions)
**Normalization**: Maintain original string while supporting normalized comparison

### References

- @README.md
- @DEVELOPMENT.md
- @CONTRIBUTING.md
