package rpm

import (
	"fmt"
	"strings"
)

// VersionRange represents an RPM version range with standard comparison operators
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single RPM version constraint
type constraint struct {
	operator string
	version  *Version
}

// NewVersionRange creates a new RPM version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	original := rangeStr
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	constraints, err := parseRPMConstraints(e, rangeStr)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    original,
	}, nil
}

// parseRPMConstraints parses RPM constraint syntax
func parseRPMConstraints(e *Ecosystem, rangeStr string) ([]*constraint, error) {
	// Handle multiple constraints separated by spaces, commas, or both (AND logic)
	// Normalize by replacing commas with spaces, then split by whitespace
	parts := strings.Fields(strings.ReplaceAll(rangeStr, ",", " "))

	var constraints []*constraint

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		constraint, err := parseRPMConstraint(e, part)
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

// parseRPMConstraint parses a single constraint
func parseRPMConstraint(e *Ecosystem, constraintStr string) (*constraint, error) {
	constraintStr = strings.TrimSpace(constraintStr)

	// RPM supports standard comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(constraintStr, op) {
			versionStr := strings.TrimSpace(constraintStr[len(op):])
			if versionStr == "" {
				return nil, fmt.Errorf("constraint %s requires version", op)
			}

			// Parse the version string once and store it
			version, err := e.NewVersion(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version in constraint: %v", err)
			}

			return &constraint{operator: op, version: version}, nil
		}
	}

	// Default to exact match - parse the version
	version, err := e.NewVersion(constraintStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version in constraint: %v", err)
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
		if !satisfiesRPMConstraint(version, c) {
			return false
		}
	}

	return true
}

// satisfiesRPMConstraint checks if a version satisfies a single constraint
func satisfiesRPMConstraint(version *Version, c *constraint) bool {
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
