package github

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	original    string
	prefix      string // v, release-, etc.
	major       int
	minor       int
	patch       int
	qualifier   string // alpha, beta, rc, etc.
	number      int    // for -beta.1, -rc.2, etc.
	isDateBased bool   // for date-based versions like 2024.01.15
}

var (
	// GitHub version pattern supports common GitHub repository formats:
	// - Semantic with v prefix: v1.0.0, v2.3.1-beta
	// - Semantic without prefix: 1.0.0, 2.3.1-alpha
	// - Date-based versions: 2024.01.15, v1.20240115.1.0
	// - Custom release patterns: release-1.5, 1.0.0-rc.1
	githubVersionPattern = regexp.MustCompile(`^(v|release-|rel-)?(\d+)\.(\d+)\.(\d+)(?:[-.]([A-Za-z]+)\.?(\d*))?$`)

	// Date-based pattern for versions like 2024.01.15 or v1.20240115.1.0
	githubDatePattern = regexp.MustCompile(`^(v)?(\d{4})\.(\d{1,2})\.(\d{1,2})$`)
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

	// Try date-based pattern first
	if matches := githubDatePattern.FindStringSubmatch(trimmed); matches != nil {
		return parseDateBasedVersion(trimmed, matches)
	}

	// Try semantic version pattern
	matches := githubVersionPattern.FindStringSubmatch(trimmed)
	if matches == nil {
		return nil, fmt.Errorf("invalid GitHub version format: %s", trimmed)
	}

	return parseSemanticVersion(trimmed, matches)
}

func parseDateBasedVersion(original string, matches []string) (*Version, error) {
	prefix := matches[1]
	year, _ := strconv.Atoi(matches[2])
	month, _ := strconv.Atoi(matches[3])
	day, _ := strconv.Atoi(matches[4])

	// Validate date components
	if month < 1 || month > 12 {
		return nil, fmt.Errorf("invalid month in date-based version: %d", month)
	}
	if day < 1 || day > 31 {
		return nil, fmt.Errorf("invalid day in date-based version: %d", day)
	}

	return &Version{
		original:    original,
		prefix:      prefix,
		major:       year,
		minor:       month,
		patch:       day,
		isDateBased: true,
	}, nil
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

	// Parse optional qualifier (alpha, beta, rc, etc.)
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
		original:    original,
		prefix:      prefix,
		major:       major,
		minor:       minor,
		patch:       patch,
		qualifier:   qualifier,
		number:      number,
		isDateBased: false,
	}, nil
}

func (v *Version) Compare(other *Version) int {
	// Handle mixed date-based and semantic comparison
	if v.isDateBased != other.isDateBased {
		// For mixed types, semantic versions are considered "newer" format
		// so they come after date-based versions
		if v.isDateBased {
			return -1 // date-based < semantic
		}
		return 1 // semantic > date-based
	}

	// Both are same type, compare major.minor.patch
	if v.major != other.major {
		return compareInt(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInt(v.minor, other.minor)
	}
	if v.patch != other.patch {
		return compareInt(v.patch, other.patch)
	}

	// For date-based versions, if major.minor.patch are equal, they're equal
	if v.isDateBased {
		return 0
	}

	// For semantic versions, compare qualifiers
	return compareQualifiers(v.qualifier, v.number, other.qualifier, other.number)
}

func (v *Version) String() string {
	return v.original
}

// compareQualifiers compares GitHub version qualifiers following common precedence
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

// getQualifierPrecedence returns precedence value for GitHub qualifiers
// Lower values have lower precedence (come first in ordering)
func getQualifierPrecedence(qualifier string) int {
	switch qualifier {
	case "dev":
		return 0
	case "alpha":
		return 1
	case "beta":
		return 2
	case "rc":
		return 3
	case "snapshot":
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
