package cran

import (
	"fmt"
	"strings"
)

// VersionRange represents a CRAN version range with CRAN-specific syntax support
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single CRAN version constraint
type constraint struct {
	operator string
	version  string
}

// NewVersionRange creates a new CRAN version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	original := rangeStr
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	constraints, err := parseConstraints(rangeStr)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    original,
	}, nil
}

// parseConstraints parses CRAN constraint syntax
func parseConstraints(rangeStr string) ([]*constraint, error) {
	// Handle multiple constraints separated by comma (AND logic)
	parts := strings.Split(rangeStr, ",")
	var constraints []*constraint

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		constraint, err := parseConstraint(part)
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

// parseConstraint parses a single constraint
func parseConstraint(constraintStr string) (*constraint, error) {
	constraintStr = strings.TrimSpace(constraintStr)

	// CRAN supports standard comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(constraintStr, op) {
			version := strings.TrimSpace(constraintStr[len(op):])
			if version == "" {
				return nil, fmt.Errorf("constraint %s requires version", op)
			}
			// Validate that the version part is valid
			ecosystem := &Ecosystem{}
			if _, err := ecosystem.NewVersion(version); err != nil {
				return nil, fmt.Errorf("invalid version in constraint %s: %w", constraintStr, err)
			}
			return &constraint{operator: op, version: version}, nil
		}
	}

	// Default to exact match - validate the version
	ecosystem := &Ecosystem{}
	if _, err := ecosystem.NewVersion(constraintStr); err != nil {
		return nil, fmt.Errorf("invalid version in constraint %s: %w", constraintStr, err)
	}
	return &constraint{operator: "=", version: constraintStr}, nil
}

// String returns the string representation of the version range
func (vr *VersionRange) String() string {
	return vr.original
}

// Contains checks if a version satisfies this range
func (vr *VersionRange) Contains(version *Version) bool {
	ecosystem := &Ecosystem{}

	// All constraints must be satisfied (AND logic)
	for _, c := range vr.constraints {
		if !satisfiesConstraint(version, c, ecosystem) {
			return false
		}
	}

	return true
}

// satisfiesConstraint checks if a version satisfies a single constraint
func satisfiesConstraint(version *Version, c *constraint, ecosystem *Ecosystem) bool {
	constraintVersion, err := ecosystem.NewVersion(c.version)
	if err != nil {
		return false
	}

	cmp := version.Compare(constraintVersion)

	switch c.operator {
	case "=":
		return cmp == 0
	case "!=":
		return cmp != 0
	case ">":
		return cmp > 0
	case ">=":
		return cmp >= 0
	case "<":
		return cmp < 0
	case "<=":
		return cmp <= 0
	default:
		return false
	}
}
