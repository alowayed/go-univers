package vers

import (
	"fmt"

	"github.com/alowayed/go-univers/pkg/ecosystem/nuget"
)

// nugetContains implements VERS constraint checking for NuGet ecosystem
func nugetContains(constraints []string, version string) (bool, error) {
	e := &nuget.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToNugetRanges converts an interval to NuGet range syntax
func intervalToNugetRanges(interval interval) []string {
	switch {
	case interval.exact != "":
		// Handle exact matches - NuGet uses bracket notation for exact matches
		return []string{fmt.Sprintf("[%s]", interval.exact)}
	case interval.exclude != "":
		// Exclusions are handled separately, not as NuGet ranges
		return []string{} // Return empty - excludes handled in contains function
	case interval.lower != "" && interval.upper == "":
		// Lower bound only - use unbounded range [version,) for inclusive, comma-separated constraint for exclusive
		if interval.lowerInclusive {
			return []string{fmt.Sprintf("[%s,)", interval.lower)}
		}
		// NuGet doesn't support (version,) syntax, use comma-separated constraint
		return []string{fmt.Sprintf(">%s,", interval.lower)}
	case interval.upper != "" && interval.lower == "":
		// Upper bound only - use unbounded range (,version] for inclusive, comma-separated constraint for exclusive
		if interval.upperInclusive {
			return []string{fmt.Sprintf("(,%s]", interval.upper)}
		}
		// NuGet doesn't support (,version) syntax, use comma-separated constraint
		return []string{fmt.Sprintf("<%s,", interval.upper)}
	case interval.lower != "" && interval.upper != "":
		// Both bounds - use comma-separated constraints
		lowerOp := ">"
		if interval.lowerInclusive {
			lowerOp = ">="
		}
		upperOp := "<"
		if interval.upperInclusive {
			upperOp = "<="
		}
		return []string{fmt.Sprintf("%s%s,%s%s", lowerOp, interval.lower, upperOp, interval.upper)}
	}

	// Empty interval
	return []string{}
}
