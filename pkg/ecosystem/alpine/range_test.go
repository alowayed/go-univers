package alpine

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	ecosystem := &Ecosystem{}

	tests := []struct {
		name    string
		input   string
		want    *VersionRange
		wantErr bool
	}{
		{
			name:  "exact version",
			input: "1.2.3",
			want: &VersionRange{
				constraints: []*constraint{{operator: "=", version: "1.2.3"}},
				original:    "1.2.3",
			},
		},
		{
			name:  "greater than",
			input: ">1.2.3",
			want: &VersionRange{
				constraints: []*constraint{{operator: ">", version: "1.2.3"}},
				original:    ">1.2.3",
			},
		},
		{
			name:  "greater than or equal",
			input: ">=1.2.3",
			want: &VersionRange{
				constraints: []*constraint{{operator: ">=", version: "1.2.3"}},
				original:    ">=1.2.3",
			},
		},
		{
			name:  "less than",
			input: "<2.0.0",
			want: &VersionRange{
				constraints: []*constraint{{operator: "<", version: "2.0.0"}},
				original:    "<2.0.0",
			},
		},
		{
			name:  "less than or equal",
			input: "<=2.0.0",
			want: &VersionRange{
				constraints: []*constraint{{operator: "<=", version: "2.0.0"}},
				original:    "<=2.0.0",
			},
		},
		{
			name:  "not equal",
			input: "!=1.5.0",
			want: &VersionRange{
				constraints: []*constraint{{operator: "!=", version: "1.5.0"}},
				original:    "!=1.5.0",
			},
		},
		{
			name:  "multiple constraints",
			input: ">=1.0.0 <2.0.0",
			want: &VersionRange{
				constraints: []*constraint{
					{operator: ">=", version: "1.0.0"},
					{operator: "<", version: "2.0.0"},
				},
				original: ">=1.0.0 <2.0.0",
			},
		},
		{
			name:  "complex range",
			input: ">=1.2.0 <2.0.0 !=1.5.0",
			want: &VersionRange{
				constraints: []*constraint{
					{operator: ">=", version: "1.2.0"},
					{operator: "<", version: "2.0.0"},
					{operator: "!=", version: "1.5.0"},
				},
				original: ">=1.2.0 <2.0.0 !=1.5.0",
			},
		},
		{
			name:  "version with suffix",
			input: ">=1.2.3_alpha",
			want: &VersionRange{
				constraints: []*constraint{{operator: ">=", version: "1.2.3_alpha"}},
				original:    ">=1.2.3_alpha",
			},
		},
		{
			name:  "version with build",
			input: ">=1.2.3-r1",
			want: &VersionRange{
				constraints: []*constraint{{operator: ">=", version: "1.2.3-r1"}},
				original:    ">=1.2.3-r1",
			},
		},

		// Error cases
		{name: "empty string", input: "", wantErr: true},
		{name: "whitespace only", input: "   ", wantErr: true},
		{name: "incomplete operator", input: ">", wantErr: true},
		{name: "incomplete operator >=", input: ">=", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersionRange(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewVersionRange() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("NewVersionRange() unexpected error: %v", err)
				return
			}

			if got.String() != tt.want.original {
				t.Errorf("NewVersionRange().String() = %q, want %q", got.String(), tt.want.original)
			}

			if len(got.constraints) != len(tt.want.constraints) {
				t.Errorf("NewVersionRange() constraints length = %d, want %d", len(got.constraints), len(tt.want.constraints))
				return
			}

			for i, constraint := range got.constraints {
				if constraint.operator != tt.want.constraints[i].operator {
					t.Errorf("NewVersionRange() constraint[%d].operator = %q, want %q", i, constraint.operator, tt.want.constraints[i].operator)
				}
				if constraint.version != tt.want.constraints[i].version {
					t.Errorf("NewVersionRange() constraint[%d].version = %q, want %q", i, constraint.version, tt.want.constraints[i].version)
				}
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	ecosystem := &Ecosystem{}

	tests := []struct {
		name     string
		rangeStr string
		version  string
		want     bool
	}{
		// Exact matches
		{name: "exact match", rangeStr: "1.2.3", version: "1.2.3", want: true},
		{name: "exact match failure", rangeStr: "1.2.3", version: "1.2.4", want: false},

		// Greater than
		{name: "greater than satisfied", rangeStr: ">1.2.3", version: "1.2.4", want: true},
		{name: "greater than not satisfied", rangeStr: ">1.2.3", version: "1.2.3", want: false},
		{name: "greater than not satisfied smaller", rangeStr: ">1.2.3", version: "1.2.2", want: false},

		// Greater than or equal
		{name: "gte satisfied equal", rangeStr: ">=1.2.3", version: "1.2.3", want: true},
		{name: "gte satisfied greater", rangeStr: ">=1.2.3", version: "1.2.4", want: true},
		{name: "gte not satisfied", rangeStr: ">=1.2.3", version: "1.2.2", want: false},

		// Less than
		{name: "less than satisfied", rangeStr: "<2.0.0", version: "1.9.9", want: true},
		{name: "less than not satisfied", rangeStr: "<2.0.0", version: "2.0.0", want: false},
		{name: "less than not satisfied greater", rangeStr: "<2.0.0", version: "2.0.1", want: false},

		// Less than or equal
		{name: "lte satisfied equal", rangeStr: "<=2.0.0", version: "2.0.0", want: true},
		{name: "lte satisfied smaller", rangeStr: "<=2.0.0", version: "1.9.9", want: true},
		{name: "lte not satisfied", rangeStr: "<=2.0.0", version: "2.0.1", want: false},

		// Not equal
		{name: "not equal satisfied", rangeStr: "!=1.5.0", version: "1.4.9", want: true},
		{name: "not equal not satisfied", rangeStr: "!=1.5.0", version: "1.5.0", want: false},

		// Multiple constraints (AND logic)
		{name: "range satisfied", rangeStr: ">=1.0.0 <2.0.0", version: "1.5.0", want: true},
		{name: "range not satisfied lower", rangeStr: ">=1.0.0 <2.0.0", version: "0.9.9", want: false},
		{name: "range not satisfied upper", rangeStr: ">=1.0.0 <2.0.0", version: "2.0.0", want: false},
		{name: "range with exclusion satisfied", rangeStr: ">=1.0.0 <2.0.0 !=1.5.0", version: "1.4.9", want: true},
		{name: "range with exclusion not satisfied", rangeStr: ">=1.0.0 <2.0.0 !=1.5.0", version: "1.5.0", want: false},

		// Alpine-specific version formats
		{name: "suffix version in range", rangeStr: ">=1.2.0_alpha", version: "1.2.0_beta", want: true},
		{name: "suffix version not in range", rangeStr: ">=1.2.0_beta", version: "1.2.0_alpha", want: false},
		{name: "build version in range", rangeStr: ">=1.2.3-r1", version: "1.2.3-r2", want: true},
		{name: "build version not in range", rangeStr: ">=1.2.3-r2", version: "1.2.3-r1", want: false},
		{name: "letter version comparison", rangeStr: ">1.1a", version: "1.1b", want: true},
		{name: "letter version comparison not satisfied", rangeStr: ">1.1b", version: "1.1a", want: false},

		// Edge cases
		{name: "release vs alpha", rangeStr: ">1.0.0_alpha", version: "1.0.0", want: true},
		{name: "alpha vs release", rangeStr: ">1.0.0", version: "1.0.0_alpha", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			versionRange, err := ecosystem.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange(%q) error: %v", tt.rangeStr, err)
			}

			version, err := ecosystem.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion(%q) error: %v", tt.version, err)
			}

			got := versionRange.Contains(version)
			if got != tt.want {
				t.Errorf("VersionRange.Contains(%q in %q) = %t, want %t", tt.version, tt.rangeStr, got, tt.want)
			}
		})
	}
}
