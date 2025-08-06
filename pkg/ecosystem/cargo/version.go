package cargo

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches Cargo version strings following SemVer 2.0 specification
// Cargo strictly follows SemVer 2.0: MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]
var versionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// Version represents a Cargo (Rust) package version following SemVer 2.0
type Version struct {
	major      int
	minor      int
	patch      int
	prerelease string
	build      string
	original   string
}

// NewVersion creates a new Cargo version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	// Trim whitespace first
	version = strings.TrimSpace(version)
	
	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid Cargo version: %s", original)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	prerelease := matches[4] // can be empty
	build := matches[5]      // can be empty

	return &Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		prerelease: prerelease,
		build:      build,
		original:   original,
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another Cargo version following SemVer 2.0 rules
func (v *Version) Compare(other *Version) int {
	// 1. Compare major.minor.patch numerically
	if cmp := compareInt(v.major, other.major); cmp != 0 {
		return cmp
	}
	if cmp := compareInt(v.minor, other.minor); cmp != 0 {
		return cmp
	}
	if cmp := compareInt(v.patch, other.patch); cmp != 0 {
		return cmp
	}

	// 2. Handle prerelease according to SemVer 2.0:
	// - Version without prerelease > version with prerelease
	// - Compare prerelease identifiers lexically if both have prerelease
	if v.prerelease == "" && other.prerelease == "" {
		return 0 // Both are release versions with same major.minor.patch
	}
	if v.prerelease == "" && other.prerelease != "" {
		return 1 // Release version > prerelease version
	}
	if v.prerelease != "" && other.prerelease == "" {
		return -1 // Prerelease version < release version
	}

	// Both have prerelease, compare them according to SemVer 2.0 rules
	return comparePrereleaseIdentifiers(v.prerelease, other.prerelease)
}

// comparePrereleaseIdentifiers compares prerelease identifiers according to SemVer 2.0
func comparePrereleaseIdentifiers(a, b string) int {
	// Split by dots and compare each identifier
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := max(len(aParts), len(bParts))

	for i := 0; i < maxLen; i++ {
		var aPart, bPart string
		
		// Missing parts are considered smaller
		if i >= len(aParts) {
			return -1 // a has fewer parts, so a < b
		}
		if i >= len(bParts) {
			return 1 // b has fewer parts, so a > b
		}
		
		aPart = aParts[i]
		bPart = bParts[i]

		// Try to parse as integers
		aNum, aIsNum := tryParseInt(aPart)
		bNum, bIsNum := tryParseInt(bPart)

		if aIsNum && bIsNum {
			// Both are numeric, compare numerically
			if cmp := compareInt(aNum, bNum); cmp != 0 {
				return cmp
			}
		} else if aIsNum && !bIsNum {
			// Numeric identifiers always have lower precedence than non-numeric
			return -1
		} else if !aIsNum && bIsNum {
			// Non-numeric identifiers always have higher precedence than numeric
			return 1
		} else {
			// Both are non-numeric, compare lexically
			if cmp := strings.Compare(aPart, bPart); cmp != 0 {
				return cmp
			}
		}
	}

	return 0 // All identifiers are equal
}

// tryParseInt attempts to parse a string as an integer
func tryParseInt(s string) (int, bool) {
	num, err := strconv.Atoi(s)
	return num, err == nil
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

