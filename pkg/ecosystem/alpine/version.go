package alpine

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches Alpine version strings
// Format: number{.number}...{letter}{_suffix{number}}...{~hash}{-r#}
var versionPattern = regexp.MustCompile(`^(\d+(?:\.\d+)*)([a-z]?)((?:_[a-z]+\d*)*)(\~[a-f0-9]+)?(-r\d+)?$`)

// unknownSuffixPrecedence is the precedence value assigned to unknown suffixes
// Unknown suffixes are ordered after all known suffixes but before each other lexicographically
const unknownSuffixPrecedence = 1000

// Version represents an Alpine Linux package version
type Version struct {
	numeric    []numericComponent // numeric components: 1.2.3 (with leading zero info)
	letter     string             // optional letter after numeric: a, b, etc.
	suffixes   []suffix           // suffixes: _alpha1, _beta, etc.
	hash       string             // commit hash: ~abc123...
	build      int                // build component: -r1, -r2, etc.
	original   string             // original version string
}

// suffix represents a version suffix like _alpha1, _beta, etc.
type suffix struct {
	name   string // alpha, beta, pre, rc, cvs, svn, git, hg, p
	number int    // optional number after suffix name
}

// Suffix precedence order (lower index = lower precedence)
// Based on Alpine apk-tools version comparison logic
var suffixOrder = map[string]int{
	"alpha": 0,
	"beta":  1,
	"pre":   2,
	"rc":    3,
	"":      4, // no suffix (release)
	"cvs":   5,
	"svn":   6,
	"git":   7,
	"hg":    8,
	"p":     9,
}

// NewVersion creates a new Alpine version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)
	
	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}
	
	// Parse using regex
	matches := versionPattern.FindStringSubmatch(version)
	if matches == nil {
		// Check if this might be a valid version that just doesn't match our regex
		// Only allow versions that contain at least some digits
		hasDigits := strings.ContainsAny(version, "0123456789")
		
		if !hasDigits {
			return nil, fmt.Errorf("invalid Alpine version: %s", original)
		}
		
		// If version has digits but doesn't match standard pattern, create a special "string-only" version
		// This handles cases like "1.0bc" mentioned in the test data comment "# invalid. do string sort"
		return &Version{
			numeric:  nil,
			letter:   "",
			suffixes: nil,
			hash:     "",
			build:    0,
			original: original,
		}, nil
	}
	
	numericPart := matches[1]
	letterPart := matches[2]
	suffixPart := matches[3]
	hashPart := matches[4]
	buildPart := matches[5]
	
	// Parse numeric components
	numeric, err := parseNumericComponents(numericPart)
	if err != nil {
		return nil, fmt.Errorf("invalid numeric components in version %s: %v", original, err)
	}
	
	// Parse suffixes
	suffixes, err := parseSuffixes(suffixPart)
	if err != nil {
		return nil, fmt.Errorf("invalid suffixes in version %s: %v", original, err)
	}
	
	// Parse hash (remove ~ prefix)
	hash := ""
	if hashPart != "" {
		hash = hashPart[1:] // remove ~
	}
	
	// Parse build component (remove -r prefix)
	build := 0
	if buildPart != "" {
		buildStr := buildPart[2:] // remove -r
		build, err = strconv.Atoi(buildStr)
		if err != nil {
			return nil, fmt.Errorf("invalid build component in version %s: %v", original, err)
		}
	}
	
	return &Version{
		numeric:  numeric,
		letter:   letterPart,
		suffixes: suffixes,
		hash:     hash,
		build:    build,
		original: original,
	}, nil
}

// numericComponent represents a numeric component with leading zero information
type numericComponent struct {
	value       int    // The actual numeric value
	originalStr string // The original string representation (to detect leading zeros)
}

// parseNumericComponents parses numeric components like "1.2.3"
func parseNumericComponents(s string) ([]numericComponent, error) {
	if s == "" {
		return nil, fmt.Errorf("empty numeric components")
	}
	
	parts := strings.Split(s, ".")
	numeric := make([]numericComponent, len(parts))
	
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid numeric component: %s", part)
		}
		numeric[i] = numericComponent{
			value:       num,
			originalStr: part,
		}
	}
	
	return numeric, nil
}

// parseSuffixes parses suffix components like "_alpha1_beta"
func parseSuffixes(s string) ([]suffix, error) {
	if s == "" {
		return nil, nil
	}
	
	// Remove leading underscore and split by underscore
	s = strings.TrimPrefix(s, "_")
	parts := strings.Split(s, "_")
	
	var suffixes []suffix
	suffixRegex := regexp.MustCompile(`^([a-z]+)(\d*)$`)
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		matches := suffixRegex.FindStringSubmatch(part)
		if matches == nil {
			return nil, fmt.Errorf("invalid suffix format: %s", part)
		}
		
		name := matches[1]
		numberStr := matches[2]
		
		// Allow unknown suffixes - they will be treated as having a very high precedence
		// This handles cases like "_foo" which should be compared lexicographically
		
		number := 0
		if numberStr != "" {
			var err error
			number, err = strconv.Atoi(numberStr)
			if err != nil {
				return nil, fmt.Errorf("invalid suffix number: %s", numberStr)
			}
		}
		
		suffixes = append(suffixes, suffix{
			name:   name,
			number: number,
		})
	}
	
	return suffixes, nil
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another Alpine version
func (v *Version) Compare(other *Version) int {
	// Handle invalid versions (no numeric components) - use string comparison
	if v.numeric == nil || other.numeric == nil {
		return strings.Compare(v.original, other.original)
	}
	
	// 1. Compare numeric components (leading zeros are ignored - use actual numeric values)
	numericCmp := compareNumericArraysNumeric(v.numeric, other.numeric)
	if numericCmp != 0 {
		return numericCmp
	}
	
	// 2. Compare letters (empty letter comes before any letter)
	letterCmp := compareLetters(v.letter, other.letter)
	if letterCmp != 0 {
		return letterCmp
	}
	
	// 3. Compare suffixes
	suffixCmp := compareSuffixArrays(v.suffixes, other.suffixes)
	if suffixCmp != 0 {
		return suffixCmp
	}
	
	// 4. Compare hashes (lexicographically)
	hashCmp := strings.Compare(v.hash, other.hash)
	if hashCmp != 0 {
		return hashCmp
	}
	
	// 5. Compare build components
	return compareInt(v.build, other.build)
}



// compareLetters compares optional letters
func compareLetters(a, b string) int {
	// Empty letter (no letter) comes before any letter
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return -1
	}
	if b == "" {
		return 1
	}
	
	// Lexical comparison for letters
	return strings.Compare(a, b)
}

// compareSuffixArrays compares suffix arrays
func compareSuffixArrays(a, b []suffix) int {
	if len(a) == 0 && len(b) == 0 {
		return 0
	}
	
	// Handle cases where one side has no suffixes (release version)
	// Compare against empty suffix (weight 4)
	if len(a) == 0 {
		// Compare release vs first suffix of b
		return compareSuffixes(suffix{name: "", number: 0}, b[0])
	}
	if len(b) == 0 {
		// Compare first suffix of a vs release
		return compareSuffixes(a[0], suffix{name: "", number: 0})
	}
	
	// Compare suffix by suffix up to the minimum length
	minLen := min(len(a), len(b))
	
	for i := 0; i < minLen; i++ {
		cmp := compareSuffixes(a[i], b[i])
		if cmp != 0 {
			return cmp
		}
	}
	
	// If all compared suffixes are equal, the longer array is "smaller"
	// This means "alpha_pre" < "alpha" (more suffixes = less stable)
	return compareInt(len(b), len(a))
}

// compareSuffixes compares two individual suffixes
func compareSuffixes(a, b suffix) int {
	// Get suffix precedence order, defaulting to a high value for unknown suffixes
	aOrder, aExists := suffixOrder[a.name]
	bOrder, bExists := suffixOrder[b.name]
	
	// Unknown suffixes get a higher precedence (come after known suffixes)
	if !aExists {
		aOrder = unknownSuffixPrecedence
	}
	if !bExists {
		bOrder = unknownSuffixPrecedence
	}
	
	// If both are unknown, compare lexicographically by name
	if !aExists && !bExists {
		nameCmp := strings.Compare(a.name, b.name)
		if nameCmp != 0 {
			return nameCmp
		}
		// If same name, compare numbers
		return compareInt(a.number, b.number)
	}
	
	// Compare by suffix precedence order first
	orderCmp := compareInt(aOrder, bOrder)
	if orderCmp != 0 {
		return orderCmp
	}
	
	// If same suffix type, compare numbers
	return compareInt(a.number, b.number)
}


// compareNumericArraysNumeric compares numeric arrays using Alpine's rules
func compareNumericArraysNumeric(a, b []numericComponent) int {
	maxLen := max(len(a), len(b))
	
	for i := range maxLen {
		var aComp, bComp numericComponent
		
		if i < len(a) {
			aComp = a[i]
		} else {
			aComp = numericComponent{value: 0, originalStr: "0"}
		}
		if i < len(b) {
			bComp = b[i]
		} else {
			bComp = numericComponent{value: 0, originalStr: "0"}
		}
		
		var cmp int
		if i == 0 {
			// Major component: always compare numerically (ignore leading zeros)
			cmp = compareInt(aComp.value, bComp.value)
		} else {
			// Minor/patch components: if either has leading zeros, use string comparison
			if hasLeadingZero(aComp.originalStr) || hasLeadingZero(bComp.originalStr) {
				cmp = strings.Compare(aComp.originalStr, bComp.originalStr)
			} else {
				// Both have no leading zeros, use numeric comparison
				cmp = compareInt(aComp.value, bComp.value)
			}
		}
		
		if cmp != 0 {
			return cmp
		}
	}
	
	return 0
}

// hasLeadingZero checks if a numeric string has leading zeros
func hasLeadingZero(s string) bool {
	return len(s) > 1 && s[0] == '0'
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
