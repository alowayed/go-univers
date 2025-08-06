package pypi

import (
	"fmt"
	"strings"
)

// VersionRange represents a PyPI version range with PEP 440 syntax support
type VersionRange struct {
	constraints []*constraint
	original    string
}

// NewVersionRange creates a new PyPI version range from a specifier string
func (e *Ecosystem) NewVersionRange(specifier string) (*VersionRange, error) {
	specifier = strings.TrimSpace(specifier)
	if specifier == "" {
		return nil, fmt.Errorf("empty specifier string")
	}

	constraints, err := parseSpecifier(specifier)
	if err != nil {
		return nil, err
	}

	return &VersionRange{
		constraints: constraints,
		original:    specifier,
	}, nil
}

// parseSpecifier parses PyPI version specifiers
func parseSpecifier(specifier string) ([]*constraint, error) {
	// Handle comma-separated constraints (AND logic)
	if strings.Contains(specifier, ",") {
		parts := strings.Split(specifier, ",")
		var allConstraints []*constraint
		for _, part := range parts {
			constraints, err := parseSpecifier(strings.TrimSpace(part))
			if err != nil {
				return nil, err
			}
			allConstraints = append(allConstraints, constraints...)
		}
		return allConstraints, nil
	}

	// Parse single constraint
	return parseSingleConstraint(specifier)
}

// parseSingleConstraint parses a single PyPI constraint
func parseSingleConstraint(con string) ([]*constraint, error) {
	con = strings.TrimSpace(con)

	// Handle PEP 440 operators in correct order
	operators := []string{"===", "~=", "==", "!=", "<=", ">=", "<", ">"}
	for _, op := range operators {
		if strings.HasPrefix(con, op) {
			version := strings.TrimSpace(con[len(op):])
			if version == "" {
				return nil, fmt.Errorf("empty version after operator '%s'", op)
			}

			// Handle compatible release operator (~=)
			if op == "~=" {
				return parseCompatibleRelease(version)
			}

			// Handle wildcard in equality (==1.2.* or !=1.2.*)
			if (op == "==" || op == "!=") && strings.HasSuffix(version, ".*") {
				return parseWildcardConstraint(op, version)
			}

			return []*constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to equality
	return []*constraint{{operator: "==", version: con}}, nil
}

// parseCompatibleRelease handles the ~= operator
func parseCompatibleRelease(version string) ([]*constraint, error) {
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		return nil, err
	}

	// ~=2.2 is equivalent to >=2.2, <3.0
	if len(v.release) == 1 {
		upperVersion := fmt.Sprintf("%d.0", v.release[0]+1)
		return []*constraint{
			{operator: ">=", version: version},
			{operator: "<", version: upperVersion},
		}, nil
	}

	// ~=1.4.2 is equivalent to >=1.4.2, <1.5.0
	if len(v.release) >= 2 {
		upperVersion := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1]+1)
		return []*constraint{
			{operator: ">=", version: version},
			{operator: "<", version: upperVersion},
		}, nil
	}

	return []*constraint{{operator: ">=", version: version}}, nil
}

// parseWildcardConstraint handles wildcard constraints like ==1.2.* or !=1.2.*
func parseWildcardConstraint(operator, version string) ([]*constraint, error) {
	// Remove the .* suffix
	baseVersion := strings.TrimSuffix(version, ".*")

	e := &Ecosystem{}
	v, err := e.NewVersion(baseVersion + ".0")
	if err != nil {
		return nil, err
	}

	if operator == "==" {
		// ==1.2.* means >=1.2.0, <1.3.0
		if len(v.release) >= 2 {
			lowerBound := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1])
			upperBound := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1]+1)
			return []*constraint{
				{operator: ">=", version: lowerBound},
				{operator: "<", version: upperBound},
			}, nil
		}

		// ==1.* means >=1.0.0, <2.0.0
		if len(v.release) >= 1 {
			lowerBound := fmt.Sprintf("%d.0.0", v.release[0])
			upperBound := fmt.Sprintf("%d.0.0", v.release[0]+1)
			return []*constraint{
				{operator: ">=", version: lowerBound},
				{operator: "<", version: upperBound},
			}, nil
		}
	}

	if operator == "!=" {
		// !=1.2.* means <1.2.0 or >=1.3.0
		if len(v.release) >= 2 {
			lowerBound := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1])
			upperBound := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1]+1)
			return []*constraint{
				{operator: "<", version: lowerBound},
				{operator: ">=", version: upperBound},
			}, nil
		}
	}

	return nil, fmt.Errorf("unsupported wildcard constraint: %s%s", operator, version)
}

// String returns the string representation of the range
func (pr *VersionRange) String() string {
	return pr.original
}

// Contains checks if a version is within this range
func (pr *VersionRange) Contains(version *Version) bool {
	// All constraints must be satisfied (AND logic)
	for _, constraint := range pr.constraints {
		if !constraint.matches(version) {
			return false
		}
	}
	return true
}

// Constraint represents a single PyPI version constraint
type constraint struct {
	operator string
	version  string
}

// matches checks if the given version matches this constraint
func (c *constraint) matches(version *Version) bool {
	// Handle arbitrary equality (===)
	if c.operator == "===" {
		return version.String() == c.version
	}

	e := &Ecosystem{}
	constraintVersion, err := e.NewVersion(c.version)
	if err != nil {
		return false
	}

	comparison := version.Compare(constraintVersion)

	switch c.operator {
	case "==":
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

// compareReleaseVersions returns -1, 0, or 1 comparing version arrays element by element
func compareReleaseVersions(a, b []int) int {
	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	for i := 0; i < maxLen; i++ {
		aVal := 0
		bVal := 0
		if i < len(a) {
			aVal = a[i]
		}
		if i < len(b) {
			bVal = b[i]
		}

		if aVal != bVal {
			return compareInt(aVal, bVal)
		}
	}

	return 0
}

// comparePrereleases returns -1, 0, or 1 where no prerelease > prerelease, alpha < beta < rc
func comparePrereleases(aPre string, aNum int, bPre string, bNum int) int {
	// No prerelease has higher precedence than prerelease
	if aPre == "" && bPre == "" {
		return 0
	}
	if aPre == "" {
		return 1
	}
	if bPre == "" {
		return -1
	}

	// Compare prerelease types: alpha < beta < rc
	aType := normalizePrereleaseType(aPre)
	bType := normalizePrereleaseType(bPre)

	if aType != bType {
		return compareInt(aType, bType)
	}

	// Same type, compare numbers
	return compareInt(aNum, bNum)
}

// normalizePrereleaseType returns numeric priority: alpha=1, beta=2, rc=3
func normalizePrereleaseType(preType string) int {
	switch strings.ToLower(preType) {
	case "a", "alpha":
		return 1
	case "b", "beta":
		return 2
	case "c", "rc":
		return 3
	default:
		return 0
	}
}

// comparePostReleases returns -1, 0, or 1 where -1 means no post-release (lower precedence)
func comparePostReleases(a, b int) int {
	// -1 means no post-release
	if a == -1 && b == -1 {
		return 0
	}
	if a == -1 {
		return -1
	}
	if b == -1 {
		return 1
	}
	return compareInt(a, b)
}

// compareDevReleases returns -1, 0, or 1 where -1 means no dev release (higher precedence)
func compareDevReleases(a, b int) int {
	// -1 means no dev release
	if a == -1 && b == -1 {
		return 0
	}
	if a == -1 {
		return 1
	}
	if b == -1 {
		return -1
	}
	return compareInt(a, b)
}

// compareInt returns -1 if a < b, 0 if a == b, 1 if a > b
func compareInt(a, b int) int {
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

