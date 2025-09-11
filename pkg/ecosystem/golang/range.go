package golang

import (
	"fmt"
	"strings"
)

// VersionRange represents a Go module version range
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single version constraint
type constraint struct {
	operator string
	version  string
}

// NewVersionRange creates a new Go module version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	// For Go modules, ranges are typically simple comparisons
	// Common patterns: >=v1.2.3, >v1.2.3, <v2.0.0, <=v1.9.9, v1.2.3
	constraints, err := parseGoRange(rangeStr)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    rangeStr,
	}, nil
}

// parseGoRange parses a Go module version range
func parseGoRange(rangeStr string) ([]*constraint, error) {
	// Handle space-separated constraints
	if strings.Contains(rangeStr, " ") {
		parts := strings.Fields(rangeStr)
		var constraints []*constraint
		for _, part := range parts {
			partConstraints, err := parseSingleGoConstraint(part)
			if err != nil {
				return nil, err
			}
			constraints = append(constraints, partConstraints...)
		}
		return constraints, nil
	}

	return parseSingleGoConstraint(rangeStr)
}

// parseSingleGoConstraint parses a single Go version constraint
func parseSingleGoConstraint(c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			version := strings.TrimSpace(c[len(op):])
			return []*constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to exact match
	return []*constraint{{operator: "=", version: c}}, nil
}

// String returns the string representation of the range
func (gr *VersionRange) String() string {
	return gr.original
}

// Contains checks if a version is within this range
func (gr *VersionRange) Contains(version *Version) bool {
	for _, constraint := range gr.constraints {
		if !constraint.matches(version) {
			return false
		}
	}
	return true
}

// matches checks if a version matches this constraint
func (c *constraint) matches(version *Version) bool {
	e := &Ecosystem{}
	constraintVersion, err := e.NewVersion(c.version)
	if err != nil {
		return false
	}

	switch c.operator {
	case "=", "==":
		return version.Compare(constraintVersion) == 0
	case "!=":
		return version.Compare(constraintVersion) != 0
	case ">":
		return version.Compare(constraintVersion) > 0
	case ">=":
		return version.Compare(constraintVersion) >= 0
	case "<":
		return version.Compare(constraintVersion) < 0
	case "<=":
		return version.Compare(constraintVersion) <= 0
	default:
		return false
	}
}
