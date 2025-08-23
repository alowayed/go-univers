package vers

import (
	"fmt"
	"strings"

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
	var lowerConstraint, upperConstraint string

	if interval.lower != "" {
		op := ">"
		if interval.lowerInclusive {
			op = ">="
		}
		lowerConstraint = fmt.Sprintf("%s%s", op, ensureVPrefix(interval.lower))
	}

	if interval.upper != "" {
		op := "<"
		if interval.upperInclusive {
			op = "<="
		}
		upperConstraint = fmt.Sprintf("%s%s", op, ensureVPrefix(interval.upper))
	}

	if lowerConstraint != "" && upperConstraint != "" {
		return []string{fmt.Sprintf("%s %s", lowerConstraint, upperConstraint)}
	} else if lowerConstraint != "" {
		return []string{lowerConstraint}
	} else if upperConstraint != "" {
		return []string{upperConstraint}
	}

	return []string{}
}

// ensureVPrefix ensures a version string has the v prefix required for Go modules
func ensureVPrefix(version string) string {
	if version == "" {
		return version
	}
	if !strings.HasPrefix(version, "v") {
		return "v" + version
	}
	return version
}
