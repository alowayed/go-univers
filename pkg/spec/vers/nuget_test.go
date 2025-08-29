package vers

import (
	"testing"
)

// TestContains_NuGet tests VERS functionality specifically for the NuGet ecosystem
func TestContains_NuGet(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "nuget simple range - contained",
			versRange: "vers:nuget/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget simple range - not contained",
			versRange: "vers:nuget/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget exact match",
			versRange: "vers:nuget/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget exact match - not equal",
			versRange: "vers:nuget/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget lower bound only",
			versRange: "vers:nuget/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget upper bound only",
			versRange: "vers:nuget/<=2.0.0",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget greater than",
			versRange: "vers:nuget/>1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget greater than - not satisfied",
			versRange: "vers:nuget/>1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget less than",
			versRange: "vers:nuget/<2.0.0",
			version:   "1.9.9",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget less than - not satisfied",
			versRange: "vers:nuget/<2.0.0",
			version:   "2.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget not equal - satisfied",
			versRange: "vers:nuget/!=1.0.0",
			version:   "1.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget not equal - not satisfied",
			versRange: "vers:nuget/!=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget multiple constraints - AND logic",
			versRange: "vers:nuget/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.2.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget multiple constraints - excluded",
			versRange: "vers:nuget/>=1.0.0|<=2.0.0|!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget prerelease version - contained",
			versRange: "vers:nuget/>=1.0.0-alpha",
			version:   "1.0.0-alpha.2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget prerelease version - not contained",
			versRange: "vers:nuget/>=1.0.0-beta",
			version:   "1.0.0-alpha.2",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget prerelease vs release",
			versRange: "vers:nuget/>=1.0.0",
			version:   "1.0.0-alpha",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget release vs prerelease",
			versRange: "vers:nuget/>=1.0.0-alpha",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget revision version - contained",
			versRange: "vers:nuget/>=1.0.0.1",
			version:   "1.0.0.2",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget revision version - not contained",
			versRange: "vers:nuget/>=1.0.0.2",
			version:   "1.0.0.1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget revision version range",
			versRange: "vers:nuget/>=1.0.0.1|<2.0.0.0",
			version:   "1.5.0.5",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget revision version range - not contained",
			versRange: "vers:nuget/>=1.0.0.1|<2.0.0.0",
			version:   "2.0.0.1",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget prerelease with revision",
			versRange: "vers:nuget/>=1.0.0.1-alpha|<2.0.0.0",
			version:   "1.5.0.2-beta",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget prerelease with revision - not contained",
			versRange: "vers:nuget/>=1.0.0.1-beta|<2.0.0.0",
			version:   "1.0.0.1-alpha",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget complex version range",
			versRange: "vers:nuget/>=1.0.0.1|<=2.0.0.0|!=1.5.0.0",
			version:   "1.2.0.5",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget complex version range - excluded",
			versRange: "vers:nuget/>=1.0.0.1|<=2.0.0.0|!=1.5.0.0",
			version:   "1.5.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "nuget star constraint - matches all",
			versRange: "vers:nuget/*",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget star constraint - matches prerelease",
			versRange: "vers:nuget/*",
			version:   "1.0.0-alpha",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget star constraint - matches revision",
			versRange: "vers:nuget/*",
			version:   "1.0.0.1",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "nuget star constraint - matches complex",
			versRange: "vers:nuget/*",
			version:   "1.0.0.1-alpha.2",
			want:      true,
			wantErr:   false,
		},
		// Error cases
		{
			name:      "nuget invalid version",
			versRange: "vers:nuget/>=1.0.0",
			version:   "1.0.0@invalid",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "nuget invalid constraint version",
			versRange: "vers:nuget/>=1..0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
