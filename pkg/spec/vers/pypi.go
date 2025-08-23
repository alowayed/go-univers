package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
)

// No regex needed - we can parse the version string more directly

// pypiContains implements VERS constraint checking for PyPI ecosystem
// with PEP 440 prerelease exclusion logic
func pypiContains(constraints []string, version string) (bool, error) {
	e := &pypi.Ecosystem{}

	// Parse the version to check if it's a prerelease
	v, err := e.NewVersion(version)
	if err != nil {
		return false, err
	}

	// Check if version is a prerelease (has prerelease or dev components)
	isPrerelease := isPyPIPrerelease(v)

	// If it's a prerelease, check if any constraint explicitly includes prereleases
	if isPrerelease && !constraintsIncludePrerelease(constraints) {
		return false, nil
	}

	return contains(e, constraints, version)
}

// constraintsIncludePrerelease checks if any constraint explicitly includes prerelease versions
func constraintsIncludePrerelease(constraints []string) bool {
	for _, constraint := range constraints {
		// If constraint contains prerelease markers, then prereleases are explicitly allowed
		if containsPrereleaseMarkers(constraint) {
			return true
		}
	}
	return false
}

// isPyPIPrerelease checks if a PyPI version has prerelease or dev components
func isPyPIPrerelease(v *pypi.Version) bool {
	// Since we can't access private fields directly, check the string representation
	// But be careful to avoid false positives from local versions
	vStr := v.String()

	// Split on '+' to isolate the main version from local version identifier
	parts := strings.Split(vStr, "+")
	mainVersion := parts[0]

	// Check for prerelease markers in the main version only
	return containsPrereleaseMarkers(mainVersion)
}

// containsPrereleaseMarkers checks if a version string contains PEP 440 prerelease markers
func containsPrereleaseMarkers(versionStr string) bool {
	// PEP 440 prerelease markers can appear directly attached to version numbers
	// e.g., "1.5.0b1", "1.5.0rc1", "1.5.0.dev1"

	versionStr = strings.ToLower(versionStr)

	// Define prerelease markers in order of length (longest first to avoid partial matches)
	markers := []string{"alpha", "beta", "dev", "rc", "a", "b"}

	for _, marker := range markers {
		// Look for the marker in the version string
		if idx := strings.Index(versionStr, marker); idx >= 0 {
			// Check that marker is preceded by a digit or dot and followed by digits or end/+
			if idx > 0 && (versionStr[idx-1] >= '0' && versionStr[idx-1] <= '9' || versionStr[idx-1] == '.') {
				afterMarker := idx + len(marker)
				if afterMarker >= len(versionStr) {
					return true // marker at end
				}
				// Check what comes after the marker
				next := versionStr[afterMarker]
				if (next >= '0' && next <= '9') || next == '+' || next == '.' {
					return true
				}
			}
		}
	}

	return false
}

// intervalToPypiRanges converts an interval to PyPI range syntax
func intervalToPypiRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("==%s", interval.exact)}
	}

	// Exclusions are handled separately, not as PyPI ranges
	if interval.exclude != "" {
		return []string{} // Return empty - excludes handled in contains function
	}

	// Handle regular intervals with bounds
	var parts []string
	if interval.lower != "" {
		op := ">"
		if interval.lowerInclusive {
			op = ">="
		}
		parts = append(parts, fmt.Sprintf("%s%s", op, interval.lower))
	}
	if interval.upper != "" {
		op := "<"
		if interval.upperInclusive {
			op = "<="
		}
		parts = append(parts, fmt.Sprintf("%s%s", op, interval.upper))
	}

	if len(parts) > 0 {
		return []string{strings.Join(parts, ", ")}
	}

	// Empty interval
	return []string{}
}
