package npm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents an NPM package version following semantic versioning
type Version struct {
	major      int
	minor      int
	patch      int
	prerelease string
	build      string
	original   string
}

// versionPattern matches NPM version strings
var versionPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// NewVersion creates a new NPM version from a string
func NewVersion(version string) (*Version, error) {
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

// isValid checks if the version is valid
func (v *Version) isValid() bool {
	return v.major >= 0 && v.minor >= 0 && v.patch >= 0
}

// Normalize returns the normalized form of the version
func (v *Version) Normalize() string {
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

// satisfies checks if this version satisfies the given constraint
func (v *Version) satisfies(constraint *constraint) bool {
	return constraint.matches(v)
}
