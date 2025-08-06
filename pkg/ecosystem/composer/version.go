package composer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Composer version patterns - matches Composer version specification
var (
	// Match dev versions: dev-branch
	devVersionPattern = regexp.MustCompile(`^dev-(.+)$`)
	// Match standard semantic versions with optional stability suffixes
	// Examples: 1.2.3, 1.2.3-alpha, 1.2.3-alpha.1, 1.2.3-RC1
	semanticVersionPattern = regexp.MustCompile(`^(?:v?)(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:\.(\d+))?(?:-(alpha|beta|RC|a|b|rc)(?:\.?(\d+))?)?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)
	// Match branch names that are not semantic versions (must contain letters)
	branchNamePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_./]*$`)
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
	major         int
	minor         int
	patch         int
	extra         int    // Fourth component (for .NET-style versions)
	stability     int    // Stability level (dev, alpha, beta, RC, stable)
	stabilityNum  int    // Stability version number (e.g., alpha.2)
	build         string // Build metadata
	isDev         bool   // True for dev- prefixed versions
	devBranch     string // Branch name for dev versions
	original      string
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

		// Parse stability suffix
		if matches[5] != "" {
			stabilityStr := strings.ToLower(matches[5])
			if stability, exists := stabilityMap[stabilityStr]; exists {
				v.stability = stability
			} else {
				v.stability = stabilityStable
			}

			// Parse stability number (alpha.1, beta.2, RC.3)
			if matches[6] != "" {
				stabilityNum, err := strconv.Atoi(matches[6])
				if err != nil {
					return nil, fmt.Errorf("invalid stability number: %s", matches[6])
				}
				v.stabilityNum = stabilityNum
			}
		} else {
			v.stability = stabilityStable
		}

		// Parse build metadata
		if matches[7] != "" {
			v.build = matches[7]
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

// isBranchName checks if a string looks like a valid branch name
func isBranchName(version string) bool {
	// Only accept very common branch name patterns to avoid false positives
	commonBranches := []string{"main", "master", "develop", "trunk"}
	for _, branch := range commonBranches {
		if version == branch {
			return true
		}
	}
	
	// Accept patterns like "feature-xyz", "bugfix-123", etc.
	if strings.HasPrefix(version, "feature-") || strings.HasPrefix(version, "bugfix-") ||
		strings.HasPrefix(version, "hotfix-") || strings.HasPrefix(version, "release-") {
		return true
	}
	
	return false
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// normalize returns the normalized form of the version for comparison
func (v *Version) normalize() string {
	if v.isDev {
		if v.devBranch != "" {
			return fmt.Sprintf("dev-%s", v.devBranch)
		}
		return "dev"
	}

	result := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if v.extra > 0 {
		result += fmt.Sprintf(".%d", v.extra)
	}

	// Add stability suffix
	if v.stability != stabilityStable {
		var stabilityStr string
		switch v.stability {
		case stabilityAlpha:
			stabilityStr = "alpha"
		case stabilityBeta:
			stabilityStr = "beta"
		case stabilityRC:
			stabilityStr = "RC"
		case stabilityDev:
			stabilityStr = "dev"
		}

		if v.stabilityNum > 0 {
			result += fmt.Sprintf("-%s.%d", stabilityStr, v.stabilityNum)
		} else {
			result += fmt.Sprintf("-%s", stabilityStr)
		}
	}

	if v.build != "" {
		result += "+" + v.build
	}

	return result
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