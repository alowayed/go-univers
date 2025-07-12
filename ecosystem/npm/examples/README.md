# NPM Version Examples

This directory contains comprehensive examples demonstrating the NPM semantic version parsing capabilities of the go-univers library.

## Examples

### `example-simple/`
Demonstrates basic NPM version operations:
- Version parsing with various formats (v-prefix, equals-prefix, prerelease, build metadata)
- Version comparison and sorting
- Simple range matching with common patterns

**Run:** `cd example-simple && go run main.go`

### `example-complex-ranges/`
Shows advanced version range operations:
- Complex range patterns (caret, tilde, X-ranges, hyphen ranges)
- OR logic with `||` operator
- Prerelease version handling in ranges
- Version sorting with complex prerelease ordering

**Run:** `cd example-complex-ranges && go run main.go`

### `example-edge-cases/`
Demonstrates edge cases and error handling:
- Version parsing edge cases (whitespace, invalid formats, zero-padding)
- Version range parsing edge cases (malformed ranges, invalid characters)
- Zero-version caret range special cases (`^0.x.y` behavior)
- Build metadata handling (ignored in comparisons)
- Prerelease boundary cases

**Run:** `cd example-edge-cases && go run main.go`

### `example-sorting/`
Demonstrates native Go sorting capabilities for NPM versions:
- Sorting version slices using `slices.SortFunc()` with the existing `Compare()` method
- Ascending and descending sort order
- Stable sorting for equal versions
- Semantic version ordering rules (normal vs prerelease, numeric vs string comparison)

**Run:** `cd example-sorting && go run main.go`

## Key Features Demonstrated

- **Version Parsing**: Handles v-prefixes, equals-prefixes, prerelease, build metadata
- **Range Parsing**: Supports caret (`^`), tilde (`~`), X-ranges, hyphen ranges, OR logic
- **Comparison**: Semantic version comparison with proper prerelease ordering
- **Sorting**: Native Go sorting using `slices.SortFunc()` with existing comparison logic
- **Edge Cases**: Zero-version handling, whitespace normalization, error handling
- **NPM Semver Compliance**: Follows NPM semantic versioning specification

## Usage Pattern

```go
import (
    "slices"
    "github.com/alowayed/go-univers/ecosystem/npm"
)

// Parse version
version, err := npm.NewVersion("1.2.3-alpha.1")

// Parse range
versionRange, err := npm.NewVersionRange("^1.0.0 || >=2.0.0-alpha")

// Check if version satisfies range
satisfies := versionRange.Contains(version)

// Sort versions
versions := []*npm.Version{v1, v2, v3}
slices.SortFunc(versions, (*npm.Version).Compare)
```

These examples showcase over 70 test cases covering comprehensive edge cases identified by comparing with Google's semver library tests.