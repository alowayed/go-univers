package cargo

import (
	"fmt"
	"strings"
)

// VersionRange represents a Cargo version range with Cargo-specific syntax support
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single Cargo version constraint
type constraint struct {
	operator string
	version  *Version
}

// NewVersionRange creates a new Cargo version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	original := rangeStr
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	constraints, err := parseConstraints(rangeStr, e)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    original,
	}, nil
}

// parseConstraints parses Cargo constraint syntax
func parseConstraints(rangeStr string, ecosystem *Ecosystem) ([]*constraint, error) {
	// Handle multiple constraints separated by commas (AND logic)
	parts := strings.Split(rangeStr, ",")
	var constraints []*constraint

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		constraint, err := parseConstraint(part, ecosystem)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, constraint)
	}

	if len(constraints) == 0 {
		return nil, fmt.Errorf("no valid constraints found")
	}

	return constraints, nil
}

// parseConstraint parses a single Cargo constraint
func parseConstraint(constraintStr string, ecosystem *Ecosystem) (*constraint, error) {
	constraintStr = strings.TrimSpace(constraintStr)

	// Handle caret constraints: ^1.2.3
	if strings.HasPrefix(constraintStr, "^") {
		version := strings.TrimSpace(constraintStr[1:])
		parsedVersion, err := ecosystem.NewVersion(version)
		if err != nil {
			return nil, fmt.Errorf("invalid version in caret constraint: %v", err)
		}
		return &constraint{operator: "^", version: parsedVersion}, nil
	}

	// Handle tilde constraints: ~1.2.3
	if strings.HasPrefix(constraintStr, "~") {
		version := strings.TrimSpace(constraintStr[1:])
		parsedVersion, err := ecosystem.NewVersion(version)
		if err != nil {
			return nil, fmt.Errorf("invalid version in tilde constraint: %v", err)
		}
		return &constraint{operator: "~", version: parsedVersion}, nil
	}

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(constraintStr, op) {
			version := strings.TrimSpace(constraintStr[len(op):])
			if version == "" {
				return nil, fmt.Errorf("constraint %s requires version", op)
			}
			parsedVersion, err := ecosystem.NewVersion(version)
			if err != nil {
				return nil, fmt.Errorf("invalid version in %s constraint: %v", op, err)
			}
			return &constraint{operator: op, version: parsedVersion}, nil
		}
	}

	// Handle wildcard patterns: 1.2.*
	if strings.Contains(constraintStr, "*") {
		return parseWildcardConstraint(constraintStr, ecosystem)
	}

	// Default to exact match
	parsedVersion, err := ecosystem.NewVersion(constraintStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version in exact constraint: %v", err)
	}
	return &constraint{operator: "=", version: parsedVersion}, nil
}

// parseWildcardConstraint handles wildcard patterns like 1.2.*
func parseWildcardConstraint(constraintStr string, ecosystem *Ecosystem) (*constraint, error) {
	// Replace * with 0 for parsing, then handle specially in Contains
	normalized := strings.ReplaceAll(constraintStr, "*", "0")
	parsedVersion, err := ecosystem.NewVersion(normalized)
	if err != nil {
		return nil, fmt.Errorf("invalid wildcard constraint: %v", err)
	}
	return &constraint{operator: "*", version: parsedVersion}, nil
}

// String returns the string representation of the version range
func (vr *VersionRange) String() string {
	return vr.original
}

// Contains checks if a version satisfies this range
func (vr *VersionRange) Contains(version *Version) bool {
	// All constraints must be satisfied (AND logic)
	for _, c := range vr.constraints {
		if !satisfiesConstraint(version, c) {
			return false
		}
	}

	return true
}

// satisfiesConstraint checks if a version satisfies a single constraint
func satisfiesConstraint(version *Version, c *constraint) bool {
	switch c.operator {
	case "=":
		return version.Compare(c.version) == 0
	case "!=":
		return version.Compare(c.version) != 0
	case ">":
		return version.Compare(c.version) > 0
	case ">=":
		return version.Compare(c.version) >= 0
	case "<":
		return version.Compare(c.version) < 0
	case "<=":
		return version.Compare(c.version) <= 0
	case "^":
		return satisfiesCaretConstraint(version, c.version)
	case "~":
		return satisfiesTildeConstraint(version, c.version)
	case "*":
		return satisfiesWildcardConstraint(version, c.version)
	default:
		return false
	}
}

// satisfiesCaretConstraint checks if version satisfies caret constraint (^1.2.3)
// Caret allows changes that do not modify the left-most non-zero digit
func satisfiesCaretConstraint(version, constraint *Version) bool {
	// Must be >= constraint version
	if version.Compare(constraint) < 0 {
		return false
	}

	// Major version must be the same
	if version.major != constraint.major {
		return false
	}

	// If major > 0, minor and patch can be anything >= constraint
	if constraint.major > 0 {
		return true
	}

	// If major == 0, minor must be the same
	if version.minor != constraint.minor {
		return false
	}

	// If major == 0 and minor > 0, patch can be anything >= constraint
	if constraint.minor > 0 {
		return true
	}

	// If major == 0 and minor == 0, patch must be the same
	return version.patch == constraint.patch
}

// satisfiesTildeConstraint checks if version satisfies tilde constraint (~1.2.3)
// Tilde allows patch-level changes if a minor version is specified
// ~1.2.3 := >=1.2.3 <1.(2+1).0 := >=1.2.3 <1.3.0
// ~1.2 := >=1.2.0 <1.(2+1).0 := >=1.2.0 <1.3.0
// ~1 := >=1.0.0 <(1+1).0.0 := >=1.0.0 <2.0.0
func satisfiesTildeConstraint(version, constraint *Version) bool {
	// Must be >= constraint version
	if version.Compare(constraint) < 0 {
		return false
	}

	// Major version must be the same
	if version.major != constraint.major {
		return false
	}

	// Minor version must be the same
	if version.minor != constraint.minor {
		return false
	}

	// Patch can be anything >= constraint patch
	return true
}

// satisfiesWildcardConstraint checks if version satisfies wildcard constraint (1.2.*)
func satisfiesWildcardConstraint(version, constraint *Version) bool {
	// The original constraint string determines what parts to match
	// For now, implement basic wildcard matching
	
	// Major must match
	if version.major != constraint.major {
		return false
	}

	// Minor must match (assuming constraint was like 1.2.*)
	if version.minor != constraint.minor {
		return false
	}

	// Patch can be anything
	return true
}