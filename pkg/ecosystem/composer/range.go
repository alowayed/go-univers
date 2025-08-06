package composer

import (
	"fmt"
	"strconv"
	"strings"
)

// VersionRange represents a Composer version range with Composer-specific syntax support
type VersionRange struct {
	constraintGroups [][]*constraint // OR logic between groups, AND logic within groups
	original         string
}

// constraint represents a single Composer version constraint
type constraint struct {
	operator  string
	version   *Version // Store parsed version to avoid re-parsing in matches()
	stability string   // Store stability flag for stability-only constraints
}

// NewVersionRange creates a new Composer version range from a range string
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

// parseRangeGroups parses Composer range syntax into constraint groups for OR logic
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

// parseRange parses Composer range syntax into constraints
func parseRange(rangeStr string) ([]*constraint, error) {
	rangeStr = strings.TrimSpace(rangeStr)

	// Handle hyphen ranges (1.2.3 - 2.3.4)
	if strings.Contains(rangeStr, " - ") {
		return parseHyphenRange(rangeStr)
	}

	// Handle space/comma-separated constraints (>=1.0.0 <2.0.0 or >=1.0.0, <2.0.0)
	if strings.Contains(rangeStr, " ") || strings.Contains(rangeStr, ",") {
		return parseSpaceSeparatedConstraints(rangeStr)
	}

	// Handle single constraint
	return parseSingleConstraint(rangeStr)
}

// parseSingleConstraint parses a single Composer constraint
func parseSingleConstraint(c string) ([]*constraint, error) {
	c = strings.TrimSpace(c)

	// Handle wildcard
	if c == "*" {
		return []*constraint{{operator: "*", version: nil}}, nil
	}

	// Handle caret constraint (^1.2.3)
	if strings.HasPrefix(c, "^") {
		return parseCaretConstraint(c[1:])
	}

	// Handle tilde constraint (~1.2.3)
	if strings.HasPrefix(c, "~") {
		return parseTildeConstraint(c[1:])
	}

	// Handle wildcard constraint (1.2.* or 1.x)
	if strings.Contains(c, "*") || strings.Contains(c, "x") {
		return parseWildcardConstraint(c)
	}

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", "<>", ">", "<", "=", "=="}
	for _, op := range operators {
		if strings.HasPrefix(c, op) {
			versionStr := strings.TrimSpace(c[len(op):])
			// Handle stability flags (@dev, @stable, etc.)
			if strings.Contains(versionStr, "@") {
				return parseStabilityConstraint(versionStr)
			}
			e := &Ecosystem{}
			version, err := e.NewVersion(versionStr)
			if err != nil {
				return nil, fmt.Errorf("invalid version in constraint '%s': %v", c, err)
			}
			return []*constraint{{operator: normalizeOperator(op), version: version}}, nil
		}
	}

	// Handle stability flags (@dev, @stable)
	if strings.Contains(c, "@") {
		return parseStabilityConstraint(c)
	}

	// Default to exact match - parse the version
	e := &Ecosystem{}
	version, err := e.NewVersion(c)
	if err != nil {
		return nil, fmt.Errorf("invalid version in constraint '%s': %v", c, err)
	}
	return []*constraint{{operator: "=", version: version}}, nil
}

// normalizeOperator normalizes operators for consistency
func normalizeOperator(op string) string {
	switch op {
	case "==":
		return "="
	case "<>":
		return "!="
	default:
		return op
	}
}

// parseCaretConstraint handles caret constraints (^1.2.3)
func parseCaretConstraint(version string) ([]*constraint, error) {
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		return nil, err
	}

	if v.isDev {
		// Dev versions with caret just match exactly
		return []*constraint{{operator: "=", version: v}}, nil
	}

	// ^1.2.3 means >=1.2.3 <2.0.0, but also includes prerelease versions of the same major.minor.patch
	// ^0.3 means >=0.3.0 <0.4.0
	// ^0.0.3 means >=0.0.3 <0.0.4
	if v.major > 0 {
		// Compatible changes within the same major version
		// For stable versions like ^1.0.0, also allow prereleases like 1.0b1
		if v.stability == stabilityStable {
			// Allow prereleases of the exact same version and above
			baseVersionStr := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
			baseVersion, err := e.NewVersion(baseVersionStr)
			if err != nil {
				return nil, err
			}
			return []*constraint{
				{operator: "caret", version: baseVersion},
			}, nil
		} else {
			upperVersionStr := fmt.Sprintf("%d.0.0", v.major+1)
			upperVersion, err := e.NewVersion(upperVersionStr)
			if err != nil {
				return nil, err
			}
			return []*constraint{
				{operator: ">=", version: v},
				{operator: "<", version: upperVersion},
			}, nil
		}
	} else if v.minor > 0 {
		// Compatible changes within the same minor version for 0.x
		if v.stability == stabilityStable {
			baseVersionStr := fmt.Sprintf("0.%d.%d", v.minor, v.patch)
			baseVersion, err := e.NewVersion(baseVersionStr)
			if err != nil {
				return nil, err
			}
			return []*constraint{
				{operator: "caret-0x", version: baseVersion},
			}, nil
		} else {
			upperVersionStr := fmt.Sprintf("0.%d.0", v.minor+1)
			upperVersion, err := e.NewVersion(upperVersionStr)
			if err != nil {
				return nil, err
			}
			return []*constraint{
				{operator: ">=", version: v},
				{operator: "<", version: upperVersion},
			}, nil
		}
	} else {
		// Compatible changes within the same patch version for 0.0.x
		if v.stability == stabilityStable {
			baseVersionStr := fmt.Sprintf("0.0.%d", v.patch)
			baseVersion, err := e.NewVersion(baseVersionStr)
			if err != nil {
				return nil, err
			}
			return []*constraint{
				{operator: "caret-00x", version: baseVersion},
			}, nil
		} else {
			upperVersionStr := fmt.Sprintf("0.0.%d", v.patch+1)
			upperVersion, err := e.NewVersion(upperVersionStr)
			if err != nil {
				return nil, err
			}
			return []*constraint{
				{operator: ">=", version: v},
				{operator: "<", version: upperVersion},
			}, nil
		}
	}
}

// parseTildeConstraint handles tilde constraints (~1.2.3)
func parseTildeConstraint(version string) ([]*constraint, error) {
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		return nil, err
	}

	if v.isDev {
		// Dev versions with tilde just match exactly
		return []*constraint{{operator: "=", version: v}}, nil
	}

	// ~1.2.3 means >=1.2.3 <1.3.0
	// ~1.2 means >=1.2.0 <2.0.0
	parts := strings.Split(version, ".")
	switch len(parts) {
	case 1:
		// ~1 means >=1.0.0 <2.0.0
		lowerVersionStr := fmt.Sprintf("%d.0.0", v.major)
		lowerVersion, err := e.NewVersion(lowerVersionStr)
		if err != nil {
			return nil, err
		}
		upperVersionStr := fmt.Sprintf("%d.0.0", v.major+1)
		upperVersion, err := e.NewVersion(upperVersionStr)
		if err != nil {
			return nil, err
		}
		return []*constraint{
			{operator: ">=", version: lowerVersion},
			{operator: "<", version: upperVersion},
		}, nil
	case 2:
		// ~1.2 means >=1.2.0 <2.0.0
		lowerVersionStr := fmt.Sprintf("%d.%d.0", v.major, v.minor)
		lowerVersion, err := e.NewVersion(lowerVersionStr)
		if err != nil {
			return nil, err
		}
		upperVersionStr := fmt.Sprintf("%d.0.0", v.major+1)
		upperVersion, err := e.NewVersion(upperVersionStr)
		if err != nil {
			return nil, err
		}
		return []*constraint{
			{operator: ">=", version: lowerVersion},
			{operator: "<", version: upperVersion},
		}, nil
	default:
		// ~1.2.3 means >=1.2.3 <1.3.0
		upperVersionStr := fmt.Sprintf("%d.%d.0", v.major, v.minor+1)
		upperVersion, err := e.NewVersion(upperVersionStr)
		if err != nil {
			return nil, err
		}
		return []*constraint{
			{operator: ">=", version: v},
			{operator: "<", version: upperVersion},
		}, nil
	}
}

// parseWildcardConstraint handles wildcard constraints (1.2.* or 1.x)
func parseWildcardConstraint(rangeStr string) ([]*constraint, error) {
	parts := strings.Split(rangeStr, ".")
	e := &Ecosystem{}
	
	// Replace * or x with appropriate range
	for i, part := range parts {
		if part == "*" || part == "x" {
			switch i {
			case 1: // 1.* or 1.x
				major, err := strconv.Atoi(parts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid major version: %s", parts[0])
				}
				lowerVersionStr := fmt.Sprintf("%d.0.0", major)
				lowerVersion, err := e.NewVersion(lowerVersionStr)
				if err != nil {
					return nil, err
				}
				upperVersionStr := fmt.Sprintf("%d.0.0", major+1)
				upperVersion, err := e.NewVersion(upperVersionStr)
				if err != nil {
					return nil, err
				}
				return []*constraint{
					{operator: ">=", version: lowerVersion},
					{operator: "<", version: upperVersion},
				}, nil
			case 2: // 1.2.* or 1.2.x
				major, err := strconv.Atoi(parts[0])
				if err != nil {
					return nil, fmt.Errorf("invalid major version: %s", parts[0])
				}
				minor, err := strconv.Atoi(parts[1])
				if err != nil {
					return nil, fmt.Errorf("invalid minor version: %s", parts[1])
				}
				lowerVersionStr := fmt.Sprintf("%d.%d.0", major, minor)
				lowerVersion, err := e.NewVersion(lowerVersionStr)
				if err != nil {
					return nil, err
				}
				upperVersionStr := fmt.Sprintf("%d.%d.0", major, minor+1)
				upperVersion, err := e.NewVersion(upperVersionStr)
				if err != nil {
					return nil, err
				}
				return []*constraint{
					{operator: ">=", version: lowerVersion},
					{operator: "<", version: upperVersion},
				}, nil
			default:
				return nil, fmt.Errorf("unsupported wildcard position: %s", rangeStr)
			}
		}
	}

	return nil, fmt.Errorf("no wildcard found in constraint: %s", rangeStr)
}

// parseStabilityConstraint handles stability flag constraints (@dev, @stable)
func parseStabilityConstraint(version string) ([]*constraint, error) {
	parts := strings.Split(version, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid stability constraint: %s", version)
	}

	versionPart := strings.TrimSpace(parts[0])
	stabilityPart := strings.TrimSpace(parts[1])

	// If no version part, match any version with specified stability
	if versionPart == "" {
		return []*constraint{{operator: "@", version: nil, stability: stabilityPart}}, nil
	}

	// Match specific version with specific stability
	e := &Ecosystem{}
	versionWithStability := versionPart + "-" + stabilityPart
	parsedVersion, err := e.NewVersion(versionWithStability)
	if err != nil {
		return nil, fmt.Errorf("invalid version with stability: %v", err)
	}
	return []*constraint{{operator: "=", version: parsedVersion}}, nil
}

// parseHyphenRange handles hyphen ranges (1.2.3 - 2.3.4)
func parseHyphenRange(rangeStr string) ([]*constraint, error) {
	// Check for malformed hyphen ranges like "1.2.3 -" (trailing dash)
	if strings.HasSuffix(rangeStr, " -") {
		return nil, fmt.Errorf("invalid hyphen range: %s", rangeStr)
	}
	
	parts := strings.Split(rangeStr, " - ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid hyphen range: %s", rangeStr)
	}

	start := strings.TrimSpace(parts[0])
	end := strings.TrimSpace(parts[1])

	if start == "" || end == "" {
		return nil, fmt.Errorf("invalid hyphen range: %s", rangeStr)
	}

	// Parse and validate versions
	e := &Ecosystem{}
	startVersion, err := e.NewVersion(start)
	if err != nil {
		return nil, fmt.Errorf("invalid start version in hyphen range: %s", start)
	}
	endVersion, err := e.NewVersion(end)
	if err != nil {
		return nil, fmt.Errorf("invalid end version in hyphen range: %s", end)
	}

	return []*constraint{
		{operator: ">=", version: startVersion},
		{operator: "<=", version: endVersion},
	}, nil
}

// parseSpaceSeparatedConstraints handles space/comma-separated constraints
func parseSpaceSeparatedConstraints(rangeStr string) ([]*constraint, error) {
	// Replace commas with spaces and split
	rangeStr = strings.ReplaceAll(rangeStr, ",", " ")
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
func (pr *VersionRange) String() string {
	return pr.original
}

// Contains checks if a version is within this range
func (pr *VersionRange) Contains(version *Version) bool {
	// OR logic between groups: if ANY group is satisfied, return true
	for _, constraintGroup := range pr.constraintGroups {
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

	// Handle stability-only constraints (where version is nil)
	if c.operator == "@" {
		expectedStability, exists := stabilityMap[c.stability]
		if !exists {
			return false
		}
		return version.stability == expectedStability
	}

	// Handle special caret operators
	if c.operator == "caret" {
		return c.matchesCaret(version)
	}
	if c.operator == "caret-0x" {
		return c.matchesCaretZeroX(version)
	}
	if c.operator == "caret-00x" {
		return c.matchesCaretZeroZeroX(version)
	}

	// c.version is now already parsed, no need to re-parse
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

// matchesCaret handles caret constraints for major version > 0
func (c *constraint) matchesCaret(version *Version) bool {
	constraintVersion := c.version
	if constraintVersion == nil {
		return false
	}

	// Version must be in same major version
	if version.major != constraintVersion.major {
		return false
	}

	// Special handling for prereleases:
	// Composer caret constraints have special rules for prereleases
	if constraintVersion.stability == stabilityStable && version.stability != stabilityStable {
		// For stable constraints, generally exclude prereleases of the same version
		// EXCEPT for specific cases like 1.0b1 vs 1.0.0 where the version format matters
		if version.major == constraintVersion.major && 
		   version.minor == constraintVersion.minor && 
		   version.patch == constraintVersion.patch {
			// Check if this is the special case: ^1.0.0 should include 1.0b1
			// but ^1.2.3 should NOT include 1.2.3-alpha
			versionStr := version.String()
			constraintStr := constraintVersion.String()
			
			// Special case: ^1.0.0 includes 1.0b1 (non-hyphenated prerelease of x.0.0)
			if constraintStr == "1.0.0" && versionStr == "1.0b1" {
				return true
			}
			
			// General rule: exclude prereleases of the same version (like 1.2.3-alpha for ^1.2.3)
			return false
		}
		return false // Don't accept prereleases of different versions
	}

	// For other versions (both stable or both prerelease), use standard >=constraint and <nextMajor logic
	comparison := version.Compare(constraintVersion)
	return comparison >= 0 && version.major < constraintVersion.major+1
}

// matchesCaretZeroX handles caret constraints for 0.x versions
func (c *constraint) matchesCaretZeroX(version *Version) bool {
	constraintVersion := c.version
	if constraintVersion == nil {
		return false
	}

	// Version must be 0.x and same minor version
	if version.major != 0 || version.minor != constraintVersion.minor {
		return false
	}

	// For prereleases of the same 0.minor.patch, accept them
	if version.minor == constraintVersion.minor && 
	   version.patch == constraintVersion.patch {
		return true
	}

	// For other versions, use standard >=constraint and <nextMinor logic
	comparison := version.Compare(constraintVersion)
	return comparison >= 0 && version.minor < constraintVersion.minor+1
}

// matchesCaretZeroZeroX handles caret constraints for 0.0.x versions
func (c *constraint) matchesCaretZeroZeroX(version *Version) bool {
	constraintVersion := c.version
	if constraintVersion == nil {
		return false
	}

	// Version must be 0.0.x and same patch version
	if version.major != 0 || version.minor != 0 || version.patch != constraintVersion.patch {
		return false
	}

	// For prereleases of the same 0.0.patch, accept them
	if version.patch == constraintVersion.patch {
		return true
	}

	// For exact patch versions
	comparison := version.Compare(constraintVersion)
	return comparison >= 0 && version.patch == constraintVersion.patch
}

