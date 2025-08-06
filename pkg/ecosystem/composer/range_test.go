package composer

import (
	"testing"
)

func TestNewVersionRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid ranges
		{"exact version", "1.2.3", false},
		{"comparison operators", ">=1.2.3", false},
		{"caret constraint", "^1.2.3", false},
		{"tilde constraint", "~1.2.3", false},
		{"wildcard constraint", "1.2.*", false},
		{"hyphen range", "1.2.3 - 2.3.4", false},
		{"space separated", ">=1.0.0 <2.0.0", false},
		{"comma separated", ">=1.0.0, <2.0.0", false},
		{"OR logic", "1.x || 2.x", false},
		{"stability constraint", "1.2.3@dev", false},
		{"stability only", "@stable", false},
		{"wildcard", "*", false},

		// Invalid ranges
		{"empty string", "", true},
		{"invalid version in caret", "^invalid", true},
		{"invalid hyphen range", "1.2.3 - ", true},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVersionRangeContains(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		version  string
		want     bool
	}{
		// Exact version matches
		{"exact match", "1.2.3", "1.2.3", true},
		{"exact no match", "1.2.3", "1.2.4", false},

		// Comparison operators
		{"greater than equal - match", ">=1.2.3", "1.2.3", true},
		{"greater than equal - higher", ">=1.2.3", "1.2.4", true},
		{"greater than equal - lower", ">=1.2.3", "1.2.2", false},
		{"less than - match", "<2.0.0", "1.9.9", true},
		{"less than - no match", "<2.0.0", "2.0.0", false},
		{"not equal - different", "!=1.2.3", "1.2.4", true},
		{"not equal - same", "!=1.2.3", "1.2.3", false},

		// Caret constraints
		{"caret major match", "^1.2.3", "1.2.3", true},
		{"caret minor bump", "^1.2.3", "1.3.0", true},
		{"caret patch bump", "^1.2.3", "1.2.4", true},
		{"caret major bump no match", "^1.2.3", "2.0.0", false},
		{"caret zero major", "^0.2.3", "0.2.4", true},
		{"caret zero major minor bump", "^0.2.3", "0.3.0", false},
		{"caret zero minor", "^0.0.3", "0.0.3", true},
		{"caret zero minor patch bump", "^0.0.3", "0.0.4", false},

		// Tilde constraints
		{"tilde patch match", "~1.2.3", "1.2.3", true},
		{"tilde patch bump", "~1.2.3", "1.2.4", true},
		{"tilde minor bump no match", "~1.2.3", "1.3.0", false},
		{"tilde two component", "~1.2", "1.2.0", true},
		{"tilde two component patch", "~1.2", "1.2.5", true},
		{"tilde two component major", "~1.2", "2.0.0", false},

		// Wildcard constraints
		{"wildcard minor match", "1.2.*", "1.2.0", true},
		{"wildcard minor patch", "1.2.*", "1.2.5", true},
		{"wildcard minor no match", "1.2.*", "1.3.0", false},
		{"wildcard major match", "1.*", "1.0.0", true},
		{"wildcard major minor", "1.*", "1.5.0", true},
		{"wildcard major no match", "1.*", "2.0.0", false},

		// Hyphen ranges
		{"hyphen range start", "1.2.3 - 2.3.4", "1.2.3", true},
		{"hyphen range end", "1.2.3 - 2.3.4", "2.3.4", true},
		{"hyphen range middle", "1.2.3 - 2.3.4", "1.5.0", true},
		{"hyphen range before", "1.2.3 - 2.3.4", "1.2.2", false},
		{"hyphen range after", "1.2.3 - 2.3.4", "2.3.5", false},

		// Combined constraints (AND logic)
		{"combined constraints match", ">=1.0.0 <2.0.0", "1.5.0", true},
		{"combined constraints low", ">=1.0.0 <2.0.0", "0.9.0", false},
		{"combined constraints high", ">=1.0.0 <2.0.0", "2.0.0", false},
		{"comma separated match", ">=1.0.0, <2.0.0", "1.5.0", true},
		{"comma separated no match", ">=1.0.0, <2.0.0", "2.0.0", false},

		// OR logic
		{"OR match first", "1.x || 2.x", "1.2.3", true},
		{"OR match second", "1.x || 2.x", "2.3.4", true},
		{"OR no match", "1.x || 2.x", "3.0.0", false},

		// Stability constraints
		{"stability exact match", "1.2.3-alpha", "1.2.3-alpha", true},
		{"stability no match", "1.2.3-alpha", "1.2.3-beta", false},
		{"stability stable vs pre", "1.2.3", "1.2.3-alpha", false},

		// Dev versions
		{"dev version match", "dev-main", "dev-main", true},
		{"dev version no match", "dev-main", "dev-feature", false},
		{"dev vs stable", ">=1.0.0", "dev-main", false},

		// Wildcard all
		{"wildcard all", "*", "1.2.3", true},
		{"wildcard all dev", "*", "dev-main", true},

		// Edge cases with stability
		{"prerelease in caret", "^1.2.3", "1.2.3-alpha", false}, // Pre-releases don't satisfy caret
		{"prerelease exact", "1.2.3-alpha", "1.2.3-alpha", true},

		// Advanced constraint combinations from real-world data
		{"complex OR constraint", "1.0.* || >=2.0.0,<3.0.0", "1.0.5", true},
		{"complex OR constraint no match", "1.0.* || >=2.0.0,<3.0.0", "1.1.0", false},
		{"multiple exclusions", ">=1.0.0,<2.0.0,!=1.2.0,!=1.3.0", "1.1.0", true},
		{"multiple exclusions match excluded", ">=1.0.0,<2.0.0,!=1.2.0,!=1.3.0", "1.2.0", false},

		// pl version constraints
		{"pl version in range", ">=1.0.0", "1.0pl1", true},
		{"pl version caret", "^1.0.0", "1.0pl1", true},

		// Alternative stability suffix formats
		{"alpha without hyphen in range", ">=1.0.0-alpha", "1.0a1", true},
		{"beta without hyphen in caret", "^1.0.0", "1.0b1", true},

		// Complex prerelease ranges
		{"prerelease range", ">=1.0.0-alpha,<1.0.0", "1.0.0-beta", true},
		{"prerelease range boundary", ">=1.0.0-alpha,<1.0.0", "1.0.0", false},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange() error = %v", err)
			}

			v, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion() error = %v", err)
			}

			got := r.Contains(v)
			if got != tt.want {
				t.Errorf("Contains(%s, %s) = %v, want %v", tt.rangeStr, tt.version, got, tt.want)
			}
		})
	}
}

func TestVersionRangeString(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
	}{
		{"exact version", "1.2.3"},
		{"caret constraint", "^1.2.3"},
		{"tilde constraint", "~1.2.3"},
		{"comparison", ">=1.2.3"},
		{"combined", ">=1.0.0 <2.0.0"},
		{"OR logic", "1.x || 2.x"},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange() error = %v", err)
			}
			if got := r.String(); got != tt.rangeStr {
				t.Errorf("String() = %v, want %v", got, tt.rangeStr)
			}
		})
	}
}

func TestCaretConstraintEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		version  string
		want     bool
	}{
		// Major version 0 special cases
		{"caret 0.x.y allows patch", "^0.2.3", "0.2.4", true},
		{"caret 0.x.y blocks minor", "^0.2.3", "0.3.0", false},
		{"caret 0.0.x exact only", "^0.0.3", "0.0.3", true},
		{"caret 0.0.x blocks patch", "^0.0.3", "0.0.4", false},

		// Regular caret behavior
		{"caret 1.x.y allows minor", "^1.2.3", "1.3.0", true},
		{"caret 1.x.y allows patch", "^1.2.3", "1.2.4", true},
		{"caret 1.x.y blocks major", "^1.2.3", "2.0.0", false},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange() error = %v", err)
			}

			v, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion() error = %v", err)
			}

			got := r.Contains(v)
			if got != tt.want {
				t.Errorf("Contains(%s, %s) = %v, want %v", tt.rangeStr, tt.version, got, tt.want)
			}
		})
	}
}

func TestTildeConstraintVariations(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		version  string
		want     bool
	}{
		// Different tilde formats
		{"tilde one component", "~1", "1.5.0", true},
		{"tilde one component major bump", "~1", "2.0.0", false},
		{"tilde two component", "~1.2", "1.2.5", true},
		{"tilde two component major bump", "~1.2", "2.0.0", false},
		{"tilde three component", "~1.2.3", "1.2.4", true},
		{"tilde three component minor bump", "~1.2.3", "1.3.0", false},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange() error = %v", err)
			}

			v, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion() error = %v", err)
			}

			got := r.Contains(v)
			if got != tt.want {
				t.Errorf("Contains(%s, %s) = %v, want %v", tt.rangeStr, tt.version, got, tt.want)
			}
		})
	}
}
