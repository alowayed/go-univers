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

# NuGet range checking with bracket notation and comma-separated constraints
univers nuget contains "[1.0.0,2.0.0]" "1.5.0"     # → true (inclusive range)
univers nuget contains "[1.0.0,)" "2.0.0"          # → true (unbounded range)
univers nuget contains ">=1.0.0,<2.0.0" "1.5.0"    # → true (comma-separated)
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

### Alpine

#### Version Formats
```go
e := &alpine.Ecosystem{}

// Basic versions
e.NewVersion("1.2.3")

// Versions with letters
e.NewVersion("1.2.3a")         // Letter suffix
e.NewVersion("2.3.0b")         // Letter suffix

// Versions with suffixes
e.NewVersion("1.2.3_alpha")    // Alpha release
e.NewVersion("1.3_alpha2")     // Alpha with number
e.NewVersion("1.2.3_beta")     // Beta release
e.NewVersion("1.2.3_pre")      // Pre-release
e.NewVersion("1.2.3_rc")       // Release candidate
e.NewVersion("0.1.0_alpha_pre2") // Multiple suffixes

// Versions with build components
e.NewVersion("1.0.4-r3")       // Build revision
e.NewVersion("20050718-r2")    // Date-based with build

// Versions with hash components
e.NewVersion("1.2.3~abc123")   // Commit hash
e.NewVersion("1.2.3~abc123-r1") // Hash with build

// Complex versions
e.NewVersion("2.3.0b-r1")      // Letter and build
e.NewVersion("1.2.3a_beta2-r5") // Letter, suffix, and build
```

#### Range Operators
```go
e := &alpine.Ecosystem{}

// Equality and inequality
e.NewVersionRange("1.2.3")         // Exact match
e.NewVersionRange("= 1.2.3")       // Explicit equals
e.NewVersionRange("!= 1.2.3")      // Not equal

// Comparison operators
e.NewVersionRange(">= 1.2.3")      // Greater than or equal
e.NewVersionRange("> 1.2.3")       // Greater than
e.NewVersionRange("<= 1.2.3")      // Less than or equal
e.NewVersionRange("< 1.2.3")       // Less than

// Multiple constraints (AND logic)
e.NewVersionRange(">= 1.0.0 < 2.0.0")      // Range constraint
e.NewVersionRange(">= 1.2.3 < 2.0.0 != 1.5.0") // With exclusion

// Alpine-specific version formats in ranges
e.NewVersionRange(">= 1.2.0_alpha")        // Suffix versions
e.NewVersionRange(">= 1.2.3-r1")           // Build versions
e.NewVersionRange("> 1.1a")                // Letter versions
```

### Cargo

#### Version Formats
```go
e := &cargo.Ecosystem{}

// Basic SemVer 2.0 versions
e.NewVersion("1.2.3")

// Versions with prerelease identifiers
e.NewVersion("1.0.0-alpha")         // Alpha release
e.NewVersion("1.0.0-beta.1")        // Beta with number
e.NewVersion("1.0.0-rc.1")          // Release candidate
e.NewVersion("2.0.0-alpha.beta.1")  // Complex prerelease

// Versions with build metadata
e.NewVersion("1.0.0+build.1")       // Build metadata
e.NewVersion("1.0.0-alpha+build")   // Prerelease with build

// Complex versions with all components
e.NewVersion("1.0.0-beta.1+build.20230101")
```

#### Range Operators
```go
e := &cargo.Ecosystem{}

// Equality and inequality
e.NewVersionRange("1.2.3")         // Exact match
e.NewVersionRange("=1.2.3")        // Explicit equals
e.NewVersionRange("!=1.2.3")       // Not equal

// Comparison operators
e.NewVersionRange(">=1.2.3")       // Greater than or equal
e.NewVersionRange(">1.2.3")        // Greater than
e.NewVersionRange("<=1.2.3")       // Less than or equal
e.NewVersionRange("<2.0.0")        // Less than

// Caret constraints (compatible within major)
e.NewVersionRange("^1.2.3")        // >=1.2.3 <2.0.0
e.NewVersionRange("^0.2.3")        // >=0.2.3 <0.3.0 (special 0.x behavior)
e.NewVersionRange("^0.0.3")        // =0.0.3 (special 0.0.x behavior)

// Tilde constraints (compatible within minor)
e.NewVersionRange("~1.2.3")        // >=1.2.3 <1.3.0
e.NewVersionRange("~1.2")          // >=1.2.0 <1.3.0

// Wildcard matching
e.NewVersionRange("1.2.*")         // >=1.2.0 <1.3.0

// Multiple constraints (AND logic)
e.NewVersionRange(">=1.0.0, <2.0.0")      // Range constraint
e.NewVersionRange(">=1.2.3, <2.0.0, !=1.5.0") // With exclusion
```

### Maven

#### Version Formats
```go
e := &maven.Ecosystem{}

// Basic versions
e.NewVersion("1.2.3")

// Versions with qualifiers
e.NewVersion("1.2.3-alpha")     // Alpha release
e.NewVersion("1.2.3-beta")      // Beta release
e.NewVersion("1.2.3-milestone") // Milestone release
e.NewVersion("1.2.3-rc")        // Release candidate
e.NewVersion("1.2.3-snapshot")  // Snapshot release
e.NewVersion("1.2.3-sp")        // Service pack

// Normalized qualifiers (equivalent to release)
e.NewVersion("1.2.3-ga")        // General availability (same as 1.2.3)
e.NewVersion("1.2.3-final")     // Final release (same as 1.2.3)
e.NewVersion("1.2.3-release")   // Release (same as 1.2.3)

// Qualifier shortcuts
e.NewVersion("1.2.3-a")         // Short for alpha
e.NewVersion("1.2.3-b")         // Short for beta
e.NewVersion("1.2.3-m")         // Short for milestone
```

#### Range Operators
```go
e := &maven.Ecosystem{}

// Exact version match
e.NewVersionRange("[1.2.3]")

// Inclusive ranges
e.NewVersionRange("[1.0.0,2.0.0]")  // >=1.0.0 and <=2.0.0

// Exclusive ranges
e.NewVersionRange("(1.0.0,2.0.0)")  // >1.0.0 and <2.0.0

// Mixed inclusive/exclusive
e.NewVersionRange("[1.0.0,2.0.0)")  // >=1.0.0 and <2.0.0
e.NewVersionRange("(1.0.0,2.0.0]")  // >1.0.0 and <=2.0.0

// Unbounded ranges
e.NewVersionRange("[1.0.0,)")       // >=1.0.0
e.NewVersionRange("(,2.0.0]")       // <=2.0.0
e.NewVersionRange("(,2.0.0)")       // <2.0.0

// Simple version (equivalent to exact match)
e.NewVersionRange("1.2.3")          // Same as [1.2.3]
```

### RubyGems

#### Version Formats
```go
e := &gem.Ecosystem{}

// Basic versions
e.NewVersion("1.2.3")

// Versions with prerelease identifiers
e.NewVersion("1.2.3-alpha")      // Alpha release
e.NewVersion("1.2.3-beta")       // Beta release  
e.NewVersion("1.2.3-rc1")        // Release candidate
e.NewVersion("1.2.3.pre")        // Pre-release format
e.NewVersion("2.0.0.rc1")        // RC with numbers

// Build metadata
e.NewVersion("1.0.0+build.1")    // Build metadata
e.NewVersion("1.0.0-alpha+build") // Prerelease with build

// Complex versions
e.NewVersion("1.0.0-beta.1")     // Complex prerelease
e.NewVersion("v1.0.0")           // With v prefix
```

#### Range Operators
```go
e := &gem.Ecosystem{}

// Equality and inequality
e.NewVersionRange("1.2.3")         // Exact match
e.NewVersionRange("= 1.2.3")       // Explicit equals
e.NewVersionRange("!= 1.2.3")      // Not equal

// Comparison operators
e.NewVersionRange(">= 1.2.3")      // Greater than or equal
e.NewVersionRange("> 1.2.3")       // Greater than
e.NewVersionRange("<= 1.2.3")      // Less than or equal
e.NewVersionRange("< 1.2.3")       // Less than

// Pessimistic constraint (twiddle-wakka)
e.NewVersionRange("~> 1.2.3")      // >= 1.2.3, < 1.3.0
e.NewVersionRange("~> 1.2")        // >= 1.2.0, < 2.0.0
e.NewVersionRange("~> 1")          // >= 1.0.0, < 2.0.0

// Multiple constraints (AND logic)
e.NewVersionRange("~> 1.2.3, >= 1.2.5")    // Pessimistic with minimum
e.NewVersionRange(">= 1.0.0, < 2.0.0")     // Range constraint
e.NewVersionRange("~> 2.0, != 2.1.0")      // Pessimistic with exclusion
```

### NPM

```go
e := &npm.Ecosystem{}

// Exact versions
e.NewVersionRange("1.2.3")

// Comparison operators  
e.NewVersionRange(">=1.2.3")
e.NewVersionRange("<2.0.0")

// Caret ranges (compatible within major version)
e.NewVersionRange("^1.2.3")  // >=1.2.3 <2.0.0

// Tilde ranges (compatible within minor version)  
e.NewVersionRange("~1.2.3")  // >=1.2.3 <1.3.0

// X-ranges (wildcard matching)
e.NewVersionRange("1.x")     // >=1.0.0 <2.0.0
e.NewVersionRange("1.2.x")   // >=1.2.0 <1.3.0

// Hyphen ranges
e.NewVersionRange("1.2.3 - 2.3.4")  // >=1.2.3 <=2.3.4

// Multiple constraints
e.NewVersionRange(">=1.0.0 <2.0.0")

// OR logic
e.NewVersionRange("1.x || 2.x")
```

### NuGet

#### Version Formats
```go
e := &nuget.Ecosystem{}

// Basic versions
e.NewVersion("1.2.3")

// Versions with v prefix
e.NewVersion("v1.2.3")

// Versions with revision (.NET 4th component)
e.NewVersion("1.2.3.4")

// Versions with prerelease identifiers
e.NewVersion("1.0.0-alpha")         // Alpha release
e.NewVersion("1.0.0-beta.1")        // Beta with number
e.NewVersion("1.0.0-rc.1")          // Release candidate
e.NewVersion("2.0.0-alpha.beta.1")  // Complex prerelease

// Versions with build metadata
e.NewVersion("1.0.0+build.1")       // Build metadata
e.NewVersion("1.0.0-alpha+build")   // Prerelease with build

// Complex versions with revision, prerelease, and build
e.NewVersion("1.2.3.4-beta.1+build.20230101")

// Major only (defaults to .0.0)
e.NewVersion("1")

// Major and minor only (defaults to .0)
e.NewVersion("1.2")
```

#### Range Operators
```go
e := &nuget.Ecosystem{}

// Exact version match
e.NewVersionRange("[1.2.3]")

// Inclusive ranges
e.NewVersionRange("[1.0.0,2.0.0]")  // >=1.0.0 and <=2.0.0

// Exclusive ranges
e.NewVersionRange("(1.0.0,2.0.0)")  // >1.0.0 and <2.0.0

// Mixed inclusive/exclusive
e.NewVersionRange("[1.0.0,2.0.0)")  // >=1.0.0 and <2.0.0
e.NewVersionRange("(1.0.0,2.0.0]")  // >1.0.0 and <=2.0.0

// Unbounded ranges
e.NewVersionRange("[1.0.0,)")       // >=1.0.0
e.NewVersionRange("(,2.0.0]")       // <=2.0.0
e.NewVersionRange("(,2.0.0)")       // <2.0.0

// Comma-separated constraints (AND logic)
e.NewVersionRange(">=1.0.0,<2.0.0") // >=1.0.0 AND <2.0.0
e.NewVersionRange(">=1.0.0,<2.0.0,!=1.5.0") // With exclusion

// Minimum version (default behavior)
e.NewVersionRange("1.0.0")          // >=1.0.0
```

### PyPI

#### Version Formats
```go
e := &pypi.Ecosystem{}

// Basic versions
e.NewVersion("1.2.3")

// Versions with epochs
e.NewVersion("2!1.2.3")

// Pre-releases
e.NewVersion("1.2.3a1")    // Alpha
e.NewVersion("1.2.3b2")    // Beta  
e.NewVersion("1.2.3rc1")   // Release candidate

// Post-releases
e.NewVersion("1.2.3.post1")

// Development releases
e.NewVersion("1.2.3.dev1")

// Local versions
e.NewVersion("1.2.3+local.1")

// Complex versions
e.NewVersion("2!1.2.3a1.post1.dev1+local.1")
```

#### Range Operators
```go
e := &pypi.Ecosystem{}

// Equality and inequality
e.NewVersionRange("==1.2.3")
e.NewVersionRange("!=1.2.3")

// Comparison operators
e.NewVersionRange(">=1.2.3")
e.NewVersionRange("<2.0.0")

// Compatible release (tilde-equals)
e.NewVersionRange("~=1.2.3")  // >=1.2.3, <1.3.0

// Wildcard matching
e.NewVersionRange("==1.2.*")  // >=1.2.0, <1.3.0
e.NewVersionRange("!=1.3.*")  // <1.3.0 or >=1.4.0

// Arbitrary equality (string matching)
e.NewVersionRange("===1.2.3")

// Multiple constraints (AND logic)
e.NewVersionRange(">=1.0.0, <2.0.0, !=1.5.0")
```

### Go

#### Version Formats
```go
e := &gomod.Ecosystem{}

// Basic semantic versions
e.NewVersion("v1.2.3")
e.NewVersion("1.2.3")  // Automatically prefixed with 'v'

// Versions with prerelease
e.NewVersion("v1.2.3-beta")
e.NewVersion("v1.2.3-alpha.1")
e.NewVersion("v1.2.3-rc.1")

// Versions with build metadata
e.NewVersion("v1.2.3+build.1")

// Complex versions with prerelease and build
e.NewVersion("v1.2.3-beta.1+build.20230101")

// Pseudo-versions (generated by Go tools)
e.NewVersion("v1.0.0-20170915032832-14c0d48ead0c")        // Pattern 1: no base version
e.NewVersion("v1.2.3-beta.0.20170915032832-14c0d48ead0c") // Pattern 2: prerelease base
e.NewVersion("v1.2.4-0.20170915032832-14c0d48ead0c")      // Pattern 3: release base
```

#### Range Operators
```go
e := &gomod.Ecosystem{}

// Equality and inequality
e.NewVersionRange("v1.2.3")       // Exact match
e.NewVersionRange("!=v1.2.3")     // Not equal

// Comparison operators
e.NewVersionRange(">=v1.2.3")     // Greater than or equal
e.NewVersionRange(">v1.2.3")      // Greater than
e.NewVersionRange("<=v1.2.3")     // Less than or equal
e.NewVersionRange("<v2.0.0")      // Less than

// Multiple constraints (AND logic)
e.NewVersionRange(">=v1.2.3 <v2.0.0")           // Range constraint
e.NewVersionRange(">=v1.0.0 <v2.0.0 !=v1.5.0")  // With exclusion
```

#### Pseudo-Version Support

```go
e := &gomod.Ecosystem{}

// Pseudo-versions are automatically recognized and parsed
pseudo, _ := e.NewVersion("v1.0.0-20170915032832-14c0d48ead0c")

// Pseudo-versions compare correctly with regular versions
regular, _ := e.NewVersion("v1.0.0")
fmt.Println(pseudo.Compare(regular)) // -1 (pseudo-versions are pre-release)

// Pseudo-versions can be used in ranges
r, _ := e.NewVersionRange(">=v1.0.0-20170915032832-14c0d48ead0c")
```

## Related Projects

- [aboutcode-org/univers](https://github.com/aboutcode-org/univers) - The original Python implementation
- [Package URL specification](https://github.com/package-url/purl-spec) - Standard for package identification
