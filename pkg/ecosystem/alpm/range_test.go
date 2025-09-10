package alpm

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
			input: "1.0.0-1",
		},
		{
			name:  "greater than",
			input: ">1.0.0-1",
		},
		{
			name:  "greater than or equal",
			input: ">=1.0.0-1",
		},
		{
			name:  "less than",
			input: "<2.0.0-1",
		},
		{
			name:  "less than or equal",
			input: "<=1.5.0-1",
		},
		{
			name:  "explicit equal",
			input: "=1.0.0-1",
		},
		{
			name:  "version with epoch",
			input: ">=2:1.0.0-1",
		},
		{
			name:  "version without release",
			input: ">=1.2.3",
		},
		{
			name:  "range with multiple constraints",
			input: ">=1.0.0-1 <2.0.0-1",
		},
		{
			name:  "complex range",
			input: ">=1.2.3-1 <1.3.0-1",
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
			input:   ">=1.2@-1",
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
			rangeStr: "1.0.0-1",
			version:  "1.0.0-1",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "1.0.0-1",
			version:  "1.0.0-2",
			want:     false,
		},
		// Greater than
		{
			name:     "greater than - true",
			rangeStr: ">1.0.0-1",
			version:  "1.0.0-2",
			want:     true,
		},
		{
			name:     "greater than - false equal",
			rangeStr: ">1.0.0-1",
			version:  "1.0.0-1",
			want:     false,
		},
		{
			name:     "greater than - false less",
			rangeStr: ">1.0.1-1",
			version:  "1.0.0-1",
			want:     false,
		},
		// Greater than or equal
		{
			name:     "greater than or equal - true greater",
			rangeStr: ">=1.0.0-1",
			version:  "1.0.0-2",
			want:     true,
		},
		{
			name:     "greater than or equal - true equal",
			rangeStr: ">=1.0.0-1",
			version:  "1.0.0-1",
			want:     true,
		},
		{
			name:     "greater than or equal - false less",
			rangeStr: ">=1.0.1-1",
			version:  "1.0.0-1",
			want:     false,
		},
		// Less than
		{
			name:     "less than - true",
			rangeStr: "<2.0.0-1",
			version:  "1.0.0-1",
			want:     true,
		},
		{
			name:     "less than - false equal",
			rangeStr: "<1.0.0-1",
			version:  "1.0.0-1",
			want:     false,
		},
		{
			name:     "less than - false greater",
			rangeStr: "<1.0.0-1",
			version:  "1.0.0-2",
			want:     false,
		},
		// Less than or equal
		{
			name:     "less than or equal - true less",
			rangeStr: "<=1.0.1-1",
			version:  "1.0.0-1",
			want:     true,
		},
		{
			name:     "less than or equal - true equal",
			rangeStr: "<=1.0.0-1",
			version:  "1.0.0-1",
			want:     true,
		},
		{
			name:     "less than or equal - false greater",
			rangeStr: "<=1.0.0-1",
			version:  "1.0.0-2",
			want:     false,
		},
		// Epoch handling
		{
			name:     "epoch comparison - greater",
			rangeStr: ">=1:1.0.0-1",
			version:  "2:0.5.0-1",
			want:     true,
		},
		{
			name:     "epoch comparison - less",
			rangeStr: "<2:1.0.0-1",
			version:  "1:2.0.0-1",
			want:     true,
		},
		{
			name:     "epoch overrides version",
			rangeStr: ">=2:1.0-1",
			version:  "3:0.9-1",
			want:     true,
		},
		// Range constraints
		{
			name:     "range - version in range",
			rangeStr: ">=1.0.0-1 <2.0.0-1",
			version:  "1.5.0-1",
			want:     true,
		},
		{
			name:     "range - version at lower bound",
			rangeStr: ">=1.0.0-1 <2.0.0-1",
			version:  "1.0.0-1",
			want:     true,
		},
		{
			name:     "range - version at upper bound",
			rangeStr: ">=1.0.0-1 <2.0.0-1",
			version:  "2.0.0-1",
			want:     false,
		},
		{
			name:     "range - version below range",
			rangeStr: ">=1.0.0-1 <2.0.0-1",
			version:  "0.9.0-1",
			want:     false,
		},
		{
			name:     "range - version above range",
			rangeStr: ">=1.0.0-1 <2.0.0-1",
			version:  "2.1.0-1",
			want:     false,
		},
		// Real ALMP scenarios
		{
			name:     "Linux kernel range",
			rangeStr: ">=6.1.0-1",
			version:  "6.1.1-1",
			want:     true,
		},
		{
			name:     "Firefox version range",
			rangeStr: ">=108.0.0-1 <109.0.0-1",
			version:  "108.0.2-1",
			want:     true,
		},
		{
			name:     "glibc release range",
			rangeStr: ">=2.36-5 <=2.36-10",
			version:  "2.36-6",
			want:     true,
		},
		{
			name:     "texlive with epoch",
			rangeStr: ">=1:2022.0-1",
			version:  "1:2022.62885-17",
			want:     true,
		},
		// Alphanumeric version comparisons
		{
			name:     "alpha version range",
			rangeStr: ">=1.0a-1",
			version:  "1.0b-1",
			want:     true,
		},
		{
			name:     "beta to rc range",
			rangeStr: ">=1.0beta-1 <1.0-1",
			version:  "1.0rc-1",
			want:     true,
		},
		{
			name:     "pre-release range",
			rangeStr: ">=2.0pre1-1 <2.0-1",
			version:  "2.0rc1-1",
			want:     true,
		},
		// Version without release
		{
			name:     "no release number - exact",
			rangeStr: "1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "no release number - range",
			rangeStr: ">=1.2.0 <1.3.0",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "mixed release formats",
			rangeStr: ">=1.0.0 <1.0.0-10",
			version:  "1.0.0-5",
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
			input: ">=1.0.0-1",
		},
		{
			name:  "complex range",
			input: ">=1.0.0-1 <2.0.0-1",
		},
		{
			name:  "exact version",
			input: "1.0.0-1",
		},
		{
			name:  "version with epoch",
			input: ">=2:1.0.0-1",
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
