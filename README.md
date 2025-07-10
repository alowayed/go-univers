# go-univers

A Go library for version comparison and range parsing across multiple software package ecosystems.

`go-univers` provides a type-safe, ecosystem-specific approach to handling version ranges and constraints. Unlike other libraries that mix version formats, go-univers ensures compile-time safety by keeping each ecosystem (npm, PyPI, etc.) in separate packages.

## Features

- **Type-safe version handling** - Impossible to accidentally mix npm and PyPI versions
- **Ecosystem-specific parsing** - Native support for each package manager's versioning rules
- **Comprehensive range syntax** - Full support for complex version constraints
- **PEP 440 compliant** - Complete PyPI version parsing including epochs, prereleases, dev releases
- **NPM semver compatible** - Supports all npm range operators (^, ~, x-ranges, etc.)
- **Well-tested** - Comprehensive table-driven tests following Go best practices

## Installation

```bash
go get github.com/alowayed/go-univers
```

## Quick Start

### NPM Versions

```go
package main

import (
    "fmt"
    "github.com/alowayed/go-univers/ecosystem/npm"
)

func main() {
    // Parse NPM versions
    v1, _ := npm.NewVersion("1.2.3")
    v2, _ := npm.NewVersion("1.2.4-alpha")
    
    // Compare versions
    if v1.Compare(v2) < 0 {
        fmt.Println("v1 is older than v2")
    }
    
    // Parse version ranges with NPM-specific syntax
    range1, _ := npm.NewVersionRange("^1.2.0")        // Caret range
    range2, _ := npm.NewVersionRange("~1.2.0")        // Tilde range  
    range3, _ := npm.NewVersionRange("1.x")           // X-range
    range4, _ := npm.NewVersionRange(">=1.0.0 <2.0.0") // Multiple constraints
    
    // Check if version satisfies range
    if range1.Contains(v1) {
        fmt.Println("v1 satisfies ^1.2.0")
    }
}
```

### PyPI Versions

```go
package main

import (
    "fmt"
    "github.com/alowayed/go-univers/ecosystem/pypi"
)

func main() {
    // Parse PyPI versions (PEP 440 compliant)
    v1, _ := pypi.NewVersion("1.2.3")
    v2, _ := pypi.NewVersion("2!1.2.3a1.post1.dev1+local.1") // Complex version
    
    // Compare versions
    if v1.Compare(v2) < 0 {
        fmt.Println("v1 is older than v2")
    }
    
    // Parse version ranges with PEP 440 syntax
    range1, _ := pypi.NewVersionRange("~=1.2.3")           // Compatible release
    range2, _ := pypi.NewVersionRange(">=1.0.0, <2.0.0")   // Multiple constraints
    range3, _ := pypi.NewVersionRange("==1.2.*")           // Wildcard matching
    range4, _ := pypi.NewVersionRange("!= 1.3.0")          // Exclusion
    
    // Check if version satisfies range
    if range1.Contains(v1) {
        fmt.Println("v1 satisfies ~=1.2.3")
    }
}
```

## Architecture

go-univers uses a **type-safe, ecosystem-isolated architecture** that prevents accidental cross-ecosystem version mixing:

```
go-univers/
├── univers.go                    # Universal interfaces (documentation)
└── ecosystem/
    ├── npm/                      # NPM ecosystem package
    │   ├── npm.go               # Version, VersionRange, Constraint types
    │   └── npm_test.go          # Comprehensive tests
    └── pypi/                     # PyPI ecosystem package
        ├── pypi.go              # Version, VersionRange, Constraint types  
        └── pypi_test.go         # Comprehensive tests
```

### Key Design Principles

1. **Type Safety**: Each ecosystem has its own types (`npm.Version`, `pypi.Version`)
2. **No Cross-Ecosystem Mixing**: Compile-time errors prevent mixing versions from different ecosystems
3. **Ecosystem-Specific Logic**: Each package implements the exact rules for its ecosystem
4. **Self-Contained**: No shared state or registries that could cause confusion

### Type Safety Example

```go
// ✅ Type-safe - all same ecosystem
npmVer := npm.NewVersion("1.2.3")
npmRange := npm.NewVersionRange("^1.2.0")
npmRange.Contains(npmVer) // This works

// ❌ Compile error - cannot mix ecosystems
pypiVer := pypi.NewVersion("1.2.3") 
npmRange.Contains(pypiVer) // Compile error!
```

## Supported Ecosystems

| Ecosystem | Package | Version Format | Range Syntax |
|-----------|---------|----------------|--------------|
| **NPM** | `ecosystem/npm` | Semantic Versioning | `^1.2.3`, `~1.2.3`, `1.x`, `>=1.0.0 <2.0.0` |
| **PyPI** | `ecosystem/pypi` | PEP 440 | `~=1.2.3`, `>=1.0.0,<2.0.0`, `==1.2.*` |

## NPM Version Syntax

go-univers supports the complete NPM semver range syntax:

```go
// Exact versions
npm.NewVersionRange("1.2.3")

// Comparison operators  
npm.NewVersionRange(">=1.2.3")
npm.NewVersionRange("<2.0.0")

// Caret ranges (compatible within major version)
npm.NewVersionRange("^1.2.3")  // >=1.2.3 <2.0.0

// Tilde ranges (compatible within minor version)  
npm.NewVersionRange("~1.2.3")  // >=1.2.3 <1.3.0

// X-ranges (wildcard matching)
npm.NewVersionRange("1.x")     // >=1.0.0 <2.0.0
npm.NewVersionRange("1.2.x")   // >=1.2.0 <1.3.0

// Hyphen ranges
npm.NewVersionRange("1.2.3 - 2.3.4")  // >=1.2.3 <=2.3.4

// Multiple constraints
npm.NewVersionRange(">=1.0.0 <2.0.0")

// OR logic
npm.NewVersionRange("1.x || 2.x")
```

## PyPI Version Syntax

go-univers is fully PEP 440 compliant and supports all Python versioning features:

### Version Formats
```go
// Basic versions
pypi.NewVersion("1.2.3")

// Versions with epochs
pypi.NewVersion("2!1.2.3")

// Pre-releases
pypi.NewVersion("1.2.3a1")    // Alpha
pypi.NewVersion("1.2.3b2")    // Beta  
pypi.NewVersion("1.2.3rc1")   // Release candidate

// Post-releases
pypi.NewVersion("1.2.3.post1")

// Development releases
pypi.NewVersion("1.2.3.dev1")

// Local versions
pypi.NewVersion("1.2.3+local.1")

// Complex versions
pypi.NewVersion("2!1.2.3a1.post1.dev1+local.1")
```

### Range Operators
```go
// Equality and inequality
pypi.NewVersionRange("==1.2.3")
pypi.NewVersionRange("!=1.2.3")

// Comparison operators
pypi.NewVersionRange(">=1.2.3")
pypi.NewVersionRange("<2.0.0")

// Compatible release (tilde-equals)
pypi.NewVersionRange("~=1.2.3")  // >=1.2.3, <1.3.0

// Wildcard matching
pypi.NewVersionRange("==1.2.*")  // >=1.2.0, <1.3.0
pypi.NewVersionRange("!=1.3.*")  // <1.3.0 or >=1.4.0

// Arbitrary equality (string matching)
pypi.NewVersionRange("===1.2.3")

// Multiple constraints (AND logic)
pypi.NewVersionRange(">=1.0.0, <2.0.0, !=1.5.0")
```

## Common Patterns

### Version Validation
```go
version, err := npm.NewVersion("1.2.3")
if err != nil {
    log.Fatal("Invalid version:", err)
}

if !version.IsValid() {
    log.Fatal("Version failed validation")
}
```

### Range Checking
```go
// Check if a version satisfies multiple ranges
version := npm.NewVersion("1.5.0")
ranges := []string{">=1.0.0", "<2.0.0", "!=1.3.0"}

satisfiesAll := true
for _, rangeStr := range ranges {
    versionRange, _ := npm.NewVersionRange(rangeStr)
    if !versionRange.Contains(version) {
        satisfiesAll = false
        break
    }
}
```

### Version Normalization
```go
version := pypi.NewVersion("01.02.03")
normalized := version.Normalize() // "1.2.3"
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests for a specific ecosystem:
```bash
go test ./ecosystem/npm
go test ./ecosystem/pypi
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add comprehensive tests for any new functionality
4. Ensure all tests pass: `go test ./...`
5. Submit a pull request

When adding new ecosystems:
1. Create a new package under `ecosystem/`
2. Implement the core interfaces defined in `univers.go`
3. Add comprehensive table-driven tests
4. Update this README with examples

## License

[View License](LICENSE)

## Related Projects

- [aboutcode-org/univers](https://github.com/aboutcode-org/univers) - The original Python implementation
- [Package URL specification](https://github.com/package-url/purl-spec) - Standard for package identification