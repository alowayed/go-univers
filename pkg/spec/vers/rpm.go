package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/rpm"
)

// rpmContains implements VERS constraint checking for RPM ecosystem
func rpmContains(constraints []string, version string) (bool, error) {
	e := &rpm.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToRpmRanges converts an interval to RPM range syntax
func intervalToRpmRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", interval.exact)}
	}

	// Exclusions are handled separately, not as RPM ranges
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
		// RPM supports comma-separated constraints like gem and cargo
		return []string{strings.Join(parts, ",")}
	}

	// Empty interval
	return []string{}
}
