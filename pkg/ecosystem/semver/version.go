package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches SemVer 2.0.0 version strings
// Group 1: major, Group 2: minor, Group 3: patch
// Group 4: prerelease (optional), Group 5: build metadata (optional)
var versionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z\-]+(?:\.[0-9A-Za-z\-]+)*))?(?:\+([0-9A-Za-z\-]+(?:\.[0-9A-Za-z\-]+)*))?$`)

// Version represents a Semantic Version 2.0.0
type Version struct {
	major      int
	minor      int
	patch      int
	prerelease string
	build      string
	original   string
}

// NewVersion creates a new SemVer version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid semantic version: %s", original)
	}

	// Parse major version
	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	// Check for leading zeros (not allowed in SemVer 2.0)
	if len(matches[1]) > 1 && matches[1][0] == '0' {
		return nil, fmt.Errorf("major version cannot have leading zeros: %s", matches[1])
	}

	// Parse minor version
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	if len(matches[2]) > 1 && matches[2][0] == '0' {
		return nil, fmt.Errorf("minor version cannot have leading zeros: %s", matches[2])
	}

	// Parse patch version
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	if len(matches[3]) > 1 && matches[3][0] == '0' {
		return nil, fmt.Errorf("patch version cannot have leading zeros: %s", matches[3])
	}

	prerelease := matches[4]
	build := matches[5]

	// Validate prerelease identifiers
	if prerelease != "" {
		if err := validatePrerelease(prerelease); err != nil {
			return nil, fmt.Errorf("invalid prerelease: %v", err)
		}
	}

	// Validate build metadata identifiers
	if build != "" {
		if err := validateBuildMetadata(build); err != nil {
			return nil, fmt.Errorf("invalid build metadata: %v", err)
		}
	}

	return &Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		prerelease: prerelease,
		build:      build,
		original:   original,
	}, nil
}

// validatePrerelease validates prerelease identifiers according to SemVer 2.0
func validatePrerelease(prerelease string) error {
	parts := strings.Split(prerelease, ".")
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("empty prerelease identifier")
		}

		// Check for valid characters (alphanumerics and hyphens only)
		if !regexp.MustCompile(`^[0-9A-Za-z\-]+$`).MatchString(part) {
			return fmt.Errorf("invalid characters in prerelease identifier: %s", part)
		}

		// Numeric identifiers must not have leading zeros
		if regexp.MustCompile(`^[0-9]+$`).MatchString(part) {
			if len(part) > 1 && part[0] == '0' {
				return fmt.Errorf("numeric prerelease identifier cannot have leading zeros: %s", part)
			}
		}
	}
	return nil
}

// validateBuildMetadata validates build metadata identifiers according to SemVer 2.0
func validateBuildMetadata(build string) error {
	parts := strings.Split(build, ".")
	for _, part := range parts {
		if part == "" {
			return fmt.Errorf("empty build metadata identifier")
		}

		// Check for valid characters (alphanumerics and hyphens only)
		if !regexp.MustCompile(`^[0-9A-Za-z\-]+$`).MatchString(part) {
			return fmt.Errorf("invalid characters in build metadata identifier: %s", part)
		}
	}
	return nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another SemVer version
// Returns -1 if this < other, 0 if this == other, 1 if this > other
func (v *Version) Compare(other *Version) int {
	// Compare major.minor.patch first
	if v.major != other.major {
		return compareInt(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInt(v.minor, other.minor)
	}
	if v.patch != other.patch {
		return compareInt(v.patch, other.patch)
	}

	// Compare prerelease according to SemVer 2.0 rules
	// Build metadata is ignored for version precedence
	return comparePrerelease(v.prerelease, other.prerelease)
}

// comparePrerelease compares prerelease versions according to SemVer 2.0 specification
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
		aIsNum := regexp.MustCompile(`^[0-9]+$`).MatchString(aPart)
		bIsNum := regexp.MustCompile(`^[0-9]+$`).MatchString(bPart)

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
			// Both are non-numeric, compare lexically (ASCII sort order)
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
