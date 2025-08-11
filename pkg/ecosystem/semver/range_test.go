package semver

import (
	"strings"
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		// Valid ranges
		{"exact match", "1.2.3", false},
		{"greater than", ">1.2.3", false},
		{"greater than or equal", ">=1.2.3", false},
		{"less than", "<1.2.3", false},
		{"less than or equal", "<=1.2.3", false},
		{"not equal", "!=1.2.3", false},
		{"equal explicit", "=1.2.3", false},
		{"wildcard", "*", false},
		{"comma separated", ">=1.0.0,<2.0.0", false},
		{"space separated", ">=1.0.0 <2.0.0", false},
		{"multiple constraints", ">=1.0.0,<2.0.0,!=1.5.0", false},
		{"prerelease constraints", ">=1.0.0-alpha", false},
		{"build metadata constraints", ">=1.0.0+build", false},
		{"valid with whitespace", "  >=1.2.3  ", false},

		// Invalid ranges
		{"empty string", "", true},
		{"only whitespace", "   ", true},
		{"invalid operator", "~1.2.3", true},
		{"missing version after operator", ">=", true},
		{"invalid version", ">invalid", true},
		{"empty constraint in comma list", ">=1.0.0,,<2.0.0", false}, // Should skip empty
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.NewVersionRange(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewVersionRange(%q) expected error, got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("NewVersionRange(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result.String() != strings.TrimSpace(tt.input) {
				t.Errorf("NewVersionRange(%q).String() = %q, want %q", tt.input, result.String(), strings.TrimSpace(tt.input))
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		name        string
		rangeStr    string
		version     string
		expected    bool
		expectError bool // For invalid range or version
	}{
		// Exact matches
		{"exact match - equal", "1.2.3", "1.2.3", true, false},
		{"exact match - not equal", "1.2.3", "1.2.4", false, false},
		{"explicit equal - match", "=1.2.3", "1.2.3", true, false},
		{"explicit equal - no match", "=1.2.3", "1.2.4", false, false},

		// Wildcards
		{"wildcard matches everything", "*", "1.2.3", true, false},
		{"wildcard matches prerelease", "*", "1.2.3-alpha", true, false},

		// Greater than
		{"greater than - true", ">1.2.3", "1.2.4", true, false},
		{"greater than - false equal", ">1.2.3", "1.2.3", false, false},
		{"greater than - false less", ">1.2.3", "1.2.2", false, false},
		{"greater than - major increment", ">1.2.3", "2.0.0", true, false},
		{"greater than with prerelease", ">1.0.0-alpha", "1.0.0-beta", true, false},
		{"greater than with prerelease vs normal", ">1.0.0-alpha", "1.0.0", true, false},

		// Greater than or equal
		{"greater than or equal - true greater", ">=1.2.3", "1.2.4", true, false},
		{"greater than or equal - true equal", ">=1.2.3", "1.2.3", true, false},
		{"greater than or equal - false", ">=1.2.3", "1.2.2", false, false},

		// Less than
		{"less than - true", "<1.2.3", "1.2.2", true, false},
		{"less than - false equal", "<1.2.3", "1.2.3", false, false},
		{"less than - false greater", "<1.2.3", "1.2.4", false, false},
		{"less than with prerelease", "<1.0.0", "1.0.0-alpha", true, false},

		// Less than or equal
		{"less than or equal - true less", "<=1.2.3", "1.2.2", true, false},
		{"less than or equal - true equal", "<=1.2.3", "1.2.3", true, false},
		{"less than or equal - false", "<=1.2.3", "1.2.4", false, false},

		// Not equal
		{"not equal - true", "!=1.2.3", "1.2.4", true, false},
		{"not equal - false", "!=1.2.3", "1.2.3", false, false},

		// Multiple constraints (comma-separated)
		{"comma range - within", ">=1.0.0,<2.0.0", "1.5.0", true, false},
		{"comma range - at lower bound", ">=1.0.0,<2.0.0", "1.0.0", true, false},
		{"comma range - at upper bound", ">=1.0.0,<2.0.0", "2.0.0", false, false},
		{"comma range - below range", ">=1.0.0,<2.0.0", "0.9.0", false, false},
		{"comma range - above range", ">=1.0.0,<2.0.0", "2.0.1", false, false},
		{"comma range with exclusion", ">=1.0.0,<2.0.0,!=1.5.0", "1.4.0", true, false},
		{"comma range with exclusion - excluded", ">=1.0.0,<2.0.0,!=1.5.0", "1.5.0", false, false},

		// Multiple constraints (space-separated)
		{"space range - within", ">=1.0.0 <2.0.0", "1.5.0", true, false},
		{"space range - at lower bound", ">=1.0.0 <2.0.0", "1.0.0", true, false},
		{"space range - at upper bound", ">=1.0.0 <2.0.0", "2.0.0", false, false},
		{"space range - below range", ">=1.0.0 <2.0.0", "0.9.0", false, false},
		{"space range - above range", ">=1.0.0 <2.0.0", "2.0.1", false, false},

		// Prerelease handling
		{"prerelease in range", ">=1.0.0-alpha <1.0.0", "1.0.0-beta", true, false},
		{"prerelease vs normal", ">=1.0.0-alpha <=1.0.0", "1.0.0", true, false},
		{"prerelease comparison", ">=1.0.0-alpha.1 <1.0.0-beta", "1.0.0-alpha.2", true, false},

		// Build metadata (should be ignored in comparisons)
		{"build metadata ignored in range", ">=1.0.0+build1", "1.0.0+build2", true, false},
		{"build metadata ignored in version", ">=1.0.0", "1.0.0+build", true, false},

		// Edge cases
		{"zero version", ">=0.0.0", "0.0.1", true, false},
		{"large version numbers", ">=999.999.999", "999.999.999", true, false},

		// Complex scenarios
		{"complex prerelease range", ">=1.0.0-alpha.1 <1.0.0-beta.1", "1.0.0-alpha.2", true, false},
		{"multiple complex constraints", ">=1.0.0,<2.0.0,!=1.2.3,>=1.1.0", "1.1.5", true, false},
		{"multiple complex constraints - excluded", ">=1.0.0,<2.0.0,!=1.2.3,>=1.1.0", "1.2.3", false, false},
		{"multiple complex constraints - below minimum", ">=1.0.0,<2.0.0,!=1.2.3,>=1.1.0", "1.0.5", false, false},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				if !tt.expectError {
					t.Fatalf("NewVersionRange(%q) unexpected error: %v", tt.rangeStr, err)
				}
				return
			}

			v, err := e.NewVersion(tt.version)
			if err != nil {
				if !tt.expectError {
					t.Fatalf("NewVersion(%q) unexpected error: %v", tt.version, err)
				}
				return
			}

			result := vr.Contains(v)
			if result != tt.expected {
				t.Errorf("VersionRange(%q).Contains(%q) = %v, want %v", tt.rangeStr, tt.version, result, tt.expected)
			}
		})
	}
}

func TestVersionRange_String(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.3", "1.2.3"},
		{">1.2.3", ">1.2.3"},
		{">=1.0.0,<2.0.0", ">=1.0.0,<2.0.0"},
		{">=1.0.0 <2.0.0", ">=1.0.0 <2.0.0"},
		{"*", "*"},
		{"  >=1.0.0  ", ">=1.0.0"}, // Trimmed input
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			vr, err := e.NewVersionRange(tt.input)
			if err != nil {
				t.Fatalf("NewVersionRange(%q) unexpected error: %v", tt.input, err)
			}

			result := vr.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}
