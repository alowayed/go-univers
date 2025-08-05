package nuget

import (
	"fmt"
	"strings"
)

// VersionRange represents a NuGet version range with NuGet-specific syntax support
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single NuGet version constraint
type constraint struct {
	operator string
	version  string
}

// NewVersionRange creates a new NuGet version range from a range string
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

// parseRange parses NuGet range syntax into constraints
func parseRange(rangeStr string) ([]*constraint, error) {
	// Trim whitespace
	rangeStr = strings.TrimSpace(rangeStr)

	// Check for bracket/paren syntax first
	if (strings.HasPrefix(rangeStr, "[") || strings.HasPrefix(rangeStr, "(")) && 
	   (strings.HasSuffix(rangeStr, "]") || strings.HasSuffix(rangeStr, ")")) {
		
		// Check for empty brackets/parens
		if rangeStr == "[]" || rangeStr == "()" {
			return nil, fmt.Errorf("empty range expression: %s", rangeStr)
		}
		
		// Handle exact version match [1.0.0]
		if strings.HasPrefix(rangeStr, "[") && strings.HasSuffix(rangeStr, "]") && !strings.Contains(rangeStr, ",") {
			version := strings.TrimSpace(rangeStr[1 : len(rangeStr)-1])
			if version == "" {
				return nil, fmt.Errorf("empty version in exact match: %s", rangeStr)
			}
			return []*constraint{{operator: "=", version: version}}, nil
		}

		// Handle inclusive ranges [1.0.0,2.0.0]
		if strings.HasPrefix(rangeStr, "[") && strings.HasSuffix(rangeStr, "]") && strings.Contains(rangeStr, ",") {
			return parseInclusiveRange(rangeStr)
		}

		// Handle exclusive ranges (1.0.0,2.0.0)
		if strings.HasPrefix(rangeStr, "(") && strings.HasSuffix(rangeStr, ")") && strings.Contains(rangeStr, ",") {
			return parseExclusiveRange(rangeStr)
		}

		// Handle mixed ranges [1.0.0,2.0.0) or (1.0.0,2.0.0]
		if ((strings.HasPrefix(rangeStr, "[") && strings.HasSuffix(rangeStr, ")")) ||
			(strings.HasPrefix(rangeStr, "(") && strings.HasSuffix(rangeStr, "]"))) && strings.Contains(rangeStr, ",") {
			return parseMixedRange(rangeStr)
		}

		// Handle unbounded ranges [1.0.0,) or (,2.0.0]  
		if strings.Contains(rangeStr, ",") {
			return parseUnboundedRange(rangeStr)
		}
	}

	// Handle multiple constraints separated by commas (NuGet allows comma-separated constraints)
	if strings.Contains(rangeStr, ",") {
		return parseCommaSeparatedConstraints(rangeStr)
	}

	// Handle single constraint (minimum version)
	return []*constraint{{operator: ">=", version: rangeStr}}, nil
}

// parseInclusiveRange handles inclusive ranges [1.0.0,2.0.0]
func parseInclusiveRange(rangeStr string) ([]*constraint, error) {
	content := rangeStr[1 : len(rangeStr)-1] // Remove [ and ]
	parts := strings.Split(content, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid inclusive range: %s", rangeStr)
	}

	start := strings.TrimSpace(parts[0])
	end := strings.TrimSpace(parts[1])

	return []*constraint{
		{operator: ">=", version: start},
		{operator: "<=", version: end},
	}, nil
}

// parseExclusiveRange handles exclusive ranges (1.0.0,2.0.0)
func parseExclusiveRange(rangeStr string) ([]*constraint, error) {
	content := rangeStr[1 : len(rangeStr)-1] // Remove ( and )
	parts := strings.Split(content, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid exclusive range: %s", rangeStr)
	}

	start := strings.TrimSpace(parts[0])
	end := strings.TrimSpace(parts[1])

	return []*constraint{
		{operator: ">", version: start},
		{operator: "<", version: end},
	}, nil
}

// parseMixedRange handles mixed ranges [1.0.0,2.0.0) or (1.0.0,2.0.0]
func parseMixedRange(rangeStr string) ([]*constraint, error) {
	content := rangeStr[1 : len(rangeStr)-1] // Remove brackets
	parts := strings.Split(content, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid mixed range: %s", rangeStr)
	}

	start := strings.TrimSpace(parts[0])
	end := strings.TrimSpace(parts[1])

	// Check if this is actually an unbounded range
	if start == "" && end != "" {
		// (,2.0.0] or (,2.0.0)
		if strings.HasSuffix(rangeStr, "]") {
			return []*constraint{{operator: "<=", version: end}}, nil
		} else {
			return []*constraint{{operator: "<", version: end}}, nil
		}
	} else if start != "" && end == "" {
		// [1.0.0,) or (1.0.0,)
		if strings.HasPrefix(rangeStr, "[") {
			return []*constraint{{operator: ">=", version: start}}, nil
		} else {
			return []*constraint{{operator: ">", version: start}}, nil
		}
	}

	// Both start and end are non-empty, handle as normal mixed range
	var constraints []*constraint

	// Start constraint
	if strings.HasPrefix(rangeStr, "[") {
		constraints = append(constraints, &constraint{operator: ">=", version: start})
	} else {
		constraints = append(constraints, &constraint{operator: ">", version: start})
	}

	// End constraint
	if strings.HasSuffix(rangeStr, "]") {
		constraints = append(constraints, &constraint{operator: "<=", version: end})
	} else {
		constraints = append(constraints, &constraint{operator: "<", version: end})
	}

	return constraints, nil
}

// parseUnboundedRange handles unbounded ranges [1.0.0,) or (,2.0.0]
func parseUnboundedRange(rangeStr string) ([]*constraint, error) {
	if len(rangeStr) < 3 {
		return nil, fmt.Errorf("invalid unbounded range: %s", rangeStr)
	}
	
	content := rangeStr[1 : len(rangeStr)-1] // Remove brackets
	parts := strings.Split(content, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid unbounded range: %s", rangeStr)
	}

	start := strings.TrimSpace(parts[0])
	end := strings.TrimSpace(parts[1])

	if start == "" && end != "" {
		// (,2.0.0] or (,2.0.0)
		if strings.HasSuffix(rangeStr, "]") {
			return []*constraint{{operator: "<=", version: end}}, nil
		} else {
			return []*constraint{{operator: "<", version: end}}, nil
		}
	} else if start != "" && end == "" {
		// [1.0.0,) or (1.0.0,)
		if strings.HasPrefix(rangeStr, "[") {
			return []*constraint{{operator: ">=", version: start}}, nil
		} else {
			return []*constraint{{operator: ">", version: start}}, nil
		}
	}

	return nil, fmt.Errorf("invalid unbounded range: %s", rangeStr)
}

// parseCommaSeparatedConstraints handles comma-separated constraints
func parseCommaSeparatedConstraints(rangeStr string) ([]*constraint, error) {
	// Reject malformed bracket/paren expressions that fall through to here
	if (strings.HasPrefix(rangeStr, "[") && !strings.HasSuffix(rangeStr, "]")) ||
	   (strings.HasPrefix(rangeStr, "(") && !strings.HasSuffix(rangeStr, ")")) ||
	   rangeStr == "[]" || rangeStr == "()" {
		return nil, fmt.Errorf("malformed range expression: %s", rangeStr)
	}
	
	parts := strings.Split(rangeStr, ",")
	var constraints []*constraint

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Parse each part as a single constraint
		partConstraints, err := parseSingleConstraint(part)
		if err != nil {
			return nil, err
		}
		constraints = append(constraints, partConstraints...)
	}

	if len(constraints) == 0 {
		return nil, fmt.Errorf("no valid constraints found in: %s", rangeStr)
	}

	return constraints, nil
}

// parseSingleConstraint parses a single NuGet constraint
func parseSingleConstraint(c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			version := strings.TrimSpace(c[len(op):])
			return []*constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to minimum version (>=)
	return []*constraint{{operator: ">=", version: c}}, nil
}

// String returns the string representation of the range
func (nr *VersionRange) String() string {
	return nr.original
}

// Contains checks if a version is within this range
func (nr *VersionRange) Contains(version *Version) bool {
	// AND logic: ALL constraints must be satisfied
	for _, constraint := range nr.constraints {
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