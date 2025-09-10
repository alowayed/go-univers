package github

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:  "exact version",
			input: "v1.0.0",
		},
		{
			name:  "greater than",
			input: ">v1.0.0",
		},
		{
			name:  "greater than or equal",
			input: ">=v1.0.0",
		},
		{
			name:  "less than",
			input: "<v2.0.0",
		},
		{
			name:  "less than or equal",
			input: "<=v1.5.0",
		},
		{
			name:  "explicit equal",
			input: "=v1.0.0",
		},
		{
			name:  "range with multiple constraints",
			input: ">=v1.0.0 <v2.0.0",
		},
		{
			name:  "range without v prefix",
			input: ">=1.0.0 <=1.5.0",
		},
		{
			name:  "date-based range",
			input: ">=2024.01.01",
		},
		{
			name:  "release prefix range",
			input: ">=release-1.0.0",
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
			name:    "invalid version in range",
			input:   ">=v1.2",
			wantErr: true,
		},
	}

	ecosystem := &Ecosystem{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.String() != tt.input {
				t.Errorf("VersionRange.String() = %v, want %v", got.String(), tt.input)
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
			name:     "exact match with v prefix",
			rangeStr: "v1.0.0",
			version:  "v1.0.0",
			want:     true,
		},
		{
			name:     "exact match without prefix",
			rangeStr: "1.0.0",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "v1.0.0",
			version:  "v1.0.1",
			want:     false,
		},
		// Greater than
		{
			name:     "greater than - true",
			rangeStr: ">v1.0.0",
			version:  "v1.0.1",
			want:     true,
		},
		{
			name:     "greater than - false equal",
			rangeStr: ">v1.0.0",
			version:  "v1.0.0",
			want:     false,
		},
		{
			name:     "greater than - false less",
			rangeStr: ">v1.0.1",
			version:  "v1.0.0",
			want:     false,
		},
		// Greater than or equal
		{
			name:     "greater than or equal - true greater",
			rangeStr: ">=v1.0.0",
			version:  "v1.0.1",
			want:     true,
		},
		{
			name:     "greater than or equal - true equal",
			rangeStr: ">=v1.0.0",
			version:  "v1.0.0",
			want:     true,
		},
		{
			name:     "greater than or equal - false less",
			rangeStr: ">=v1.0.1",
			version:  "v1.0.0",
			want:     false,
		},
		// Less than
		{
			name:     "less than - true",
			rangeStr: "<v2.0.0",
			version:  "v1.0.0",
			want:     true,
		},
		{
			name:     "less than - false equal",
			rangeStr: "<v1.0.0",
			version:  "v1.0.0",
			want:     false,
		},
		{
			name:     "less than - false greater",
			rangeStr: "<v1.0.0",
			version:  "v1.0.1",
			want:     false,
		},
		// Less than or equal
		{
			name:     "less than or equal - true less",
			rangeStr: "<=v1.0.1",
			version:  "v1.0.0",
			want:     true,
		},
		{
			name:     "less than or equal - true equal",
			rangeStr: "<=v1.0.0",
			version:  "v1.0.0",
			want:     true,
		},
		{
			name:     "less than or equal - false greater",
			rangeStr: "<=v1.0.0",
			version:  "v1.0.1",
			want:     false,
		},
		// Range constraints
		{
			name:     "range - version in range",
			rangeStr: ">=v1.0.0 <v2.0.0",
			version:  "v1.5.0",
			want:     true,
		},
		{
			name:     "range - version at lower bound",
			rangeStr: ">=v1.0.0 <v2.0.0",
			version:  "v1.0.0",
			want:     true,
		},
		{
			name:     "range - version at upper bound",
			rangeStr: ">=v1.0.0 <v2.0.0",
			version:  "v2.0.0",
			want:     false,
		},
		{
			name:     "range - version below range",
			rangeStr: ">=v1.0.0 <v2.0.0",
			version:  "v0.9.0",
			want:     false,
		},
		{
			name:     "range - version above range",
			rangeStr: ">=v1.0.0 <v2.0.0",
			version:  "v2.1.0",
			want:     false,
		},
		// Prefix handling
		{
			name:     "mixed prefix in range and version",
			rangeStr: ">=v1.0.0",
			version:  "1.0.1",
			want:     true,
		},
		{
			name:     "release prefix range",
			rangeStr: ">=release-1.0.0",
			version:  "release-1.5.0",
			want:     true,
		},
		// Qualifiers
		{
			name:     "pre-release in range",
			rangeStr: ">=v1.0.0-beta",
			version:  "v1.0.0",
			want:     true,
		},
		{
			name:     "release vs pre-release constraint",
			rangeStr: ">=v1.0.0",
			version:  "v1.0.0-beta",
			want:     false,
		},
		// Date-based versions
		{
			name:     "date-based range",
			rangeStr: ">=2024.01.01",
			version:  "2024.06.15",
			want:     true,
		},
		{
			name:     "date-based exact",
			rangeStr: "2024.01.15",
			version:  "2024.01.15",
			want:     true,
		},
		{
			name:     "date-based vs semantic (date < semantic)",
			rangeStr: ">=2024.01.01",
			version:  "v1.0.0",
			want:     true,
		},
		// Real GitHub scenarios
		{
			name:     "Go version range",
			rangeStr: ">=v1.20.0 <v1.22.0",
			version:  "v1.21.0",
			want:     true,
		},
		{
			name:     "Kubernetes version",
			rangeStr: ">=v1.28.0",
			version:  "v1.29.0",
			want:     true,
		},
		{
			name:     "Pre-release constraint",
			rangeStr: ">=v2.0.0-beta",
			version:  "v2.0.0-rc.1",
			want:     true,
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
				t.Errorf("VersionRange.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionRange_String(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "simple range",
			input: ">=v1.0.0",
		},
		{
			name:  "complex range",
			input: ">=v1.0.0 <v2.0.0",
		},
		{
			name:  "exact version",
			input: "v1.0.0",
		},
		{
			name:  "date-based range",
			input: ">=2024.01.01",
		},
	}

	ecosystem := &Ecosystem{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := ecosystem.NewVersionRange(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse range %s: %v", tt.input, err)
			}

			if got := vr.String(); got != tt.input {
				t.Errorf("VersionRange.String() = %v, want %v", got, tt.input)
			}
		})
	}
}
