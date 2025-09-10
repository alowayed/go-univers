package hex

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	original      string
	major         int
	minor         int
	patch         int
	preRelease    []string // pre-release identifiers
	buildMetadata string   // build metadata (ignored in comparisons)
}

var (
	// Simple Hex version pattern - much more readable and maintainable
	// Format: MAJOR.MINOR.PATCH with optional pre-release and build metadata
	// Examples: 1.0.0, 2.1.3-alpha.1, 1.0.0+build.1
	hexVersionPattern = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-([a-zA-Z0-9\-\.]+))?(?:\+([a-zA-Z0-9\-\.]+))?$`)
	// Partial version for pessimistic operator: MAJOR.MINOR
	hexPartialVersionPattern = regexp.MustCompile(`^(\d+)\.(\d+)$`)
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

	// Try full semantic version pattern first
	matches := hexVersionPattern.FindStringSubmatch(trimmed)
	if matches != nil {
		return parseSemanticVersion(trimmed, matches)
	}

	// Try partial version pattern (for pessimistic operator)
	partialMatches := hexPartialVersionPattern.FindStringSubmatch(trimmed)
	if partialMatches != nil {
		return parsePartialVersion(trimmed, partialMatches)
	}

	return nil, fmt.Errorf("invalid Hex version format: %s", trimmed)
}

func parseSemanticVersion(original string, matches []string) (*Version, error) {
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

	// Parse optional pre-release identifiers
	var preRelease []string
	if matches[4] != "" {
		preRelease = strings.Split(matches[4], ".")
		// Validate pre-release identifiers
		for _, identifier := range preRelease {
			if identifier == "" {
				return nil, fmt.Errorf("invalid pre-release identifier: empty identifier")
			}
		}
	}

	// Parse optional build metadata
	buildMetadata := ""
	if len(matches) > 5 && matches[5] != "" {
		buildMetadata = matches[5]
	}

	return &Version{
		original:      original,
		major:         major,
		minor:         minor,
		patch:         patch,
		preRelease:    preRelease,
		buildMetadata: buildMetadata,
	}, nil
}

func parsePartialVersion(original string, matches []string) (*Version, error) {
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

	// Partial versions default to patch 0
	return &Version{
		original: original,
		major:    major,
		minor:    minor,
		patch:    0,
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

	// Build metadata is ignored in comparisons (SemVer 2.0 rule)
	// Pre-release versions have lower precedence than normal versions
	return comparePreRelease(v.preRelease, other.preRelease)
}

func (v *Version) String() string {
	return v.original
}

// comparePreRelease compares pre-release versions following SemVer 2.0 rules
func comparePreRelease(pr1, pr2 []string) int {
	// No pre-release (release) is higher than any pre-release
	if len(pr1) == 0 && len(pr2) == 0 {
		return 0
	}
	if len(pr1) == 0 {
		return 1 // Release > any pre-release
	}
	if len(pr2) == 0 {
		return -1 // Any pre-release < release
	}

	// Both have pre-release, compare identifiers one by one
	maxLen := len(pr1)
	if len(pr2) > maxLen {
		maxLen = len(pr2)
	}

	for i := 0; i < maxLen; i++ {
		if i >= len(pr1) {
			return -1 // pr1 has fewer identifiers, so it's less
		}
		if i >= len(pr2) {
			return 1 // pr2 has fewer identifiers, so it's less
		}

		cmp := comparePreReleaseIdentifier(pr1[i], pr2[i])
		if cmp != 0 {
			return cmp
		}
	}

	return 0
}

// comparePreReleaseIdentifier compares individual pre-release identifiers
func comparePreReleaseIdentifier(id1, id2 string) int {
	// Try to parse as integers first
	num1, err1 := strconv.Atoi(id1)
	num2, err2 := strconv.Atoi(id2)

	if err1 == nil && err2 == nil {
		// Both are numbers, compare numerically
		return compareInt(num1, num2)
	}
	if err1 == nil {
		// id1 is number, id2 is not: number < string
		return -1
	}
	if err2 == nil {
		// id2 is number, id1 is not: string > number
		return 1
	}

	// Both are strings, compare lexicographically
	if id1 < id2 {
		return -1
	}
	if id1 > id2 {
		return 1
	}
	return 0
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
