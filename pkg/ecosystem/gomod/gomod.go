package gomod

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Version represents a Go module version following semantic versioning with Go-specific extensions
type Version struct {
	major      int
	minor      int
	patch      int
	prerelease string
	build      string
	pseudo     *PseudoVersion
	original   string
}

// PseudoVersion represents a Go pseudo-version
type PseudoVersion struct {
	baseVersion string
	timestamp   time.Time
	revision    string
}

// VersionRange represents a Go module version range
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single version constraint
type constraint struct {
	operator string
	version  string
}

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

// NewVersion creates a new Go module version from a string
func NewVersion(version string) (*Version, error) {
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
			pseudo:   &pseudo.PseudoVersion,
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
	PseudoVersion
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
			PseudoVersion
		}{
			major: major,
			minor: 0,
			patch: 0,
			PseudoVersion: PseudoVersion{
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
			PseudoVersion
		}{
			major: major,
			minor: minor,
			patch: patch,
			PseudoVersion: PseudoVersion{
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
			PseudoVersion
		}{
			major: major,
			minor: minor,
			patch: patch,
			PseudoVersion: PseudoVersion{
				baseVersion: fmt.Sprintf("v%d.%d.%d", major, minor, patch-1),
				timestamp:   timestamp,
				revision:    matches[5],
			},
		}, nil
	}
	
	return nil, fmt.Errorf("not a pseudo-version")
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Normalize returns the normalized form of the version
func (v *Version) Normalize() string {
	if v.pseudo != nil {
		return v.normalizePseudo()
	}
	
	result := fmt.Sprintf("v%d.%d.%d", v.major, v.minor, v.patch)
	if v.prerelease != "" {
		result += "-" + v.prerelease
	}
	if v.build != "" {
		result += "+" + v.build
	}
	return result
}

// normalizePseudo returns the normalized form of a pseudo-version
func (v *Version) normalizePseudo() string {
	if v.pseudo == nil {
		return v.Normalize()
	}
	
	timestamp := v.pseudo.timestamp.Format("20060102150405")
	
	// Determine which pattern this pseudo-version follows
	if v.minor == 0 && v.patch == 0 && !strings.Contains(v.pseudo.baseVersion, "-") {
		// Pattern 1: vX.0.0-yyyymmddhhmmss-abcdefabcdef
		return fmt.Sprintf("v%d.0.0-%s-%s", v.major, timestamp, v.pseudo.revision)
	} else if strings.Contains(v.pseudo.baseVersion, "-") {
		// Pattern 2: vX.Y.Z-pre.0.yyyymmddhhmmss-abcdefabcdef
		parts := strings.Split(v.pseudo.baseVersion, "-")
		pre := strings.Join(parts[1:], "-")
		return fmt.Sprintf("v%d.%d.%d-%s.0.%s-%s", v.major, v.minor, v.patch, pre, timestamp, v.pseudo.revision)
	} else {
		// Pattern 3: vX.Y.(Z+1)-0.yyyymmddhhmmss-abcdefabcdef
		return fmt.Sprintf("v%d.%d.%d-0.%s-%s", v.major, v.minor, v.patch, timestamp, v.pseudo.revision)
	}
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

// isPseudo returns true if this is a pseudo-version
func (v *Version) isPseudo() bool {
	return v.pseudo != nil
}

// NewVersionRange creates a new Go module version range from a range string
func NewVersionRange(rangeStr string) (*VersionRange, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}
	
	// For Go modules, ranges are typically simple comparisons
	// Common patterns: >=v1.2.3, >v1.2.3, <v2.0.0, <=v1.9.9, v1.2.3
	constraints, err := parseGoRange(rangeStr)
	if err != nil {
		return nil, err
	}
	
	return &VersionRange{
		constraints: constraints,
		original:    rangeStr,
	}, nil
}

// parseGoRange parses a Go module version range
func parseGoRange(rangeStr string) ([]*constraint, error) {
	// Handle space-separated constraints
	if strings.Contains(rangeStr, " ") {
		parts := strings.Fields(rangeStr)
		var constraints []*constraint
		for _, part := range parts {
			partConstraints, err := parseSingleGoConstraint(part)
			if err != nil {
				return nil, err
			}
			constraints = append(constraints, partConstraints...)
		}
		return constraints, nil
	}
	
	return parseSingleGoConstraint(rangeStr)
}

// parseSingleGoConstraint parses a single Go version constraint
func parseSingleGoConstraint(c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)
	
	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			version := strings.TrimSpace(c[len(op):])
			return []*constraint{{operator: op, version: version}}, nil
		}
	}
	
	// Default to exact match
	return []*constraint{{operator: "=", version: c}}, nil
}

// String returns the string representation of the range
func (gr *VersionRange) String() string {
	return gr.original
}

// Contains checks if a version is within this range
func (gr *VersionRange) Contains(version *Version) bool {
	for _, constraint := range gr.constraints {
		if !constraint.matches(version) {
			return false
		}
	}
	return true
}

// matches checks if a version matches this constraint
func (c *constraint) matches(version *Version) bool {
	constraintVersion, err := NewVersion(c.version)
	if err != nil {
		return false
	}
	
	switch c.operator {
	case "=", "==":
		return version.Compare(constraintVersion) == 0
	case "!=":
		return version.Compare(constraintVersion) != 0
	case ">":
		return version.Compare(constraintVersion) > 0
	case ">=":
		return version.Compare(constraintVersion) >= 0
	case "<":
		return version.Compare(constraintVersion) < 0
	case "<=":
		return version.Compare(constraintVersion) <= 0
	default:
		return false
	}
}

// Helper functions
func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

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