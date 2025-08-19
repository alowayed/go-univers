package conan

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		// Valid ranges - simple constraints
		{"exact version", "1.2.3", false},
		{"greater than", ">1.2.3", false},
		{"greater than or equal", ">=1.2.3", false},
		{"less than", "<2.0.0", false},
		{"less than or equal", "<=2.0.0", false},
		{"tilde constraint", "~1.2.3", false},
		{"caret constraint", "^1.2.3", false},

		// Valid ranges - multiple constraints
		{"comma separated", ">=1.2.0, <2.0.0", false},
		{"space separated", ">=1.2.0 <2.0.0", false},
		{"mixed separators", ">=1.2.0, <2.0.0 !=1.5.0", false},

		// Valid ranges - OR logic
		{"or logic", ">1.0.0 <2.0.0 || ^3.2.0", false},
		{"complex or", ">=1.0.0 <1.5.0 || >=2.0.0 <2.5.0", false},

		// Valid ranges - with letters and extended formats
		{"version with letter", ">=1.2.3a", false},
		{"openssl style", ">=1.0.2n", false},
		{"many parts", ">=1.2.3.4.5", false},

		// Valid ranges - whitespace handling
		{"whitespace around operators", " >= 1.2.3 ", false},
		{"extra whitespace", "  >=  1.2.3  ,  <  2.0.0  ", false},

		// Invalid ranges
		{"empty string", "", true},
		{"only whitespace", "   ", true},
		{"invalid operator", "~=1.2.3", true},
		{"invalid version in constraint", ">invalid@version", true},
		{"only operator", ">", true},
		{"malformed constraint", ">=1.2.3<", true},
		{"double operators", ">>1.2.3", true},
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

			if result == nil {
				t.Errorf("NewVersionRange(%q) returned nil result", tt.input)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		version  string
		expected bool
	}{
		// Exact matches
		{"exact match", "1.2.3", "1.2.3", true},
		{"exact no match", "1.2.3", "1.2.4", false},

		// Greater than constraints
		{"greater than - true", ">1.2.3", "1.2.4", true},
		{"greater than - false (equal)", ">1.2.3", "1.2.3", false},
		{"greater than - false (less)", ">1.2.3", "1.2.2", false},

		// Greater than or equal constraints
		{"greater than or equal - true (greater)", ">=1.2.3", "1.2.4", true},
		{"greater than or equal - true (equal)", ">=1.2.3", "1.2.3", true},
		{"greater than or equal - false", ">=1.2.3", "1.2.2", false},

		// Less than constraints
		{"less than - true", "<2.0.0", "1.9.9", true},
		{"less than - false (equal)", "<2.0.0", "2.0.0", false},
		{"less than - false (greater)", "<2.0.0", "2.0.1", false},

		// Less than or equal constraints
		{"less than or equal - true (less)", "<=2.0.0", "1.9.9", true},
		{"less than or equal - true (equal)", "<=2.0.0", "2.0.0", true},
		{"less than or equal - false", "<=2.0.0", "2.0.1", false},

		// Tilde constraints
		{"tilde - patch level change allowed", "~1.2.3", "1.2.5", true},
		{"tilde - patch level equal", "~1.2.3", "1.2.3", true},
		{"tilde - minor change not allowed", "~1.2.3", "1.3.0", false},
		{"tilde - major change not allowed", "~1.2.3", "2.0.0", false},
		{"tilde - less than base", "~1.2.3", "1.2.2", false},

		// Caret constraints
		{"caret - compatible change allowed", "^1.2.3", "1.5.0", true},
		{"caret - equal version", "^1.2.3", "1.2.3", true},
		{"caret - major change not allowed", "^1.2.3", "2.0.0", false},
		{"caret - less than base", "^1.2.3", "1.2.2", false},
		{"caret - zero major", "^0.2.3", "0.2.5", true},
		{"caret - zero major, minor change not allowed", "^0.2.3", "0.3.0", false},

		// Multiple constraints (AND logic)
		{"multiple constraints - both satisfied", ">=1.2.0, <2.0.0", "1.5.0", true},
		{"multiple constraints - first not satisfied", ">=1.2.0, <2.0.0", "1.1.0", false},
		{"multiple constraints - second not satisfied", ">=1.2.0, <2.0.0", "2.1.0", false},
		{"three constraints - all satisfied", ">=1.0.0, <2.0.0, !=1.5.0", "1.4.0", true},

		// Extended Conan version formats
		{"version with letter", ">=1.2.3a", "1.2.3b", true},
		{"version with letter - false", ">=1.2.3b", "1.2.3a", false},
		{"openssl style version", ">=1.0.2n", "1.0.2o", true},
		{"many parts version", ">=1.2.3.4.5", "1.2.3.4.6", true},

		// With prerelease versions
		{"prerelease in range", ">=1.0.0-alpha", "1.0.0-beta", true},
		{"prerelease vs release", ">1.0.0-alpha", "1.0.0", true},
		{"release vs prerelease", "<1.0.0", "1.0.0-alpha", true},

		// Edge cases
		{"single digit version", ">=1", "2", true},
		{"single digit version - false", ">=2", "1", false},
		{"missing parts treated as zero", ">=1.2", "1.2.0", true},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vRange, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange(%q) unexpected error: %v", tt.rangeStr, err)
			}

			version, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion(%q) unexpected error: %v", tt.version, err)
			}

			result := vRange.Contains(version)
			if result != tt.expected {
				t.Errorf("Range(%q).Contains(%q) = %v, want %v", tt.rangeStr, tt.version, result, tt.expected)
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
		{">=1.2.0, <2.0.0", ">=1.2.0, <2.0.0"},
		{"~1.2.3", "~1.2.3"},
		{"^1.2.3", "^1.2.3"},
		{"  >=1.2.3  ", "  >=1.2.3  "}, // Should preserve original
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			vRange, err := e.NewVersionRange(tt.input)
			if err != nil {
				t.Fatalf("NewVersionRange(%q) unexpected error: %v", tt.input, err)
			}

			result := vRange.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test complex scenarios that might be specific to Conan
func TestVersionRange_ConanSpecificScenarios(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		version  string
		expected bool
	}{
		// Test tilde behavior with different part counts
		{"tilde with two parts", "~1.2", "1.2.5", true},
		{"tilde with two parts - minor change", "~1.2", "1.3.0", false},
		{"tilde with one part", "~1", "1.9.9", true},
		{"tilde with one part - major change", "~1", "2.0.0", false},

		// Test caret behavior with zeros
		{"caret with leading zero", "^0.0.3", "0.0.4", true},
		{"caret with leading zero - patch change", "^0.0.3", "0.1.0", false},
		{"caret all zeros", "^0.0.0", "0.0.1", true},

		// Test mixed numeric/alphabetic version parts
		{"mixed parts in constraint", ">=1.2a.3", "1.2b.3", true},
		{"mixed parts - numeric vs alpha", ">=1.2.3", "1.2.3a", true},
		{"mixed parts - alpha vs numeric", ">=1.2.3a", "1.2.4", false},

		// Complex constraints with extended versions
		{"complex constraint with letters", ">=1.0.2n, <1.1.0", "1.0.2o", true},
		{"complex constraint with many parts", ">=1.2.3.4, <1.2.4", "1.2.3.5", true},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vRange, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange(%q) unexpected error: %v", tt.rangeStr, err)
			}

			version, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion(%q) unexpected error: %v", tt.version, err)
			}

			result := vRange.Contains(version)
			if result != tt.expected {
				t.Errorf("Range(%q).Contains(%q) = %v, want %v", tt.rangeStr, tt.version, result, tt.expected)
			}
		})
	}
}
