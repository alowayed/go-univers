# go-univers: mostly universal version and version ranges comparison and conversion

[![Go](https://github.com/alowayed/go-univers/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/alowayed/go-univers/actions/workflows/go.yml)

A Go library to:
1. Parse and compare versions.
2. Parse ranges and check if it contains a version.
3. Sort versions.

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
    "github.com/alowayed/go-univers/pkg/spec/vers"
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
    
    // VERS range checking
    result, _ := vers.Contains("vers:npm/>=1.2.0|<=2.0.0", "1.5.0")
    fmt.Printf("VERS result: %t\n", result) // true
}
```

## Supported Ecosystems

| Ecosystem | Package | [VERS versioning scheme](https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst#some-of-the-known-versioning-schemes) |
|-----------|---------|-----------|
| **Alpine** | `pkg/ecosystem/alpine` | `alpine` ✅ |
| **Apache** | ❌ | `apache` ❌ |
| **Arch Linux (ALPM)** | ❌ | `alpm` ❌ |
| **Cargo** | `pkg/ecosystem/cargo` | `cargo` ✅ |
| **Conan** | `pkg/ecosystem/conan` | [`conan` ❌](https://github.com/alowayed/go-univers/issues/59) |
| **Composer** | `pkg/ecosystem/composer` | [`composer` ❌](https://github.com/alowayed/go-univers/issues/54) |
| **CRAN** | `pkg/ecosystem/cran` | ❌ |
| **Debian** | `pkg/ecosystem/debian` | `deb` ✅ |
| **Gentoo** | `pkg/ecosystem/gentoo` | [`ebuild` ❌](https://github.com/alowayed/go-univers/issues/70) |
| **GitHub** | ❌ | `github` ❌ |
| **Go** | `pkg/ecosystem/gomod` | `golang` ✅ |
| **Hex (Elixir)** | ❌ | `hex` ❌ |
| **Intdot** | ❌ | `intdot` ❌ |
| **Mattermost** | ❌ | `mattermost` ❌ |
| **Maven** | `pkg/ecosystem/maven` | `maven` ✅ |
| **Mozilla** | ❌ | `mozilla` ❌ |
| **Nginx** | ❌ | `nginx` ❌ |
| **NPM** | `pkg/ecosystem/npm` | `npm` ✅ |
| **NuGet** | `pkg/ecosystem/nuget` | `nuget` ✅ |
| **OpenSSL** | ❌ | `openssl` ❌ |
| **PyPI** | `pkg/ecosystem/pypi` | `pypi` ✅ |
| **RPM** | `pkg/ecosystem/rpm` | `rpm` ✅ |
| **RubyGems** | `pkg/ecosystem/gem` | `gem` ✅ |
| **SemVer** | `pkg/ecosystem/semver` | `generic` ✅ |

## CLI

go-univers provides a command-line interface for version operations:

```bash
# Build the CLI binary
go build -o univers ./cmd
```

The CLI follows the pattern: `univers <ecosystem|spec> <command> [args]`

### Examples

```bash
# Compare versions (outputs -1, 0, or 1)
univers npm compare "1.2.3" "1.2.4"           # → -1 (first < second)
univers pypi compare "2.0.0" "1.9.9"          # → 1 (first > second)
univers semver compare "1.2.3" "1.2.3"        # → 0 (equal)

# Sort versions in ascending order
univers gem sort "2.0.0" "1.0.0-alpha" "1.0.0"
# → "1.0.0-alpha" "1.0.0" "2.0.0"

# Check if version satisfies range (outputs true/false)
univers cargo contains "^1.2.0" "1.2.5"       # → true
univers maven contains "[1.0.0,2.0.0]" "1.5.0" # → true
univers vers contains "vers:npm/>=1.2.0|<=2.0.0" "1.5.0" # → true
univers vers contains "vers:alpine/>=1.2.0-r5" "1.2.1-r3" # → true
```

## Documentation

- **[CONTRIBUTING.md](./CONTRIBUTING.md)** - Contribution guidelines and architecture details
- **[DEVELOPMENT.md](./DEVELOPMENT.md)** - Extended development documentation
- **Individual ecosystem documentation** - See `pkg/ecosystem/<ecosystem>/` directories for detailed examples

## Related Projects

- [aboutcode-org/univers](https://github.com/aboutcode-org/univers) - The original Python implementation
- [Package URL specification](https://github.com/package-url/purl-spec) - Standard for package identification
