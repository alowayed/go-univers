package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/semver"
)

// semverContains implements VERS constraint checking for generic SemVer ecosystem
func semverContains(constraints []string, version string) (bool, error) {
	e := &semver.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToSemverRanges converts an interval to SemVer range syntax
func intervalToSemverRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", interval.exact)}
	}

	// Exclusions are handled separately, not as semver ranges
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
		// SemVer supports both space and comma-separated, using space-separated like npm
		return []string{strings.Join(parts, " ")}
	}

	// Empty interval
	return []string{}
}
