package vers

import (
	"fmt"

	"github.com/alowayed/go-univers/pkg/ecosystem/maven"
)

// mavenContains implements VERS constraint checking for Maven ecosystem
func mavenContains(constraints []string, version string) (bool, error) {
	e := &maven.Ecosystem{}
	return contains(e, constraints, version)
}

// intervalToMavenRanges converts an interval to Maven range syntax
func intervalToMavenRanges(interval interval) []string {
	// Handle exact matches
	if interval.exact != "" {
		return []string{fmt.Sprintf("[%s]", interval.exact)}
	}

	// Exclusions are now handled separately, not as Maven ranges
	if interval.exclude != "" {
		return []string{} // Return empty - excludes handled in contains function
	}

	// Handle regular intervals with bounds
	lowerBracket := "["
	if !interval.lowerInclusive {
		lowerBracket = "("
	}

	upperBracket := "]"
	if !interval.upperInclusive {
		upperBracket = ")"
	}

	if interval.lower != "" && interval.upper != "" {
		// Both bounds: [lower,upper]
		return []string{fmt.Sprintf("%s%s,%s%s", lowerBracket, interval.lower, interval.upper, upperBracket)}
	} else if interval.lower != "" {
		// Only lower bound: [lower,)
		return []string{fmt.Sprintf("%s%s,)", lowerBracket, interval.lower)}
	} else if interval.upper != "" {
		// Only upper bound: (,upper]
		return []string{fmt.Sprintf("(,%s%s", interval.upper, upperBracket)}
	}

	// Empty interval
	return []string{}
}
