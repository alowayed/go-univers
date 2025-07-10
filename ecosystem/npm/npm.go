package npm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents an NPM package version following semantic versioning
type Version struct {
	major      int
	minor      int
	patch      int
	prerelease string
	build      string
	original   string
}

// versionPattern matches NPM version strings
var versionPattern = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$`)

// NewVersion creates a new NPM version from a string
func NewVersion(version string) (*Version, error) {
	// Remove leading v or =
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "=")
	
	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid NPM version: %s", version)
	}

	major, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid major version: %s", matches[1])
	}

	minor, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid minor version: %s", matches[2])
	}

	patch, err := strconv.Atoi(matches[3])
	if err != nil {
		return nil, fmt.Errorf("invalid patch version: %s", matches[3])
	}

	return &Version{
		major:      major,
		minor:      minor,
		patch:      patch,
		prerelease: matches[4],
		build:      matches[5],
		original:   version,
	}, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// IsValid checks if the version is valid
func (v *Version) IsValid() bool {
	return v.major >= 0 && v.minor >= 0 && v.patch >= 0
}

// Normalize returns the normalized form of the version
func (v *Version) Normalize() string {
	result := fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	if v.prerelease != "" {
		result += "-" + v.prerelease
	}
	if v.build != "" {
		result += "+" + v.build
	}
	return result
}

// Compare compares this version with another NPM version
func (v *Version) Compare(other *Version) int {
	// Compare major.minor.patch
	if v.major != other.major {
		return compareInt(v.major, other.major)
	}
	if v.minor != other.minor {
		return compareInt(v.minor, other.minor)
	}
	if v.patch != other.patch {
		return compareInt(v.patch, other.patch)
	}

	// Compare prerelease according to semver rules
	return comparePrerelease(v.prerelease, other.prerelease)
}

// Satisfies checks if this version satisfies the given constraint
func (v *Version) Satisfies(constraint *Constraint) bool {
	return constraint.Matches(v)
}

// VersionRange represents an NPM version range with NPM-specific syntax support
type VersionRange struct {
	constraints []*Constraint
	original    string
}

// NewVersionRange creates a new NPM version range from a range string
func NewVersionRange(rangeStr string) (*VersionRange, error) {
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

// parseRange parses NPM range syntax into constraints
func parseRange(rangeStr string) ([]*Constraint, error) {
	// Handle OR logic (||)
	if strings.Contains(rangeStr, "||") {
		parts := strings.Split(rangeStr, "||")
		var allConstraints []*Constraint
		for _, part := range parts {
			constraints, err := parseRange(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			allConstraints = append(allConstraints, constraints...)
		}
		return allConstraints, nil
	}

	// Handle hyphen ranges (1.2.3 - 2.3.4)
	if strings.Contains(rangeStr, " - ") {
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
func parseSingleConstraint(constraint string) ([]*Constraint, error) {
	constraint = strings.TrimSpace(constraint)

	// Handle wildcard
	if constraint == "*" {
		return []*Constraint{{operator: "*", version: "*"}}, nil
	}

	// Handle caret range (^1.2.3)
	if strings.HasPrefix(constraint, "^") {
		return parseCaretRange(constraint[1:])
	}

	// Handle tilde range (~1.2.3)
	if strings.HasPrefix(constraint, "~") {
		return parseTildeRange(constraint[1:])
	}

	// Handle x-range (1.x, 1.2.x)
	if strings.Contains(constraint, "x") || strings.Contains(constraint, "X") {
		return parseXRange(constraint)
	}

	// Handle comparison operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(constraint, op) {
			version := strings.TrimSpace(constraint[len(op):])
			return []*Constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to exact match
	return []*Constraint{{operator: "=", version: constraint}}, nil
}

// parseCaretRange handles caret ranges (^1.2.3)
func parseCaretRange(version string) ([]*Constraint, error) {
	v, err := NewVersion(version)
	if err != nil {
		return nil, err
	}

	// ^1.2.3 means >=1.2.3 <2.0.0
	return []*Constraint{
		{operator: ">=", version: v.Normalize()},
		{operator: "<", version: fmt.Sprintf("%d.0.0", v.major+1)},
	}, nil
}

// parseTildeRange handles tilde ranges (~1.2.3)
func parseTildeRange(version string) ([]*Constraint, error) {
	v, err := NewVersion(version)
	if err != nil {
		return nil, err
	}

	// ~1.2.3 means >=1.2.3 <1.3.0
	return []*Constraint{
		{operator: ">=", version: v.Normalize()},
		{operator: "<", version: fmt.Sprintf("%d.%d.0", v.major, v.minor+1)},
	}, nil
}

// parseXRange handles x-ranges (1.x, 1.2.x)
func parseXRange(rangeStr string) ([]*Constraint, error) {
	parts := strings.Split(rangeStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid x-range: %s", rangeStr)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid major version in x-range: %s", parts[0])
	}

	// 1.x means >=1.0.0 <2.0.0
	if len(parts) == 2 && (parts[1] == "x" || parts[1] == "X") {
		return []*Constraint{
			{operator: ">=", version: fmt.Sprintf("%d.0.0", major)},
			{operator: "<", version: fmt.Sprintf("%d.0.0", major+1)},
		}, nil
	}

	// 1.2.x means >=1.2.0 <1.3.0
	if len(parts) == 3 && (parts[2] == "x" || parts[2] == "X") {
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid minor version in x-range: %s", parts[1])
		}
		return []*Constraint{
			{operator: ">=", version: fmt.Sprintf("%d.%d.0", major, minor)},
			{operator: "<", version: fmt.Sprintf("%d.%d.0", major, minor+1)},
		}, nil
	}

	return nil, fmt.Errorf("unsupported x-range format: %s", rangeStr)
}

// parseHyphenRange handles hyphen ranges (1.2.3 - 2.3.4)
func parseHyphenRange(rangeStr string) ([]*Constraint, error) {
	parts := strings.Split(rangeStr, " - ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid hyphen range: %s", rangeStr)
	}

	return []*Constraint{
		{operator: ">=", version: strings.TrimSpace(parts[0])},
		{operator: "<=", version: strings.TrimSpace(parts[1])},
	}, nil
}

// parseSpaceSeparatedConstraints handles space-separated constraints (>=1.0.0 <2.0.0)
func parseSpaceSeparatedConstraints(rangeStr string) ([]*Constraint, error) {
	parts := strings.Fields(rangeStr)
	var constraints []*Constraint

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
	// All constraints must be satisfied (AND logic)
	for _, constraint := range nr.constraints {
		if !constraint.Matches(version) {
			return false
		}
	}
	return true
}

// Constraints returns all constraints in this range
func (nr *VersionRange) Constraints() []*Constraint {
	return nr.constraints
}

// IsEmpty returns true if the range contains no valid versions
func (nr *VersionRange) IsEmpty() bool {
	return len(nr.constraints) == 0
}

// Constraint represents a single NPM version constraint
type Constraint struct {
	operator string
	version  string
}

// String returns the string representation of the constraint
func (c *Constraint) String() string {
	if c.operator == "=" {
		return c.version
	}
	return c.operator + c.version
}

// Operator returns the constraint operator
func (c *Constraint) Operator() string {
	return c.operator
}

// Version returns the version part of the constraint
func (c *Constraint) Version() string {
	return c.version
}

// Matches checks if the given version matches this constraint
func (c *Constraint) Matches(version *Version) bool {
	if c.operator == "*" {
		return true
	}

	constraintVersion, err := NewVersion(c.version)
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

// Helper functions

func comparePrerelease(a, b string) int {
	// No prerelease has higher precedence than prerelease
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1
	}
	if b == "" {
		return -1
	}

	// Split by dots and compare each part
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}

	for i := 0; i < maxLen; i++ {
		var aPart, bPart string
		if i < len(aParts) {
			aPart = aParts[i]
		}
		if i < len(bParts) {
			bPart = bParts[i]
		}

		// Missing part has lower precedence
		if aPart == "" && bPart != "" {
			return -1
		}
		if aPart != "" && bPart == "" {
			return 1
		}

		// Try to parse as numbers
		aNum, aIsNum := parseNum(aPart)
		bNum, bIsNum := parseNum(bPart)

		if aIsNum && bIsNum {
			if aNum != bNum {
				return compareInt(aNum, bNum)
			}
		} else if aIsNum {
			return -1 // Numeric identifiers have lower precedence
		} else if bIsNum {
			return 1
		} else {
			// Both are strings, compare lexically
			if aPart != bPart {
				return strings.Compare(aPart, bPart)
			}
		}
	}

	return 0
}

func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func parseNum(s string) (int, bool) {
	if num, err := strconv.Atoi(s); err == nil {
		return num, true
	}
	return 0, false
}