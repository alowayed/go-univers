package npm

import (
	"fmt"
	"strconv"
	"strings"
)

// VersionRange represents an NPM version range with NPM-specific syntax support
type VersionRange struct {
	constraintGroups [][]*constraint // OR logic between groups, AND logic within groups
	original         string
}

// constraint represents a single NPM version constraint
type constraint struct {
	operator string
	version  string
}

// NewVersionRange creates a new NPM version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	constraintGroups, err := parseRangeGroups(rangeStr)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraintGroups: constraintGroups,
		original:         rangeStr,
	}, nil
}

// parseRangeGroups parses NPM range syntax into constraint groups for OR logic
func parseRangeGroups(rangeStr string) ([][]*constraint, error) {
	// Handle OR logic (||) - each OR'd part becomes a separate group
	if strings.Contains(rangeStr, "||") {
		parts := strings.Split(rangeStr, "||")
		var constraintGroups [][]*constraint
		for _, part := range parts {
			constraints, err := parseRange(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			constraintGroups = append(constraintGroups, constraints)
		}
		return constraintGroups, nil
	}

	// Single group (no OR logic)
	constraints, err := parseRange(rangeStr)
	if err != nil {
		return nil, err
	}
	return [][]*constraint{constraints}, nil
}

// parseRange parses NPM range syntax into constraints
func parseRange(rangeStr string) ([]*constraint, error) {
	// Trim whitespace and remove parentheses
	rangeStr = strings.TrimSpace(rangeStr)
	rangeStr = strings.TrimPrefix(rangeStr, "(")
	rangeStr = strings.TrimSuffix(rangeStr, ")")

	// Handle hyphen ranges (1.2.3 - 2.3.4)
	// Also catch malformed hyphen ranges like "1.2.3 -" or "- 1.2.3"
	if strings.Contains(rangeStr, " - ") || strings.HasSuffix(rangeStr, " -") || strings.HasPrefix(rangeStr, "- ") {
		return parseHyphenRange(rangeStr)
	}

	// Handle space-separated constraints (>=1.0.0 <2.0.0)
	if strings.Contains(rangeStr, " ") && !strings.HasPrefix(rangeStr, "^") && !strings.HasPrefix(rangeStr, "~") {
		return parseSpaceSeparatedConstraints(rangeStr)
	}

	// Handle single constraint
	return parseSingleConstraint(rangeStr)
}

// parseSingleConstraint parses a single NPM constraint
func parseSingleConstraint(c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Check for invalid characters
	if strings.ContainsAny(c, "@#$%&!") {
		return nil, fmt.Errorf("invalid characters in constraint: %s", c)
	}

	// Handle wildcard
	if c == "*" {
		return []*constraint{{operator: "*", version: "*"}}, nil
	}

	// Handle caret range (^1.2.3)
	if strings.HasPrefix(c, "^") {
		return parseCaretRange(c[1:])
	}

	// Handle tilde range (~1.2.3)
	if strings.HasPrefix(c, "~") {
		return parseTildeRange(c[1:])
	}

	// Handle x-range (1.x, 1.2.x)
	if strings.Contains(c, "x") || strings.Contains(c, "X") {
		return parseXRange(c)
	}

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

// parseCaretRange handles caret ranges (^1.2.3)
func parseCaretRange(version string) ([]*constraint, error) {
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		return nil, err
	}

	// Special rules for caret ranges with zero versions
	if v.major == 0 {
		if v.minor == 0 {
			// ^0.0.3 means >=0.0.3 <0.0.4 (only patch changes)
			return []*constraint{
				{operator: ">=", version: v.normalize()},
				{operator: "<", version: fmt.Sprintf("0.0.%d", v.patch+1)},
			}, nil
		}
		// ^0.2.3 means >=0.2.3 <0.3.0-0 (patch and minor changes, excludes prereleases from next minor)
		return []*constraint{
			{operator: ">=", version: v.normalize()},
			{operator: "<", version: fmt.Sprintf("0.%d.0-0", v.minor+1)},
		}, nil
	}

	// ^1.2.3 means >=1.2.3 <2.0.0-0 (excludes prereleases from next major)
	return []*constraint{
		{operator: ">=", version: v.normalize()},
		{operator: "<", version: fmt.Sprintf("%d.0.0-0", v.major+1)},
	}, nil
}

// parseTildeRange handles tilde ranges (~1.2.3)
func parseTildeRange(version string) ([]*constraint, error) {
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		return nil, err
	}

	// ~1.2.3 means >=1.2.3 <1.3.0-0 (excludes prereleases from next minor)
	return []*constraint{
		{operator: ">=", version: v.normalize()},
		{operator: "<", version: fmt.Sprintf("%d.%d.0-0", v.major, v.minor+1)},
	}, nil
}

// parseXRange handles x-ranges (1.x, 1.2.x)
func parseXRange(rangeStr string) ([]*constraint, error) {
	parts := strings.Split(rangeStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid x-range: %s", rangeStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version in x-range: %s", parts[0])
	}

	// 1.x means >=1.0.0-0 <2.0.0-0 (includes prereleases in range, excludes prereleases from next major)
	if len(parts) == 2 && (parts[1] == "x" || parts[1] == "X") {
		return []*constraint{
			{operator: ">=", version: fmt.Sprintf("%d.0.0-0", major)},
			{operator: "<", version: fmt.Sprintf("%d.0.0-0", major+1)},
		}, nil
	}

	// 1.2.x means >=1.2.0-0 <1.3.0-0 (includes prereleases in range, excludes prereleases from next minor)
	if len(parts) == 3 && (parts[2] == "x" || parts[2] == "X") {
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version in x-range: %s", parts[1])
		}
		return []*constraint{
			{operator: ">=", version: fmt.Sprintf("%d.%d.0-0", major, minor)},
			{operator: "<", version: fmt.Sprintf("%d.%d.0-0", major, minor+1)},
		}, nil
	}

	return nil, fmt.Errorf("unsupported x-range format: %s", rangeStr)
}

// parseHyphenRange handles hyphen ranges (1.2.3 - 2.3.4)
func parseHyphenRange(rangeStr string) ([]*constraint, error) {
	parts := strings.Split(rangeStr, " - ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid hyphen range: %s", rangeStr)
	}

	start := strings.TrimSpace(parts[0])
	end := strings.TrimSpace(parts[1])

	// Both parts must be non-empty and valid versions
	if start == "" || end == "" {
		return nil, fmt.Errorf("invalid hyphen range: %s", rangeStr)
	}

	// Validate that both parts are valid versions
	e := &Ecosystem{}
	if _, err := e.NewVersion(start); err != nil {
		return nil, fmt.Errorf("invalid start version in hyphen range: %s", start)
	}
	if _, err := e.NewVersion(end); err != nil {
		return nil, fmt.Errorf("invalid end version in hyphen range: %s", end)
	}

	return []*constraint{
		{operator: ">=", version: start},
		{operator: "<=", version: end},
	}, nil
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
func (nr *VersionRange) String() string {
	return nr.original
}

// Contains checks if a version is within this range
func (nr *VersionRange) Contains(version *Version) bool {
	// OR logic between groups: if ANY group is satisfied, return true
	for _, constraintGroup := range nr.constraintGroups {
		// AND logic within group: ALL constraints in this group must be satisfied
		groupSatisfied := true
		for _, constraint := range constraintGroup {
			if !constraint.matches(version) {
				groupSatisfied = false
				break
			}
		}
		if groupSatisfied {
			return true
		}
	}
	return false
}

// matches checks if the given version matches this constraint
func (c *constraint) matches(version *Version) bool {
	if c.operator == "*" {
		return true
	}

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
