package npm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches NPM version strings
var versionPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// Version represents an NPM package version following semantic versioning
type Version struct {
	major      int
	minor      int
	patch      int
	prerelease string
	build      string
	original   string
}

// NewVersion creates a new NPM version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	// Trim whitespace first
	version = strings.TrimSpace(version)
	// Remove leading v or = (only one v)
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "=")

	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid NPM version: %s", original)
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

	return &Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		prerelease: matches[4],
		build:      matches[5],
		original:   strings.TrimSpace(original),
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// normalize returns the normalized form of the version
func (v *Version) normalize() string {
	result := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if v.prerelease != "" {
		result += "-" + v.prerelease
	}
	if v.build != "" {
		result += "+" + v.build
	}
	return result
}

// Compare compares this version with another NPM version
func (v *Version) Compare(other *Version) int {
	// Compare major.minor.patch
	if v.major != other.major {
		return compareInt(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInt(v.minor, other.minor)
	}
	if v.patch != other.patch {
		return compareInt(v.patch, other.patch)
	}

	// Compare prerelease according to semver rules
	return comparePrerelease(v.prerelease, other.prerelease)
}

func comparePrerelease(a, b string) int {
	// No prerelease has higher precedence than prerelease
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}

	// Split by dots and compare each part
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

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

		// Missing part has lower precedence
		if aPart == "" && bPart != "" {
			return -1
		}
		if aPart != "" && bPart == "" {
			return 1
		}

		// Try to parse as numbers
		aNum, aIsNum := parseNum(aPart)
		bNum, bIsNum := parseNum(bPart)

		if aIsNum && bIsNum {
			if aNum != bNum {
				return compareInt(aNum, bNum)
			}
		} else if aIsNum {
			return -1 // Numeric identifiers have lower precedence
		} else if bIsNum {
			return 1
		} else {
			// Both are strings, compare lexically
			if aPart != bPart {
				return strings.Compare(aPart, bPart)
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

// parseNum returns the integer value and true if s is a valid number, otherwise 0 and false
func parseNum(s string) (int, bool) {
	if num, err := strconv.Atoi(s); err == nil {
		return num, true
	}
	return 0, false
}
