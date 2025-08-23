package vers

import (
	"testing"
)

// TestContains_NPM tests VERS functionality specifically for the NPM ecosystem
func TestContains_NPM(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "npm simple range - contained",
			versRange: "vers:npm/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm simple range - not contained",
			versRange: "vers:npm/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm exact match",
			versRange: "vers:npm/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm exact match - not equal",
			versRange: "vers:npm/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm lower bound only",
			versRange: "vers:npm/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm upper bound only",
			versRange: "vers:npm/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm with prerelease versions",
			versRange: "vers:npm/>=1.0.0-alpha|<=2.0.0",
			version:   "1.5.0-beta",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm complex range - multiple intervals",
			versRange: "vers:npm/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm complex range - second interval",
			versRange: "vers:npm/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "2.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm complex range - outside intervals",
			versRange: "vers:npm/>=1.0.0|<=1.5.0|>=2.0.0|<=2.5.0",
			version:   "1.8.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm exclude constraint",
			versRange: "vers:npm/!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm exclude constraint - allowed version",
			versRange: "vers:npm/!=1.5.0",
			version:   "1.6.0",
			want:      true,
			wantErr:   false,
		},
		// Additional edge case tests for != operator in NPM
		{
			name:      "npm exclude prerelease version",
			versRange: "vers:npm/!=1.0.0-alpha.1",
			version:   "1.0.0-alpha.1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm exclude prerelease - different prerelease allowed",
			versRange: "vers:npm/!=1.0.0-alpha.1",
			version:   "1.0.0-alpha.2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm multiple excludes with prerelease",
			versRange: "vers:npm/!=1.0.0|!=1.0.0-alpha",
			version:   "1.0.0-alpha",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "npm star constraint",
			versRange: "vers:npm/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm multiple lower bounds - should use most restrictive",
			versRange: "vers:npm/>=1.0.0|>=2.0.0|<=3.0.0",
			version:   "1.5.0",
			want:      false, // Should use most restrictive lower bound (>=2.0.0)
			wantErr:   false,
		},
		{
			name:      "npm multiple upper bounds - should use most restrictive",
			versRange: "vers:npm/>=1.0.0|<=3.0.0|<=2.0.0",
			version:   "2.5.0",
			want:      false, // Should use most restrictive upper bound (<=2.0.0)
			wantErr:   false,
		},
		{
			name:      "npm greater than constraint",
			versRange: "vers:npm/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm less than constraint",
			versRange: "vers:npm/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "npm invalid version in NPM",
			versRange: "vers:npm/>=1.0.0",
			version:   "invalid-npm-version",
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
