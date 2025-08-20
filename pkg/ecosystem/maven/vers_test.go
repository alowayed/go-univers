package maven

import (
	"testing"
)

func TestParseVers(t *testing.T) {
	tests := []struct {
		name        string
		versString  string
		testVersion string
		wantContains bool
		wantErr     bool
	}{
		{
			name:         "valid range with >=",
			versString:   "vers:maven/>=1.0.0",
			testVersion:  "1.5.0",
			wantContains: true,
			wantErr:      false,
		},
		{
			name:         "valid range with <=",
			versString:   "vers:maven/<=2.0.0",
			testVersion:  "1.5.0",
			wantContains: true,
			wantErr:      false,
		},
		{
			name:         "valid range with multiple constraints",
			versString:   "vers:maven/>=1.0.0|<=2.0.0",
			testVersion:  "1.5.0",
			wantContains: true,
			wantErr:      false,
		},
		{
			name:         "version outside range",
			versString:   "vers:maven/>=2.0.0|<=3.0.0",
			testVersion:  "1.0.0",
			wantContains: false,
			wantErr:      false,
		},
		{
			name:         "exact match with =",
			versString:   "vers:maven/=1.5.0",
			testVersion:  "1.5.0",
			wantContains: true,
			wantErr:      false,
		},
		{
			name:         "exact match fails",
			versString:   "vers:maven/=1.5.0",
			testVersion:  "1.6.0",
			wantContains: false,
			wantErr:      false,
		},
		{
			name:         "complex range example from issue",
			versString:   "vers:maven/>=1.0.0-beta1|<=1.7.5|>=7.0.0-M1|<=7.0.7",
			testVersion:  "1.1.0",
			wantContains: true,
			wantErr:      false,
		},
		{
			name:       "wrong ecosystem",
			versString: "vers:npm/>=1.0.0",
			wantErr:    true,
		},
		{
			name:       "invalid VERS format",
			versString: "not-vers-format",
			wantErr:    true,
		},
		{
			name:         "invalid version in constraint - caught at Contains time",
			versString:   "vers:maven/>=invalid-version",
			testVersion:  "1.0.0",
			wantContains: false, // Should not contain due to invalid constraint
			wantErr:      false, // ParseVers succeeds, but Contains fails gracefully
		},
		{
			name:         "!= operator - supported in VERS algorithm",
			versString:   "vers:maven/>=1.0.0|!=1.5.0",
			testVersion:  "1.5.0",
			wantContains: false, // 1.5.0 is excluded by !=
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := ParseVers(tt.versString)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseVers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return // Don't test Contains if we expected an error
			}

			// Test Contains functionality
			if tt.testVersion != "" {
				e := &Ecosystem{}
				version, err := e.NewVersion(tt.testVersion)
				if err != nil {
					t.Errorf("Failed to create test version: %v", err)
					return
				}

				got := vr.Contains(version)
				if got != tt.wantContains {
					t.Errorf("Contains() = %v, want %v for version %s in range %s", got, tt.wantContains, tt.testVersion, tt.versString)
				}
			}
		})
	}
}

func TestParseVers_OriginalString(t *testing.T) {
	versString := "vers:maven/>=1.0.0|<=2.0.0"
	vr, err := ParseVers(versString)
	if err != nil {
		t.Errorf("ParseVers() error = %v", err)
		return
	}

	if vr.String() != versString {
		t.Errorf("String() = %v, want %v", vr.String(), versString)
	}
}