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
	version  *Version
}

// NewVersionRange creates a new CRAN version range from a range string
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

// parseConstraints parses CRAN constraint syntax
func parseConstraints(rangeStr string, e *Ecosystem) ([]*constraint, error) {
	// Handle multiple constraints separated by comma (AND logic)
	parts := strings.Split(rangeStr, ",")
	var constraints []*constraint

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		constraint, err := parseConstraint(part, e)
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
func parseConstraint(constraintStr string, e *Ecosystem) (*constraint, error) {
	constraintStr = strings.TrimSpace(constraintStr)

	// CRAN supports standard comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(constraintStr, op) {
			versionStr := strings.TrimSpace(constraintStr[len(op):])
			if versionStr == "" {
				return nil, fmt.Errorf("constraint %s requires version", op)
			}
			// Parse and store the version object
			version, err := e.NewVersion(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version in constraint %s: %w", constraintStr, err)
			}
			return &constraint{operator: op, version: version}, nil
		}
	}

	// Default to exact match - parse and store the version
	version, err := e.NewVersion(constraintStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version in constraint %s: %w", constraintStr, err)
	}
	return &constraint{operator: "=", version: version}, nil
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
	cmp := version.Compare(c.version)

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
