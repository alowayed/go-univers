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
	"sort"
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

	remaining := versString[len("vers:"):]
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

	// Sort constraints by version for consistent processing
	e := &maven.Ecosystem{}
	sort.Slice(versConstraints, func(i, j int) bool {
		vi, errI := e.NewVersion(versConstraints[i].version)
		vj, errJ := e.NewVersion(versConstraints[j].version)
		if errI != nil || errJ != nil {
			// Fallback to string comparison if version parsing fails
			return versConstraints[i].version < versConstraints[j].version
		}
		return vi.Compare(vj) < 0
	})

	// Apply constraint simplification for redundant constraints
	simplifiedConstraints := removeRedundantConstraints(versConstraints)

	// Convert simplified constraints to Maven ranges using VERS containment logic
	ranges, err := convertConstraintsToMavenRanges(simplifiedConstraints)
	if err != nil {
		return nil, err
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

// removeRedundantConstraints removes only truly redundant constraints that have the same version
func removeRedundantConstraints(constraints []versConstraint) []versConstraint {
	if len(constraints) <= 1 {
		return constraints
	}

	var result []versConstraint
	seen := make(map[string]versConstraint)

	for _, constraint := range constraints {
		key := constraint.version + constraint.operator

		// Skip exact duplicates
		if existing, exists := seen[key]; exists && existing.operator == constraint.operator {
			continue
		}

		// For same version, apply basic redundancy rules
		versionKey := constraint.version
		if existing, exists := seen[versionKey]; exists {
			// If we have >=X and >X for same version, keep >=X (more inclusive)
			if constraint.operator == ">=" && existing.operator == ">" {
				// Replace with more inclusive constraint
				for i, r := range result {
					if r.version == versionKey && r.operator == ">" {
						result[i] = constraint
						break
					}
				}
				seen[key] = constraint
				continue
			}
			// If we have <=X and <X for same version, keep <=X (more inclusive)
			if constraint.operator == "<=" && existing.operator == "<" {
				// Replace with more inclusive constraint
				for i, r := range result {
					if r.version == versionKey && r.operator == "<" {
						result[i] = constraint
						break
					}
				}
				seen[key] = constraint
				continue
			}
		}

		result = append(result, constraint)
		seen[key] = constraint
		seen[versionKey] = constraint
	}

	return result
}

// Helper functions to categorize operators
func isLowerBound(op string) bool {
	return op == ">" || op == ">="
}

func isUpperBound(op string) bool {
	return op == "<" || op == "<="
}

// convertConstraintsToMavenRanges converts simplified VERS constraints to Maven ranges
func convertConstraintsToMavenRanges(constraints []versConstraint) ([]*maven.VersionRange, error) {
	e := &maven.Ecosystem{}
	var ranges []*maven.VersionRange

	// Handle exact matches first
	for _, constraint := range constraints {
		if constraint.operator == "=" {
			rangeStr := fmt.Sprintf("[%s]", constraint.version)
			r, err := e.NewVersionRange(rangeStr)
			if err != nil {
				return nil, fmt.Errorf("failed to create Maven range '%s': %w", rangeStr, err)
			}
			ranges = append(ranges, r)
		}
	}

	// Group non-equality constraints into intervals
	var nonEqualConstraints []versConstraint
	for _, constraint := range constraints {
		if constraint.operator != "=" && constraint.operator != "!=" {
			nonEqualConstraints = append(nonEqualConstraints, constraint)
		}
	}

	// Process constraints pairwise to create intervals
	intervals := createIntervalsFromConstraints(nonEqualConstraints)

	for _, interval := range intervals {
		rangeStr := convertIntervalToMavenRange(interval)
		if rangeStr != "" {
			r, err := e.NewVersionRange(rangeStr)
			if err != nil {
				return nil, fmt.Errorf("failed to create Maven range '%s': %w", rangeStr, err)
			}
			ranges = append(ranges, r)
		}
	}

	return ranges, nil
}

// createIntervalsFromConstraints creates intervals from non-equality constraints
// This follows the VERS specification for processing constraints pairwise
func createIntervalsFromConstraints(constraints []versConstraint) []interval {
	var intervals []interval

	// Process constraints pairwise as specified in VERS containment checking
	for i := 0; i < len(constraints); i += 2 {
		if i+1 < len(constraints) {
			// Create interval from pair of constraints
			first := constraints[i]
			second := constraints[i+1]

			var lower, upper string
			var lowerInclusive, upperInclusive bool

			// Determine which constraint is lower bound and which is upper bound
			if isLowerBound(first.operator) && isUpperBound(second.operator) {
				lower = first.version
				lowerInclusive = first.operator == ">="
				upper = second.version
				upperInclusive = second.operator == "<="
			} else if isUpperBound(first.operator) && isLowerBound(second.operator) {
				lower = second.version
				lowerInclusive = second.operator == ">="
				upper = first.version
				upperInclusive = first.operator == "<="
			} else {
				// Both same type - create separate unbounded intervals
				intervals = append(intervals, createUnboundedInterval(first))
				intervals = append(intervals, createUnboundedInterval(second))
				continue
			}

			intervals = append(intervals, interval{
				lower:          lower,
				lowerInclusive: lowerInclusive,
				upper:          upper,
				upperInclusive: upperInclusive,
			})
		} else {
			// Single constraint - create unbounded interval
			intervals = append(intervals, createUnboundedInterval(constraints[i]))
		}
	}

	return intervals
}

// createUnboundedInterval creates an unbounded interval from a single constraint
func createUnboundedInterval(constraint versConstraint) interval {
	if isLowerBound(constraint.operator) {
		return interval{
			lower:          constraint.version,
			lowerInclusive: constraint.operator == ">=",
		}
	} else {
		return interval{
			upper:          constraint.version,
			upperInclusive: constraint.operator == "<=",
		}
	}
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
}

// convertIntervalToMavenRange converts an interval to Maven range syntax
func convertIntervalToMavenRange(interval interval) string {
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
