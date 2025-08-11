package semver

import (
	"fmt"
	"strings"
)

// VersionRange represents a SemVer version range with standard comparison operators
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single SemVer version constraint
type constraint struct {
	operator string
	version  *Version
}

// NewVersionRange creates a new SemVer version range from a range string
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

// parseRange parses SemVer range syntax into constraints
func parseRange(rangeStr string) ([]*constraint, error) {
	// Handle comma-separated constraints (>=1.0.0,<2.0.0)
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

// parseSingleConstraint parses a single SemVer constraint
func parseSingleConstraint(c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Handle wildcard
	if c == "*" {
		return []*constraint{{operator: "*", version: nil}}, nil
	}

	// Handle comparison operators (order matters - check longer operators first)
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			versionStr := strings.TrimSpace(c[len(op):])
			if versionStr == "" {
				return nil, fmt.Errorf("missing version after operator %s", op)
			}

			e := &Ecosystem{}
			version, err := e.NewVersion(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version %s: %v", versionStr, err)
			}

			return []*constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to exact match
	e := &Ecosystem{}
	version, err := e.NewVersion(c)
	if err != nil {
		return nil, fmt.Errorf("invalid version %s: %v", c, err)
	}

	return []*constraint{{operator: "=", version: version}}, nil
}

// parseCommaSeparatedConstraints handles comma-separated constraints (>=1.0.0,<2.0.0)
func parseCommaSeparatedConstraints(rangeStr string) ([]*constraint, error) {
	parts := strings.Split(rangeStr, ",")
	var constraints []*constraint

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		partConstraints, err := parseSingleConstraint(part)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, partConstraints...)
	}

	if len(constraints) == 0 {
		return nil, fmt.Errorf("no valid constraints found")
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
func (sr *VersionRange) String() string {
	return sr.original
}

// Contains checks if a version is within this range
func (sr *VersionRange) Contains(version *Version) bool {
	// ALL constraints must be satisfied (AND logic)
	for _, constraint := range sr.constraints {
		if !constraint.matches(version) {
			return false
		}
	}
	return true
}

// matches checks if the given version matches this constraint
func (c *constraint) matches(version *Version) bool {
	// Wildcard matches everything
	if c.operator == "*" {
		return true
	}

	if c.version == nil {
		return false
	}

	comparison := version.Compare(c.version)

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
