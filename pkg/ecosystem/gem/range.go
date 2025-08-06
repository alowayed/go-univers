package gem

import (
	"fmt"
	"strings"
)

// VersionRange represents a Ruby Gem version range with Gem-specific syntax support
type VersionRange struct {
	constraints []*constraint
	original    string
}

// constraint represents a single Ruby Gem version constraint
type constraint struct {
	operator string
	version  string
}

// NewVersionRange creates a new Ruby Gem version range from a range string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	original := rangeStr
	rangeStr = strings.TrimSpace(rangeStr)
	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	constraints, err := parseConstraints(rangeStr)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    original,
	}, nil
}

// parseConstraints parses Ruby Gem constraint syntax
func parseConstraints(rangeStr string) ([]*constraint, error) {
	// Handle multiple constraints separated by commas
	parts := strings.Split(rangeStr, ",")
	var constraints []*constraint

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		constraint, err := parseConstraint(part)
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
func parseConstraint(constraintStr string) (*constraint, error) {
	constraintStr = strings.TrimSpace(constraintStr)

	// Pessimistic constraint (~>)
	if strings.HasPrefix(constraintStr, "~>") {
		version := strings.TrimSpace(constraintStr[2:])
		if version == "" {
			return nil, fmt.Errorf("pessimistic constraint requires version")
		}
		return &constraint{operator: "~>", version: version}, nil
	}

	// Other operators
	operators := []string{">=", "<=", "!=", ">", "<", "="}
	for _, op := range operators {
		if strings.HasPrefix(constraintStr, op) {
			version := strings.TrimSpace(constraintStr[len(op):])
			if version == "" {
				return nil, fmt.Errorf("constraint %s requires version", op)
			}
			return &constraint{operator: op, version: version}, nil
		}
	}

	// Default to exact match
	return &constraint{operator: "=", version: constraintStr}, nil
}

// String returns the string representation of the version range
func (vr *VersionRange) String() string {
	return vr.original
}

// Contains checks if a version satisfies this range
func (vr *VersionRange) Contains(version *Version) bool {
	ecosystem := &Ecosystem{}

	// All constraints must be satisfied (AND logic)
	for _, c := range vr.constraints {
		if !satisfiesConstraint(version, c, ecosystem) {
			return false
		}
	}

	return true
}

// satisfiesConstraint checks if a version satisfies a single constraint
func satisfiesConstraint(version *Version, c *constraint, ecosystem *Ecosystem) bool {
	constraintVersion, err := ecosystem.NewVersion(c.version)
	if err != nil {
		return false
	}

	cmp := version.Compare(constraintVersion)

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
	case "~>":
		return satisfiesPessimistic(version, constraintVersion)
	default:
		return false
	}
}

// satisfiesPessimistic implements the Ruby Gem pessimistic constraint (~>)
func satisfiesPessimistic(version, constraint *Version) bool {
	// ~> 1.2.3 means >= 1.2.3 and < 1.3.0
	// ~> 1.2 means >= 1.2.0 and < 2.0.0

	// Must be >= constraint version
	if version.Compare(constraint) < 0 {
		return false
	}

	// Get the numeric parts of both version and constraint for comparison
	versionNumeric, _ := version.splitNumericAndPrerelease()
	constraintNumeric, constraintPrerelease := constraint.splitNumericAndPrerelease()

	// For range calculations, we need to understand the original precision
	// Count numeric segments from the original constraint string
	constraintStr := constraint.String()
	mainPart := constraintStr
	if dashIndex := strings.Index(constraintStr, "-"); dashIndex != -1 {
		mainPart = constraintStr[:dashIndex]
	}
	originalSegments := strings.Split(mainPart, ".")
	numericSegments := len(originalSegments)

	// For pessimistic constraints, all segments except the last must match exactly
	numSegmentsToCheck := numericSegments - 1

	// Special case: single segment constraint (~> 1)
	if numericSegments == 1 {
		numSegmentsToCheck = 1
	}

	// Special case: constraint has prerelease (~> 1.0.0-alpha)
	// When constraint has prerelease, all numeric segments must match exactly
	if len(constraintPrerelease) > 0 {
		numSegmentsToCheck = numericSegments
	}

	// Check that the required segments match exactly
	for i := 0; i < numSegmentsToCheck; i++ {
		var vSeg, cSeg int
		if i < len(versionNumeric) {
			vSeg = versionNumeric[i].numValue
		}
		if i < len(constraintNumeric) {
			cSeg = constraintNumeric[i].numValue
		}

		if vSeg != cSeg {
			return false
		}
	}

	return true
}
