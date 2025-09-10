package mattermost

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
	// Constraint pattern for individual Mattermost version constraints
	constraintPattern = regexp.MustCompile(`^(>=|<=|>|<|=)?(.+)$`)
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

	// Parse constraints by splitting on spaces
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
	// Split by spaces to handle multiple constraints
	parts := strings.Fields(rangeStr)
	if len(parts) == 0 {
		return nil, fmt.Errorf("no constraints found")
	}

	var constraints []*constraint

	for _, part := range parts {
		constraint, err := parseConstraint(part, ecosystem)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, constraint)
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
