package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/gem"
)

// gemContains implements VERS constraint checking for RubyGems ecosystem
func gemContains(constraints []string, version string) (bool, error) {
	e := &gem.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToGemRanges converts an interval to RubyGems range syntax
func intervalToGemRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("=%s", interval.exact)}
	}

	// Exclusions are handled separately, not as gem ranges
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
		return []string{strings.Join(parts, ",")}
	}

	// Empty interval
	return []string{}
}
