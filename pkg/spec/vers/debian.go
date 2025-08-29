package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/debian"
)

// debianContains implements VERS constraint checking for Debian ecosystem
func debianContains(constraints []string, version string) (bool, error) {
	e := &debian.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToDebianRanges converts an interval to Debian range syntax
func intervalToDebianRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", interval.exact)}
	}

	// Exclusions are handled separately, not as Debian ranges
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
		// Debian supports comma-separated constraints like rpm, gem and cargo
		return []string{strings.Join(parts, ",")}
	}

	// Empty interval
	return []string{}
}
