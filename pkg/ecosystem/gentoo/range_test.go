package gentoo

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		rangeStr string
		wantErr  bool
	}{
		// Valid ranges
		{">=1.0.0", false},
		{"<2.0.0", false},
		{"!=1.5.0", false},
		{">1.0.0", false},
		{"<=2.0.0", false},
		{"=1.2.3", false},
		{"1.2.3", false}, // defaults to exact match
		{">=1.0.0, <2.0.0", false},
		{">=1.0.0 <2.0.0", false},
		{">1.0_alpha, <=2.0_beta", false},
		{">=1.0-r1, <2.0-r5", false},

		// Invalid ranges
		{"", true},
		{">=", true},
		{">", true},
		{"<", true},
		{"<=", true},
		{"!=", true},
		{"= ", true},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.rangeStr, func(t *testing.T) {
			r, err := e.NewVersionRange(tt.rangeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange(%q) error = %v, wantErr %v", tt.rangeStr, err, tt.wantErr)
				return
			}
			if !tt.wantErr && r.String() != tt.rangeStr {
				t.Errorf("NewVersionRange(%q).String() = %v, want %v", tt.rangeStr, r.String(), tt.rangeStr)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		rangeStr string
		version  string
		want     bool
	}{
		// Basic comparisons
		{">=1.0.0", "1.0.0", true},
		{">=1.0.0", "1.0.1", true},
		{">=1.0.0", "0.9.9", false},
		{">1.0.0", "1.0.0", false},
		{">1.0.0", "1.0.1", true},
		{"<2.0.0", "1.9.9", true},
		{"<2.0.0", "2.0.0", false},
		{"<=2.0.0", "2.0.0", true},
		{"<=2.0.0", "2.0.1", false},
		{"=1.2.3", "1.2.3", true},
		{"=1.2.3", "1.2.4", false},
		{"!=1.2.3", "1.2.4", true},
		{"!=1.2.3", "1.2.3", false},

		// Suffix handling
		{">=1.0_alpha", "1.0_alpha", true},
		{">=1.0_alpha", "1.0_beta", true},
		{">=1.0_alpha", "1.0", true},
		{">=1.0_alpha", "0.9", false},
		{">1.0_alpha", "1.0_alpha", false},
		{">1.0_alpha", "1.0_beta", true},
		{"<1.0", "1.0_alpha", true},
		{"<1.0", "1.0_rc", true},
		{"<1.0", "1.0", false},

		// Revision handling
		{">=1.0-r1", "1.0", false}, // 1.0 is 1.0-r0
		{">=1.0-r1", "1.0-r1", true},
		{">=1.0-r1", "1.0-r2", true},
		{">1.0", "1.0-r1", true},
		{"<1.0-r2", "1.0-r1", true},
		{"<1.0-r2", "1.0-r2", false},

		// Multiple constraints (AND logic)
		{">=1.0.0, <2.0.0", "1.5.0", true},
		{">=1.0.0, <2.0.0", "0.9.9", false},
		{">=1.0.0, <2.0.0", "2.0.0", false},
		{">=1.0.0 <2.0.0", "1.5.0", true},
		{">=1.0.0 <2.0.0", "0.9.9", false},
		{">=1.0.0 <2.0.0", "2.0.0", false},
		{">1.0_alpha, <=2.0_beta", "1.0", true},
		{">1.0_alpha, <=2.0_beta", "1.0_alpha", false},
		{">1.0_alpha, <=2.0_beta", "2.0_beta", true},
		{">1.0_alpha, <=2.0_beta", "2.0", false}, // 2.0 > 2.0_beta

		// Complex version constraints
		{">=1.0a_alpha1-r2, <2.0b_rc3-r5", "1.5c_beta2-r3", true},
		{">=1.0a_alpha1-r2, <2.0b_rc3-r5", "1.0a_alpha1-r1", false},
		{">=1.0a_alpha1-r2, <2.0b_rc3-r5", "2.0b_rc3-r5", false},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.rangeStr+"_contains_"+tt.version, func(t *testing.T) {
			r, err := e.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange(%q) failed: %v", tt.rangeStr, err)
			}
			v, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion(%q) failed: %v", tt.version, err)
			}

			got := r.Contains(v)
			if got != tt.want {
				t.Errorf("Range %q contains %q = %v, want %v", tt.rangeStr, tt.version, got, tt.want)
			}
		})
	}
}
