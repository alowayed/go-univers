package apache

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	original  string
	major     int
	minor     int
	patch     int
	qualifier string
	number    int // for RC1, RC2, etc.
}

var (
	// Apache version pattern supports common Apache project formats:
	// - Basic semantic: 2.4.41, 9.0.45, 8.5.75
	// - Release candidates: 2.4.41-RC1, 9.0.0-RC2
	// - Beta/Alpha releases: 2.4.0-beta, 3.0.0-alpha
	// - Development versions: 2.5.0-dev
	apacheVersionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([A-Za-z]+)(\d*))?$`)
)

func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	if version == "" {
		return nil, fmt.Errorf("version string cannot be empty")
	}

	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return nil, fmt.Errorf("version string cannot be empty or only whitespace")
	}

	matches := apacheVersionPattern.FindStringSubmatch(trimmed)
	if matches == nil {
		return nil, fmt.Errorf("invalid Apache version format: %s", trimmed)
	}

	// Parse major version
	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	// Parse minor version
	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	// Parse patch version
	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	// Parse optional qualifier (RC, beta, alpha, dev, etc.)
	qualifier := ""
	number := 0
	if matches[4] != "" {
		qualifier = strings.ToLower(matches[4])
		if matches[5] != "" {
			number, err = strconv.Atoi(matches[5])
			if err != nil {
				return nil, fmt.Errorf("invalid qualifier number: %s", matches[5])
			}
		}
	}

	return &Version{
		original:  version,
		major:     major,
		minor:     minor,
		patch:     patch,
		qualifier: qualifier,
		number:    number,
	}, nil
}

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

	// Compare qualifiers
	return compareQualifiers(v.qualifier, v.number, other.qualifier, other.number)
}

func (v *Version) String() string {
	return v.original
}

// compareQualifiers compares Apache version qualifiers
func compareQualifiers(q1 string, n1 int, q2 string, n2 int) int {
	// No qualifier (release) is higher than any qualifier
	if q1 == "" && q2 == "" {
		return 0
	}
	if q1 == "" {
		return 1 // Release > any qualifier
	}
	if q2 == "" {
		return -1 // Any qualifier < release
	}

	// Both have qualifiers - compare by precedence
	p1 := getQualifierPrecedence(q1)
	p2 := getQualifierPrecedence(q2)

	if p1 != p2 {
		return compareInt(p1, p2)
	}

	// Same qualifier type, compare numbers
	return compareInt(n1, n2)
}

// getQualifierPrecedence returns precedence value for Apache qualifiers
// Lower values have lower precedence (come first in ordering)
func getQualifierPrecedence(qualifier string) int {
	switch qualifier {
	case "alpha":
		return 1
	case "beta":
		return 2
	case "rc":
		return 3
	case "dev":
		return 4
	default:
		return 99 // Unknown qualifiers come last
	}
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
