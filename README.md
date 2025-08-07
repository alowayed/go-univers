# go-univers: mostly universal version and version ranges comparison and conversion

[![Go](https://github.com/alowayed/go-univers/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/alowayed/go-univers/actions/workflows/go.yml)

A Go library to:
1. Parse and compare versions.
2. Parse ranges and check if it contains a version.
3. Sort versions.

## Supported Ecosystems

| Ecosystem | Package | Version Format | Range Syntax |
|-----------|---------|----------------|--------------|
| **Alpine** | `pkg/ecosystem/alpine` | Alpine Package Versioning | `>=1.2.3`, `<2.0.0`, `!=1.5.0` |
| **Cargo** | `pkg/ecosystem/cargo` | SemVer 2.0 | `^1.2.3`, `~1.2.3`, `>=1.0.0`, `1.2.*` |
| **Composer** | `pkg/ecosystem/composer` | Composer Versioning | `^1.2.3`, `~1.2.3`, `1.2.*`, `>=1.0.0,<2.0.0` |
| **Go** | `pkg/ecosystem/gomod` | Go Module Versioning | `>=v1.2.3`, `<v2.0.0`, `!=v1.3.0` |
| **Maven** | `pkg/ecosystem/maven` | Maven Versioning | `[1.0.0]`, `[1.0.0,2.0.0]`, `(1.0.0,)` |
| **NPM** | `pkg/ecosystem/npm` | Semantic Versioning | `^1.2.3`, `~1.2.3`, `1.x`, `>=1.0.0 <2.0.0` |
| **NuGet** | `pkg/ecosystem/nuget` | SemVer 2.0 + .NET Extensions | `[1.0.0]`, `[1.0.0,2.0.0]`, `>=1.0.0,<2.0.0` |
| **PyPI** | `pkg/ecosystem/pypi` | PEP 440 | `~=1.2.3`, `>=1.0.0,<2.0.0`, `==1.2.*` |
| **RubyGems** | `pkg/ecosystem/gem` | Ruby Gem Versioning | `~> 1.2.3`, `>= 1.0.0`, `!= 1.5.0` |

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
    // Create ecosystem instance
    e := &npm.Ecosystem{}
    
    // Parse versions.
    v1, _ := e.NewVersion("1.2.3")
    v2, _ := e.NewVersion("1.2.4-alpha")
    
    // Compare versions.
    if v1.Compare(v2) < 0 {
        fmt.Println("v1 is older than v2")
    }
    
    // Parse version ranges.
    r1, _ := e.NewVersionRange("^1.2.0")        // Caret range
    r2, _ := e.NewVersionRange("~1.2.0")        // Tilde range  
    r3, _ := e.NewVersionRange("1.x")           // X-range
    r4, _ := e.NewVersionRange(">=1.0.0 <2.0.0") // Multiple constraints
    
    // Check if version satisfies range
    if r1.Contains(v1) {
        fmt.Println("v1 satisfies ^1.2.0")
    }
    
    // Sort versions.
    versions := []*npm.Version{v2, v1}
    slices.SortFunc(versions, (*npm.Version).Compare)
    fmt.Printf("Sorted: %+v\n", versions) // {v1, v2}
}
```

## CLI

go-univers provides a command-line interface for version operations.

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

# Compare Alpine versions with suffix handling
univers alpine compare "1.0.0_alpha" "1.0.0"  # → -1 (alpha < release)
univers alpine compare "2.0.0" "1.9.9"         # → 1 (first > second)

# Compare Ruby Gem versions with prerelease handling
univers gem compare "1.0.0-alpha" "1.0.0"  # → -1 (prerelease < release)
univers gem compare "2.0.0" "1.9.9"        # → 1 (first > second)

# Compare Cargo versions with SemVer 2.0 compliance
univers cargo compare "1.0.0-alpha" "1.0.0"  # → -1 (prerelease < release)
univers cargo compare "1.2.3" "1.2.4"        # → -1 (first < second)

# Compare Composer versions with stability flags
univers composer compare "1.2.3-alpha" "1.2.3"  # → -1 (prerelease < stable)
univers composer compare "2.0.0" "1.9.9"         # → 1 (first > second)

# Compare NuGet versions with SemVer 2.0 and .NET extensions
univers nuget compare "1.0.0-alpha" "1.0.0"  # → -1 (prerelease < release)
univers nuget compare "1.2.3.4" "1.2.3"      # → 1 (revision > no revision)
```

#### Sort Versions
```bash
# Sort Alpine versions with proper suffix ordering
univers alpine sort "2.0.0" "1.0.0_alpha" "1.0.0"
# → "1.0.0_alpha" "1.0.0" "2.0.0"

# Sort Go module versions including pseudo-versions
univers go sort "v2.0.0" "v1.2.3" "v1.0.0-20170915032832-14c0d48ead0c"
# → v1.0.0-20170915032832-14c0d48ead0c, v1.2.3, v2.0.0

# Sort Ruby Gem versions with proper prerelease ordering
univers gem sort "2.0.0" "1.0.0-alpha" "1.0.0"
# → "1.0.0-alpha" "1.0.0" "2.0.0"

# Sort Cargo versions with SemVer 2.0 prerelease rules
univers cargo sort "1.0.0" "1.0.0-beta.1" "1.0.0-beta.11" "1.0.0-alpha"
# → "1.0.0-alpha" "1.0.0-beta.1" "1.0.0-beta.11" "1.0.0"

# Sort Composer versions with stability ordering (dev < alpha < beta < RC < stable)
univers composer sort "1.2.3" "1.2.3-beta" "dev-main" "1.2.3-alpha"
# → "dev-main" "1.2.3-alpha" "1.2.3-beta" "1.2.3"

# Sort NuGet versions with SemVer 2.0 prerelease and revision handling
univers nuget sort "1.0.0" "1.0.0-beta" "1.0.0.1" "1.0.0-alpha"
# → "1.0.0-alpha" "1.0.0-beta" "1.0.0" "1.0.0.1"
```

#### Check Range Satisfaction
```bash
# Alpine range checking  
univers alpine contains ">=1.2.0" "1.2.5"     # → true
univers alpine contains "<2.0.0" "1.9.9"      # → true

# PyPI range checking
univers pypi contains "~=1.2.0" "1.2.5"   # → true
univers pypi contains "==1.2.*" "1.2.5"   # → true

# Go module range checking
univers go contains ">=v1.2.0 <v2.0.0" "v1.5.0"  # → true
univers go contains "<v1.9.0" "v2.0.0"    # → false

# Ruby Gem pessimistic constraint checking
univers gem contains "~> 1.2.0" "1.2.5"  # → true (patch increment allowed)
univers gem contains "~> 1.2.0" "1.3.0"  # → false (minor increment not allowed)

# Cargo constraint checking with SemVer 2.0 caret/tilde ranges
univers cargo contains "^1.2.0" "1.2.5"   # → true (compatible within major)
univers cargo contains "^1.2.0" "2.0.0"   # → false (major increment not allowed)
univers cargo contains "~1.2.0" "1.2.5"   # → true (patch increment allowed)

# Composer constraint checking with caret, tilde, and wildcard ranges
univers composer contains "^1.2.0" "1.3.0"   # → true (compatible within major)
univers composer contains "~1.2.0" "1.2.5"   # → true (compatible within minor)
univers composer contains "1.2.*" "1.2.9"    # → true (wildcard match)

# NuGet range checking with bracket notation and comma-separated constraints
univers nuget contains "[1.0.0,2.0.0]" "1.5.0"     # → true (inclusive range)
univers nuget contains "[1.0.0,)" "2.0.0"          # → true (unbounded range)
univers nuget contains ">=1.0.0,<2.0.0" "1.5.0"    # → true (comma-separated)
```

## Architecture

go-univers uses a **type-safe, ecosystem-isolated architecture** that prevents accidental cross-ecosystem version mixing. Each ecosystem (npm, pypi, go, etc.) has its own `Version` and `VersionRange` types, eliminating the common bug of accidentally comparing versions from different package managers.

See [DEVELOPMENT.md](./DEVELOPMENT.md) for detailed architecture documentation.

## Ecosystem Examples

Each ecosystem has its own version and range syntax. Here are key patterns:

### NPM (Semantic Versioning)
```go
e := &npm.Ecosystem{}

// Parse versions and ranges
v1, _ := e.NewVersion("1.2.3")
r1, _ := e.NewVersionRange("^1.2.0")        // Caret: >=1.2.0 <2.0.0
r2, _ := e.NewVersionRange("~1.2.0")        // Tilde: >=1.2.0 <1.3.0
r3, _ := e.NewVersionRange(">=1.0.0 <2.0.0") // Multiple constraints

// Check version against range
fmt.Println(r1.Contains(v1)) // true
```

### PyPI (PEP 440)
```go
e := &pypi.Ecosystem{}

// Complex version formats
v1, _ := e.NewVersion("1.2.3a1")           // Alpha release
v2, _ := e.NewVersion("2!1.2.3.post1")     // Epoch and post-release
v3, _ := e.NewVersion("1.2.3+local.1")     // Local version

// Range operators
r1, _ := e.NewVersionRange("~=1.2.3")      // Compatible release
r2, _ := e.NewVersionRange("==1.2.*")      // Wildcard matching
r3, _ := e.NewVersionRange(">=1.0.0, <2.0.0, !=1.5.0") // Multiple constraints
```

### Go Modules
```go
e := &gomod.Ecosystem{}

// Regular and pseudo-versions
v1, _ := e.NewVersion("v1.2.3")
v2, _ := e.NewVersion("v1.2.3-beta")
pseudo, _ := e.NewVersion("v1.0.0-20170915032832-14c0d48ead0c")

// Range constraints
r1, _ := e.NewVersionRange(">=v1.2.3 <v2.0.0")
```

For complete syntax documentation of all ecosystems, see the [Supported Ecosystems](#supported-ecosystems) table above.

## Related Projects

- [aboutcode-org/univers](https://github.com/aboutcode-org/univers) - The original Python implementation
- [Package URL specification](https://github.com/package-url/purl-spec) - Standard for package identification