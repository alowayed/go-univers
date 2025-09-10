package hex

import (
	"fmt"
	"regexp"
	"strings"
)

type VersionRange struct {
	original    string
	constraints []*constraint
}

type constraint struct {
	operator string
	version  *Version
}

var (
	// Constraint pattern for Hex version constraints
	// Supports: >=, <=, >, <, =, ~> (pessimistic operator)
	constraintPattern = regexp.MustCompile(`^(>=|<=|>|<|=|~>)?(.+)$`)
)

func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	if rangeStr == "" {
		return nil, fmt.Errorf("range string cannot be empty")
	}

	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(rangeStr)
	if trimmed == "" {
		return nil, fmt.Errorf("range string cannot be empty or only whitespace")
	}

	// Parse constraints by splitting on spaces or "and" keywords
	constraints, err := parseConstraints(trimmed, e)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		original:    rangeStr,
		constraints: constraints,
	}, nil
}

func parseConstraints(rangeStr string, ecosystem *Ecosystem) ([]*constraint, error) {
	// Split by spaces and "and" keywords to handle multiple constraints
	parts := strings.Fields(rangeStr)
	if len(parts) == 0 {
		return nil, fmt.Errorf("no constraints found")
	}

	var constraints []*constraint

	for _, part := range parts {
		// Skip "and" keywords
		if strings.ToLower(part) == "and" {
			continue
		}

		constraint, err := parseConstraint(part, ecosystem)
		if err != nil {
			return nil, err
		}

		// Handle pessimistic operator (~>) by converting to range
		if constraint.operator == "~>" {
			pessimisticConstraints := expandPessimisticConstraint(constraint)
			constraints = append(constraints, pessimisticConstraints...)
		} else {
			constraints = append(constraints, constraint)
		}
	}

	return constraints, nil
}

func parseConstraint(constraintStr string, ecosystem *Ecosystem) (*constraint, error) {
	matches := constraintPattern.FindStringSubmatch(constraintStr)
	if matches == nil {
		return nil, fmt.Errorf("invalid constraint format: %s", constraintStr)
	}

	operator := matches[1]
	versionStr := strings.TrimSpace(matches[2])

	// Default operator is "=" (exact match)
	if operator == "" {
		operator = "="
	}

	// Parse the version
	version, err := ecosystem.NewVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version in constraint: %s: %v", versionStr, err)
	}

	return &constraint{
		operator: operator,
		version:  version,
	}, nil
}

// expandPessimisticConstraint converts ~> operator to equivalent >= and < constraints
// Standard pessimistic operator behavior:
// ~> 1.2.3 means >= 1.2.3 and < 1.3.0 (increment minor)
// ~> 1.2.0 means >= 1.2.0 and < 1.3.0 (increment minor)
// ~> 1.14 means >= 1.14.0 and < 2.0.0 (increment major for partial version)
func expandPessimisticConstraint(c *constraint) []*constraint {
	v := c.version

	// Create >= constraint with full version (ensuring patch is present)
	lowerVersion := &Version{
		original: fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch),
		major:    v.major,
		minor:    v.minor,
		patch:    v.patch,
	}
	geConstraint := &constraint{
		operator: ">=",
		version:  lowerVersion,
	}

	// Determine upper bound
	var upperBound *Version
	if strings.Count(v.original, ".") == 1 {
		// Partial version like "1.14" - increment major
		upperBound = &Version{
			original: fmt.Sprintf("%d.0.0", v.major+1),
			major:    v.major + 1,
			minor:    0,
			patch:    0,
		}
	} else {
		// Full version like "1.2.3" or "1.2.0" - increment minor
		upperBound = &Version{
			original: fmt.Sprintf("%d.%d.0", v.major, v.minor+1),
			major:    v.major,
			minor:    v.minor + 1,
			patch:    0,
		}
	}

	ltConstraint := &constraint{
		operator: "<",
		version:  upperBound,
	}

	return []*constraint{geConstraint, ltConstraint}
}

func (r *VersionRange) Contains(version *Version) bool {
	// All constraints must be satisfied
	for _, constraint := range r.constraints {
		if !constraint.matches(version) {
			return false
		}
	}
	return true
}

func (r *VersionRange) String() string {
	return r.original
}

func (c *constraint) matches(version *Version) bool {
	cmp := version.Compare(c.version)

	switch c.operator {
	case "=":
		return cmp == 0
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
