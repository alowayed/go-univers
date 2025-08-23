package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
)

// pypiContains implements VERS constraint checking for PyPI ecosystem
func pypiContains(constraints []string, version string) (bool, error) {
	e := &pypi.Ecosystem{}
	return contains(e, constraints, version)
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
