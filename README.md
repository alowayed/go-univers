# go-univers: mostly universal version and version ranges comparison and conversion

[![Go](https://github.com/alowayed/go-univers/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/alowayed/go-univers/actions/workflows/go.yml)

A Go library to:
1. Parse and compare versions.
2. Parse ragnes and check if it contains a version.
3. Sort versions.

## Supported Ecosystems

| Ecosystem | Package | Version Format | Range Syntax |
|-----------|---------|----------------|--------------|
| **NPM** | `ecosystem/npm` | Semantic Versioning | `^1.2.3`, `~1.2.3`, `1.x`, `>=1.0.0 <2.0.0` |
| **PyPI** | `ecosystem/pypi` | PEP 440 | `~=1.2.3`, `>=1.0.0,<2.0.0`, `==1.2.*` |
| **Go** | `ecosystem/gomod` | Go Module Versioning | `>=v1.2.3`, `<v2.0.0`, `!=v1.3.0` |

## Installation

```bash
go get github.com/alowayed/go-univers
```

## Quick Start

```go
package main

import (
    "fmt"
    "slices"
    "github.com/alowayed/go-univers/pkg/ecosystem/npm"
)

func main() {
    // Parse versions.
    v1, _ := npm.NewVersion("1.2.3")
    v2, _ := npm.NewVersion("1.2.4-alpha")
    
    // Compare versions.
    if v1.Compare(v2) < 0 {
        fmt.Println("v1 is older than v2")
    }
    
    // Parse version ranges.
    range1, _ := npm.NewVersionRange("^1.2.0")        // Caret range
    range2, _ := npm.NewVersionRange("~1.2.0")        // Tilde range  
    range3, _ := npm.NewVersionRange("1.x")           // X-range
    range4, _ := npm.NewVersionRange(">=1.0.0 <2.0.0") // Multiple constraints
    
    // Check if version satisfies range
    if range1.Contains(v1) {
        fmt.Println("v1 satisfies ^1.2.0")
    }
    
    // Sort versions.
    versions := []*npm.Version{v2, v1}
    slices.SortFunc(versions, (*npm.Version).Compare)
    fmt.Printf("Sorted: %+v\n", versions) // {v1, v2}
}
```

## CLI

go-univers aprovides a command-line interface for version operations.

### Building the CLI

```bash
# Build the CLI binary
go build -o univers ./cmd
```

### CLI Usage

The CLI follows the pattern: `univers <ecosystem> <command> [args]`

#### Compare Versions
```bash
# Compare two NPM versions (outputs -1, 0, or 1)
univers npm compare "1.2.3" "1.2.4"     # → -1 (first < second)
univers npm compare "2.0.0" "1.9.9"     # → 1 (first > second)
univers npm compare "1.2.3" "1.2.3"     # → 0 (equal)
```

#### Sort Versions
```bash
# Sort Go module versions including pseudo-versions
univers go sort "v2.0.0" "v1.2.3" "v1.0.0-20170915032832-14c0d48ead0c"
# → v1.0.0-20170915032832-14c0d48ead0c, v1.2.3, v2.0.0
```

#### Check Range Satisfaction
```bash
# PyPI range checking
univers pypi contains "~=1.2.0" "1.2.5"   # → true
univers pypi contains "==1.2.*" "1.2.5"   # → true

# Go module range checking
univers go contains ">=v1.2.0 <v2.0.0" "v1.5.0"  # → true
univers go contains "<v1.9.0" "v2.0.0"    # → false
```

## Architecture

go-univers uses a **type-safe, ecosystem-isolated architecture** that prevents accidental cross-ecosystem version mixing:

```
go-univers/
├── cmd/
│   └── main.go                 # CLI application
└── pkg/
    ├── univers.go              # Universal interfaces (documentation)
    └── ecosystem/
        ├── gomod/              # Go ecosystem package
        │   ├── gomod.go        # Version, VersionRange, Constraint types
        │   └── gomod_test.go   # Comprehensive tests
        ├── npm/                # NPM ecosystem package
        │   ├── npm.go          # Version, VersionRange, Constraint types
        │   └── npm_test.go     # Comprehensive tests
        │
        └── [OTHER ECOSYSTEMS]
```

## Ecosystem specific

### NPM

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

## PyPI

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

## Go

### Version Formats
```go
// Basic semantic versions
gomod.NewVersion("v1.2.3")
gomod.NewVersion("1.2.3")  // Automatically prefixed with 'v'

// Versions with prerelease
gomod.NewVersion("v1.2.3-beta")
gomod.NewVersion("v1.2.3-alpha.1")
gomod.NewVersion("v1.2.3-rc.1")

// Versions with build metadata
gomod.NewVersion("v1.2.3+build.1")

// Complex versions with prerelease and build
gomod.NewVersion("v1.2.3-beta.1+build.20230101")

// Pseudo-versions (generated by Go tools)
gomod.NewVersion("v1.0.0-20170915032832-14c0d48ead0c")        // Pattern 1: no base version
gomod.NewVersion("v1.2.3-beta.0.20170915032832-14c0d48ead0c") // Pattern 2: prerelease base
gomod.NewVersion("v1.2.4-0.20170915032832-14c0d48ead0c")      // Pattern 3: release base
```

### Range Operators
```go
// Equality and inequality
gomod.NewVersionRange("v1.2.3")       // Exact match
gomod.NewVersionRange("!=v1.2.3")     // Not equal

// Comparison operators
gomod.NewVersionRange(">=v1.2.3")     // Greater than or equal
gomod.NewVersionRange(">v1.2.3")      // Greater than
gomod.NewVersionRange("<=v1.2.3")     // Less than or equal
gomod.NewVersionRange("<v2.0.0")      // Less than

// Multiple constraints (AND logic)
gomod.NewVersionRange(">=v1.2.3 <v2.0.0")           // Range constraint
gomod.NewVersionRange(">=v1.0.0 <v2.0.0 !=v1.5.0")  // With exclusion
```

### Pseudo-Version Support

```go
// Pseudo-versions are automatically recognized and parsed
pseudo, _ := gomod.NewVersion("v1.0.0-20170915032832-14c0d48ead0c")

// Pseudo-versions compare correctly with regular versions
regular, _ := gomod.NewVersion("v1.0.0")
fmt.Println(pseudo.Compare(regular)) // -1 (pseudo-versions are pre-release)

// Pseudo-versions can be used in ranges
range1, _ := gomod.NewVersionRange(">=v1.0.0-20170915032832-14c0d48ead0c")
```

## Related Projects

- [aboutcode-org/univers](https://github.com/aboutcode-org/univers) - The original Python implementation
- [Package URL specification](https://github.com/package-url/purl-spec) - Standard for package identification
