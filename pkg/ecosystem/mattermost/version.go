package mattermost

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	original  string
	prefix    string // v, etc.
	major     int
	minor     int
	patch     int
	qualifier string // esr, rc, etc.
	number    int    // for rc1, rc2, etc.
}

var (
	// Mattermost version pattern supports:
	// - Semantic with v prefix: v8.1.5, v10.12.0
	// - Semantic without prefix: 8.1.5, 10.12.0
	// - ESR versions: v8.1.5-esr, 8.1.5-esr
	// - Release candidates: v8.1.0-rc1, v10.12.0-rc2
	mattermostVersionPattern = regexp.MustCompile(`^(v)?(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-(esr|rc)(\d*))?$`)
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

	// Parse Mattermost version pattern
	matches := mattermostVersionPattern.FindStringSubmatch(trimmed)
	if matches == nil {
		return nil, fmt.Errorf("invalid Mattermost version format: %s", trimmed)
	}

	return parseSemanticVersion(trimmed, matches)
}

func parseSemanticVersion(original string, matches []string) (*Version, error) {
	prefix := matches[1]

	// Parse major version
	major, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[2])
	}

	// Parse minor version
	minor, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[3])
	}

	// Parse patch version
	patch, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[4])
	}

	// Parse optional qualifier (esr, rc, etc.)
	qualifier := ""
	number := 0
	if matches[5] != "" {
		qualifier = strings.ToLower(matches[5])
		if matches[6] != "" {
			number, err = strconv.Atoi(matches[6])
			if err != nil {
				return nil, fmt.Errorf("invalid qualifier number: %s", matches[6])
			}
		}
	}

	return &Version{
		original:  original,
		prefix:    prefix,
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

// compareQualifiers compares Mattermost version qualifiers following precedence
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

// getQualifierPrecedence returns precedence value for Mattermost qualifiers
// Lower values have lower precedence (come first in ordering)
func getQualifierPrecedence(qualifier string) int {
	switch qualifier {
	case "rc":
		return 1
	case "esr":
		return 2
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
