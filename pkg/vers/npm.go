package vers

import (
	"fmt"

	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
)

// npmContains implements VERS constraint checking for NPM ecosystem
func npmContains(constraints []string, version string) (bool, error) {
	e := &npm.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToNpmRanges converts an interval to NPM range syntax
func intervalToNpmRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", interval.exact)}
	}

	// Exclusions are handled separately, not as NPM ranges
	if interval.exclude != "" {
		return []string{} // Return empty - excludes handled in contains function
	}

	// Handle regular intervals with bounds
	var constraints []string

	if interval.lower != "" && interval.upper != "" {
		// Both bounds: >=lower <=upper
		if interval.lowerInclusive {
			constraints = append(constraints, fmt.Sprintf(">=%s", interval.lower))
		} else {
			constraints = append(constraints, fmt.Sprintf(">%s", interval.lower))
		}
		if interval.upperInclusive {
			constraints = append(constraints, fmt.Sprintf("<=%s", interval.upper))
		} else {
			constraints = append(constraints, fmt.Sprintf("<%s", interval.upper))
		}
		// Return as space-separated constraint for NPM
		return []string{fmt.Sprintf("%s %s", constraints[0], constraints[1])}
	} else if interval.lower != "" {
		// Only lower bound: >=lower
		if interval.lowerInclusive {
			return []string{fmt.Sprintf(">=%s", interval.lower)}
		} else {
			return []string{fmt.Sprintf(">%s", interval.lower)}
		}
	} else if interval.upper != "" {
		// Only upper bound: <=upper
		if interval.upperInclusive {
			return []string{fmt.Sprintf("<=%s", interval.upper)}
		} else {
			return []string{fmt.Sprintf("<%s", interval.upper)}
		}
	}

	// Empty interval
	return []string{}
}
