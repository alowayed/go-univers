package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/nuget"
)

// nugetContains implements VERS constraint checking for NuGet ecosystem
func nugetContains(constraints []string, version string) (bool, error) {
	e := &nuget.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToNugetRanges converts an interval to NuGet range syntax
func intervalToNugetRanges(interval interval) []string {
	// Handle exact matches - NuGet uses bracket notation for exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("[%s]", interval.exact)}
	}

	// Exclusions are handled separately, not as NuGet ranges
	if interval.exclude != "" {
		return []string{} // Return empty - excludes handled in contains function
	}

	// Handle regular intervals with bounds
	// NuGet has special handling for single bounds using unbounded ranges
	if interval.lower != "" && interval.upper == "" {
		// Lower bound only - use unbounded range [version,) for inclusive, comma-separated constraint for exclusive
		if interval.lowerInclusive {
			return []string{fmt.Sprintf("[%s,)", interval.lower)}
		} else {
			// NuGet doesn't support (version,) syntax, use comma-separated constraint
			return []string{fmt.Sprintf(">%s,", interval.lower)}
		}
	}

	if interval.upper != "" && interval.lower == "" {
		// Upper bound only - use unbounded range (,version] for inclusive, comma-separated constraint for exclusive
		if interval.upperInclusive {
			return []string{fmt.Sprintf("(,%s]", interval.upper)}
		} else {
			// NuGet doesn't support (,version) syntax, use comma-separated constraint
			return []string{fmt.Sprintf("<%s,", interval.upper)}
		}
	}

	if interval.lower != "" && interval.upper != "" {
		// Both bounds - use comma-separated constraints
		var parts []string
		op := ">"
		if interval.lowerInclusive {
			op = ">="
		}
		parts = append(parts, fmt.Sprintf("%s%s", op, interval.lower))

		op = "<"
		if interval.upperInclusive {
			op = "<="
		}
		parts = append(parts, fmt.Sprintf("%s%s", op, interval.upper))

		return []string{strings.Join(parts, ",")}
	}

	// Empty interval
	return []string{}
}
