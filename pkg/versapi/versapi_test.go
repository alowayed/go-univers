package versapi

import (
	"testing"
)

func TestContains(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
		{
			name:      "maven simple range - contained",
			versRange: "vers:maven/>=1.0.0|<=2.0.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven simple range - not contained",
			versRange: "vers:maven/>=2.0.0|<=3.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "maven complex range from issue",
			versRange: "vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7",
			version:   "1.1.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven exact match",
			versRange: "vers:maven/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven with exclusion",
			versRange: "vers:maven/>=1.0.0|!=1.5.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "invalid VERS format",
			versRange: "not-vers-format",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "unsupported ecosystem",
			versRange: "vers:unsupported/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "invalid version",
			versRange: "vers:maven/>=1.0.0",
			version:   "invalid-version",
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

func TestCompare(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version1  string
		version2  string
		want      int
		wantErr   bool
	}{
		{
			name:      "maven versions - first less than second",
			versRange: "vers:maven/>=1.0.0",
			version1:  "1.0.0",
			version2:  "2.0.0",
			want:      -1,
			wantErr:   false,
		},
		{
			name:      "maven versions - equal",
			versRange: "vers:maven/>=1.0.0",
			version1:  "1.5.0",
			version2:  "1.5.0",
			want:      0,
			wantErr:   false,
		},
		{
			name:      "maven versions - first greater than second",
			versRange: "vers:maven/>=1.0.0",
			version1:  "2.0.0",
			version2:  "1.0.0",
			want:      1,
			wantErr:   false,
		},
		{
			name:      "invalid VERS format",
			versRange: "not-vers-format",
			version1:  "1.0.0",
			version2:  "2.0.0",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "unsupported ecosystem",
			versRange: "vers:unsupported/>=1.0.0",
			version1:  "1.0.0",
			version2:  "2.0.0",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "invalid first version",
			versRange: "vers:maven/>=1.0.0",
			version1:  "invalid-version",
			version2:  "2.0.0",
			want:      0,
			wantErr:   true,
		},
		{
			name:      "invalid second version",
			versRange: "vers:maven/>=1.0.0",
			version1:  "1.0.0",
			version2:  "invalid-version",
			want:      0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compare(tt.versRange, tt.version1, tt.version2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}