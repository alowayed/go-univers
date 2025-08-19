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
| **Conan** | `pkg/ecosystem/conan` | Extended SemVer | `^1.2.3`, `~1.2.3`, `>=1.0.0`, `<2.0.0` |
| **Composer** | `pkg/ecosystem/composer` | Composer Versioning | `^1.2.3`, `~1.2.3`, `1.2.*`, `>=1.0.0,<2.0.0` |
| **CRAN** | `pkg/ecosystem/cran` | R Package Versioning | `>=1.2.3`, `<2.0.0`, `!=1.5.0` |
| **Debian** | `pkg/ecosystem/debian` | Debian Package Versioning | `>=1.2.3`, `<2.0.0`, `!=1.5.0`, `>>1.0`, `<<2.0` |
| **Gentoo** | `pkg/ecosystem/gentoo` | Gentoo Package Versioning | `>=1.2.3`, `<2.0.0`, `!=1.5.0` |
| **Go** | `pkg/ecosystem/gomod` | Go Module Versioning | `>=v1.2.3`, `<v2.0.0`, `!=v1.3.0` |
| **Maven** | `pkg/ecosystem/maven` | Maven Versioning | `[1.0.0]`, `[1.0.0,2.0.0]`, `(1.0.0,)` |
| **NPM** | `pkg/ecosystem/npm` | Semantic Versioning | `^1.2.3`, `~1.2.3`, `1.x`, `>=1.0.0 <2.0.0` |
| **NuGet** | `pkg/ecosystem/nuget` | SemVer 2.0 + .NET Extensions | `[1.0.0]`, `[1.0.0,2.0.0]`, `>=1.0.0,<2.0.0` |
| **PyPI** | `pkg/ecosystem/pypi` | PEP 440 | `~=1.2.3`, `>=1.0.0,<2.0.0`, `==1.2.*` |
| **RPM** | `pkg/ecosystem/rpm` | RPM Package Versioning | `>=1.2.3`, `<2.0.0`, `!=1.5.0` |
| **RubyGems** | `pkg/ecosystem/gem` | Ruby Gem Versioning | `~> 1.2.3`, `>= 1.0.0`, `!= 1.5.0` |
| **SemVer** | `pkg/ecosystem/semver` | Semantic Versioning 2.0.0 | `>=1.2.3`, `<2.0.0`, `!=1.5.0` |

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

# Compare Debian versions with epoch and revision handling
univers debian compare "1:1.0-1" "2.0"       # → 1 (epoch 1 > epoch 0)
univers debian compare "1.0~beta" "1.0"      # → -1 (tilde sorts before release)

# Compare Gentoo versions with suffix and revision handling
univers gentoo compare "1.0_alpha" "1.0"     # → -1 (alpha < release)
univers gentoo compare "1.0-r1" "1.0"        # → 1 (revision > no revision)
univers gentoo compare "2.3a" "2.3b"         # → -1 (letter a < letter b)

# Compare Ruby Gem versions with prerelease handling
univers gem compare "1.0.0-alpha" "1.0.0"  # → -1 (prerelease < release)
univers gem compare "2.0.0" "1.9.9"        # → 1 (first > second)

# Compare Cargo versions with SemVer 2.0 compliance
univers cargo compare "1.0.0-alpha" "1.0.0"  # → -1 (prerelease < release)
univers cargo compare "1.2.3" "1.2.4"        # → -1 (first < second)

# Compare Conan versions with extended semantic versioning
univers conan compare "1.2.3" "1.2.4"        # → -1 (first < second)
univers conan compare "1.0.2n" "1.0.2m"      # → 1 (letter n > letter m)
univers conan compare "1.2.3-alpha" "1.2.3"  # → -1 (prerelease < release)

# Compare Composer versions with stability flags
univers composer compare "1.2.3-alpha" "1.2.3"  # → -1 (prerelease < stable)
univers composer compare "2.0.0" "1.9.9"         # → 1 (first > second)

# Compare CRAN versions with numeric component handling
univers cran compare "1.2.3" "1.2.4"     # → -1 (first < second)
univers cran compare "1.2" "1.2.0"       # → -1 (fewer components < more components)
univers cran compare "1-2-3" "1.2.3"     # → 0 (dashes normalized to periods)

# Compare NuGet versions with SemVer 2.0 and .NET extensions
univers nuget compare "1.0.0-alpha" "1.0.0"  # → -1 (prerelease < release)
univers nuget compare "1.2.3.4" "1.2.3"      # → 1 (revision > no revision)

# Compare RPM versions with epoch, version, and release handling
univers rpm compare "1:1.2.3-4" "2.0.0"      # → 1 (epoch 1 > epoch 0)
univers rpm compare "1.0~beta" "1.0"         # → -1 (tilde sorts before release)
univers rpm compare "2.0.0-1.el8" "2.0.0-2.el8" # → -1 (release 1 < release 2)

# Compare SemVer versions with strict semantic versioning 2.0.0 rules
univers semver compare "1.0.0-alpha" "1.0.0"      # → -1 (prerelease < release)
univers semver compare "1.0.0-alpha.1" "1.0.0-alpha.beta" # → -1 (numeric < non-numeric)
univers semver compare "1.0.0+build1" "1.0.0+build2" # → 0 (build metadata ignored)
```

#### Sort Versions
```bash
# Sort Alpine versions with proper suffix ordering
univers alpine sort "2.0.0" "1.0.0_alpha" "1.0.0"
# → "1.0.0_alpha" "1.0.0" "2.0.0"

# Sort CRAN versions with component length handling
univers cran sort "1.2.3" "1.2" "1.2.4" "1.10"
# → "1.2" "1.2.3" "1.2.4" "1.10"

# Sort Debian versions with epoch, tilde, and revision handling
univers debian sort "1:1.0" "1.0~beta" "1.0" "1.0-1"
# → "1.0~beta" "1.0" "1.0-1" "1:1.0"

# Sort Gentoo versions with suffix and revision handling
univers gentoo sort "2.0" "1.0_alpha" "1.0" "1.0-r1" "1.0_beta"
# → "1.0_alpha" "1.0_beta" "1.0" "1.0-r1" "2.0"

# Sort Go module versions including pseudo-versions
univers go sort "v2.0.0" "v1.2.3" "v1.0.0-20170915032832-14c0d48ead0c"
# → v1.0.0-20170915032832-14c0d48ead0c, v1.2.3, v2.0.0

# Sort Ruby Gem versions with proper prerelease ordering
univers gem sort "2.0.0" "1.0.0-alpha" "1.0.0"
# → "1.0.0-alpha" "1.0.0" "2.0.0"

# Sort Cargo versions with SemVer 2.0 prerelease rules
univers cargo sort "1.0.0" "1.0.0-beta.1" "1.0.0-beta.11" "1.0.0-alpha"
# → "1.0.0-alpha" "1.0.0-beta.1" "1.0.0-beta.11" "1.0.0"

# Sort Conan versions with extended semantic versioning and letter handling
univers conan sort "1.2.3" "1.0.2n" "1.0.2m" "1.2.3-alpha" "1.2.3.4"
# → "1.0.2m" "1.0.2n" "1.2.3-alpha" "1.2.3" "1.2.3.4"

# Sort Composer versions with stability ordering (dev < alpha < beta < RC < stable)
univers composer sort "1.2.3" "1.2.3-beta" "dev-main" "1.2.3-alpha"
# → "dev-main" "1.2.3-alpha" "1.2.3-beta" "1.2.3"

# Sort NuGet versions with SemVer 2.0 prerelease and revision handling
univers nuget sort "1.0.0" "1.0.0-beta" "1.0.0.1" "1.0.0-alpha"
# → "1.0.0-alpha" "1.0.0-beta" "1.0.0" "1.0.0.1"

# Sort RPM versions with epoch, tilde, and release handling
univers rpm sort "1:1.0" "1.0~beta" "1.0" "1.0-1" "2.0.0-1.el8"
# → "1.0~beta" "1.0" "1.0-1" "2.0.0-1.el8" "1:1.0"

# Sort SemVer versions with strict semantic versioning 2.0.0 precedence
univers semver sort "1.0.0" "1.0.0-beta.2" "1.0.0-beta.11" "1.0.0-alpha" "1.0.0-rc.1"
# → "1.0.0-alpha" "1.0.0-beta.2" "1.0.0-beta.11" "1.0.0-rc.1" "1.0.0"
```

#### Check Range Satisfaction
```bash
# Alpine range checking  
univers alpine contains ">=1.2.0" "1.2.5"     # → true
univers alpine contains "<2.0.0" "1.9.9"      # → true

# Debian range checking with epoch and revision support
univers debian contains ">=1:1.0" "1:1.5-1"   # → true
univers debian contains ">>1.0" "1.0~beta"    # → false (tilde < 1.0)

# Gentoo range checking with suffix and revision support
univers gentoo contains ">=1.0" "1.0_beta"    # → false (beta < release)
univers gentoo contains ">=1.0_alpha" "1.0"   # → true (release > alpha)
univers gentoo contains ">=1.0, <2.0" "1.5-r2" # → true (within range)

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

# Conan range checking with extended semantic versioning
univers conan contains ">=1.2.0" "1.2.5"     # → true
univers conan contains "^1.2.0" "1.5.0"      # → true (compatible within major)
univers conan contains "~1.2.3" "1.2.5"      # → true (patch increment allowed)
univers conan contains ">=1.0.2m" "1.0.2n"   # → true (letter n > letter m)

# Composer constraint checking with caret, tilde, and wildcard ranges
univers composer contains "^1.2.0" "1.3.0"   # → true (compatible within major)
univers composer contains "~1.2.0" "1.2.5"   # → true (compatible within minor)
univers composer contains "1.2.*" "1.2.9"    # → true (wildcard match)

# CRAN range checking with standard operators
univers cran contains ">=1.2.0" "1.2.5"      # → true
univers cran contains ">=1.2.0, <2.0.0" "1.5.0"  # → true (multiple constraints)
univers cran contains "!=1.2.3" "1.2.4"      # → true (not equal constraint)

# NuGet range checking with bracket notation and comma-separated constraints
univers nuget contains "[1.0.0,2.0.0]" "1.5.0"     # → true (inclusive range)
univers nuget contains "[1.0.0,)" "2.0.0"          # → true (unbounded range)
univers nuget contains ">=1.0.0,<2.0.0" "1.5.0"    # → true (comma-separated)

# RPM range checking with epoch and release support
univers rpm contains ">=1:1.0.0" "1:1.5.0-1"      # → true (epoch and version match)
univers rpm contains ">=1.0.0 <2.0.0" "1.5.0-1.el8" # → true (within range)
univers rpm contains ">1.0~beta" "1.0"            # → true (release > tilde)

# SemVer range checking with standard comparison operators
univers semver contains ">=1.0.0,<2.0.0" "1.5.0"        # → true (within range)
univers semver contains ">=1.0.0-alpha" "1.0.0-beta"     # → true (beta > alpha)
univers semver contains "!=1.5.0" "1.5.0+build"         # → false (build metadata ignored)
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

### CRAN (R Package Versioning)
```go
e := &cran.Ecosystem{}

// Numeric component versions
v1, _ := e.NewVersion("1.2")            // Two components
v2, _ := e.NewVersion("1.2.3")          // Three components
v3, _ := e.NewVersion("1-2-3")          // Dashes normalized to periods

// Range operators
r1, _ := e.NewVersionRange(">=1.2.0")     // Standard comparison
r2, _ := e.NewVersionRange(">=1.0, <2.0") // Multiple constraints
r3, _ := e.NewVersionRange("!=1.5.0")     // Not equal constraint

// Check version against range
fmt.Println(r1.Contains(v2)) // true
```

### Conan (C/C++ Package Manager)
```go
e := &conan.Ecosystem{}

// Extended semantic versioning with flexible parts
v1, _ := e.NewVersion("1.2.3")              // Standard semver
v2, _ := e.NewVersion("1.0.2n")             // OpenSSL-style with letter
v3, _ := e.NewVersion("1.2.3.4.5")          // Extended multi-part version
v4, _ := e.NewVersion("1.2.3-alpha")        // With prerelease
v5, _ := e.NewVersion("2.1b-beta+build")    // Complex version with all components

// Range constraints
r1, _ := e.NewVersionRange(">=1.2.0")       // Standard comparison
r2, _ := e.NewVersionRange("^1.2.3")        // Caret: compatible within major
r3, _ := e.NewVersionRange("~1.2.3")        // Tilde: patch-level changes
r4, _ := e.NewVersionRange(">=1.0.0, <2.0.0") // Multiple constraints

// Check version against range
fmt.Println(r2.Contains(v1)) // true
```

### Debian Package Versioning
```go
e := &debian.Ecosystem{}

// Epoch, upstream, and revision components
v1, _ := e.NewVersion("1:2.3.4-5")         // Epoch 1, upstream 2.3.4, revision 5
v2, _ := e.NewVersion("1.0~beta1")          // Tilde for pre-release
v3, _ := e.NewVersion("1.0+dfsg-1ubuntu1")  // Complex revision

// Range constraints
r1, _ := e.NewVersionRange(">=1:1.0-1")     // Epoch-aware ranges
r2, _ := e.NewVersionRange(">>1.0, <<2.0") // Debian-specific operators
r3, _ := e.NewVersionRange(">=1.0, !=1.5") // Multiple constraints
```

### Gentoo Package Versioning
```go
e := &gentoo.Ecosystem{}

// Version formats with suffixes and revisions
v1, _ := e.NewVersion("1.2.3")              // Basic version
v2, _ := e.NewVersion("2.0a")               // Version with letter suffix
v3, _ := e.NewVersion("1.0_alpha1")         // Version with alpha suffix
v4, _ := e.NewVersion("1.5_beta2-r3")       // Complex version with beta suffix and revision

// Range constraints
r1, _ := e.NewVersionRange(">=1.0.0")       // Standard comparison
r2, _ := e.NewVersionRange(">=1.0, <2.0")  // Multiple constraints
r3, _ := e.NewVersionRange(">=1.0_alpha, !=1.2") // Suffix-aware constraints

// Check version against range
fmt.Println(r1.Contains(v3)) // true
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

### RPM Package Versioning
```go
e := &rpm.Ecosystem{}

// Epoch, version, and release components
v1, _ := e.NewVersion("1:2.3.4-5.el8")       // Epoch 1, version 2.3.4, release 5.el8
v2, _ := e.NewVersion("1.0~beta1")            // Tilde for pre-release
v3, _ := e.NewVersion("2.4.37-1.fc34.x86_64") // Fedora release with architecture

// Range constraints
r1, _ := e.NewVersionRange(">=1:1.0.0")      // Epoch-aware ranges
r2, _ := e.NewVersionRange(">=1.0.0 <2.0.0") // Standard comparison operators
r3, _ := e.NewVersionRange(">=1.0.0, !=1.5.0") // Multiple constraints
```

### SemVer (Semantic Versioning 2.0.0)
```go
e := &semver.Ecosystem{}

// Standard semantic versioning formats
v1, _ := e.NewVersion("1.2.3")                  // Basic version
v2, _ := e.NewVersion("1.0.0-alpha.1")          // Prerelease
v3, _ := e.NewVersion("1.0.0+build.1")          // Build metadata
v4, _ := e.NewVersion("1.0.0-beta.2+exp.sha.1") // Both prerelease and build

// Range constraints with standard comparison operators
r1, _ := e.NewVersionRange(">=1.0.0,<2.0.0")    // Comma-separated constraints
r2, _ := e.NewVersionRange(">=1.0.0 <2.0.0")    // Space-separated constraints
r3, _ := e.NewVersionRange("!=1.5.0")           // Not equal constraint

// Check version against range (build metadata ignored in comparisons)
fmt.Println(r1.Contains(v1)) // true
```

For complete syntax documentation of all ecosystems, see the [Supported Ecosystems](#supported-ecosystems) table above.

## Related Projects

- [aboutcode-org/univers](https://github.com/aboutcode-org/univers) - The original Python implementation
- [Package URL specification](https://github.com/package-url/purl-spec) - Standard for package identification