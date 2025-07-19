package maven

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type Version struct {
	original string
	elements []element
}

type element struct {
	value    interface{} // string or int
	isNumber bool
}

func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	if version == "" {
		return nil, fmt.Errorf("version string cannot be empty")
	}

	// Trim leading and trailing whitespace
	trimmed := strings.TrimSpace(version)
	if trimmed == "" {
		return nil, fmt.Errorf("version string cannot be empty or only whitespace")
	}

	// Basic validation - Maven versions should contain at least one digit or known qualifier
	if !isValidMavenVersion(trimmed) {
		return nil, fmt.Errorf("invalid Maven version format: %s", trimmed)
	}

	elements := parseVersionString(trimmed)
	
	return &Version{
		original: version,
		elements: elements,
	}, nil
}

func isValidMavenVersion(version string) bool {
	// Maven versions should contain at least one digit or be a known qualifier
	hasDigit := false
	hasKnownQualifier := false
	
	for _, r := range version {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	
	// Check for known qualifiers
	lower := strings.ToLower(version)
	knownQualifiers := []string{"alpha", "beta", "milestone", "rc", "snapshot", "ga", "final", "release", "sp"}
	for _, qualifier := range knownQualifiers {
		if strings.Contains(lower, qualifier) {
			hasKnownQualifier = true
			break
		}
	}
	
	// Also accept single-letter qualifiers
	if len(version) == 1 && (version == "a" || version == "b" || version == "m") {
		hasKnownQualifier = true
	}
	
	return hasDigit || hasKnownQualifier
}

func (v *Version) Compare(other *Version) int {
	// Compare elements one by one
	maxLen := len(v.elements)
	if len(other.elements) > maxLen {
		maxLen = len(other.elements)
	}
	
	for i := 0; i < maxLen; i++ {
		var elem1, elem2 element
		
		// Get element or use "null" element if past end
		if i < len(v.elements) {
			elem1 = v.elements[i]
		} else {
			elem1 = element{value: 0, isNumber: true} // null element
		}
		
		if i < len(other.elements) {
			elem2 = other.elements[i]
		} else {
			elem2 = element{value: 0, isNumber: true} // null element
		}
		
		cmp := compareElements(elem1, elem2)
		if cmp != 0 {
			return cmp
		}
	}
	
	return 0 // versions are equal
}

func compareElements(e1, e2 element) int {
	// If both are numbers, compare numerically
	if e1.isNumber && e2.isNumber {
		n1 := e1.value.(int)
		n2 := e2.value.(int)
		if n1 < n2 {
			return -1
		}
		if n1 > n2 {
			return 1
		}
		return 0
	}
	
	// If one is number and other is string, number comes first (unless string is empty/release)
	if e1.isNumber && !e2.isNumber {
		s2 := e2.value.(string)
		if s2 == "" {
			// number vs empty string: empty string (release) is greater
			return -1
		}
		if s2 == "sp" {
			// number vs sp: sp is greater
			return -1
		}
		// number vs other qualifier: number is greater
		return 1
	}
	
	if !e1.isNumber && e2.isNumber {
		s1 := e1.value.(string)
		if s1 == "" {
			// empty string (release) vs number: empty string is greater
			return 1
		}
		if s1 == "sp" {
			// sp vs number: sp is greater
			return 1
		}
		// other qualifier vs number: number is greater
		return -1
	}
	
	// Both are strings - compare by qualifier order
	s1 := e1.value.(string)
	s2 := e2.value.(string)
	
	order1, exists1 := qualifierOrder[s1]
	order2, exists2 := qualifierOrder[s2]
	
	// Unknown qualifiers come after known qualifiers
	if !exists1 && !exists2 {
		// Both unknown - lexicographic comparison
		if s1 < s2 {
			return -1
		}
		if s1 > s2 {
			return 1
		}
		return 0
	}
	
	if !exists1 {
		return 1 // unknown qualifier comes after known
	}
	
	if !exists2 {
		return -1 // known qualifier comes before unknown
	}
	
	// Both are known qualifiers
	if order1 < order2 {
		return -1
	}
	if order1 > order2 {
		return 1
	}
	return 0
}

func (v *Version) String() string {
	return v.original
}

// qualifierOrder defines the precedence of Maven qualifiers
var qualifierOrder = map[string]int{
	"alpha":     1,
	"a":         1,
	"beta":      2,
	"b":         2,
	"milestone": 3,
	"m":         3,
	"rc":        4,
	"cr":        4,
	"snapshot":  5,
	"":          6, // release version (no qualifier)
	"ga":        6,
	"final":     6,
	"release":   6,
	"sp":        7,
}

func parseVersionString(version string) []element {
	var elements []element
	
	// Split by common separators and transitions
	parts := tokenize(version)
	
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		// Normalize qualifiers
		normalized := normalizeQualifier(part)
		
		// Try to parse as number
		if num, err := strconv.Atoi(normalized); err == nil {
			elements = append(elements, element{value: num, isNumber: true})
		} else {
			elements = append(elements, element{value: normalized, isNumber: false})
		}
	}
	
	// Trim trailing null elements (0, "", "final", "ga")
	elements = trimTrailingNulls(elements)
	
	return elements
}

func tokenize(version string) []string {
	var tokens []string
	var current strings.Builder
	
	for i, r := range version {
		switch {
		case r == '.' || r == '-':
			// Add current token if not empty
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		case i > 0:
			prev := rune(version[i-1])
			// Check for transitions between digits and letters
			if (unicode.IsDigit(prev) && unicode.IsLetter(r)) ||
				(unicode.IsLetter(prev) && unicode.IsDigit(r)) {
				// Add current token and start new one
				if current.Len() > 0 {
					tokens = append(tokens, current.String())
					current.Reset()
				}
			}
			current.WriteRune(r)
		default:
			current.WriteRune(r)
		}
	}
	
	// Add final token
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	
	return tokens
}

func normalizeQualifier(s string) string {
	lower := strings.ToLower(s)
	
	// Handle qualifier shortcuts
	switch lower {
	case "a":
		return "alpha"
	case "b":
		return "beta"
	case "m":
		return "milestone"
	case "cr":
		return "rc"
	case "ga", "final", "release":
		return ""
	}
	
	return lower
}

func trimTrailingNulls(elements []element) []element {
	// Remove trailing elements that are equivalent to "null"
	for len(elements) > 0 {
		last := elements[len(elements)-1]
		if isNullElement(last) {
			elements = elements[:len(elements)-1]
		} else {
			break
		}
	}
	return elements
}

func isNullElement(e element) bool {
	if e.isNumber {
		return e.value.(int) == 0
	}
	str := e.value.(string)
	return str == "" || str == "final" || str == "ga" || str == "release"
}