package cran

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid ranges
		{
			name:  "exact version",
			input: "1.2.3",
		},
		{
			name:  "greater than",
			input: ">1.2.3",
		},
		{
			name:  "greater than or equal",
			input: ">=1.2.3",
		},
		{
			name:  "less than",
			input: "<1.2.3",
		},
		{
			name:  "less than or equal",
			input: "<=1.2.3",
		},
		{
			name:  "not equal",
			input: "!=1.2.3",
		},
		{
			name:  "multiple constraints",
			input: ">=1.2.0, <2.0.0",
		},
		{
			name:  "multiple constraints with spaces",
			input: ">= 1.2.0 , < 2.0.0",
		},
		{
			name:  "complex range",
			input: ">=1.2.0, <2.0.0, !=1.5.0",
		},
		{
			name:  "range with dashes",
			input: ">=1-2-0",
		},
		{
			name:  "range with whitespace",
			input: "  >= 1.2.3  ",
		},
		// Error cases
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "invalid operator",
			input:   "~>1.2.3",
			wantErr: true,
		},
		{
			name:    "operator without version",
			input:   ">=",
			wantErr: true,
		},
		{
			name:    "invalid version in range",
			input:   ">=1.abc",
			wantErr: true,
		},
		{
			name:    "empty constraint in list",
			input:   ">=1.2.3, , <2.0.0",
			wantErr: false, // Empty constraints are skipped
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ecosystem.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		version  string
		want     bool
	}{
		// Exact matches
		{
			name:     "exact match",
			rangeStr: "1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "exact match with dashes",
			rangeStr: "1-2-3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "1.2.3",
			version:  "1.2.4",
			want:     false,
		},
		// Greater than
		{
			name:     "greater than satisfied",
			rangeStr: ">1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than not satisfied",
			rangeStr: ">1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "greater than not satisfied - less",
			rangeStr: ">1.2.3",
			version:  "1.2.2",
			want:     false,
		},
		// Greater than or equal
		{
			name:     "greater than or equal satisfied - greater",
			rangeStr: ">=1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than or equal satisfied - equal",
			rangeStr: ">=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "greater than or equal not satisfied",
			rangeStr: ">=1.2.3",
			version:  "1.2.2",
			want:     false,
		},
		// Less than
		{
			name:     "less than satisfied",
			rangeStr: "<1.2.3",
			version:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than not satisfied",
			rangeStr: "<1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "less than not satisfied - greater",
			rangeStr: "<1.2.3",
			version:  "1.2.4",
			want:     false,
		},
		// Less than or equal
		{
			name:     "less than or equal satisfied - less",
			rangeStr: "<=1.2.3",
			version:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than or equal satisfied - equal",
			rangeStr: "<=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "less than or equal not satisfied",
			rangeStr: "<=1.2.3",
			version:  "1.2.4",
			want:     false,
		},
		// Not equal
		{
			name:     "not equal satisfied",
			rangeStr: "!=1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "not equal not satisfied",
			rangeStr: "!=1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		// Multiple constraints (AND logic)
		{
			name:     "multiple constraints satisfied",
			rangeStr: ">=1.2.0, <2.0.0",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "multiple constraints not satisfied - too low",
			rangeStr: ">=1.2.0, <2.0.0",
			version:  "1.1.0",
			want:     false,
		},
		{
			name:     "multiple constraints not satisfied - too high",
			rangeStr: ">=1.2.0, <2.0.0",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "multiple constraints satisfied at boundary",
			rangeStr: ">=1.2.0, <2.0.0",
			version:  "1.2.0",
			want:     true,
		},
		{
			name:     "complex constraints satisfied",
			rangeStr: ">=1.2.0, <2.0.0, !=1.5.0",
			version:  "1.4.0",
			want:     true,
		},
		{
			name:     "complex constraints not satisfied - excluded version",
			rangeStr: ">=1.2.0, <2.0.0, !=1.5.0",
			version:  "1.5.0",
			want:     false,
		},
		// Version component length differences
		{
			name:     "different component lengths - version longer",
			rangeStr: ">=1.2",
			version:  "1.2.1",
			want:     true,
		},
		{
			name:     "different component lengths - constraint longer",
			rangeStr: ">=1.2.1",
			version:  "1.2",
			want:     false,
		},
		{
			name:     "exact match different lengths",
			rangeStr: "1.2.0",
			version:  "1.2",
			want:     false,
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := ecosystem.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("Failed to parse range %s: %v", tt.rangeStr, err)
			}

			v, err := ecosystem.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.version, err)
			}

			got := vr.Contains(v)
			if got != tt.want {
				t.Errorf("VersionRange.Contains() = %v, want %v (range=%s, version=%s)", got, tt.want, tt.rangeStr, tt.version)
			}
		})
	}
}
