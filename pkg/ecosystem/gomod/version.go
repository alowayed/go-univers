package gomod

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Regular expressions for Go version parsing
var (
	// Standard semantic version: v1.2.3, v1.2.3-beta, v1.2.3+build
	semverPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

	// Pseudo-version patterns
	// vX.0.0-yyyymmddhhmmss-abcdefabcdef (no base version)
	pseudoPattern1 = regexp.MustCompile(`^v(\d+)\.0\.0-(\d{14})-([a-f0-9]{12})$`)

	// vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef (pre-release base)
	pseudoPattern2 = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)-([^.]+)\.0\.(\d{14})-([a-f0-9]{12})$`)

	// vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef (release base)
	pseudoPattern3 = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)-0\.(\d{14})-([a-f0-9]{12})$`)
)

// Version represents a Go module version following semantic versioning with Go-specific extensions
type Version struct {
	major      int
	minor      int
	patch      int
	prerelease string
	build      string
	pseudo     *pseudoVersion
	original   string
}

// pseudoVersion represents a Go pseudo-version
type pseudoVersion struct {
	baseVersion string
	timestamp   time.Time
	revision    string
}

// NewVersion creates a new Go module version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Ensure version starts with 'v'
	if !strings.HasPrefix(version, "v") {
		version = "v" + version
	}

	// Try to parse as pseudo-version first
	if pseudo, err := parsePseudoVersion(version); err == nil {
		return &Version{
			major:    pseudo.major,
			minor:    pseudo.minor,
			patch:    pseudo.patch,
			pseudo:   &pseudo.pseudoVersion,
			original: original,
		}, nil
	}

	// Try to parse as standard semantic version
	matches := semverPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid Go module version: %s", original)
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
		original:   original,
	}, nil
}

// parsePseudoVersion attempts to parse a pseudo-version
func parsePseudoVersion(version string) (*struct {
	major, minor, patch int
	pseudoVersion
}, error) {
	// Pattern 1: vX.0.0-yyyymmddhhmmss-abcdefabcdef
	if matches := pseudoPattern1.FindStringSubmatch(version); matches != nil {
		major, _ := strconv.Atoi(matches[1])
		timestamp, err := time.Parse("20060102150405", matches[2])
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp in pseudo-version: %s", matches[2])
		}

		return &struct {
			major, minor, patch int
			pseudoVersion
		}{
			major: major,
			minor: 0,
			patch: 0,
			pseudoVersion: pseudoVersion{
				baseVersion: fmt.Sprintf("v%d.0.0", major),
				timestamp:   timestamp,
				revision:    matches[3],
			},
		}, nil
	}

	// Pattern 2: vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef
	if matches := pseudoPattern2.FindStringSubmatch(version); matches != nil {
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])
		timestamp, err := time.Parse("20060102150405", matches[5])
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp in pseudo-version: %s", matches[5])
		}

		return &struct {
			major, minor, patch int
			pseudoVersion
		}{
			major: major,
			minor: minor,
			patch: patch,
			pseudoVersion: pseudoVersion{
				baseVersion: fmt.Sprintf("v%d.%d.%d-%s", major, minor, patch, matches[4]),
				timestamp:   timestamp,
				revision:    matches[6],
			},
		}, nil
	}

	// Pattern 3: vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef
	if matches := pseudoPattern3.FindStringSubmatch(version); matches != nil {
		major, _ := strconv.Atoi(matches[1])
		minor, _ := strconv.Atoi(matches[2])
		patch, _ := strconv.Atoi(matches[3])
		timestamp, err := time.Parse("20060102150405", matches[4])
		if err != nil {
			return nil, fmt.Errorf("invalid timestamp in pseudo-version: %s", matches[4])
		}

		return &struct {
			major, minor, patch int
			pseudoVersion
		}{
			major: major,
			minor: minor,
			patch: patch,
			pseudoVersion: pseudoVersion{
				baseVersion: fmt.Sprintf("v%d.%d.%d", major, minor, patch-1),
				timestamp:   timestamp,
				revision:    matches[5],
			},
		}, nil
	}

	return nil, fmt.Errorf("not a pseudo-version")
}

// Compare compares this version with another Go module version
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

	// Handle pseudo-version comparison
	if v.pseudo != nil && other.pseudo != nil {
		return v.pseudo.timestamp.Compare(other.pseudo.timestamp)
	}
	if v.pseudo != nil && other.pseudo == nil {
		// Pseudo-versions are pre-release, so they come before releases
		if other.prerelease == "" {
			return -1
		}
		// Compare with prerelease
		return comparePrerelease("pseudo", other.prerelease)
	}
	if v.pseudo == nil && other.pseudo != nil {
		if v.prerelease == "" {
			return 1
		}
		return comparePrerelease(v.prerelease, "pseudo")
	}

	// Compare prerelease according to semver rules
	return comparePrerelease(v.prerelease, other.prerelease)
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
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

// comparePrerelease returns -1, 0, or 1 comparing prereleases where empty string (release) > any prerelease
func comparePrerelease(a, b string) int {
	// No prerelease (release) has higher precedence than prerelease
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}

	// Special handling for pseudo-versions
	if a == "pseudo" && b != "pseudo" {
		return -1
	}
	if a != "pseudo" && b == "pseudo" {
		return 1
	}
	if a == "pseudo" && b == "pseudo" {
		return 0
	}

	// Lexicographic comparison for prereleases
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

