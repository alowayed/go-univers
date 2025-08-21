package vers

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
			name:      "maven exact match",
			versRange: "vers:maven/=1.5.0",
			version:   "1.5.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven exact match - not equal",
			versRange: "vers:maven/=1.5.0",
			version:   "1.6.0",
			want:      false,
			wantErr:   false,
		},
		{
			name:      "maven lower bound only",
			versRange: "vers:maven/>=1.0.0",
			version:   "2.0.0",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven upper bound only",
			versRange: "vers:maven/<=2.0.0",
			version:   "1.0.0",
			want:      true,
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
			name:      "maven complex range - second interval",
			versRange: "vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7",
			version:   "7.0.5",
			want:      true,
			wantErr:   false,
		},
		{
			name:      "maven complex range - outside intervals",
			versRange: "vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7",
			version:   "5.0.0",
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
		{
			name:      "invalid constraint format",
			versRange: "vers:maven/~1.0.0",
			version:   "1.0.0",
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
