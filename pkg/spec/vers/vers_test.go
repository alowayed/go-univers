package vers

import (
	"testing"
)

// TestContains tests general VERS functionality including format validation,
// parsing errors, and cross-ecosystem behavior.
func TestContains(t *testing.T) {
	tests := []struct {
		name      string
		versRange string
		version   string
		want      bool
		wantErr   bool
	}{
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
		// VERS validation edge cases
		{
			name:      "empty ecosystem name",
			versRange: "vers:/>1.0.0",
			version:   "1.1.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "empty constraints",
			versRange: "vers:maven/",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "missing colon",
			versRange: "versmaven/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "missing slash separator",
			versRange: "vers:maven>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "uppercase ecosystem should fail per VERS spec",
			versRange: "vers:MAVEN/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true, // VERS spec requires lowercase
		},
		{
			name:      "multiple stars should fail",
			versRange: "vers:maven/*|*",
			version:   "1.0.0",
			want:      false,
			wantErr:   true, // VERS spec allows only one star
		},
		{
			name:      "mixed case ecosystem should fail",
			versRange: "vers:Maven/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		// Additional ecosystem name validation tests
		{
			name:      "ecosystem with hyphen should fail",
			versRange: "vers:my-ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with underscore should fail",
			versRange: "vers:my_ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with period should fail",
			versRange: "vers:my.ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with special characters should fail",
			versRange: "vers:my@ecosystem/>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "ecosystem with numbers should pass validation but fail as unsupported",
			versRange: "vers:ecosystem2/>=1.0.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   true, // Should fail due to unsupported ecosystem, not validation
		},
		{
			name:      "ecosystem with only numbers should pass validation but fail as unsupported",
			versRange: "vers:123/>=1.0.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   true, // Should fail due to unsupported ecosystem, not validation
		},
		{
			name:      "ecosystem with mixed letters and numbers should pass validation but fail as unsupported",
			versRange: "vers:maven2/>=1.0.0",
			version:   "1.5.0",
			want:      false,
			wantErr:   true, // Should fail due to unsupported ecosystem, not validation
		},
		// Additional VERS specification compliance tests
		{
			name:      "non-printable ASCII character should fail",
			versRange: "vers:maven/>=1.0.0\x00",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "tab character should fail",
			versRange: "vers:maven/>=1.0.0\t",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "newline character should fail",
			versRange: "vers:maven/>=1.0.0\n",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "star with other non-empty constraints should fail",
			versRange: "vers:maven/*|>=1.0.0",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "star with whitespace-only constraints should pass",
			versRange: "vers:maven/*| | ",
			version:   "1.0.0",
			want:      true,
			wantErr:   false,
		},
		// Additional edge case tests for != operator validation
		{
			name:      "empty version after != operator should fail",
			versRange: "vers:maven/!=",
			version:   "1.0.0",
			want:      false,
			wantErr:   true,
		},
		{
			name:      "whitespace-only version after != operator should fail",
			versRange: "vers:maven/!= ",
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
