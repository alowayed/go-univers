package semver

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
		// Valid versions
		{
			name:      "basic version",
			input:     "1.2.3",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "", build: "", original: "1.2.3"},
		},
		{
			name:      "version with prerelease",
			input:     "1.2.3-alpha",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "alpha", build: "", original: "1.2.3-alpha"},
		},
		{
			name:      "version with build metadata",
			input:     "1.2.3+build.1",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "", build: "build.1", original: "1.2.3+build.1"},
		},
		{
			name:      "version with prerelease and build metadata",
			input:     "1.2.3-alpha.1+build.1",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "alpha.1", build: "build.1", original: "1.2.3-alpha.1+build.1"},
		},
		{
			name:      "zero versions",
			input:     "0.0.0",
			wantError: false,
			expected:  &Version{major: 0, minor: 0, patch: 0, prerelease: "", build: "", original: "0.0.0"},
		},
		{
			name:      "large numbers",
			input:     "999.999.999",
			wantError: false,
			expected:  &Version{major: 999, minor: 999, patch: 999, prerelease: "", build: "", original: "999.999.999"},
		},
		{
			name:      "complex prerelease",
			input:     "1.2.3-alpha.beta.1.2.3",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "alpha.beta.1.2.3", build: "", original: "1.2.3-alpha.beta.1.2.3"},
		},
		{
			name:      "prerelease with hyphens",
			input:     "1.2.3-alpha-1",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "alpha-1", build: "", original: "1.2.3-alpha-1"},
		},
		{
			name:      "build metadata with hyphens",
			input:     "1.2.3+build-1",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "", build: "build-1", original: "1.2.3+build-1"},
		},
		{
			name:      "whitespace handling",
			input:     "  1.2.3  ",
			wantError: false,
			expected:  &Version{major: 1, minor: 2, patch: 3, prerelease: "", build: "", original: "  1.2.3  "},
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
			name:      "leading zeros in major",
			input:     "01.2.3",
			wantError: true,
		},
		{
			name:      "leading zeros in minor",
			input:     "1.02.3",
			wantError: true,
		},
		{
			name:      "leading zeros in patch",
			input:     "1.2.03",
			wantError: true,
		},
		{
			name:      "negative major",
			input:     "-1.2.3",
			wantError: true,
		},
		{
			name:      "negative minor",
			input:     "1.-2.3",
			wantError: true,
		},
		{
			name:      "negative patch",
			input:     "1.2.-3",
			wantError: true,
		},
		{
			name:      "missing minor",
			input:     "1",
			wantError: true,
		},
		{
			name:      "missing patch",
			input:     "1.2",
			wantError: true,
		},
		{
			name:      "too many components",
			input:     "1.2.3.4",
			wantError: true,
		},
		{
			name:      "non-numeric major",
			input:     "a.2.3",
			wantError: true,
		},
		{
			name:      "non-numeric minor",
			input:     "1.a.3",
			wantError: true,
		},
		{
			name:      "non-numeric patch",
			input:     "1.2.a",
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
			name:      "leading zeros in numeric prerelease",
			input:     "1.2.3-alpha.01",
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

			if result.major != tt.expected.major {
				t.Errorf("NewVersion(%q).major = %v, want %v", tt.input, result.major, tt.expected.major)
			}
			if result.minor != tt.expected.minor {
				t.Errorf("NewVersion(%q).minor = %v, want %v", tt.input, result.minor, tt.expected.minor)
			}
			if result.patch != tt.expected.patch {
				t.Errorf("NewVersion(%q).patch = %v, want %v", tt.input, result.patch, tt.expected.patch)
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
		// Basic comparisons
		{"equal versions", "1.2.3", "1.2.3", 0},
		{"major version difference", "2.0.0", "1.9.9", 1},
		{"major version difference (less)", "1.9.9", "2.0.0", -1},
		{"minor version difference", "1.2.0", "1.1.9", 1},
		{"minor version difference (less)", "1.1.9", "1.2.0", -1},
		{"patch version difference", "1.2.3", "1.2.2", 1},
		{"patch version difference (less)", "1.2.2", "1.2.3", -1},

		// Prerelease comparisons
		{"normal version vs prerelease", "1.0.0", "1.0.0-alpha", 1},
		{"prerelease vs normal version", "1.0.0-alpha", "1.0.0", -1},
		{"prerelease identifiers - numeric", "1.0.0-1", "1.0.0-2", -1},
		{"prerelease identifiers - numeric vs alpha", "1.0.0-1", "1.0.0-alpha", -1},
		{"prerelease identifiers - alpha vs numeric", "1.0.0-alpha", "1.0.0-1", 1},
		{"prerelease identifiers - alpha vs beta", "1.0.0-alpha", "1.0.0-beta", -1},
		{"prerelease identifiers - beta vs alpha", "1.0.0-beta", "1.0.0-alpha", 1},
		{"prerelease identifiers - different lengths", "1.0.0-alpha.1", "1.0.0-alpha", 1},
		{"prerelease identifiers - shorter vs longer", "1.0.0-alpha", "1.0.0-alpha.1", -1},

		// SemVer spec examples
		{"semver spec example 1", "1.0.0-alpha", "1.0.0-alpha.1", -1},
		{"semver spec example 2", "1.0.0-alpha.1", "1.0.0-alpha.beta", -1},
		{"semver spec example 3", "1.0.0-alpha.beta", "1.0.0-beta", -1},
		{"semver spec example 4", "1.0.0-beta", "1.0.0-beta.2", -1},
		{"semver spec example 5", "1.0.0-beta.2", "1.0.0-beta.11", -1},
		{"semver spec example 6", "1.0.0-beta.11", "1.0.0-rc.1", -1},
		{"semver spec example 7", "1.0.0-rc.1", "1.0.0", -1},

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
		{"1.2.3-alpha", "1.2.3-alpha"},
		{"1.2.3+build", "1.2.3+build"},
		{"1.2.3-alpha+build", "1.2.3-alpha+build"},
		{"  1.2.3  ", "  1.2.3  "}, // Should preserve original with whitespace
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
