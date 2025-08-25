package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/cargo"
)

// cargoContains implements VERS constraint checking for Cargo (Rust) ecosystem
func cargoContains(constraints []string, version string) (bool, error) {
	e := &cargo.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToCargoRanges converts an interval to Cargo range syntax
func intervalToCargoRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", interval.exact)}
	}

	// Exclusions are handled separately, not as cargo ranges
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
		// Cargo uses comma-separated constraints like gem
		return []string{strings.Join(parts, ",")}
	}

	// Empty interval
	return []string{}
}
