package composer

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Composer version patterns - matches Composer version specification
var (
	// Match dev versions: dev-branch
	devVersionPattern = regexp.MustCompile(`^dev-(.+)$`)
	// Match standard semantic versions with optional stability suffixes
	// Examples: 1.2.3, 1.2.3-alpha, 1.2.3-alpha.1, 1.2.3-RC1, 1.0a1, 1.0pl1
	// Capture groups: (1)major (2)minor (3)patch (4)extra (5)extra2 (6)stability1 (7)stabilityNum1 (8)stability2 (9)stabilityNum2 (10)build
	semanticVersionPattern = regexp.MustCompile(`^(?:v?)(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?(?:(?:-(alpha|beta|RC|a|b|rc|dev|patch)(?:\.?(\d+))?)|(?:(alpha|beta|RC|a|b|rc|dev|pl)(\d+)?))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
)

// Stability levels ordered from least to most stable
const (
	stabilityDev    = 0
	stabilityAlpha  = 1
	stabilityBeta   = 2
	stabilityRC     = 3
	stabilityStable = 4
)

// stabilityMap maps stability strings to numeric values for comparison
var stabilityMap = map[string]int{
	"dev":   stabilityDev,
	"alpha": stabilityAlpha,
	"a":     stabilityAlpha,
	"beta":  stabilityBeta,
	"b":     stabilityBeta,
	"RC":    stabilityRC,
	"rc":    stabilityRC,
}

// Version represents a Composer package version following Composer specification
type Version struct {
	major        int
	minor        int
	patch        int
	extra        int    // Fourth component (for .NET-style versions)
	stability    int    // Stability level (dev, alpha, beta, RC, stable)
	stabilityNum int    // Stability version number (e.g., alpha.2)
	build        string // Build metadata
	isDev        bool   // True for dev- prefixed versions
	devBranch    string // Branch name for dev versions
	original     string
}

// NewVersion creates a new Composer version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("invalid Composer version: %s", original)
	}

	v := &Version{
		original: original,
	}

	// Check if this is a dev version (dev-main, dev-feature-branch)
	if matches := devVersionPattern.FindStringSubmatch(version); matches != nil {
		v.isDev = true
		v.stability = stabilityDev
		v.devBranch = matches[1]
		return v, nil
	}

	// Try semantic version pattern
	if matches := semanticVersionPattern.FindStringSubmatch(version); matches != nil {
		// Parse numeric version components
		if matches[1] != "" {
			major, err := strconv.Atoi(matches[1])
			if err != nil {
				return nil, fmt.Errorf("invalid major version: %s", matches[1])
			}
			v.major = major
		}

		if matches[2] != "" {
			minor, err := strconv.Atoi(matches[2])
			if err != nil {
				return nil, fmt.Errorf("invalid minor version: %s", matches[2])
			}
			v.minor = minor
		}

		if matches[3] != "" {
			patch, err := strconv.Atoi(matches[3])
			if err != nil {
				return nil, fmt.Errorf("invalid patch version: %s", matches[3])
			}
			v.patch = patch
		}

		if matches[4] != "" {
			extra, err := strconv.Atoi(matches[4])
			if err != nil {
				return nil, fmt.Errorf("invalid extra version: %s", matches[4])
			}
			v.extra = extra
		}

		// Handle 5th component (rare but exists in some versions)
		if matches[5] != "" {
			// Just ignore additional components beyond the 4th for now
		}

		// Parse stability suffix - handle both formats
		// Format 1: -alpha.1 (matches[6] and matches[7])
		// Format 2: alpha1 (matches[8] and matches[9])
		if matches[6] != "" { // Hyphenated format: -alpha.1
			stabilityStr := strings.ToLower(matches[6])
			if stabilityStr == "patch" || stabilityStr == "pl" {
				v.stability = stabilityStable // Treat patch/pl as stable
			} else if stability, exists := stabilityMap[stabilityStr]; exists {
				v.stability = stability
			} else {
				v.stability = stabilityStable
			}

			// Parse stability number (alpha.1, beta.2, RC.3)
			if matches[7] != "" {
				stabilityNum, err := strconv.Atoi(matches[7])
				if err != nil {
					return nil, fmt.Errorf("invalid stability number: %s", matches[7])
				}
				v.stabilityNum = stabilityNum
			}
		} else if matches[8] != "" { // Direct format: alpha1
			stabilityStr := strings.ToLower(matches[8])
			if stabilityStr == "pl" {
				v.stability = stabilityStable // Treat pl as stable
			} else if stability, exists := stabilityMap[stabilityStr]; exists {
				v.stability = stability
			} else {
				v.stability = stabilityStable
			}

			// Parse stability number (alpha1, beta2, RC1)
			if matches[9] != "" {
				stabilityNum, err := strconv.Atoi(matches[9])
				if err != nil {
					return nil, fmt.Errorf("invalid stability number: %s", matches[9])
				}
				v.stabilityNum = stabilityNum
			}
		} else {
			v.stability = stabilityStable
		}

		// Parse build metadata
		if matches[10] != "" {
			v.build = matches[10]
		}

		return v, nil
	}

	// Check if this is a valid branch name (common branch patterns)
	if isBranchName(version) {
		v.isDev = true
		v.stability = stabilityDev
		v.devBranch = version
		return v, nil
	}

	return nil, fmt.Errorf("invalid Composer version: %s", original)
}

// isBranchName checks if a string looks like a valid branch name following Composer conventions
func isBranchName(version string) bool {
	// Empty strings are not valid branch names
	if version == "" {
		return false
	}

	// Don't accept strings that look too much like malformed versions
	// This prevents false positives for invalid version strings
	if strings.Contains(version, ".") && (len(strings.Split(version, ".")) >= 2) {
		// If it contains dots and looks like it could be a version, be more strict
		if !strings.Contains(version, "/") && !strings.Contains(version, "-") {
			return false
		}
	}

	// Allow most strings that look like branch names
	// Let the validation logic handle truly invalid cases

	// Common branch name patterns used in Git repositories
	commonBranches := []string{
		"main", "master", "develop", "development", "trunk",
		"stable", "staging", "production", "prod",
	}
	if slices.Contains(commonBranches, version) {
		return true
	}

	// Accept conventional Git Flow and GitHub Flow patterns
	branchPrefixes := []string{
		"feature/", "feature-", "feat/", "feat-",
		"bugfix/", "bugfix-", "fix/", "fix-",
		"hotfix/", "hotfix-", "patch/", "patch-",
		"release/", "release-", "rel/", "rel-",
		"chore/", "chore-", "docs/", "docs-", "doc/", "doc-",
		"refactor/", "refactor-", "style/", "style-",
	}

	for _, prefix := range branchPrefixes {
		if strings.HasPrefix(version, prefix) && len(version) > len(prefix) {
			return true
		}
	}

	// Accept version branches (like v1.x, 1.x-dev, etc.)
	if strings.HasSuffix(version, "-dev") && len(version) > 4 {
		return true
	}

	return false
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another Composer version following Composer rules
func (v *Version) Compare(other *Version) int {
	// Dev versions are always less than stable versions
	if v.isDev && !other.isDev {
		return -1
	}
	if !v.isDev && other.isDev {
		return 1
	}

	// Both are dev versions - compare branch names lexically
	if v.isDev && other.isDev {
		return strings.Compare(v.devBranch, other.devBranch)
	}

	// Compare major.minor.patch.extra
	if v.major != other.major {
		return compareInt(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInt(v.minor, other.minor)
	}
	if v.patch != other.patch {
		return compareInt(v.patch, other.patch)
	}
	if v.extra != other.extra {
		return compareInt(v.extra, other.extra)
	}

	// Compare stability - stable versions are higher than pre-release
	if v.stability != other.stability {
		return compareInt(v.stability, other.stability)
	}

	// Same stability level - compare stability numbers
	return compareInt(v.stabilityNum, other.stabilityNum)
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
