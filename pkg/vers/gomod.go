package vers

import (
	"fmt"

	"github.com/alowayed/go-univers/pkg/ecosystem/gomod"
)

// gomodContains implements VERS constraint checking for Go modules ecosystem using 'golang' scheme
func gomodContains(constraints []string, version string) (bool, error) {
	e := &gomod.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToGomodRanges converts an interval to Go module range syntax
func intervalToGomodRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", ensureVPrefix(interval.exact))}
	}

	// Exclusions are handled separately, not as Go module ranges
	if interval.exclude != "" {
		return []string{} // Return empty - excludes handled in contains function
	}

	// Handle regular intervals with bounds
	if interval.lower != "" && interval.upper != "" {
		// Both bounds: create a single space-separated constraint (AND logic)
		var lowerConstraint, upperConstraint string

		if interval.lowerInclusive {
			lowerConstraint = fmt.Sprintf(">=%s", ensureVPrefix(interval.lower))
		} else {
			lowerConstraint = fmt.Sprintf(">%s", ensureVPrefix(interval.lower))
		}

		if interval.upperInclusive {
			upperConstraint = fmt.Sprintf("<=%s", ensureVPrefix(interval.upper))
		} else {
			upperConstraint = fmt.Sprintf("<%s", ensureVPrefix(interval.upper))
		}

		return []string{fmt.Sprintf("%s %s", lowerConstraint, upperConstraint)}
	} else if interval.lower != "" {
		// Only lower bound
		if interval.lowerInclusive {
			return []string{fmt.Sprintf(">=%s", ensureVPrefix(interval.lower))}
		} else {
			return []string{fmt.Sprintf(">%s", ensureVPrefix(interval.lower))}
		}
	} else if interval.upper != "" {
		// Only upper bound
		if interval.upperInclusive {
			return []string{fmt.Sprintf("<=%s", ensureVPrefix(interval.upper))}
		} else {
			return []string{fmt.Sprintf("<%s", ensureVPrefix(interval.upper))}
		}
	}

	return []string{}
}

// ensureVPrefix ensures a version string has the v prefix required for Go modules
func ensureVPrefix(version string) string {
	if version == "" {
		return version
	}
	if version[0] != 'v' {
		return "v" + version
	}
	return version
}
