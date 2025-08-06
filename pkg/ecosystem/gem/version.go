package gem

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// versionPattern matches Ruby Gem version strings
// Supports: 1.2.3, 1.2.3.4, 1.2.3-alpha, 1.2.3.pre, etc.
var versionPattern = regexp.MustCompile(`^v?(\d+(?:\.\d+)*(?:\.[a-zA-Z]+\d*)*(?:-[a-zA-Z0-9.-]+)*(?:\+[a-zA-Z0-9.-]+)*)$`)

// Version represents a Ruby Gem package version
type Version struct {
	segments []segment
	original string
}

// segment represents a version segment (numeric or string)
type segment struct {
	value     string
	isNumeric bool
	numValue  int
}

// NewVersion creates a new Ruby Gem version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)
	
	// Remove leading v prefix
	version = strings.TrimPrefix(version, "v")
	
	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}
	
	// Basic validation
	if !versionPattern.MatchString("v"+version) {
		return nil, fmt.Errorf("invalid Ruby Gem version: %s", original)
	}
	
	// Canonicalize and parse segments
	canonical := canonicalizeVersion(version)
	segments, err := parseSegments(canonical)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %s: %v", original, err)
	}
	
	return &Version{
		segments: segments,
		original: original,
	}, nil
}

// canonicalizeVersion transforms version string to canonical form
func canonicalizeVersion(version string) string {
	// Handle prerelease indicators (-, +)
	parts := strings.FieldsFunc(version, func(r rune) bool {
		return r == '-' || r == '+'
	})
	
	if len(parts) == 0 {
		return version
	}
	
	// Process main version part
	main := parts[0]
	result := addDotsBetweenNumericAndAlpha(main)
	
	// Add prerelease/build parts back
	for i := 1; i < len(parts); i++ {
		if strings.Contains(version, "-"+parts[i]) {
			result += "-" + addDotsBetweenNumericAndAlpha(parts[i])
		} else {
			result += "+" + addDotsBetweenNumericAndAlpha(parts[i])
		}
	}
	
	return result
}

// addDotsBetweenNumericAndAlpha adds dots between numeric and alphabetic segments
func addDotsBetweenNumericAndAlpha(s string) string {
	if len(s) == 0 {
		return s
	}
	
	var result strings.Builder
	var prev rune
	
	for i, r := range s {
		if i > 0 {
			isCurrentNumeric := r >= '0' && r <= '9'
			isPrevNumeric := prev >= '0' && prev <= '9'
			
			// Add dot if transitioning between numeric and alpha
			if (isCurrentNumeric && !isPrevNumeric) || (!isCurrentNumeric && isPrevNumeric) {
				if prev != '.' && r != '.' {
					result.WriteRune('.')
				}
			}
		}
		result.WriteRune(r)
		prev = r
	}
	
	return result.String()
}

// parseSegments parses canonical version into segments
func parseSegments(version string) ([]segment, error) {
	var segments []segment
	
	// First handle prerelease/build separators at top level
	mainPart := version
	prereleasePart := ""
	buildPart := ""
	
	// Extract build metadata (after +)
	if plusIndex := strings.Index(version, "+"); plusIndex != -1 {
		buildPart = version[plusIndex+1:]
		mainPart = version[:plusIndex]
	}
	
	// Extract prerelease (after -)
	if dashIndex := strings.Index(mainPart, "-"); dashIndex != -1 {
		prereleasePart = mainPart[dashIndex+1:]
		mainPart = mainPart[:dashIndex]
	}
	
	// Parse main version parts (numeric segments)
	parts := strings.Split(mainPart, ".")
	for _, part := range parts {
		if part == "" {
			continue
		}
		
		// Check if this part contains letters (prerelease indicator)
		if containsLetter(part) {
			// This is a prerelease segment
			segments = append(segments, createSegment(part))
		} else {
			// This is a numeric segment
			segments = append(segments, createSegment(part))
		}
	}
	
	// Add prerelease segments
	if prereleasePart != "" {
		prereleaseParts := strings.Split(prereleasePart, ".")
		for _, part := range prereleaseParts {
			if part != "" {
				// Prerelease parts are always treated as non-numeric for comparison purposes
				segments = append(segments, segment{
					value:     strings.ToLower(part),
					isNumeric: false,
					numValue:  0,
				})
			}
		}
	}
	
	// Add build segments
	if buildPart != "" {
		buildParts := strings.Split(buildPart, ".")
		for _, part := range buildParts {
			if part != "" {
				segments = append(segments, createSegment(part))
			}
		}
	}
	
	// Remove trailing zero segments from numeric part only
	segments = removeTrailingZeros(segments)
	
	return segments, nil
}

// containsLetter checks if string contains any letter
func containsLetter(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

// createSegment creates a segment from a string part
func createSegment(part string) segment {
	if numValue, err := strconv.Atoi(part); err == nil {
		return segment{
			value:     part,
			isNumeric: true,
			numValue:  numValue,
		}
	}
	return segment{
		value:     strings.ToLower(part), // Case-insensitive comparison
		isNumeric: false,
		numValue:  0,
	}
}

// removeTrailingZeros removes trailing zero segments
func removeTrailingZeros(segments []segment) []segment {
	for len(segments) > 1 && segments[len(segments)-1].isNumeric && segments[len(segments)-1].numValue == 0 {
		segments = segments[:len(segments)-1]
	}
	return segments
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another Ruby Gem version
func (v *Version) Compare(other *Version) int {
	// First compare the numeric parts
	vNumeric, vPrerelease := v.splitNumericAndPrerelease()
	oNumeric, oPrerelease := other.splitNumericAndPrerelease()
	
	// Compare numeric parts first
	numericCmp := compareSegmentArrays(vNumeric, oNumeric)
	if numericCmp != 0 {
		return numericCmp
	}
	
	// If numeric parts are equal, compare prerelease parts
	// No prerelease > prerelease
	if len(vPrerelease) == 0 && len(oPrerelease) == 0 {
		return 0
	}
	if len(vPrerelease) == 0 {
		return 1 // release > prerelease
	}
	if len(oPrerelease) == 0 {
		return -1 // prerelease < release
	}
	
	// Both have prerelease, compare them
	return compareSegmentArrays(vPrerelease, oPrerelease)
}

// splitNumericAndPrerelease splits version into numeric and prerelease parts
func (v *Version) splitNumericAndPrerelease() ([]segment, []segment) {
	var numeric, prerelease []segment
	
	for _, seg := range v.segments {
		if seg.isNumeric {
			numeric = append(numeric, seg)
		} else {
			prerelease = append(prerelease, seg)
		}
	}
	
	return numeric, prerelease
}

// compareSegmentArrays compares two arrays of segments
func compareSegmentArrays(a, b []segment) int {
	maxLen := max(len(a), len(b))
	
	for i := range maxLen {
		var aSeg, bSeg segment
		
		if i < len(a) {
			aSeg = a[i]
		} else {
			aSeg = segment{value: "0", isNumeric: true, numValue: 0}
		}
		
		if i < len(b) {
			bSeg = b[i]
		} else {
			bSeg = segment{value: "0", isNumeric: true, numValue: 0}
		}
		
		cmp := compareSegments(aSeg, bSeg)
		if cmp != 0 {
			return cmp
		}
	}
	
	return 0
}

// compareSegments compares two version segments
func compareSegments(a, b segment) int {
	// Both numeric
	if a.isNumeric && b.isNumeric {
		return compareInt(a.numValue, b.numValue)
	}
	
	// One numeric, one string - in prerelease context, strings have precedence
	if a.isNumeric && !b.isNumeric {
		return -1
	}
	if !a.isNumeric && b.isNumeric {
		return 1
	}
	
	// Both strings - lexical comparison
	return strings.Compare(a.value, b.value)
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
