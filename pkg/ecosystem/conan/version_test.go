package conan

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		expected  *Version
	}{
		// Valid versions - basic patterns
		{
			name:      "single digit",
			input:     "1",
			wantError: false,
			expected:  &Version{parts: []string{"1"}, prerelease: "", build: "", original: "1"},
		},
		{
			name:      "two parts",
			input:     "1.2",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2"}, prerelease: "", build: "", original: "1.2"},
		},
		{
			name:      "three parts",
			input:     "1.2.3",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "", build: "", original: "1.2.3"},
		},
		{
			name:      "four parts",
			input:     "1.2.3.4",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3", "4"}, prerelease: "", build: "", original: "1.2.3.4"},
		},
		{
			name:      "many parts",
			input:     "1.2.3.4.5.6",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3", "4", "5", "6"}, prerelease: "", build: "", original: "1.2.3.4.5.6"},
		},

		// Valid versions - with letters
		{
			name:      "version with letter",
			input:     "1.2.3a",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3a"}, prerelease: "", build: "", original: "1.2.3a"},
		},
		{
			name:      "version with letter in middle",
			input:     "1.2a.3",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2a", "3"}, prerelease: "", build: "", original: "1.2a.3"},
		},
		{
			name:      "version with multiple letters",
			input:     "1.2abc.3def",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2abc", "3def"}, prerelease: "", build: "", original: "1.2abc.3def"},
		},
		{
			name:      "openssl style version",
			input:     "1.0.2n",
			wantError: false,
			expected:  &Version{parts: []string{"1", "0", "2n"}, prerelease: "", build: "", original: "1.0.2n"},
		},

		// Valid versions - with prerelease
		{
			name:      "version with prerelease",
			input:     "1.2.3-alpha",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "alpha", build: "", original: "1.2.3-alpha"},
		},
		{
			name:      "version with complex prerelease",
			input:     "1.2.3-alpha.1.2",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "alpha.1.2", build: "", original: "1.2.3-alpha.1.2"},
		},
		{
			name:      "version with prerelease containing hyphens",
			input:     "1.2.3-alpha-1",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "alpha-1", build: "", original: "1.2.3-alpha-1"},
		},

		// Valid versions - with build metadata
		{
			name:      "version with build metadata",
			input:     "1.2.3+build",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "", build: "build", original: "1.2.3+build"},
		},
		{
			name:      "version with complex build metadata",
			input:     "1.2.3+build.1.2",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "", build: "build.1.2", original: "1.2.3+build.1.2"},
		},

		// Valid versions - with both prerelease and build
		{
			name:      "version with prerelease and build",
			input:     "1.2.3-alpha+build",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "alpha", build: "build", original: "1.2.3-alpha+build"},
		},
		{
			name:      "complex version with all components",
			input:     "1.2.3.a.5-beta.1+build.456",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3", "a", "5"}, prerelease: "beta.1", build: "build.456", original: "1.2.3.a.5-beta.1+build.456"},
		},

		// Valid versions - case sensitivity and whitespace
		{
			name:      "uppercase converted to lowercase",
			input:     "1.2.3-ALPHA",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "alpha", build: "", original: "1.2.3-ALPHA"},
		},
		{
			name:      "mixed case converted to lowercase",
			input:     "1.2.3A-Beta+BUILD",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3a"}, prerelease: "beta", build: "build", original: "1.2.3A-Beta+BUILD"},
		},
		{
			name:      "whitespace trimmed",
			input:     "  1.2.3  ",
			wantError: false,
			expected:  &Version{parts: []string{"1", "2", "3"}, prerelease: "", build: "", original: "  1.2.3  "},
		},

		// Valid versions - zero components
		{
			name:      "zero version",
			input:     "0.0.0",
			wantError: false,
			expected:  &Version{parts: []string{"0", "0", "0"}, prerelease: "", build: "", original: "0.0.0"},
		},

		// Invalid versions
		{
			name:      "empty string",
			input:     "",
			wantError: true,
		},
		{
			name:      "only whitespace",
			input:     "   ",
			wantError: true,
		},
		{
			name:      "invalid characters - symbols",
			input:     "1.2.3@",
			wantError: true,
		},
		{
			name:      "invalid characters - spaces in version",
			input:     "1.2 .3",
			wantError: true,
		},
		{
			name:      "empty version part",
			input:     "1..3",
			wantError: true,
		},
		{
			name:      "trailing dot",
			input:     "1.2.3.",
			wantError: true,
		},
		{
			name:      "leading dot",
			input:     ".1.2.3",
			wantError: true,
		},
		{
			name:      "only dots",
			input:     "...",
			wantError: true,
		},
		{
			name:      "empty prerelease",
			input:     "1.2.3-",
			wantError: true,
		},
		{
			name:      "empty build metadata",
			input:     "1.2.3+",
			wantError: true,
		},
		{
			name:      "invalid characters in prerelease",
			input:     "1.2.3-alpha@beta",
			wantError: true,
		},
		{
			name:      "invalid characters in build metadata",
			input:     "1.2.3+build@1",
			wantError: true,
		},
		{
			name:      "empty prerelease identifier",
			input:     "1.2.3-alpha..beta",
			wantError: true,
		},
		{
			name:      "empty build metadata identifier",
			input:     "1.2.3+build..1",
			wantError: true,
		},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := e.NewVersion(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("NewVersion(%q) expected error, got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("NewVersion(%q) unexpected error: %v", tt.input, err)
				return
			}

			// Compare parts
			if len(result.parts) != len(tt.expected.parts) {
				t.Errorf("NewVersion(%q).parts length = %v, want %v", tt.input, len(result.parts), len(tt.expected.parts))
			} else {
				for i, part := range result.parts {
					if part != tt.expected.parts[i] {
						t.Errorf("NewVersion(%q).parts[%d] = %v, want %v", tt.input, i, part, tt.expected.parts[i])
					}
				}
			}

			if result.prerelease != tt.expected.prerelease {
				t.Errorf("NewVersion(%q).prerelease = %v, want %v", tt.input, result.prerelease, tt.expected.prerelease)
			}
			if result.build != tt.expected.build {
				t.Errorf("NewVersion(%q).build = %v, want %v", tt.input, result.build, tt.expected.build)
			}
			if result.original != tt.expected.original {
				t.Errorf("NewVersion(%q).original = %v, want %v", tt.input, result.original, tt.expected.original)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name     string
		version1 string
		version2 string
		expected int
	}{
		// Basic numeric comparisons
		{"equal versions", "1.2.3", "1.2.3", 0},
		{"major version difference", "2.0.0", "1.9.9", 1},
		{"major version difference (less)", "1.9.9", "2.0.0", -1},
		{"minor version difference", "1.2.0", "1.1.9", 1},
		{"minor version difference (less)", "1.1.9", "1.2.0", -1},
		{"patch version difference", "1.2.3", "1.2.2", 1},
		{"patch version difference (less)", "1.2.2", "1.2.3", -1},

		// Different number of parts
		{"fewer parts vs more parts (equal)", "1.2", "1.2.0", 0},
		{"fewer parts vs more parts (less)", "1.1", "1.2.0", -1},
		{"fewer parts vs more parts (greater)", "1.3", "1.2.0", 1},
		{"many parts vs few parts", "1.2.3.4.5", "1.2.3", 1},
		{"single part vs multi part", "2", "1.9.9", 1},

		// Numeric vs alphabetic parts
		{"numeric vs alphabetic", "1.2.3", "1.2.3a", -1},
		{"alphabetic vs numeric", "1.2.3a", "1.2.3", 1},
		{"alphabetic comparison", "1.2.3a", "1.2.3b", -1},
		{"alphabetic comparison (reverse)", "1.2.3b", "1.2.3a", 1},

		// Extended Conan version formats
		{"four parts comparison", "1.2.3.4", "1.2.3.5", -1},
		{"openssl style versions", "1.0.2n", "1.0.2m", 1},
		{"mixed numeric and alphabetic", "1.2a.3", "1.2b.3", -1},
		{"complex alphabetic parts", "1.2abc.3", "1.2abd.3", -1},

		// Numeric ordering in parts
		{"numeric ordering in parts", "1.2.11", "1.2.9", 1},
		{"large numbers", "1.999.999", "1.1000.0", -1},

		// Natural comparison fixes - mixed numeric/alphabetic parts
		{"numeric vs alphanumeric with same base", "1.2.10", "1.2.3a", 1},  // 10 > 3a
		{"numeric vs alphanumeric different base", "1.2.3", "1.2.10a", -1}, // 3 < 10a
		{"pure numeric vs pure alphabetic", "1.2.3", "1.2.a", -1},          // numeric < alphabetic
		{"alphanumeric vs pure numeric", "1.2.3a", "1.2.4", -1},            // 3a < 4

		// Prerelease comparisons
		{"normal version vs prerelease", "1.0.0", "1.0.0-alpha", 1},
		{"prerelease vs normal version", "1.0.0-alpha", "1.0.0", -1},
		{"prerelease identifiers - numeric", "1.0.0-1", "1.0.0-2", -1},
		{"prerelease identifiers - numeric vs alpha", "1.0.0-1", "1.0.0-alpha", -1},
		{"prerelease identifiers - alpha vs numeric", "1.0.0-alpha", "1.0.0-1", 1},
		{"prerelease identifiers - alpha vs beta", "1.0.0-alpha", "1.0.0-beta", -1},
		{"prerelease identifiers - different lengths", "1.0.0-alpha.1", "1.0.0-alpha", 1},
		{"prerelease identifiers - shorter vs longer", "1.0.0-alpha", "1.0.0-alpha.1", -1},

		// Complex prerelease comparisons
		{"complex prerelease 1", "1.0.0-alpha.beta.1", "1.0.0-alpha.beta.2", -1},
		{"complex prerelease 2", "1.0.0-alpha.1.2", "1.0.0-alpha.1.10", -1},
		{"complex prerelease 3", "1.0.0-alpha.1.beta", "1.0.0-alpha.1.1", 1},

		// Build metadata should be ignored
		{"build metadata ignored 1", "1.0.0+build1", "1.0.0+build2", 0},
		{"build metadata ignored 2", "1.0.0-alpha+build1", "1.0.0-alpha+build2", 0},
		{"build metadata vs no build", "1.0.0+build", "1.0.0", 0},

		// Zero versions
		{"zero versions", "0.0.0", "0.0.1", -1},
		{"zero major", "0.1.0", "1.0.0", -1},

		// Edge cases
		{"single digit versions", "1", "2", -1},
		{"single vs multi part", "1", "0.9", 1},
		{"letters in single part", "1a", "1b", -1},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := e.NewVersion(tt.version1)
			if err != nil {
				t.Fatalf("NewVersion(%q) unexpected error: %v", tt.version1, err)
			}

			v2, err := e.NewVersion(tt.version2)
			if err != nil {
				t.Fatalf("NewVersion(%q) unexpected error: %v", tt.version2, err)
			}

			result := v1.Compare(v2)
			if result != tt.expected {
				t.Errorf("Compare(%q, %q) = %d, want %d", tt.version1, tt.version2, result, tt.expected)
			}

			// Test symmetry (reverse comparison should give opposite result)
			reverseResult := v2.Compare(v1)
			expectedReverse := -tt.expected
			if reverseResult != expectedReverse {
				t.Errorf("Compare(%q, %q) = %d, want %d (reverse test)", tt.version2, tt.version1, reverseResult, expectedReverse)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.3", "1.2.3"},
		{"1.2.3a", "1.2.3a"},
		{"1.2.3-alpha", "1.2.3-alpha"},
		{"1.2.3+build", "1.2.3+build"},
		{"1.2.3-alpha+build", "1.2.3-alpha+build"},
		{"1.2.3.4.5", "1.2.3.4.5"},
		{"1.0.2n", "1.0.2n"},
		{"  1.2.3  ", "  1.2.3  "},                 // Should preserve original with whitespace
		{"1.2.3A-Beta+BUILD", "1.2.3A-Beta+BUILD"}, // Should preserve original case
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			v, err := e.NewVersion(tt.input)
			if err != nil {
				t.Fatalf("NewVersion(%q) unexpected error: %v", tt.input, err)
			}

			result := v.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}
