package mattermost

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
			input: "v8.1.5",
		},
		{
			name:  "greater than",
			input: ">v8.1.0",
		},
		{
			name:  "greater than or equal",
			input: ">=v8.0.0",
		},
		{
			name:  "less than",
			input: "<v9.0.0",
		},
		{
			name:  "less than or equal",
			input: "<=v8.1.5",
		},
		{
			name:  "explicit equal",
			input: "=v8.1.5",
		},
		{
			name:  "range with multiple constraints",
			input: ">=v8.0.0 <v9.0.0",
		},
		{
			name:  "range without v prefix",
			input: ">=8.0.0 <=8.1.5",
		},
		{
			name:  "ESR range",
			input: ">=v8.1.5-esr",
		},
		{
			name:  "RC range",
			input: ">=v8.1.0-rc1",
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
			input:   ">=v8.1",
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
			rangeStr: "v8.1.5",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "exact match without prefix",
			rangeStr: "8.1.5",
			version:  "8.1.5",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "v8.1.5",
			version:  "v8.1.4",
			want:     false,
		},
		// Greater than
		{
			name:     "greater than - true",
			rangeStr: ">v8.1.0",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "greater than - false equal",
			rangeStr: ">v8.1.5",
			version:  "v8.1.5",
			want:     false,
		},
		{
			name:     "greater than - false less",
			rangeStr: ">v8.1.5",
			version:  "v8.1.0",
			want:     false,
		},
		// Greater than or equal
		{
			name:     "greater than or equal - true greater",
			rangeStr: ">=v8.1.0",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "greater than or equal - true equal",
			rangeStr: ">=v8.1.5",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "greater than or equal - false less",
			rangeStr: ">=v8.1.5",
			version:  "v8.1.0",
			want:     false,
		},
		// Less than
		{
			name:     "less than - true",
			rangeStr: "<v9.0.0",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "less than - false equal",
			rangeStr: "<v8.1.5",
			version:  "v8.1.5",
			want:     false,
		},
		{
			name:     "less than - false greater",
			rangeStr: "<v8.1.0",
			version:  "v8.1.5",
			want:     false,
		},
		// Less than or equal
		{
			name:     "less than or equal - true less",
			rangeStr: "<=v8.1.5",
			version:  "v8.1.0",
			want:     true,
		},
		{
			name:     "less than or equal - true equal",
			rangeStr: "<=v8.1.5",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "less than or equal - false greater",
			rangeStr: "<=v8.1.0",
			version:  "v8.1.5",
			want:     false,
		},
		// Range constraints
		{
			name:     "range - version in range",
			rangeStr: ">=v8.0.0 <v9.0.0",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "range - version at lower bound",
			rangeStr: ">=v8.0.0 <v9.0.0",
			version:  "v8.0.0",
			want:     true,
		},
		{
			name:     "range - version at upper bound",
			rangeStr: ">=v8.0.0 <v9.0.0",
			version:  "v9.0.0",
			want:     false,
		},
		{
			name:     "range - version below range",
			rangeStr: ">=v8.0.0 <v9.0.0",
			version:  "v7.8.11",
			want:     false,
		},
		{
			name:     "range - version above range",
			rangeStr: ">=v8.0.0 <v9.0.0",
			version:  "v10.0.0",
			want:     false,
		},
		// Prefix handling
		{
			name:     "mixed prefix in range and version",
			rangeStr: ">=v8.1.0",
			version:  "8.1.5",
			want:     true,
		},
		// Qualifiers
		{
			name:     "rc in range",
			rangeStr: ">=v8.1.0-rc1",
			version:  "v8.1.0",
			want:     true,
		},
		{
			name:     "release vs rc constraint",
			rangeStr: ">=v8.1.0",
			version:  "v8.1.0-rc1",
			want:     false,
		},
		{
			name:     "ESR in range",
			rangeStr: ">=v8.1.5-esr",
			version:  "v8.2.0",
			want:     true,
		},
		{
			name:     "rc vs ESR constraint",
			rangeStr: ">=v8.1.0-rc1",
			version:  "v8.1.5-esr",
			want:     true,
		},
		// Real Mattermost scenarios
		{
			name:     "Mattermost 8.x series range",
			rangeStr: ">=v8.0.0 <v9.0.0",
			version:  "v8.1.5",
			want:     true,
		},
		{
			name:     "Mattermost ESR constraint",
			rangeStr: ">=v8.1.5-esr",
			version:  "v8.1.5-esr",
			want:     true,
		},
		{
			name:     "Mattermost RC progression",
			rangeStr: ">=v8.1.0-rc1",
			version:  "v8.1.0-rc2",
			want:     true,
		},
		{
			name:     "Cross-major version check",
			rangeStr: ">=v10.0.0",
			version:  "v10.11.2",
			want:     true,
		},
		{
			name:     "ESR vs newer major",
			rangeStr: ">=v8.1.5-esr",
			version:  "v10.0.0",
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
			input: ">=v8.0.0",
		},
		{
			name:  "complex range",
			input: ">=v8.0.0 <v9.0.0",
		},
		{
			name:  "exact version",
			input: "v8.1.5",
		},
		{
			name:  "ESR range",
			input: ">=v8.1.5-esr",
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
