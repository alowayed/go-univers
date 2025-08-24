package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/alpine"
)

// alpineContains implements VERS constraint checking for Alpine ecosystem
func alpineContains(constraints []string, version string) (bool, error) {
	e := &alpine.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToAlpineRanges converts an interval to Alpine range syntax
func intervalToAlpineRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", interval.exact)}
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
		return []string{strings.Join(parts, " ")}
	}

	// Empty interval
	return []string{}
}
