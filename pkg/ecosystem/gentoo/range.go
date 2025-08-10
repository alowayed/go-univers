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
	version  *Version
}

// NewVersionRange creates a new Gentoo version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	constraints, err := parseRange(e, rangeStr)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    rangeStr,
	}, nil
}

// parseRange parses Gentoo range syntax into constraints
func parseRange(e *Ecosystem, rangeStr string) ([]*constraint, error) {
	rangeStr = strings.TrimSpace(rangeStr)

	// Normalize separators by replacing commas with spaces, then split.
	// This handles both ">=1.0, <2.0" and ">=1.0 <2.0".
	normalized := strings.ReplaceAll(rangeStr, ",", " ")
	parts := strings.Fields(normalized)

	if len(parts) <= 1 {
		// This also handles single versions like "1.2.3" correctly.
		return parseSingleConstraint(e, rangeStr)
	}

	var constraints []*constraint
	for _, part := range parts {
		partConstraints, err := parseSingleConstraint(e, part)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, partConstraints...)
	}

	return constraints, nil
}

// parseSingleConstraint parses a single Gentoo constraint
func parseSingleConstraint(e *Ecosystem, c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			versionStr := strings.TrimSpace(c[len(op):])
			if versionStr == "" {
				return nil, fmt.Errorf("missing version after operator %s", op)
			}
			version, err := e.NewVersion(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version %s: %w", versionStr, err)
			}
			return []*constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to exact match
	version, err := e.NewVersion(c)
	if err != nil {
		return nil, fmt.Errorf("invalid version %s: %w", c, err)
	}
	return []*constraint{{operator: "=", version: version}}, nil
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
