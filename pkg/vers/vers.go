// Package vers provides VERS (Version Range Specification) parsing and utilities.
// VERS is a universal notation for expressing version ranges across different package ecosystems.
//
// VERS syntax: vers:<ecosystem>/<constraints>
// Example: vers:maven/>=1.0.0|<=2.0.0
//
// This package provides both stateless convenience functions and shared utilities
// for ecosystem-specific VERS parsing implementations.
package vers

import (
	"fmt"
	"strings"
)

// ParseVersString parses a VERS string and returns the ecosystem and constraints.
// Example: "vers:maven/>=1.0.0|<=2.0.0" returns ("maven", ">=1.0.0|<=2.0.0")
func ParseVersString(versString string) (ecosystem, constraints string, err error) {
	if !strings.HasPrefix(versString, "vers:") {
		return "", "", fmt.Errorf("invalid VERS format: must start with 'vers:'")
	}

	// Remove "vers:" prefix
	remaining := versString[5:]

	// Split on first "/" to separate ecosystem from constraints
	parts := strings.SplitN(remaining, "/", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid VERS format: missing '/' separator between ecosystem and constraints")
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

// Constraint represents a single version constraint in VERS notation.
type Constraint struct {
	Operator string // ">=", "<=", ">", "<", "=", "!="
	Version  string // The version string
}

// VersionComparator interface allows the VERS containment algorithm to work with any version type.
type VersionComparator interface {
	// Compare returns -1 if this < other, 0 if this == other, 1 if this > other
	Compare(other VersionComparator) int
	String() string
}

// ContainsVersion implements the VERS state machine algorithm for checking if a version
// is contained within a set of constraints. This implements the algorithm described in
// the VERS specification.
func ContainsVersion(version VersionComparator, constraints []Constraint, parseVersion func(string) (VersionComparator, error)) (bool, error) {
	if len(constraints) == 0 {
		return false, fmt.Errorf("no constraints provided")
	}

	// Handle single constraint case
	if len(constraints) == 1 {
		constraint := constraints[0]
		constraintVersion, err := parseVersion(constraint.Version)
		if err != nil {
			return false, fmt.Errorf("invalid constraint version '%s': %w", constraint.Version, err)
		}
		return satisfiesConstraint(version, constraint.Operator, constraintVersion), nil
	}

	// Check for explicit exclusions (!=) first
	for _, constraint := range constraints {
		if constraint.Operator == "!=" {
			constraintVersion, err := parseVersion(constraint.Version)
			if err != nil {
				return false, fmt.Errorf("invalid constraint version '%s': %w", constraint.Version, err)
			}
			if version.Compare(constraintVersion) == 0 {
				return false, nil // Version is explicitly excluded
			}
		}
	}

	// Check for explicit inclusions (=) 
	for _, constraint := range constraints {
		if constraint.Operator == "=" {
			constraintVersion, err := parseVersion(constraint.Version)
			if err != nil {
				return false, fmt.Errorf("invalid constraint version '%s': %w", constraint.Version, err)
			}
			if version.Compare(constraintVersion) == 0 {
				return true, nil // Version is explicitly included
			}
		}
	}

	// Filter out = and != constraints for interval checking
	var intervalConstraints []Constraint
	for _, constraint := range constraints {
		if constraint.Operator != "=" && constraint.Operator != "!=" {
			intervalConstraints = append(intervalConstraints, constraint)
		}
	}

	if len(intervalConstraints) == 0 {
		return false, nil // No interval constraints and no explicit match
	}

	// Sort constraints by version for interval detection
	sortedConstraints, err := sortConstraintsByVersion(intervalConstraints, parseVersion)
	if err != nil {
		return false, err
	}

	// Check intervals using the VERS state machine algorithm
	return checkIntervals(version, sortedConstraints, parseVersion)
}

// satisfiesConstraint checks if a version satisfies a single constraint
func satisfiesConstraint(version VersionComparator, operator string, constraintVersion VersionComparator) bool {
	cmp := version.Compare(constraintVersion)
	switch operator {
	case ">=":
		return cmp >= 0
	case "<=":
		return cmp <= 0
	case ">":
		return cmp > 0
	case "<":
		return cmp < 0
	case "=":
		return cmp == 0
	case "!=":
		return cmp != 0
	default:
		return false
	}
}

// sortConstraintsByVersion sorts constraints by their version values
func sortConstraintsByVersion(constraints []Constraint, parseVersion func(string) (VersionComparator, error)) ([]constraintWithVersion, error) {
	var result []constraintWithVersion
	for _, constraint := range constraints {
		version, err := parseVersion(constraint.Version)
		if err != nil {
			return nil, fmt.Errorf("invalid constraint version '%s': %w", constraint.Version, err)
		}
		result = append(result, constraintWithVersion{
			constraint: constraint,
			version:    version,
		})
	}

	// Simple bubble sort by version (could be optimized)
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].version.Compare(result[j].version) > 0 {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

type constraintWithVersion struct {
	constraint Constraint
	version    VersionComparator
}

// checkIntervals implements the VERS interval checking algorithm
func checkIntervals(version VersionComparator, sortedConstraints []constraintWithVersion, parseVersion func(string) (VersionComparator, error)) (bool, error) {
	if len(sortedConstraints) == 0 {
		return false, nil
	}

	// Check first constraint for lower bound patterns
	first := sortedConstraints[0]
	if (first.constraint.Operator == "<" || first.constraint.Operator == "<=") {
		if satisfiesConstraint(version, first.constraint.Operator, first.version) {
			return true, nil
		}
	}

	// Check last constraint for upper bound patterns  
	if len(sortedConstraints) > 1 {
		last := sortedConstraints[len(sortedConstraints)-1]
		if (last.constraint.Operator == ">" || last.constraint.Operator == ">=") {
			if satisfiesConstraint(version, last.constraint.Operator, last.version) {
				return true, nil
			}
		}
	}

	// Check for intervals between consecutive constraints
	for i := 0; i < len(sortedConstraints)-1; i++ {
		current := sortedConstraints[i]
		next := sortedConstraints[i+1]

		// Check if current is a lower bound and next is an upper bound
		if (current.constraint.Operator == ">" || current.constraint.Operator == ">=") &&
			(next.constraint.Operator == "<" || next.constraint.Operator == "<=") {

			// Check if version falls in the interval [current, next]
			satisfiesCurrent := satisfiesConstraint(version, current.constraint.Operator, current.version)
			satisfiesNext := satisfiesConstraint(version, next.constraint.Operator, next.version)

			if satisfiesCurrent && satisfiesNext {
				return true, nil
			}
		}
	}

	return false, nil
}

// ParseVersConstraints parses VERS constraint string into individual constraints.
// Example: ">=1.0.0|<=2.0.0|!=1.5.0" returns constraints with operators and versions.
func ParseVersConstraints(constraints string) ([]Constraint, error) {
	if constraints == "" {
		return nil, fmt.Errorf("empty constraints")
	}

	// Split on pipe separator
	constraintStrs := strings.Split(constraints, "|")
	result := make([]Constraint, 0, len(constraintStrs))

	for _, constraintStr := range constraintStrs {
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

// parseConstraint parses a single constraint string into operator and version.
func parseConstraint(constraintStr string) (Constraint, error) {
	// Check for two-character operators first
	if len(constraintStr) >= 2 {
		twoChar := constraintStr[:2]
		switch twoChar {
		case ">=", "<=", "!=":
			version := strings.TrimSpace(constraintStr[2:])
			if version == "" {
				return Constraint{}, fmt.Errorf("missing version after operator '%s'", twoChar)
			}
			return Constraint{
				Operator: twoChar,
				Version:  version,
			}, nil
		}
	}

	// Check for single-character operators
	if len(constraintStr) >= 1 {
		oneChar := constraintStr[:1]
		switch oneChar {
		case ">", "<", "=":
			version := strings.TrimSpace(constraintStr[1:])
			if version == "" {
				return Constraint{}, fmt.Errorf("missing version after operator '%s'", oneChar)
			}
			return Constraint{
				Operator: oneChar,
				Version:  version,
			}, nil
		}
	}

	return Constraint{}, fmt.Errorf("no valid operator found in constraint")
}

