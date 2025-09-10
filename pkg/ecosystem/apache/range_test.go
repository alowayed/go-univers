package apache

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
			input: "2.4.41",
		},
		{
			name:  "greater than",
			input: ">2.4.0",
		},
		{
			name:  "greater than or equal",
			input: ">=2.4.0",
		},
		{
			name:  "less than",
			input: "<3.0.0",
		},
		{
			name:  "less than or equal",
			input: "<=2.4.41",
		},
		{
			name:  "explicit equal",
			input: "=2.4.41",
		},
		{
			name:  "range with multiple constraints",
			input: ">=2.4.0 <3.0.0",
		},
		{
			name:  "range with inclusive bounds",
			input: ">=2.4.0 <=2.4.41",
		},
		{
			name:  "range with RC version",
			input: ">=2.4.41-RC1",
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
			input:   ">=2.4",
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
			rangeStr: "2.4.41",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "exact match explicit",
			rangeStr: "=2.4.41",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "2.4.41",
			version:  "2.4.40",
			want:     false,
		},
		// Greater than
		{
			name:     "greater than - true",
			rangeStr: ">2.4.0",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "greater than - false equal",
			rangeStr: ">2.4.41",
			version:  "2.4.41",
			want:     false,
		},
		{
			name:     "greater than - false less",
			rangeStr: ">2.4.41",
			version:  "2.4.40",
			want:     false,
		},
		// Greater than or equal
		{
			name:     "greater than or equal - true greater",
			rangeStr: ">=2.4.0",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "greater than or equal - true equal",
			rangeStr: ">=2.4.41",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "greater than or equal - false less",
			rangeStr: ">=2.4.41",
			version:  "2.4.40",
			want:     false,
		},
		// Less than
		{
			name:     "less than - true",
			rangeStr: "<3.0.0",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "less than - false equal",
			rangeStr: "<2.4.41",
			version:  "2.4.41",
			want:     false,
		},
		{
			name:     "less than - false greater",
			rangeStr: "<2.4.41",
			version:  "2.4.42",
			want:     false,
		},
		// Less than or equal
		{
			name:     "less than or equal - true less",
			rangeStr: "<=2.4.41",
			version:  "2.4.40",
			want:     true,
		},
		{
			name:     "less than or equal - true equal",
			rangeStr: "<=2.4.41",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "less than or equal - false greater",
			rangeStr: "<=2.4.41",
			version:  "2.4.42",
			want:     false,
		},
		// Range constraints
		{
			name:     "range - version in range",
			rangeStr: ">=2.4.0 <3.0.0",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "range - version at lower bound",
			rangeStr: ">=2.4.0 <3.0.0",
			version:  "2.4.0",
			want:     true,
		},
		{
			name:     "range - version at upper bound",
			rangeStr: ">=2.4.0 <3.0.0",
			version:  "3.0.0",
			want:     false,
		},
		{
			name:     "range - version below range",
			rangeStr: ">=2.4.0 <3.0.0",
			version:  "2.3.0",
			want:     false,
		},
		{
			name:     "range - version above range",
			rangeStr: ">=2.4.0 <3.0.0",
			version:  "3.1.0",
			want:     false,
		},
		{
			name:     "inclusive range",
			rangeStr: ">=2.4.0 <=2.4.41",
			version:  "2.4.41",
			want:     true,
		},
		// Qualifiers
		{
			name:     "RC version in range",
			rangeStr: ">=2.4.41-RC1",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "RC version exact match",
			rangeStr: "2.4.41-RC1",
			version:  "2.4.41-RC1",
			want:     true,
		},
		{
			name:     "release vs RC constraint",
			rangeStr: ">=2.4.41",
			version:  "2.4.41-RC1",
			want:     false,
		},
		// Real Apache scenarios
		{
			name:     "Apache HTTP Server 2.4.x",
			rangeStr: ">=2.4.0 <2.5.0",
			version:  "2.4.41",
			want:     true,
		},
		{
			name:     "Apache Tomcat 9.0.x",
			rangeStr: ">=9.0.0 <10.0.0",
			version:  "9.0.45",
			want:     true,
		},
		{
			name:     "Minimum Apache version",
			rangeStr: ">=2.4.35",
			version:  "2.4.41",
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
			input: ">=2.4.0",
		},
		{
			name:  "complex range",
			input: ">=2.4.0 <3.0.0",
		},
		{
			name:  "exact version",
			input: "2.4.41",
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
