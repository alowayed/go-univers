package gentoo

import (
	"fmt"
	"strings"
)

// VersionRange represents a Gentoo version range with Gentoo-specific syntax support
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single Gentoo version constraint
type constraint struct {
	operator string
	version  string
}

// NewVersionRange creates a new Gentoo version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	constraints, err := parseRange(rangeStr)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    rangeStr,
	}, nil
}

// parseRange parses Gentoo range syntax into constraints
func parseRange(rangeStr string) ([]*constraint, error) {
	rangeStr = strings.TrimSpace(rangeStr)

	// Handle comma-separated constraints (>=1.0.0, <2.0.0)
	if strings.Contains(rangeStr, ",") {
		return parseCommaSeparatedConstraints(rangeStr)
	}

	// Handle space-separated constraints (>=1.0.0 <2.0.0)
	if strings.Contains(rangeStr, " ") {
		return parseSpaceSeparatedConstraints(rangeStr)
	}

	// Handle single constraint
	return parseSingleConstraint(rangeStr)
}

// parseSingleConstraint parses a single Gentoo constraint
func parseSingleConstraint(c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			version := strings.TrimSpace(c[len(op):])
			if version == "" {
				return nil, fmt.Errorf("missing version after operator %s", op)
			}
			return []*constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to exact match
	return []*constraint{{operator: "=", version: c}}, nil
}

// parseCommaSeparatedConstraints handles comma-separated constraints (>=1.0.0, <2.0.0)
func parseCommaSeparatedConstraints(rangeStr string) ([]*constraint, error) {
	parts := strings.Split(rangeStr, ",")
	var constraints []*constraint

	for _, part := range parts {
		partConstraints, err := parseSingleConstraint(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, partConstraints...)
	}

	return constraints, nil
}

// parseSpaceSeparatedConstraints handles space-separated constraints (>=1.0.0 <2.0.0)
func parseSpaceSeparatedConstraints(rangeStr string) ([]*constraint, error) {
	parts := strings.Fields(rangeStr)
	var constraints []*constraint

	for _, part := range parts {
		partConstraints, err := parseSingleConstraint(part)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, partConstraints...)
	}

	return constraints, nil
}

// String returns the string representation of the range
func (gr *VersionRange) String() string {
	return gr.original
}

// Contains checks if a version is within this range
func (gr *VersionRange) Contains(version *Version) bool {
	// AND logic: ALL constraints must be satisfied
	for _, constraint := range gr.constraints {
		if !constraint.matches(version) {
			return false
		}
	}
	return true
}

// matches checks if the given version matches this constraint
func (c *constraint) matches(version *Version) bool {
	e := &Ecosystem{}
	constraintVersion, err := e.NewVersion(c.version)
	if err != nil {
		return false
	}

	comparison := version.Compare(constraintVersion)

	switch c.operator {
	case "=":
		return comparison == 0
	case "!=":
		return comparison != 0
	case "<":
		return comparison < 0
	case "<=":
		return comparison <= 0
	case ">":
		return comparison > 0
	case ">=":
		return comparison >= 0
	default:
		return false
	}
}
