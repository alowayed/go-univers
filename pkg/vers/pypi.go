package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
)

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
	isPrerelease := hasPrerelease(v) || hasDev(v)

	// If it's a prerelease, check if any constraint explicitly includes prereleases
	if isPrerelease && !constraintsIncludePrerelease(constraints) {
		return false, nil
	}

	return contains(e, constraints, version)
}

// constraintsIncludePrerelease checks if any constraint explicitly includes prerelease versions
func constraintsIncludePrerelease(constraints []string) bool {
	for _, constraint := range constraints {
		// If constraint contains prerelease markers (a, b, alpha, beta, rc, dev),
		// then prereleases are explicitly allowed
		if strings.Contains(constraint, "a") || strings.Contains(constraint, "b") ||
			strings.Contains(constraint, "alpha") || strings.Contains(constraint, "beta") ||
			strings.Contains(constraint, "rc") || strings.Contains(constraint, "dev") {
			return true
		}
	}
	return false
}

// hasPrerelease checks if a PyPI version has prerelease components
func hasPrerelease(v *pypi.Version) bool {
	// Since we can't access private fields directly, check the string representation
	// But be more careful to avoid false positives from local versions
	vStr := v.String()

	// Split on '+' to isolate the main version from local version identifier
	parts := strings.Split(vStr, "+")
	mainVersion := parts[0]

	// Check for prerelease markers in the main version only
	return strings.Contains(mainVersion, "a") || strings.Contains(mainVersion, "b") ||
		strings.Contains(mainVersion, "alpha") || strings.Contains(mainVersion, "beta") ||
		strings.Contains(mainVersion, "rc")
}

// hasDev checks if a PyPI version has dev components
func hasDev(v *pypi.Version) bool {
	// Check via string representation since we can't access private fields
	// Split on '+' to isolate the main version from local version identifier
	vStr := v.String()
	parts := strings.Split(vStr, "+")
	mainVersion := parts[0]

	return strings.Contains(mainVersion, "dev")
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
