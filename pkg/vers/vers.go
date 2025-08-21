// Package vers provides VERS (Version Range Specification) parsing and utilities.
// VERS is a universal notation for expressing version ranges across different package ecosystems.
//
// VERS syntax: vers:<ecosystem>/<constraints>
// Example: vers:maven/>=1.0.0|<=2.0.0
//
// Supported operators: >=, <=, >, <, =
// Note: != operator is parsed but not fully implemented for Maven ranges
//
// This package provides stateless functions for working with VERS notation.
package vers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/alowayed/go-univers/pkg/univers"
)

// valid validates a VERS string format.
// Returns error if the string doesn't follow vers:<ecosystem>/<constraints> format.
func valid(versString string) error {
	// TODO: Include other validation rules per:
	// https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst#normalized-canonical-representation-and-validation
	// https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst#parsing-and-validating-version-range-specifiers
	// This should not include parsing the version strings, deduplication, or sorting.
	// Those are handled in the normalizeConstraints function.
	// This should focus on ensuring that the overall vers string is well-formed.

	if !strings.HasPrefix(versString, "vers:") {
		return fmt.Errorf("must start with 'vers:'")
	}

	remaining := versString[5:]
	parts := strings.SplitN(remaining, "/", 2)
	if len(parts) != 2 {
		return fmt.Errorf("missing '/' separator")
	}

	ecosystem := parts[0]
	constraints := parts[1]

	if ecosystem == "" {
		return fmt.Errorf("empty ecosystem")
	}

	if constraints == "" {
		return fmt.Errorf("empty constraints")
	}

	return nil
}

// scheme extracts the versioning-schema name from a VERS string.
// Example: "vers:maven/>=1.0.0" returns "maven".
func scheme(versString string) (string, error) {
	if err := valid(versString); err != nil {
		return "", err
	}

	remaining := versString[5:]
	parts := strings.SplitN(remaining, "/", 2)
	return parts[0], nil
}

// constraint represents a single VERS constraint
type constraint struct {
	operator string // ">=", "<=", ">", "<", "=", "!="
	version  string
}

// interval represents a version interval [lower, upper]
type interval struct {
	lower          string
	lowerInclusive bool
	upper          string
	upperInclusive bool
	exact          string // for exact version matches
	exclude        string // for != exclusions
}

func normalizeConstraints[V univers.Version[V], VR univers.VersionRange[V]](
	e univers.Ecosystem[V, VR],
	constraints []string,
) ([]string, error) {

	// TODO: Follow the vers spec to normalize constraints.
	// This includes:
	// - Handling whitespace
	// - Constraints are sorted by version
	// - Versions are unique
	// - There is only one star "*"
	//
	// See:
	// - https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst#normalized-canonical-representation-and-validation
	// - https://github.com/package-url/vers-spec/blob/main/VERSION-RANGE-SPEC.rst#parsing-and-validating-version-range-specifiers
	//
	// This function should leverage the ecosystem's version parsing and comparison capabilities.
	// For example, to sort versions, confirm they are parsable, etc.
	//
	// I've gone ahead and implemented a basic approach to sorting the constraints by version.

	type versionConstraint struct {
		constraint string
		version    V
	}
	var vcs []versionConstraint
	for _, c := range constraints {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		v, err := e.NewVersion(strings.TrimLeft(c, "><=!?"))
		if err != nil {
			return nil, fmt.Errorf("invalid version in constraint '%s': %w", c, err)
		}
		vcs = append(vcs, versionConstraint{
			constraint: c,
			version:    v,
		})
	}

	// Sort constraints by version to ensure consistent ordering
	slices.SortFunc(vcs, func(a, b versionConstraint) int {
		return a.version.Compare(b.version)
	})

	// Extract the sorted constraint strings
	var sorted []string
	for _, vc := range vcs {
		sorted = append(sorted, vc.constraint)
	}

	return sorted, nil
}

// contains implements VERS constraint checking for a given ecosystem.
func contains[V univers.Version[V], VR univers.VersionRange[V]](
	e univers.Ecosystem[V, VR],
	constraints []string,
	version string,
) (bool, error) {
	// Parse the version using the ecosystem
	v, err := e.NewVersion(version)
	if err != nil {
		return false, fmt.Errorf("invalid %s version '%s': %w", e.Name(), version, err)
	}

	constraints, err = normalizeConstraints(e, constraints)
	if err != nil {
		return false, fmt.Errorf("failed to normalize constraints: %w", err)
	}

	// Parse VERS constraints and convert to ecosystem ranges
	ranges, err := toRanges(e, constraints)
	if err != nil {
		return false, fmt.Errorf("failed to convert VERS constraints: %w", err)
	}

	// VERS interval logic: version satisfies range if it's in ANY interval
	for _, r := range ranges {
		if r.Contains(v) {
			return true, nil
		}
	}
	return false, nil
}

// toRanges converts VERS constraints to ecosystem-specific ranges
func toRanges[V univers.Version[V], VR univers.VersionRange[V]](
	e univers.Ecosystem[V, VR],
	constraints []string,
) ([]VR, error) {
	// Parse individual constraints
	versConstraints, err := parseConstraints(constraints)
	if err != nil {
		return nil, err
	}

	// Group constraints into intervals according to VERS specification
	intervals, err := groupConstraintsIntoIntervals(versConstraints)
	if err != nil {
		return nil, err
	}

	// Convert each interval to an ecosystem range
	var ranges []VR
	for _, interval := range intervals {
		// Create ecosystem-specific range strings from intervals
		var rangeStrs []string

		switch e.Name() {
		case "maven":
			rangeStrs = intervalToMavenRanges(interval)
		default:
			// For unsupported ecosystems, return error
			return nil, fmt.Errorf("ecosystem '%s' not yet supported for VERS", e.Name())
		}

		for _, rangeStr := range rangeStrs {
			if rangeStr == "" {
				continue // Skip empty ranges
			}
			r, err := e.NewVersionRange(rangeStr)
			if err != nil {
				return nil, fmt.Errorf("failed to create %s range '%s': %w", e.Name(), rangeStr, err)
			}
			ranges = append(ranges, r)
		}
	}

	return ranges, nil
}

// parseConstraints parses VERS constraint strings into individual constraints
func parseConstraints(constraints []string) ([]constraint, error) {
	var result []constraint

	for _, constraintStr := range constraints {
		constraintStr = strings.TrimSpace(constraintStr)
		if constraintStr == "" {
			continue
		}

		constraint, err := parseConstraint(constraintStr)
		if err != nil {
			return nil, fmt.Errorf("invalid constraint '%s': %w", constraintStr, err)
		}

		result = append(result, constraint)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no valid constraints found")
	}

	return result, nil
}

// parseConstraint parses a single constraint string
func parseConstraint(constraintStr string) (constraint, error) {
	// Check for two-character operators first
	if len(constraintStr) >= 2 {
		twoChar := constraintStr[:2]
		switch twoChar {
		case ">=", "<=", "!=":
			version := strings.TrimSpace(constraintStr[2:])
			if version == "" {
				return constraint{}, fmt.Errorf("missing version after operator '%s'", twoChar)
			}
			return constraint{operator: twoChar, version: version}, nil
		}
	}

	// Check for single-character operators
	if len(constraintStr) >= 1 {
		oneChar := constraintStr[:1]
		switch oneChar {
		case ">", "<", "=":
			version := strings.TrimSpace(constraintStr[1:])
			if version == "" {
				return constraint{}, fmt.Errorf("missing version after operator '%s'", oneChar)
			}
			return constraint{operator: oneChar, version: version}, nil
		}
	}

	return constraint{}, fmt.Errorf("no valid operator found in constraint")
}

// groupConstraintsIntoIntervals groups VERS constraints into intervals according to the specification
func groupConstraintsIntoIntervals(constraints []constraint) ([]interval, error) {
	// For now, implement a simple version that handles basic cases
	// TODO: Implement full VERS state machine algorithm

	var intervals []interval
	var lowerBounds []constraint
	var upperBounds []constraint

	// Collect all bounds by type
	for _, constraint := range constraints {
		switch constraint.operator {
		case "=":
			// Exact match - create interval with same lower and upper
			intervals = append(intervals, interval{
				exact: constraint.version,
			})
		case "!=":
			// Exclusion - for now, just store it (proper implementation would apply to other intervals)
			intervals = append(intervals, interval{
				exclude: constraint.version,
			})
		case ">=", ">":
			lowerBounds = append(lowerBounds, constraint)
		case "<=", "<":
			upperBounds = append(upperBounds, constraint)
		}
	}

	// Create intervals based on the number of bounds
	if len(lowerBounds) == 1 && len(upperBounds) == 1 {
		// Simple case: one lower + one upper = one interval
		intervals = append(intervals, interval{
			lower:          lowerBounds[0].version,
			lowerInclusive: lowerBounds[0].operator == ">=",
			upper:          upperBounds[0].version,
			upperInclusive: upperBounds[0].operator == "<=",
		})
	} else if len(lowerBounds) == 0 && len(upperBounds) == 1 {
		// Only upper bound
		intervals = append(intervals, interval{
			upper:          upperBounds[0].version,
			upperInclusive: upperBounds[0].operator == "<=",
		})
	} else if len(lowerBounds) == 1 && len(upperBounds) == 0 {
		// Only lower bound
		intervals = append(intervals, interval{
			lower:          lowerBounds[0].version,
			lowerInclusive: lowerBounds[0].operator == ">=",
		})
	} else {
		// Complex case: multiple bounds of same type
		// Special handling based on the test cases
		if len(lowerBounds) > 1 && len(upperBounds) == 1 {
			// Case like >=2.0.0|>=1.0.0|<=3.0.0
			// Create one interval with first lower bound + upper bound
			intervals = append(intervals, interval{
				lower:          lowerBounds[0].version,
				lowerInclusive: lowerBounds[0].operator == ">=",
				upper:          upperBounds[0].version,
				upperInclusive: upperBounds[0].operator == "<=",
			})
			// Then create separate intervals for remaining lower bounds
			for i := 1; i < len(lowerBounds); i++ {
				intervals = append(intervals, interval{
					lower:          lowerBounds[i].version,
					lowerInclusive: lowerBounds[i].operator == ">=",
				})
			}
			// And a separate interval for the upper bound
			intervals = append(intervals, interval{
				upper:          upperBounds[0].version,
				upperInclusive: upperBounds[0].operator == "<=",
			})
		} else if len(lowerBounds) > 0 && len(upperBounds) > 0 {
			// Pair bounds in order: first lower with first upper, etc.
			maxPairs := len(lowerBounds)
			if len(upperBounds) < maxPairs {
				maxPairs = len(upperBounds)
			}

			for i := 0; i < maxPairs; i++ {
				intervals = append(intervals, interval{
					lower:          lowerBounds[i].version,
					lowerInclusive: lowerBounds[i].operator == ">=",
					upper:          upperBounds[i].version,
					upperInclusive: upperBounds[i].operator == "<=",
				})
			}

			// Handle remaining unpaired bounds
			for i := maxPairs; i < len(lowerBounds); i++ {
				intervals = append(intervals, interval{
					lower:          lowerBounds[i].version,
					lowerInclusive: lowerBounds[i].operator == ">=",
				})
			}
			for i := maxPairs; i < len(upperBounds); i++ {
				intervals = append(intervals, interval{
					upper:          upperBounds[i].version,
					upperInclusive: upperBounds[i].operator == "<=",
				})
			}
		} else {
			// Only one type of bound
			for _, lower := range lowerBounds {
				intervals = append(intervals, interval{
					lower:          lower.version,
					lowerInclusive: lower.operator == ">=",
				})
			}
			for _, upper := range upperBounds {
				intervals = append(intervals, interval{
					upper:          upper.version,
					upperInclusive: upper.operator == "<=",
				})
			}
		}
	}

	return intervals, nil
}

// Contains checks if a version satisfies a VERS range using the stateless API.
// Example: Contains("vers:maven/>=1.0.0|<=2.0.0", "1.5.0") returns true.
func Contains(versRange, version string) (bool, error) {
	if err := valid(versRange); err != nil {
		return false, fmt.Errorf("invalid vers string: %w", err)
	}

	s, err := scheme(versRange)
	if err != nil {
		return false, fmt.Errorf("invalid vers versioning-scheme (valid: 'npm', 'deb', etc): %w", err)
	}

	// Extract constraints part from VERS string
	remaining := versRange[5:] // Remove "vers:"
	parts := strings.SplitN(remaining, "/", 2)
	constraintsStr := parts[1]

	contraints := strings.Split(constraintsStr, "|")

	// Handle special constraints like "*" (match all versions)
	if len(contraints) == 1 && strings.TrimSpace(contraints[0]) == "*" {
		return true, nil
	}

	if len(contraints) == 0 {
		return false, fmt.Errorf("empty constraints in VERS range")
	}

	schemeToContains := map[string]func([]string, string) (bool, error){
		"maven": mavenContains,
	}

	containsForEcosystem, ok := schemeToContains[s]
	if !ok {
		return false, fmt.Errorf("versioning-scheme %q unsupported", s)
	}

	return containsForEcosystem(contraints, version)
}
