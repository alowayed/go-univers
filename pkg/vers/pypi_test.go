package vers

import (
	"testing"
)

// TestContains_PyPI tests VERS functionality specifically for the PyPI ecosystem
func TestContains_PyPI(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "pypi simple range - contained",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi simple range - not contained",
			versRange: "vers:pypi/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exact match",
			versRange: "vers:pypi/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi exact match - not equal",
			versRange: "vers:pypi/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi lower bound only",
			versRange: "vers:pypi/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi upper bound only",
			versRange: "vers:pypi/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with epoch versions",
			versRange: "vers:pypi/>=1!1.0.0|<=1!2.0.0",
			version:   "1!1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with prerelease versions - alpha",
			versRange: "vers:pypi/>=1.0.0a1|<=2.0.0",
			version:   "1.5.0b1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with prerelease versions - beta",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0b1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi with prerelease versions - rc",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0rc1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi with post-release versions",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0.post1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi with development release versions",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0.dev1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi with local version identifiers",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0+local.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi complex range - multiple intervals",
			versRange: "vers:pypi/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi complex range - second interval",
			versRange: "vers:pypi/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "2.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi complex range - outside intervals",
			versRange: "vers:pypi/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.8.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exclude constraint",
			versRange: "vers:pypi/!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exclude constraint - allowed version",
			versRange: "vers:pypi/!=1.5.0",
			version:   "1.6.0",
			want:      true,
			wantErr:   false,
		},
		// Additional edge case tests for != operator in PyPI
		{
			name:      "pypi exclude with local version identifier",
			versRange: "vers:pypi/!=1.0.0+local",
			version:   "1.0.0+local",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exclude with local version - different local also excluded per PEP 440",
			versRange: "vers:pypi/!=1.0.0+local1",
			version:   "1.0.0+local2",
			want:      false, // Local versions are ignored in comparison per PEP 440
			wantErr:   false,
		},
		{
			name:      "pypi exclude with local version - different base version allowed",
			versRange: "vers:pypi/!=1.0.0+local1",
			version:   "1.0.1+local1",
			want:      true, // Different base version (1.0.1 vs 1.0.0)
			wantErr:   false,
		},
		{
			name:      "pypi exclude prerelease version",
			versRange: "vers:pypi/!=1.0.0a1",
			version:   "1.0.0a1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi exclude dev version",
			versRange: "vers:pypi/!=1.0.0.dev1",
			version:   "1.0.0.dev1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi star constraint",
			versRange: "vers:pypi/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi multiple lower bounds - should use most restrictive",
			versRange: "vers:pypi/>=1.0.0|>=2.0.0|<=3.0.0",
			version:   "1.5.0",
			want:      false, // Should use most restrictive lower bound (>=2.0.0)
			wantErr:   false,
		},
		{
			name:      "pypi multiple upper bounds - should use most restrictive",
			versRange: "vers:pypi/>=1.0.0|<=3.0.0|<=2.0.0",
			version:   "2.5.0",
			want:      false, // Should use most restrictive upper bound (<=2.0.0)
			wantErr:   false,
		},
		{
			name:      "pypi greater than constraint",
			versRange: "vers:pypi/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi less than constraint",
			versRange: "vers:pypi/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi epoch precedence test",
			versRange: "vers:pypi/>=1!1.0.0",
			version:   "2.0.0",
			want:      false, // Epoch 1 > Epoch 0 (default)
			wantErr:   false,
		},
		{
			name:      "pypi epoch precedence test - satisfied",
			versRange: "vers:pypi/>=1!1.0.0",
			version:   "1!1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease vs release precedence",
			versRange: "vers:pypi/>=1.0.0",
			version:   "1.0.0a1",
			want:      false, // Prerelease < release
			wantErr:   false,
		},
		{
			name:      "pypi prerelease vs release precedence - satisfied",
			versRange: "vers:pypi/>=1.0.0a1",
			version:   "1.0.0",
			want:      true, // Release > prerelease
			wantErr:   false,
		},
		// Additional PyPI prerelease exclusion edge cases
		{
			name:      "pypi prerelease exclusion - alpha format 1",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0a1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease exclusion - alpha format 2",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0alpha1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease exclusion - beta format 1",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0beta1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease exclusion - dot dev format",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0.dev1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease inclusion - explicit alpha in constraint",
			versRange: "vers:pypi/>=1.0.0alpha1|<=2.0.0",
			version:   "1.5.0alpha1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease inclusion - explicit beta in constraint",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0beta1",
			version:   "1.5.0beta1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi prerelease inclusion - explicit dev in constraint",
			versRange: "vers:pypi/>=1.0.0.dev1|<=2.0.0",
			version:   "1.5.0.dev1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi local version with prerelease-like text - should be included",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0+alpha.build",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi local version with dev-like text - should be included",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0+dev.build",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "pypi mixed prerelease and local version",
			versRange: "vers:pypi/>=1.0.0|<=2.0.0",
			version:   "1.5.0a1+local.build",
			want:      false, // Should be excluded due to prerelease
			wantErr:   false,
		},
		{
			name:      "pypi mixed prerelease and local version - explicit inclusion",
			versRange: "vers:pypi/>=1.0.0a1|<=2.0.0",
			version:   "1.5.0a1+local.build",
			want:      true, // Should be included due to explicit prerelease in constraint
			wantErr:   false,
		},
		{
			name:      "pypi invalid version",
			versRange: "vers:pypi/>=1.0.0",
			version:   "invalid-pypi-version",
			want:      false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Contains(tt.versRange, tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
