package alpm

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// versionPattern matches ALMP version strings
// Format: [epoch:]pkgver[-pkgrel]
// Note: pkgrel is always numeric, so we split on the last hyphen followed by digits
var versionPattern = regexp.MustCompile(`^(?:(\d+):)?(.+?)(?:-(\d+))?$`)

// Version represents an ALMP package version
type Version struct {
	epoch    int    // optional epoch (defaults to 0)
	pkgver   string // package version (upstream software version)
	pkgrel   int    // optional package release number (defaults to 0)
	original string // original version string
}

// NewVersion creates a new ALMP version from a string
func (e *Ecosystem) NewVersion(version string) (*Version, error) {
	original := version
	version = strings.TrimSpace(version)

	if version == "" {
		return nil, fmt.Errorf("empty version string")
	}

	// Split on epoch first
	var epochStr string
	var versionPart string

	if colonIndex := strings.Index(version, ":"); colonIndex != -1 {
		epochStr = version[:colonIndex]
		versionPart = version[colonIndex+1:]
	} else {
		epochStr = ""
		versionPart = version
	}

	// Split version part on last hyphen followed by digits (pkgrel)
	var pkgver string
	var pkgrelStr string

	// Find the last hyphen followed by only digits
	lastValidHyphen := -1
	for i := len(versionPart) - 1; i >= 0; i-- {
		if versionPart[i] == '-' && i+1 < len(versionPart) {
			// Check if everything after this hyphen is digits
			afterHyphen := versionPart[i+1:]
			if len(afterHyphen) > 0 && isAllDigits(afterHyphen) {
				lastValidHyphen = i
				break
			}
		}
	}

	if lastValidHyphen != -1 {
		pkgver = versionPart[:lastValidHyphen]
		pkgrelStr = versionPart[lastValidHyphen+1:]
	} else {
		pkgver = versionPart
		pkgrelStr = ""
	}

	// Parse epoch (default to 0)
	epoch := 0
	if epochStr != "" {
		var err error
		epoch, err = strconv.Atoi(epochStr)
		if err != nil {
			return nil, fmt.Errorf("invalid epoch in version %s: %v", original, err)
		}
		if epoch < 0 {
			return nil, fmt.Errorf("epoch cannot be negative in version %s", original)
		}
	}

	// Validate pkgver (cannot be empty)
	if pkgver == "" {
		return nil, fmt.Errorf("pkgver cannot be empty in version %s", original)
	}

	// Validate characters in pkgver - ALMP allows alphanumerics, periods, underscores
	if err := validateALMPVersionString(pkgver, "pkgver"); err != nil {
		return nil, fmt.Errorf("invalid pkgver in %s: %v", original, err)
	}

	// Parse pkgrel (default to 0)
	pkgrel := 0
	if pkgrelStr != "" {
		var err error
		pkgrel, err = strconv.Atoi(pkgrelStr)
		if err != nil {
			return nil, fmt.Errorf("invalid pkgrel in version %s: %v", original, err)
		}
		if pkgrel < 0 {
			return nil, fmt.Errorf("pkgrel cannot be negative in version %s", original)
		}
	}

	return &Version{
		epoch:    epoch,
		pkgver:   pkgver,
		pkgrel:   pkgrel,
		original: original,
	}, nil
}

// validateALMPVersionString validates that a version string contains only allowed characters
// ALMP allows: alphanumerics, periods, underscores, plus signs, hyphens
func validateALMPVersionString(s, part string) error {
	for _, r := range s {
		if !isValidALMPVersionChar(r) {
			return fmt.Errorf("invalid character %q in %s", r, part)
		}
	}
	return nil
}

// isValidALMPVersionChar checks if a character is valid in an ALMP version string
// Based on Arch Linux PKGBUILD specification
func isValidALMPVersionChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '_' || r == '+' || r == '-'
}

// isAllDigits checks if a string contains only digits
func isAllDigits(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.original
}

// Compare compares this version with another ALMP version using vercmp rules
// Follows Arch Linux vercmp(8) algorithm:
// 1. Compare epochs first (higher epoch wins)
// 2. Compare pkgver parts using vercmp's alphanumeric rules
// 3. Compare pkgrel parts numerically
func (v *Version) Compare(other *Version) int {
	// 1. Compare epochs first
	if v.epoch != other.epoch {
		if v.epoch < other.epoch {
			return -1
		}
		return 1
	}

	// 2. Compare pkgver parts using vercmp rules
	pkgverCmp := compareALMPVersionString(v.pkgver, other.pkgver)
	if pkgverCmp != 0 {
		return pkgverCmp
	}

	// 3. Compare pkgrel parts numerically
	if v.pkgrel < other.pkgrel {
		return -1
	}
	if v.pkgrel > other.pkgrel {
		return 1
	}

	return 0
}

// compareALMPVersionString compares two ALMP version strings using vercmp rules
// This implements the Arch Linux vercmp algorithm which alternates between
// comparing non-numeric and numeric segments
func compareALMPVersionString(a, b string) int {
	i, j := 0, 0

	for i < len(a) || j < len(b) {
		// Extract non-digit segments
		iStart := i
		for i < len(a) && !unicode.IsDigit(rune(a[i])) {
			i++
		}
		aNonDigit := a[iStart:i]

		jStart := j
		for j < len(b) && !unicode.IsDigit(rune(b[j])) {
			j++
		}
		bNonDigit := b[jStart:j]

		// Compare non-digit segments with vercmp alphanumeric precedence
		nonDigitCmp := compareALMPNonDigits(aNonDigit, bNonDigit)
		if nonDigitCmp != 0 {
			return nonDigitCmp
		}

		// Extract digit segments
		iStart = i
		for i < len(a) && unicode.IsDigit(rune(a[i])) {
			i++
		}
		aDigit := a[iStart:i]

		jStart = j
		for j < len(b) && unicode.IsDigit(rune(b[j])) {
			j++
		}
		bDigit := b[jStart:j]

		// Compare digit segments numerically
		digitCmp := compareALMPDigits(aDigit, bDigit)
		if digitCmp != 0 {
			return digitCmp
		}
	}

	return 0
}

// compareALMPNonDigits compares non-digit segments using vercmp alphanumeric precedence
// vercmp precedence: 1.0a < 1.0b < 1.0beta < 1.0p < 1.0pre < 1.0rc < 1.0 < 1.0.a < 1.0.1
func compareALMPNonDigits(a, b string) int {
	// Empty strings sort before non-empty strings in vercmp
	// This handles cases like "1.0" vs "1.0a" where "1.0" should be greater
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return 1 // empty (like end of "1.0") sorts after non-empty (like "a" in "1.0a")
	}
	if b == "" {
		return -1 // non-empty sorts before empty
	}

	// Compare lexicographically - this gives us the correct alphanumeric precedence
	return strings.Compare(a, b)
}

// compareALMPDigits compares digit strings numerically
func compareALMPDigits(a, b string) int {
	// Empty string is treated as 0
	if a == "" && b == "" {
		return 0
	}
	if a == "" {
		return -1
	}
	if b == "" {
		return 1
	}

	// Convert to integers for comparison
	aNum, aErr := strconv.ParseUint(a, 10, 64)
	bNum, bErr := strconv.ParseUint(b, 10, 64)

	if aErr == nil && bErr == nil {
		if aNum < bNum {
			return -1
		}
		if aNum > bNum {
			return 1
		}
		return 0
	}

	// Fallback for very large numbers that don't fit in uint64
	// Compare by length first (longer number is larger)
	if len(a) < len(b) {
		return -1
	}
	if len(a) > len(b) {
		return 1
	}

	// If lengths are equal, string comparison works for digits
	return strings.Compare(a, b)
}
