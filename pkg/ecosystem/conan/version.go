package conan

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Package-level compiled regular expressions for performance
var (
	// versionPattern matches Conan version strings with extended semver format
	// Supports: MAJOR[.MINOR[.PATCH[.EXTRA...]]][-prerelease][+build]
	// Examples: 1, 1.2, 1.2.3, 1.2.3.a, 1.2.3-alpha, 1.2.3+build, 1.2.3.4.5
	versionPattern = regexp.MustCompile(`^([0-9a-z]+(?:\.[0-9a-z]+)*)(?:-([0-9a-z\-]+(?:\.[0-9a-z\-]+)*))?(?:\+([0-9a-z\-]+(?:\.[0-9a-z\-]+)*))?$`)

	// Patterns for validation
	versionPartPattern    = regexp.MustCompile(`^[0-9a-z]+$`)
	prereleasePartPattern = regexp.MustCompile(`^[0-9a-z\-]+$`)
	numericPattern        = regexp.MustCompile(`^[0-9]+$`)
)

// Version represents a Conan version
type Version struct {
	parts      []string // Main version parts (e.g., ["1", "2", "3", "a"])
	prerelease string   // Prerelease identifier (optional)
	build      string   // Build metadata (optional)
	original   string   // Original version string
}

// NewVersion creates a new Conan version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(strings.ToLower(version))

	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid conan version: %s", original)
	}

	// Parse main version parts
	mainVersion := matches[1]
	parts := strings.Split(mainVersion, ".")

	// Validate each part
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("empty version part in: %s", original)
		}
		if !versionPartPattern.MatchString(part) {
			return nil, fmt.Errorf("invalid version part '%s' in: %s", part, original)
		}
	}

	prerelease := matches[2]
	build := matches[3]

	// Validate prerelease identifiers
	if prerelease != "" {
		if err := validateIdentifiers(prerelease, "prerelease"); err != nil {
			return nil, fmt.Errorf("invalid prerelease in %s: %v", original, err)
		}
	}

	// Validate build metadata identifiers
	if build != "" {
		if err := validateIdentifiers(build, "build metadata"); err != nil {
			return nil, fmt.Errorf("invalid build metadata in %s: %v", original, err)
		}
	}

	return &Version{
		parts:      parts,
		prerelease: prerelease,
		build:      build,
		original:   original,
	}, nil
}

// validateIdentifiers validates prerelease or build metadata identifiers
func validateIdentifiers(identifiers, identifierType string) error {
	parts := strings.Split(identifiers, ".")
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("empty %s identifier", identifierType)
		}
		if !prereleasePartPattern.MatchString(part) {
			return fmt.Errorf("invalid characters in %s identifier: %s", identifierType, part)
		}
	}
	return nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another Conan version
// Returns -1 if this < other, 0 if this == other, 1 if this > other
func (v *Version) Compare(other *Version) int {
	// Compare main version parts first
	if result := compareVersionParts(v.parts, other.parts); result != 0 {
		return result
	}

	// Compare prerelease according to Conan rules
	// Build metadata is ignored for version precedence (like SemVer)
	return comparePrerelease(v.prerelease, other.prerelease)
}

// compareVersionParts compares version parts following Conan rules
func compareVersionParts(a, b []string) int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	for i := 0; i < maxLen; i++ {
		var aPart, bPart string
		if i < len(a) {
			aPart = a[i]
		} else {
			aPart = "0" // Missing parts are treated as 0
		}
		if i < len(b) {
			bPart = b[i]
		} else {
			bPart = "0" // Missing parts are treated as 0
		}

		// Compare parts using natural ordering
		if aPart != bPart {
			return naturalCompare(aPart, bPart)
		}
	}

	return 0
}

// comparePrerelease compares prerelease versions according to Conan rules
func comparePrerelease(a, b string) int {
	// No prerelease has higher precedence than prerelease
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1 // Normal version > prerelease version
	}
	if b == "" {
		return -1 // Prerelease version < normal version
	}

	// Both have prerelease, compare dot-separated identifiers
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	// Compare identifiers from left to right
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aPart, bPart string
		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}

		// A larger set of pre-release fields has a higher precedence than a smaller set
		if aPart == "" && bPart != "" {
			return -1
		}
		if aPart != "" && bPart == "" {
			return 1
		}

		// Both parts exist, compare them
		aIsNum := numericPattern.MatchString(aPart)
		bIsNum := numericPattern.MatchString(bPart)

		if aIsNum && bIsNum {
			// Both are numeric, compare numerically
			aNum, _ := strconv.Atoi(aPart)
			bNum, _ := strconv.Atoi(bPart)
			if aNum != bNum {
				return compareInt(aNum, bNum)
			}
		} else if aIsNum {
			// Numeric identifiers always have lower precedence than non-numeric
			return -1
		} else if bIsNum {
			// Non-numeric identifiers always have higher precedence than numeric
			return 1
		} else {
			// Both are non-numeric, compare lexically
			if aPart != bPart {
				if aPart < bPart {
					return -1
				}
				return 1
			}
		}
	}

	return 0
}

// naturalCompare compares two version parts using natural ordering
// Handles mixed alphanumeric parts correctly (e.g., "2" < "10" < "10a" < "3a")
func naturalCompare(a, b string) int {
	aIsNum := numericPattern.MatchString(a)
	bIsNum := numericPattern.MatchString(b)

	if aIsNum && bIsNum {
		// Both are numeric, compare numerically
		aNum, _ := strconv.Atoi(a)
		bNum, _ := strconv.Atoi(b)
		return compareInt(aNum, bNum)
	} else if aIsNum && !bIsNum {
		// Pure numeric vs alphanumeric: need to check if b starts with a number
		if bNum := extractLeadingNumber(b); bNum != "" {
			aNum, _ := strconv.Atoi(a)
			bNumVal, _ := strconv.Atoi(bNum)
			if aNum != bNumVal {
				return compareInt(aNum, bNumVal)
			}
			// Same numeric part, pure number comes before suffixed version
			return -1
		}
		// Pure numeric comes before pure alphabetic
		return -1
	} else if !aIsNum && bIsNum {
		// Alphanumeric vs pure numeric: check if a starts with a number
		if aNum := extractLeadingNumber(a); aNum != "" {
			aNumVal, _ := strconv.Atoi(aNum)
			bNum, _ := strconv.Atoi(b)
			if aNumVal != bNum {
				return compareInt(aNumVal, bNum)
			}
			// Same numeric part, suffixed version comes after pure number
			return 1
		}
		// Pure alphabetic comes after pure numeric
		return 1
	} else {
		// Both are alphanumeric, compare lexically
		if a < b {
			return -1
		}
		return 1
	}
}

// extractLeadingNumber extracts the leading numeric portion of a string
func extractLeadingNumber(s string) string {
	for i, r := range s {
		if r < '0' || r > '9' {
			if i == 0 {
				return ""
			}
			return s[:i]
		}
	}
	// Entire string is numeric
	return s
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
