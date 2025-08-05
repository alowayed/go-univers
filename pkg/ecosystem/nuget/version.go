package nuget

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches NuGet version strings following SemVer 2.0 with .NET extensions
var versionPattern = regexp.MustCompile(`^v?(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// Version represents a NuGet package version following SemVer 2.0 with .NET extensions
type Version struct {
	major      int
	minor      int
	patch      int
	revision   int    // .NET-specific 4th component
	prerelease string
	build      string
	original   string
}

// NewVersion creates a new NuGet version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	// Trim whitespace first
	version = strings.TrimSpace(version)
	// Remove leading v
	version = strings.TrimPrefix(version, "v")

	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid NuGet version: %s", original)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor := 0
	if matches[2] != "" {
		minor, err = strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version: %s", matches[2])
		}
	}

	patch := 0
	if matches[3] != "" {
		patch, err = strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid patch version: %s", matches[3])
		}
	}

	revision := 0
	if matches[4] != "" {
		revision, err = strconv.Atoi(matches[4])
		if err != nil {
			return nil, fmt.Errorf("invalid revision version: %s", matches[4])
		}
	}

	return &Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		revision:   revision,
		prerelease: matches[5],
		build:      matches[6],
		original:   strings.TrimSpace(original),
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}


// Compare compares this version with another NuGet version
func (v *Version) Compare(other *Version) int {
	// Compare major.minor.patch.revision
	if v.major != other.major {
		return compareInt(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInt(v.minor, other.minor)
	}
	if v.patch != other.patch {
		return compareInt(v.patch, other.patch)
	}
	if v.revision != other.revision {
		return compareInt(v.revision, other.revision)
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