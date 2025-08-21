// Package vers provides VERS (Version Range Specification) parsing and utilities.
// VERS is a universal notation for expressing version ranges across different package ecosystems.
//
// VERS syntax: vers:<ecosystem>/<constraints>
// Example: vers:maven/>=1.0.0|<=2.0.0
//
// This package provides stateless functions for working with VERS notation.
package vers

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/maven"
)

// Contains checks if a version satisfies a VERS range using the stateless API.
// Example: Contains("vers:maven/>=1.0.0|<=2.0.0", "1.5.0") returns true.
func Contains(versRange, version string) (bool, error) {
	ecosystem, constraints, err := parseVersString(versRange)
	if err != nil {
		return false, fmt.Errorf("invalid VERS range: %w", err)
	}

	switch ecosystem {
	case "maven":
		return containsMaven(constraints, version)
	default:
		return false, fmt.Errorf("ecosystem '%s' not supported", ecosystem)
	}
}

// parseVersString parses a VERS string and returns the ecosystem and constraints.
func parseVersString(versString string) (ecosystem, constraints string, err error) {
	if !strings.HasPrefix(versString, "vers:") {
		return "", "", fmt.Errorf("invalid VERS format: must start with 'vers:'")
	}

	remaining := versString[5:]
	parts := strings.SplitN(remaining, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid VERS format: missing '/' separator")
	}

	ecosystem = parts[0]
	constraints = parts[1]

	if ecosystem == "" {
		return "", "", fmt.Errorf("invalid VERS format: empty ecosystem")
	}

	if constraints == "" {
		return "", "", fmt.Errorf("invalid VERS format: empty constraints")
	}

	return ecosystem, constraints, nil
}

// containsMaven implements Contains for Maven ecosystem
func containsMaven(constraints, version string) (bool, error) {
	// Parse the version using Maven
	e := &maven.Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		return false, fmt.Errorf("invalid Maven version '%s': %w", version, err)
	}

	// Parse VERS constraints and convert to Maven ranges
	ranges, err := convertVersToMavenRanges(constraints)
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

// versConstraint represents a single VERS constraint
type versConstraint struct {
	operator string // ">=", "<=", ">", "<", "=", "!="
	version  string
}

// convertVersToMavenRanges converts VERS constraints to Maven native ranges
func convertVersToMavenRanges(constraints string) ([]*maven.VersionRange, error) {
	// Parse individual constraints
	versConstraints, err := parseVersConstraints(constraints)
	if err != nil {
		return nil, err
	}

	// Group constraints into intervals according to VERS specification
	intervals, err := groupConstraintsIntoIntervals(versConstraints)
	if err != nil {
		return nil, err
	}

	// Convert each interval to a Maven range
	var ranges []*maven.VersionRange
	e := &maven.Ecosystem{}

	for _, interval := range intervals {
		rangeStr := convertIntervalToMavenRange(interval)
		if rangeStr == "" {
			continue // Skip empty intervals (like exclusions we don't handle yet)
		}
		r, err := e.NewVersionRange(rangeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to create Maven range '%s': %w", rangeStr, err)
		}
		ranges = append(ranges, r)
	}

	return ranges, nil
}

// parseVersConstraints parses VERS constraint string into individual constraints
func parseVersConstraints(constraints string) ([]versConstraint, error) {
	if constraints == "" {
		return nil, fmt.Errorf("empty constraints")
	}

	// Split on pipe separator
	constraintStrs := strings.Split(constraints, "|")
	var result []versConstraint

	for _, constraintStr := range constraintStrs {
		constraintStr = strings.TrimSpace(constraintStr)
		if constraintStr == "" {
			continue
		}

		constraint, err := parseVersConstraint(constraintStr)
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

// parseVersConstraint parses a single constraint string
func parseVersConstraint(constraintStr string) (versConstraint, error) {
	// Check for two-character operators first
	if len(constraintStr) >= 2 {
		twoChar := constraintStr[:2]
		switch twoChar {
		case ">=", "<=", "!=":
			version := strings.TrimSpace(constraintStr[2:])
			if version == "" {
				return versConstraint{}, fmt.Errorf("missing version after operator '%s'", twoChar)
			}
			return versConstraint{operator: twoChar, version: version}, nil
		}
	}

	// Check for single-character operators
	if len(constraintStr) >= 1 {
		oneChar := constraintStr[:1]
		switch oneChar {
		case ">", "<", "=":
			version := strings.TrimSpace(constraintStr[1:])
			if version == "" {
				return versConstraint{}, fmt.Errorf("missing version after operator '%s'", oneChar)
			}
			return versConstraint{operator: oneChar, version: version}, nil
		}
	}

	return versConstraint{}, fmt.Errorf("no valid operator found in constraint")
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

// groupConstraintsIntoIntervals groups VERS constraints into intervals according to the specification
func groupConstraintsIntoIntervals(constraints []versConstraint) ([]interval, error) {
	// For now, implement a simple version that handles basic cases
	// TODO: Implement full VERS state machine algorithm

	var intervals []interval
	var currentLower, currentUpper versConstraint
	var hasLower, hasUpper bool

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
			currentLower = constraint
			hasLower = true
		case "<=", "<":
			currentUpper = constraint
			hasUpper = true
		}

		// If we have both lower and upper bounds, create an interval
		if hasLower && hasUpper {
			intervals = append(intervals, interval{
				lower:          currentLower.version,
				lowerInclusive: currentLower.operator == ">=",
				upper:          currentUpper.version,
				upperInclusive: currentUpper.operator == "<=",
			})
			hasLower, hasUpper = false, false
		}
	}

	// Handle remaining single bounds
	if hasLower && !hasUpper {
		intervals = append(intervals, interval{
			lower:          currentLower.version,
			lowerInclusive: currentLower.operator == ">=",
		})
	}
	if hasUpper && !hasLower {
		intervals = append(intervals, interval{
			upper:          currentUpper.version,
			upperInclusive: currentUpper.operator == "<=",
		})
	}

	return intervals, nil
}

// convertIntervalToMavenRange converts an interval to Maven range syntax
func convertIntervalToMavenRange(interval interval) string {
	if interval.exact != "" {
		return fmt.Sprintf("[%s]", interval.exact)
	}

	if interval.exclude != "" {
		// Maven doesn't support exclusions directly, so we'll skip these for now
		// In a full implementation, we'd need to handle this differently
		return ""
	}

	// Convert interval bounds to Maven range syntax
	lowerBracket := "["
	if !interval.lowerInclusive {
		lowerBracket = "("
	}

	upperBracket := "]"
	if !interval.upperInclusive {
		upperBracket = ")"
	}

	if interval.lower != "" && interval.upper != "" {
		return fmt.Sprintf("%s%s,%s%s", lowerBracket, interval.lower, interval.upper, upperBracket)
	} else if interval.lower != "" {
		return fmt.Sprintf("%s%s,)", lowerBracket, interval.lower)
	} else if interval.upper != "" {
		return fmt.Sprintf("(,%s%s", interval.upper, upperBracket)
	}

	return ""
}
