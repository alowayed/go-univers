// Package pypi provides functionality for working with PyPI package versions following PEP 440.
package pypi

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Version represents a PyPI package version following PEP 440
type Version struct {
	epoch       int
	release     []int
	prerelease  string
	preNumber   int
	postrelease int
	dev         int
	local       string
	original    string
}

// versionPattern matches PyPI version strings according to PEP 440
var versionPattern = regexp.MustCompile(`^(?:([0-9]+)!)?([0-9]+(?:\.[0-9]+)*?)(?:\.?(a|b|rc|alpha|beta|c)([0-9]+))?(?:\.?(post|rev|r)([0-9]+))?(?:\.?(dev)([0-9]+))?(?:\+([a-zA-Z0-9]+(?:[-_.][a-zA-Z0-9]+)*))?$`)

// NewVersion creates a new PyPI version from a string
func NewVersion(version string) (*Version, error) {
	version = strings.TrimSpace(version)
	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Parse using regex
	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		return nil, fmt.Errorf("invalid PyPI version format: %s", version)
	}

	pv := &Version{
		epoch:       0,
		postrelease: -1,
		dev:         -1,
		original:    version,
	}

	// Parse epoch (group 1)
	if matches[1] != "" {
		epoch, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid epoch: %s", matches[1])
		}
		pv.epoch = epoch
	}

	// Parse release version (group 2)
	if matches[2] == "" {
		return nil, fmt.Errorf("missing release version")
	}
	releaseParts := strings.Split(matches[2], ".")
	pv.release = make([]int, len(releaseParts))
	for i, part := range releaseParts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid release part: %s", part)
		}
		pv.release[i] = num
	}

	// Parse prerelease (groups 3, 4)
	if matches[3] != "" {
		pv.prerelease = matches[3]
		if matches[4] != "" {
			preNum, err := strconv.Atoi(matches[4])
			if err != nil {
				return nil, fmt.Errorf("invalid prerelease number: %s", matches[4])
			}
			pv.preNumber = preNum
		}
	}

	// Parse post-release (groups 5, 6)
	if matches[5] != "" {
		if matches[6] != "" {
			postNum, err := strconv.Atoi(matches[6])
			if err != nil {
				return nil, fmt.Errorf("invalid post number: %s", matches[6])
			}
			pv.postrelease = postNum
		} else {
			pv.postrelease = 0
		}
	}

	// Parse dev release (groups 7, 8)
	if matches[7] != "" {
		if matches[8] != "" {
			devNum, err := strconv.Atoi(matches[8])
			if err != nil {
				return nil, fmt.Errorf("invalid dev number: %s", matches[8])
			}
			pv.dev = devNum
		} else {
			pv.dev = 0
		}
	}

	// Parse local version (group 9)
	if matches[9] != "" {
		pv.local = matches[9]
	}

	return pv, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// IsValid checks if the version is valid
func (v *Version) IsValid() bool {
	return len(v.release) > 0
}

// Normalize returns the normalized form of the version
func (v *Version) Normalize() string {
	result := ""
	
	// Add epoch if present
	if v.epoch > 0 {
		result += fmt.Sprintf("%d!", v.epoch)
	}
	
	// Add release version
	releaseStrs := make([]string, len(v.release))
	for i, part := range v.release {
		releaseStrs[i] = strconv.Itoa(part)
	}
	result += strings.Join(releaseStrs, ".")
	
	// Add prerelease
	if v.prerelease != "" {
		result += v.prerelease + strconv.Itoa(v.preNumber)
	}
	
	// Add post-release
	if v.postrelease >= 0 {
		result += "post" + strconv.Itoa(v.postrelease)
	}
	
	// Add dev release
	if v.dev >= 0 {
		result += "dev" + strconv.Itoa(v.dev)
	}
	
	// Add local version
	if v.local != "" {
		result += "+" + v.local
	}
	
	return result
}

// Compare compares this version with another PyPI version according to PEP 440
func (v *Version) Compare(other *Version) int {
	// Compare epoch
	if v.epoch != other.epoch {
		return compareInt(v.epoch, other.epoch)
	}

	// Compare release versions
	releaseComparison := compareReleaseVersions(v.release, other.release)
	if releaseComparison != 0 {
		return releaseComparison
	}

	// Compare prerelease
	preComparison := comparePrereleases(v.prerelease, v.preNumber, other.prerelease, other.preNumber)
	if preComparison != 0 {
		return preComparison
	}

	// Compare post-release
	postComparison := comparePostReleases(v.postrelease, other.postrelease)
	if postComparison != 0 {
		return postComparison
	}

	// Compare dev release
	return compareDevReleases(v.dev, other.dev)
}

// Satisfies checks if this version satisfies the given constraint
func (v *Version) Satisfies(constraint *Constraint) bool {
	return constraint.Matches(v)
}

// VersionRange represents a PyPI version range with PEP 440 syntax support
type VersionRange struct {
	constraints []*Constraint
	original    string
}

// NewVersionRange creates a new PyPI version range from a specifier string
func NewVersionRange(specifier string) (*VersionRange, error) {
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
func parseSpecifier(specifier string) ([]*Constraint, error) {
	// Handle comma-separated constraints (AND logic)
	if strings.Contains(specifier, ",") {
		parts := strings.Split(specifier, ",")
		var allConstraints []*Constraint
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
func parseSingleConstraint(constraint string) ([]*Constraint, error) {
	constraint = strings.TrimSpace(constraint)

	// Handle PEP 440 operators in correct order
	operators := []string{"===", "~=", "==", "!=", "<=", ">=", "<", ">"}
	for _, op := range operators {
		if strings.HasPrefix(constraint, op) {
			version := strings.TrimSpace(constraint[len(op):])
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

			return []*Constraint{{operator: op, version: version}}, nil
		}
	}

	// Default to equality
	return []*Constraint{{operator: "==", version: constraint}}, nil
}

// parseCompatibleRelease handles the ~= operator
func parseCompatibleRelease(version string) ([]*Constraint, error) {
	v, err := NewVersion(version)
	if err != nil {
		return nil, err
	}

	// ~=2.2 is equivalent to >=2.2, <3.0
	if len(v.release) == 1 {
		upperVersion := fmt.Sprintf("%d.0", v.release[0]+1)
		return []*Constraint{
			{operator: ">=", version: version},
			{operator: "<", version: upperVersion},
		}, nil
	}

	// ~=1.4.2 is equivalent to >=1.4.2, <1.5.0
	if len(v.release) >= 2 {
		upperVersion := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1]+1)
		return []*Constraint{
			{operator: ">=", version: version},
			{operator: "<", version: upperVersion},
		}, nil
	}

	return []*Constraint{{operator: ">=", version: version}}, nil
}

// parseWildcardConstraint handles wildcard constraints like ==1.2.* or !=1.2.*
func parseWildcardConstraint(operator, version string) ([]*Constraint, error) {
	// Remove the .* suffix
	baseVersion := strings.TrimSuffix(version, ".*")
	
	v, err := NewVersion(baseVersion + ".0")
	if err != nil {
		return nil, err
	}

	if operator == "==" {
		// ==1.2.* means >=1.2.0, <1.3.0
		if len(v.release) >= 2 {
			lowerBound := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1])
			upperBound := fmt.Sprintf("%d.%d.0", v.release[0], v.release[1]+1)
			return []*Constraint{
				{operator: ">=", version: lowerBound},
				{operator: "<", version: upperBound},
			}, nil
		}

		// ==1.* means >=1.0.0, <2.0.0
		if len(v.release) >= 1 {
			lowerBound := fmt.Sprintf("%d.0.0", v.release[0])
			upperBound := fmt.Sprintf("%d.0.0", v.release[0]+1)
			return []*Constraint{
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
			return []*Constraint{
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
		if !constraint.Matches(version) {
			return false
		}
	}
	return true
}

// Constraints returns all constraints in this range
func (pr *VersionRange) Constraints() []*Constraint {
	return pr.constraints
}

// IsEmpty returns true if the range contains no valid versions
func (pr *VersionRange) IsEmpty() bool {
	return len(pr.constraints) == 0
}

// Constraint represents a single PyPI version constraint
type Constraint struct {
	operator string
	version  string
}

// String returns the string representation of the constraint
func (c *Constraint) String() string {
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
	// Handle arbitrary equality (===)
	if c.operator == "===" {
		return version.String() == c.version
	}

	constraintVersion, err := NewVersion(c.version)
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