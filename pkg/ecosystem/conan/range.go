package conan

import (
	"fmt"
	"regexp"
	"strings"
)

// Package-level compiled regular expressions for range parsing
var (
	// constraintPattern matches individual constraints
	// Supports: >=, >, <=, <, ~, ^, !=, exact version
	constraintPattern = regexp.MustCompile(`^\s*(>=|>|<=|<|~|\^|!=|=)?\s*([0-9a-z\.\-\+]+.*?)\s*$`)
)

// VersionRange represents a Conan version range
type VersionRange struct {
	constraints []constraint
	original    string
}

// constraint represents a single version constraint
type constraint struct {
	operator string
	version  *Version
}

// NewVersionRange creates a new Conan version range from a string
func (e *Ecosystem) NewVersionRange(rangeStr string) (*VersionRange, error) {
	original := rangeStr
	rangeStr = strings.TrimSpace(strings.ToLower(rangeStr))

	if rangeStr == "" {
		return nil, fmt.Errorf("empty range string")
	}

	// Handle OR logic (||)
	orParts := strings.Split(rangeStr, "||")
	var allConstraints []constraint

	for _, orPart := range orParts {
		orPart = strings.TrimSpace(orPart)

		// Parse AND constraints (comma or space separated)
		andParts := splitConstraints(orPart)

		for _, andPart := range andParts {
			andPart = strings.TrimSpace(andPart)
			if andPart == "" {
				continue
			}

			constraint, err := parseConstraint(andPart, e)
			if err != nil {
				return nil, fmt.Errorf("invalid constraint '%s' in range '%s': %v", andPart, original, err)
			}

			allConstraints = append(allConstraints, constraint)
		}
	}

	if len(allConstraints) == 0 {
		return nil, fmt.Errorf("no valid constraints found in range: %s", original)
	}

	return &VersionRange{
		constraints: allConstraints,
		original:    original,
	}, nil
}

// splitConstraints splits a string into individual constraints
func splitConstraints(s string) []string {
	// First split by comma to handle comma-separated constraints
	commaParts := strings.Split(s, ",")
	var result []string

	for _, part := range commaParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// For each comma-separated part, check if it contains multiple space-separated constraints
		// But be careful not to split things like " >= 1.2.3 " which should be one constraint
		spaceParts := strings.Fields(part)
		if len(spaceParts) > 1 {
			// Check if this looks like operator + version (single constraint with spaces)
			// vs multiple separate constraints
			if len(spaceParts) == 2 && isOperator(spaceParts[0]) {
				// Single constraint with spaces like ">= 1.2.3"
				result = append(result, strings.TrimSpace(part))
			} else {
				// Multiple space-separated constraints
				result = append(result, spaceParts...)
			}
		} else {
			// Single constraint
			result = append(result, part)
		}
	}

	return result
}

// isOperator checks if a string is a valid constraint operator
func isOperator(s string) bool {
	operators := []string{">=", ">", "<=", "<", "~", "^", "!=", "="}
	for _, op := range operators {
		if s == op {
			return true
		}
	}
	return false
}

// parseConstraint parses a single constraint string
func parseConstraint(constraintStr string, e *Ecosystem) (constraint, error) {
	matches := constraintPattern.FindStringSubmatch(constraintStr)
	if matches == nil {
		return constraint{}, fmt.Errorf("invalid constraint format: %s", constraintStr)
	}

	operator := matches[1]
	versionStr := matches[2]

	// Default operator is exact match
	if operator == "" {
		operator = "="
	}

	// Parse the version
	version, err := e.NewVersion(versionStr)
	if err != nil {
		return constraint{}, fmt.Errorf("invalid version in constraint: %v", err)
	}

	return constraint{
		operator: operator,
		version:  version,
	}, nil
}

// Contains checks if a version satisfies this range
func (r *VersionRange) Contains(version *Version) bool {
	if len(r.constraints) == 0 {
		return false
	}

	// Group constraints by OR logic (|| separates OR groups)
	// Within each OR group, all constraints must be satisfied (AND logic)
	constraintGroups := r.groupConstraintsByOR()

	// Check if any OR group is satisfied
	for _, group := range constraintGroups {
		if r.groupSatisfied(group, version) {
			return true
		}
	}

	return false
}

// groupConstraintsByOR groups constraints by OR logic
func (r *VersionRange) groupConstraintsByOR() [][]constraint {
	// For now, treat all constraints as a single AND group
	// This is a simplified implementation - in a full implementation,
	// we would need to track which constraints came from which OR part
	return [][]constraint{r.constraints}
}

// groupSatisfied checks if all constraints in a group are satisfied
func (r *VersionRange) groupSatisfied(group []constraint, version *Version) bool {
	for _, c := range group {
		if !r.constraintSatisfied(c, version) {
			return false
		}
	}
	return true
}

// constraintSatisfied checks if a single constraint is satisfied
func (r *VersionRange) constraintSatisfied(c constraint, version *Version) bool {
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
	case "~":
		// Tilde allows patch-level changes
		return r.tildeMatch(version, c.version)
	case "^":
		// Caret allows compatible changes
		return r.caretMatch(version, c.version)
	default:
		return false
	}
}

// tildeMatch implements tilde (~) constraint logic
func (r *VersionRange) tildeMatch(version, constraint *Version) bool {
	// Version must be >= constraint version
	if version.Compare(constraint) < 0 {
		return false
	}

	// For Conan, tilde allows changes in the last specified part
	// ~1.2.3 allows 1.2.3, 1.2.4, 1.2.5, etc. but not 1.3.0
	// ~1.2 allows 1.2.0, 1.2.1, etc. but not 1.3.0
	// ~1 allows 1.0.0, 1.1.0, 1.2.0, etc. but not 2.0.0
	if len(constraint.parts) == 0 {
		return true
	}

	// For tilde, behavior depends on number of parts in constraint:
	// ~1 means major=1 must match (anything 1.x.x)
	// ~1.2 means major=1 AND minor=2 must match (anything 1.2.x)
	// ~1.2.3 means major=1 AND minor=2 must match (anything 1.2.x)

	switch len(constraint.parts) {
	case 1: // ~1 := >=1.0.0 <2.0.0
		// Only major needs to match
		vPart := "0"
		if len(version.parts) > 0 {
			vPart = version.parts[0]
		}
		return vPart == constraint.parts[0]

	case 2: // ~1.2 := >=1.2.0 <1.3.0
		// Major and minor must match
		for i := 0; i < 2; i++ {
			vPart := "0"
			if i < len(version.parts) {
				vPart = version.parts[i]
			}
			if vPart != constraint.parts[i] {
				return false
			}
		}
		return true

	default: // ~1.2.3 := >=1.2.3 <1.3.0
		// Major and minor must match, patch can be anything
		for i := 0; i < 2; i++ {
			vPart := "0"
			if i < len(version.parts) {
				vPart = version.parts[i]
			}
			if i < len(constraint.parts) && vPart != constraint.parts[i] {
				return false
			}
		}
		return true
	}
}

// caretMatch implements caret (^) constraint logic
func (r *VersionRange) caretMatch(version, constraint *Version) bool {
	// Version must be >= constraint version
	if version.Compare(constraint) < 0 {
		return false
	}

	// For Conan, caret allows changes that don't modify the left-most non-zero digit
	// e.g., ^1.2.3 allows 1.2.3 to 1.x.x but not 2.0.0
	// ^0.2.3 allows 0.2.3 to 0.2.x but not 0.3.0
	if len(constraint.parts) == 0 {
		return true
	}

	// Caret allows compatible changes based on the highest precedence non-zero component
	// ^1.2.3 allows 1.x.x (major is significant)
	// ^0.2.3 allows 0.2.x (minor is significant since major=0)
	// ^0.0.3 allows 0.0.x (patch is significant since major=0 and minor=0)

	// Determine which components must match exactly
	if len(constraint.parts) >= 1 && constraint.parts[0] != "0" {
		// Major is non-zero, so major must match
		vPart := "0"
		if len(version.parts) > 0 {
			vPart = version.parts[0]
		}
		return vPart == constraint.parts[0]
	} else if len(constraint.parts) >= 2 && constraint.parts[1] != "0" {
		// Major is zero but minor is non-zero, so major and minor must match
		for i := 0; i < 2; i++ {
			vPart := "0"
			if i < len(version.parts) {
				vPart = version.parts[i]
			}
			cPart := "0"
			if i < len(constraint.parts) {
				cPart = constraint.parts[i]
			}
			if vPart != cPart {
				return false
			}
		}
		return true
	} else {
		// Major and minor are both zero, all parts must match except the last can vary
		// For ^0.0.3, parts 0 and 1 must be "0", part 2 can be >= 3
		for i := 0; i < len(constraint.parts)-1; i++ {
			vPart := "0"
			if i < len(version.parts) {
				vPart = version.parts[i]
			}
			if vPart != constraint.parts[i] {
				return false
			}
		}
		return true
	}
}

// String returns the string representation of the range
func (r *VersionRange) String() string {
	return r.original
}
