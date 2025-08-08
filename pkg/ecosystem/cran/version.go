package cran

import (
	"fmt"
	"math/big"
	"regexp"
	"strings"
)

// versionPattern matches CRAN version strings - at least two non-negative integers separated by . or -
var versionPattern = regexp.MustCompile(`^(\d+(?:[.-]\d+)+)$`)

// Version represents a CRAN package version
type Version struct {
	components []int
	original   string
}

// NewVersion creates a new CRAN version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	// Trim whitespace
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("invalid CRAN version: empty string")
	}

	// Check basic format
	if !versionPattern.MatchString(version) {
		return nil, fmt.Errorf("invalid CRAN version: %s", original)
	}

	// Normalize dashes to periods
	normalized := strings.ReplaceAll(version, "-", ".")

	// Split by periods
	parts := strings.Split(normalized, ".")

	// Must have at least 2 components
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid CRAN version: %s (must have at least 2 components)", original)
	}

	components := make([]int, len(parts))
	for i, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid CRAN version: %s (empty component)", original)
		}

		// Parse as big int first to validate, then convert to regular int
		bigNum := new(big.Int)
		if _, ok := bigNum.SetString(part, 10); !ok {
			return nil, fmt.Errorf("invalid CRAN version: %s (non-numeric component: %s)", original, part)
		}

		// Check if it fits in int
		if !bigNum.IsInt64() {
			return nil, fmt.Errorf("invalid CRAN version: %s (component too large: %s)", original, part)
		}

		num := bigNum.Int64()
		if num < 0 {
			return nil, fmt.Errorf("invalid CRAN version: %s (negative component: %s)", original, part)
		}

		// Convert to int for easier comparison
		if num > int64(^uint(0)>>1) { // Check if fits in int
			return nil, fmt.Errorf("invalid CRAN version: %s (component too large: %s)", original, part)
		}

		components[i] = int(num)
	}

	return &Version{
		components: components,
		original:   original,
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another CRAN version
func (v *Version) Compare(other *Version) int {
	// Compare components sequentially
	minLen := len(v.components)
	if len(other.components) < minLen {
		minLen = len(other.components)
	}

	// Compare common components
	for i := 0; i < minLen; i++ {
		if v.components[i] != other.components[i] {
			return compareInt(v.components[i], other.components[i])
		}
	}

	// If all common components are equal, longer version is greater
	return compareInt(len(v.components), len(other.components))
}

// compareInt returns -1 if a < b, 0 if a == b, 1 if a > b
func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
