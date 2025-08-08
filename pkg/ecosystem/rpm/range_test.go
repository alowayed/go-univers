package rpm

import (
	"testing"
)

func TestEcosystem_NewVersionRange(t *testing.T) {
	tests := []struct {
		name      string
		rangeStr  string
		wantCount int
		wantErr   bool
	}{
		// Single constraints
		{
			name:      "exact version",
			rangeStr:  "1.2.3",
			wantCount: 1,
		},
		{
			name:      "greater than",
			rangeStr:  ">1.2.3",
			wantCount: 1,
		},
		{
			name:      "greater than or equal",
			rangeStr:  ">=1.2.3",
			wantCount: 1,
		},
		{
			name:      "less than",
			rangeStr:  "<1.2.3",
			wantCount: 1,
		},
		{
			name:      "less than or equal",
			rangeStr:  "<=1.2.3",
			wantCount: 1,
		},
		{
			name:      "not equal",
			rangeStr:  "!=1.2.3",
			wantCount: 1,
		},
		{
			name:      "explicit equal",
			rangeStr:  "=1.2.3",
			wantCount: 1,
		},

		// Multiple constraints (space-separated)
		{
			name:      "range with spaces",
			rangeStr:  ">=1.2.0 <2.0.0",
			wantCount: 2,
		},
		{
			name:      "multiple constraints",
			rangeStr:  ">1.0.0 <=2.0.0 !=1.5.0",
			wantCount: 3,
		},

		// Multiple constraints (comma-separated)
		{
			name:      "range with commas",
			rangeStr:  ">=1.2.0,<2.0.0",
			wantCount: 2,
		},
		{
			name:      "mixed separators",
			rangeStr:  ">=1.0.0, <=2.0.0, !=1.5.0",
			wantCount: 3,
		},

		// Complex versions in ranges
		{
			name:      "range with epoch",
			rangeStr:  ">=2:1.0.0",
			wantCount: 1,
		},
		{
			name:      "range with release",
			rangeStr:  ">=1.2.3-4",
			wantCount: 1,
		},
		{
			name:      "range with full version",
			rangeStr:  ">=1:2.3.4-5",
			wantCount: 1,
		},

		// Error cases
		{
			name:     "empty range",
			rangeStr: "",
			wantErr:  true,
		},
		{
			name:     "whitespace only",
			rangeStr: "   ",
			wantErr:  true,
		},
		{
			name:     "operator without version",
			rangeStr: ">=",
			wantErr:  true,
		},
		{
			name:     "invalid version in range",
			rangeStr: ">=1.0.0@invalid",
			wantErr:  true,
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersionRange(tt.rangeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(got.constraints) != tt.wantCount {
				t.Errorf("Ecosystem.NewVersionRange() constraint count = %v, want %v", len(got.constraints), tt.wantCount)
			}

			if got.original != tt.rangeStr {
				t.Errorf("Ecosystem.NewVersionRange() original = %v, want %v", got.original, tt.rangeStr)
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
			rangeStr: "1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "exact match explicit",
			rangeStr: "=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "exact match fail",
			rangeStr: "1.2.3",
			version:  "1.2.4",
			want:     false,
		},

		// Greater than
		{
			name:     "greater than pass",
			rangeStr: ">1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than fail",
			rangeStr: ">1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "greater than fail lower",
			rangeStr: ">1.2.3",
			version:  "1.2.2",
			want:     false,
		},

		// Greater than or equal
		{
			name:     "greater than or equal pass higher",
			rangeStr: ">=1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than or equal pass equal",
			rangeStr: ">=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "greater than or equal fail",
			rangeStr: ">=1.2.3",
			version:  "1.2.2",
			want:     false,
		},

		// Less than
		{
			name:     "less than pass",
			rangeStr: "<1.2.3",
			version:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than fail equal",
			rangeStr: "<1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "less than fail higher",
			rangeStr: "<1.2.3",
			version:  "1.2.4",
			want:     false,
		},

		// Less than or equal
		{
			name:     "less than or equal pass lower",
			rangeStr: "<=1.2.3",
			version:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than or equal pass equal",
			rangeStr: "<=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "less than or equal fail",
			rangeStr: "<=1.2.3",
			version:  "1.2.4",
			want:     false,
		},

		// Not equal
		{
			name:     "not equal pass higher",
			rangeStr: "!=1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "not equal pass lower",
			rangeStr: "!=1.2.3",
			version:  "1.2.2",
			want:     true,
		},
		{
			name:     "not equal fail",
			rangeStr: "!=1.2.3",
			version:  "1.2.3",
			want:     false,
		},

		// Range constraints (space-separated)
		{
			name:     "range pass inside",
			rangeStr: ">=1.2.0 <2.0.0",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "range pass lower bound",
			rangeStr: ">=1.2.0 <2.0.0",
			version:  "1.2.0",
			want:     true,
		},
		{
			name:     "range fail upper bound",
			rangeStr: ">=1.2.0 <2.0.0",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "range fail too low",
			rangeStr: ">=1.2.0 <2.0.0",
			version:  "1.1.0",
			want:     false,
		},
		{
			name:     "range fail too high",
			rangeStr: ">=1.2.0 <2.0.0",
			version:  "2.1.0",
			want:     false,
		},

		// Range constraints (comma-separated)
		{
			name:     "comma range pass",
			rangeStr: ">=1.2.0,<2.0.0",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "comma range fail",
			rangeStr: ">=1.2.0,<2.0.0",
			version:  "2.0.0",
			want:     false,
		},

		// Complex constraints
		{
			name:     "multiple constraints pass",
			rangeStr: ">1.0.0 <=2.0.0 !=1.5.0",
			version:  "1.4.0",
			want:     true,
		},
		{
			name:     "multiple constraints fail exclusion",
			rangeStr: ">1.0.0 <=2.0.0 !=1.5.0",
			version:  "1.5.0",
			want:     false,
		},
		{
			name:     "multiple constraints fail range",
			rangeStr: ">1.0.0 <=2.0.0 !=1.5.0",
			version:  "3.0.0",
			want:     false,
		},

		// Epoch handling in ranges
		{
			name:     "epoch constraint pass",
			rangeStr: ">=2:1.0.0",
			version:  "2:1.5.0",
			want:     true,
		},
		{
			name:     "epoch constraint fail",
			rangeStr: ">=2:1.0.0",
			version:  "1:2.0.0",
			want:     false,
		},
		{
			name:     "epoch constraint no epoch vs epoch",
			rangeStr: ">=2:1.0.0",
			version:  "3.0.0",
			want:     false,
		},

		// Release handling in ranges
		{
			name:     "release constraint pass",
			rangeStr: ">=1.2.3-2",
			version:  "1.2.3-3",
			want:     true,
		},
		{
			name:     "release constraint fail",
			rangeStr: ">=1.2.3-2",
			version:  "1.2.3-1",
			want:     false,
		},
		{
			name:     "release vs no release",
			rangeStr: ">=1.2.3-1",
			version:  "1.2.3",
			want:     false,
		},

		// Tilde handling in ranges
		{
			name:     "tilde range pass",
			rangeStr: ">=1.0~beta <1.0",
			version:  "1.0~rc1",
			want:     true,
		},
		{
			name:     "tilde range fail",
			rangeStr: ">1.0~beta",
			version:  "1.0~alpha",
			want:     false,
		},

		// Real-world examples
		{
			name:     "kernel version range",
			rangeStr: ">=3.10.0-1062.el7",
			version:  "3.10.0-1127.el7",
			want:     true,
		},
		{
			name:     "glibc version range",
			rangeStr: ">=2.17-290.el7 <2.17-310.el7",
			version:  "2.17-305.el7",
			want:     true,
		},
		{
			name:     "exclude specific version",
			rangeStr: ">=1.0.0 !=1.2.3-broken",
			version:  "1.2.3-1",
			want:     true,
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := ecosystem.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange() error = %v", err)
			}

			v, err := ecosystem.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion() error = %v", err)
			}

			if got := vr.Contains(v); got != tt.want {
				t.Errorf("VersionRange.Contains() = %v, want %v (range=%s, version=%s)", got, tt.want, tt.rangeStr, tt.version)
			}
		})
	}
}

func TestVersionRange_String(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
	}{
		{"simple range", ">=1.2.3"},
		{"complex range", ">=1.0.0 <2.0.0 !=1.5.0"},
		{"comma separated", ">=1.0.0,<2.0.0"},
		{"epoch in range", ">=2:1.0.0"},
		{"release in range", ">=1.2.3-4"},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := ecosystem.NewVersionRange(tt.rangeStr)
			if err != nil {
				t.Fatalf("NewVersionRange() error = %v", err)
			}

			if got := vr.String(); got != tt.rangeStr {
				t.Errorf("VersionRange.String() = %v, want %v", got, tt.rangeStr)
			}
		})
	}
}

// Test constraint parsing edge cases
func TestParseRPMConstraints(t *testing.T) {
	tests := []struct {
		name     string
		rangeStr string
		wantOps  []string
		wantVers []string
		wantErr  bool
	}{
		{
			name:     "single constraint no operator",
			rangeStr: "1.2.3",
			wantOps:  []string{"="},
			wantVers: []string{"1.2.3"},
		},
		{
			name:     "single constraint with operator",
			rangeStr: ">=1.2.3",
			wantOps:  []string{">="},
			wantVers: []string{"1.2.3"},
		},
		{
			name:     "multiple space separated",
			rangeStr: ">=1.0.0 <2.0.0",
			wantOps:  []string{">=", "<"},
			wantVers: []string{"1.0.0", "2.0.0"},
		},
		{
			name:     "multiple comma separated",
			rangeStr: ">=1.0.0,<2.0.0",
			wantOps:  []string{">=", "<"},
			wantVers: []string{"1.0.0", "2.0.0"},
		},
		{
			name:     "mixed separators with spaces",
			rangeStr: ">=1.0.0, <2.0.0, !=1.5.0",
			wantOps:  []string{">=", "<", "!="},
			wantVers: []string{"1.0.0", "2.0.0", "1.5.0"},
		},
		{
			name:     "operator without version",
			rangeStr: ">= ",
			wantErr:  true,
		},
		{
			name:     "empty after processing",
			rangeStr: " , , ",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ecosystem := &Ecosystem{}
			constraints, err := parseRPMConstraints(ecosystem, tt.rangeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseRPMConstraints() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if len(constraints) != len(tt.wantOps) {
				t.Errorf("parseRPMConstraints() got %d constraints, want %d", len(constraints), len(tt.wantOps))
				return
			}

			for i, c := range constraints {
				if c.operator != tt.wantOps[i] {
					t.Errorf("parseRPMConstraints() constraint %d operator = %v, want %v", i, c.operator, tt.wantOps[i])
				}
				if c.version.String() != tt.wantVers[i] {
					t.Errorf("parseRPMConstraints() constraint %d version = %v, want %v", i, c.version.String(), tt.wantVers[i])
				}
			}
		})
	}
}
