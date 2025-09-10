package hex

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
			input: "1.0.0",
		},
		{
			name:  "greater than",
			input: ">1.0.0",
		},
		{
			name:  "greater than or equal",
			input: ">=1.0.0",
		},
		{
			name:  "less than",
			input: "<2.0.0",
		},
		{
			name:  "less than or equal",
			input: "<=1.5.0",
		},
		{
			name:  "explicit equal",
			input: "=1.0.0",
		},
		{
			name:  "pessimistic operator patch",
			input: "~>1.2.3",
		},
		{
			name:  "pessimistic operator minor",
			input: "~>1.2.0",
		},
		{
			name:  "pessimistic operator major",
			input: "~>1.0",
		},
		{
			name:  "range with multiple constraints",
			input: ">=1.0.0 <2.0.0",
		},
		{
			name:  "pre-release range",
			input: ">=1.0.0-alpha",
		},
		{
			name:  "complex range",
			input: ">=1.2.3 <1.3.0",
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
			input:   ">=1.2.a",
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
			name:     "exact match",
			rangeStr: "1.0.0",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "1.0.0",
			version:  "1.0.1",
			want:     false,
		},
		// Greater than
		{
			name:     "greater than - true",
			rangeStr: ">1.0.0",
			version:  "1.0.1",
			want:     true,
		},
		{
			name:     "greater than - false equal",
			rangeStr: ">1.0.0",
			version:  "1.0.0",
			want:     false,
		},
		{
			name:     "greater than - false less",
			rangeStr: ">1.0.1",
			version:  "1.0.0",
			want:     false,
		},
		// Greater than or equal
		{
			name:     "greater than or equal - true greater",
			rangeStr: ">=1.0.0",
			version:  "1.0.1",
			want:     true,
		},
		{
			name:     "greater than or equal - true equal",
			rangeStr: ">=1.0.0",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "greater than or equal - false less",
			rangeStr: ">=1.0.1",
			version:  "1.0.0",
			want:     false,
		},
		// Less than
		{
			name:     "less than - true",
			rangeStr: "<2.0.0",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "less than - false equal",
			rangeStr: "<1.0.0",
			version:  "1.0.0",
			want:     false,
		},
		{
			name:     "less than - false greater",
			rangeStr: "<1.0.0",
			version:  "1.0.1",
			want:     false,
		},
		// Less than or equal
		{
			name:     "less than or equal - true less",
			rangeStr: "<=1.0.1",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "less than or equal - true equal",
			rangeStr: "<=1.0.0",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "less than or equal - false greater",
			rangeStr: "<=1.0.0",
			version:  "1.0.1",
			want:     false,
		},
		// Pessimistic operator
		{
			name:     "pessimistic patch level - in range",
			rangeStr: "~>1.2.3",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "pessimistic patch level - at boundary",
			rangeStr: "~>1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "pessimistic patch level - out of range minor",
			rangeStr: "~>1.2.3",
			version:  "1.3.0",
			want:     false,
		},
		{
			name:     "pessimistic patch level - out of range major",
			rangeStr: "~>1.2.3",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "pessimistic minor level - in range",
			rangeStr: "~>1.2.0",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "pessimistic minor level - out of range",
			rangeStr: "~>1.2.0",
			version:  "1.3.0",
			want:     false,
		},
		{
			name:     "pessimistic major level - in range",
			rangeStr: "~>1.0.0",
			version:  "1.0.5",
			want:     true,
		},
		{
			name:     "pessimistic major level - out of range",
			rangeStr: "~>1.0.0",
			version:  "2.0.0",
			want:     false,
		},
		// Range constraints
		{
			name:     "range - version in range",
			rangeStr: ">=1.0.0 <2.0.0",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "range - version at lower bound",
			rangeStr: ">=1.0.0 <2.0.0",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "range - version at upper bound",
			rangeStr: ">=1.0.0 <2.0.0",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "range - version below range",
			rangeStr: ">=1.0.0 <2.0.0",
			version:  "0.9.0",
			want:     false,
		},
		{
			name:     "range - version above range",
			rangeStr: ">=1.0.0 <2.0.0",
			version:  "2.1.0",
			want:     false,
		},
		// Pre-release handling
		{
			name:     "pre-release in range",
			rangeStr: ">=1.0.0-alpha",
			version:  "1.0.0",
			want:     true,
		},
		{
			name:     "release vs pre-release constraint",
			rangeStr: ">=1.0.0",
			version:  "1.0.0-alpha",
			want:     false,
		},
		{
			name:     "pre-release to pre-release",
			rangeStr: ">=1.0.0-alpha",
			version:  "1.0.0-beta",
			want:     true,
		},
		// Build metadata (ignored in comparisons)
		{
			name:     "build metadata ignored",
			rangeStr: "1.0.0",
			version:  "1.0.0+build.123",
			want:     true,
		},
		{
			name:     "both have build metadata",
			rangeStr: "1.0.0+build.1",
			version:  "1.0.0+build.2",
			want:     true,
		},
		// Real Hex scenarios
		{
			name:     "Phoenix version range",
			rangeStr: "~>1.7.0",
			version:  "1.7.10",
			want:     true,
		},
		{
			name:     "Ecto version range",
			rangeStr: ">=3.0.0 <4.0.0",
			version:  "3.10.3",
			want:     true,
		},
		{
			name:     "Pre-release Phoenix",
			rangeStr: "~>1.8.0",
			version:  "1.8.1",
			want:     true,
		},
		{
			name:     "Elixir compatibility",
			rangeStr: "~>1.14",
			version:  "1.15.7",
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
			input: ">=1.0.0",
		},
		{
			name:  "complex range",
			input: ">=1.0.0 <2.0.0",
		},
		{
			name:  "exact version",
			input: "1.0.0",
		},
		{
			name:  "pessimistic range",
			input: "~>1.2.3",
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
