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
	version  *Version
}

// NewVersionRange creates a new NuGet version range from a range string
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

// parseRange parses NuGet range syntax into constraints
func parseRange(e *Ecosystem, rangeStr string) ([]*constraint, error) {
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
			versionStr := strings.TrimSpace(rangeStr[1 : len(rangeStr)-1])
			if versionStr == "" {
				return nil, fmt.Errorf("empty version in exact match: %s", rangeStr)
			}
			version, err := e.NewVersion(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version in exact match: %w", err)
			}
			return []*constraint{{operator: "=", version: version}}, nil
		}

		// Handle inclusive ranges [1.0.0,2.0.0]
		if strings.HasPrefix(rangeStr, "[") && strings.HasSuffix(rangeStr, "]") && strings.Contains(rangeStr, ",") {
			return parseInclusiveRange(e, rangeStr)
		}

		// Handle exclusive ranges (1.0.0,2.0.0)
		if strings.HasPrefix(rangeStr, "(") && strings.HasSuffix(rangeStr, ")") && strings.Contains(rangeStr, ",") {
			return parseExclusiveRange(e, rangeStr)
		}

		// Handle mixed ranges [1.0.0,2.0.0) or (1.0.0,2.0.0]
		if ((strings.HasPrefix(rangeStr, "[") && strings.HasSuffix(rangeStr, ")")) ||
			(strings.HasPrefix(rangeStr, "(") && strings.HasSuffix(rangeStr, "]"))) && strings.Contains(rangeStr, ",") {
			return parseMixedRange(e, rangeStr)
		}
	}

	// Handle multiple constraints separated by commas (NuGet allows comma-separated constraints)
	if strings.Contains(rangeStr, ",") {
		return parseCommaSeparatedConstraints(e, rangeStr)
	}

	// Handle single constraint (minimum version)
	version, err := e.NewVersion(rangeStr)
	if err != nil {
		return nil, fmt.Errorf("invalid version in minimum constraint: %w", err)
	}
	return []*constraint{{operator: ">=", version: version}}, nil
}

// parseInclusiveRange handles inclusive ranges [1.0.0,2.0.0]
func parseInclusiveRange(e *Ecosystem, rangeStr string) ([]*constraint, error) {
	content := rangeStr[1 : len(rangeStr)-1] // Remove [ and ]
	parts := strings.Split(content, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid inclusive range: %s", rangeStr)
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	startVersion, err := e.NewVersion(startStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start version in inclusive range: %w", err)
	}
	
	endVersion, err := e.NewVersion(endStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end version in inclusive range: %w", err)
	}

	return []*constraint{
		{operator: ">=", version: startVersion},
		{operator: "<=", version: endVersion},
	}, nil
}

// parseExclusiveRange handles exclusive ranges (1.0.0,2.0.0)
func parseExclusiveRange(e *Ecosystem, rangeStr string) ([]*constraint, error) {
	content := rangeStr[1 : len(rangeStr)-1] // Remove ( and )
	parts := strings.Split(content, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid exclusive range: %s", rangeStr)
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	startVersion, err := e.NewVersion(startStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start version in exclusive range: %w", err)
	}
	
	endVersion, err := e.NewVersion(endStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end version in exclusive range: %w", err)
	}

	return []*constraint{
		{operator: ">", version: startVersion},
		{operator: "<", version: endVersion},
	}, nil
}

// parseMixedRange handles mixed ranges [1.0.0,2.0.0) or (1.0.0,2.0.0] and unbounded ranges [1.0.0,) or (,2.0.0]
func parseMixedRange(e *Ecosystem, rangeStr string) ([]*constraint, error) {
	content := rangeStr[1 : len(rangeStr)-1] // Remove brackets
	parts := strings.Split(content, ",")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid mixed range: %s", rangeStr)
	}

	startStr := strings.TrimSpace(parts[0])
	endStr := strings.TrimSpace(parts[1])

	// Check if this is actually an unbounded range
	if startStr == "" && endStr != "" {
		// (,2.0.0] or (,2.0.0)
		endVersion, err := e.NewVersion(endStr)
		if err != nil {
			return nil, fmt.Errorf("invalid end version in unbounded range: %w", err)
		}
		if strings.HasSuffix(rangeStr, "]") {
			return []*constraint{{operator: "<=", version: endVersion}}, nil
		} else {
			return []*constraint{{operator: "<", version: endVersion}}, nil
		}
	} else if startStr != "" && endStr == "" {
		// [1.0.0,) or (1.0.0,)
		startVersion, err := e.NewVersion(startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start version in unbounded range: %w", err)
		}
		if strings.HasPrefix(rangeStr, "[") {
			return []*constraint{{operator: ">=", version: startVersion}}, nil
		} else {
			return []*constraint{{operator: ">", version: startVersion}}, nil
		}
	}

	// Both start and end are non-empty, handle as normal mixed range
	startVersion, err := e.NewVersion(startStr)
	if err != nil {
		return nil, fmt.Errorf("invalid start version in mixed range: %w", err)
	}
	
	endVersion, err := e.NewVersion(endStr)
	if err != nil {
		return nil, fmt.Errorf("invalid end version in mixed range: %w", err)
	}

	var constraints []*constraint

	// Start constraint
	if strings.HasPrefix(rangeStr, "[") {
		constraints = append(constraints, &constraint{operator: ">=", version: startVersion})
	} else {
		constraints = append(constraints, &constraint{operator: ">", version: startVersion})
	}

	// End constraint
	if strings.HasSuffix(rangeStr, "]") {
		constraints = append(constraints, &constraint{operator: "<=", version: endVersion})
	} else {
		constraints = append(constraints, &constraint{operator: "<", version: endVersion})
	}

	return constraints, nil
}


// parseCommaSeparatedConstraints handles comma-separated constraints
func parseCommaSeparatedConstraints(e *Ecosystem, rangeStr string) ([]*constraint, error) {
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
		partConstraints, err := parseSingleConstraint(e, part)
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
func parseSingleConstraint(e *Ecosystem, c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			versionStr := strings.TrimSpace(c[len(op):])
			version, err := e.NewVersion(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version in constraint %s: %w", c, err)
			}
			return []*constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to minimum version (>=)
	version, err := e.NewVersion(c)
	if err != nil {
		return nil, fmt.Errorf("invalid version in constraint %s: %w", c, err)
	}
	return []*constraint{{operator: ">=", version: version}}, nil
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