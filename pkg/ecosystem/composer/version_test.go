package composer

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     *Version
		wantErr  bool
	}{
		// Basic semantic versions
		{
			name:  "basic semantic version",
			input: "1.2.3",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityStable,
				original:  "1.2.3",
			},
		},
		{
			name:  "version with v prefix",
			input: "v1.2.3",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityStable,
				original:  "v1.2.3",
			},
		},
		{
			name:  "two component version",
			input: "1.2",
			want: &Version{
				major: 1, minor: 2, patch: 0,
				stability: stabilityStable,
				original:  "1.2",
			},
		},
		{
			name:  "single component version",
			input: "1",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityStable,
				original:  "1",
			},
		},

		// Versions with extra component (four parts)
		{
			name:  "four component version",
			input: "1.2.3.4",
			want: &Version{
				major: 1, minor: 2, patch: 3, extra: 4,
				stability: stabilityStable,
				original:  "1.2.3.4",
			},
		},

		// Stability versions
		{
			name:  "alpha version",
			input: "1.2.3-alpha",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityAlpha,
				original:  "1.2.3-alpha",
			},
		},
		{
			name:  "alpha version with number",
			input: "1.2.3-alpha.1",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityAlpha, stabilityNum: 1,
				original: "1.2.3-alpha.1",
			},
		},
		{
			name:  "beta version",
			input: "1.2.3-beta",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityBeta,
				original:  "1.2.3-beta",
			},
		},
		{
			name:  "beta version with number",
			input: "1.2.3-beta.2",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityBeta, stabilityNum: 2,
				original: "1.2.3-beta.2",
			},
		},
		{
			name:  "RC version",
			input: "1.2.3-RC",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityRC,
				original:  "1.2.3-RC",
			},
		},
		{
			name:  "RC version with number",
			input: "1.2.3-RC.1",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityRC, stabilityNum: 1,
				original: "1.2.3-RC.1",
			},
		},
		{
			name:  "rc lowercase version",
			input: "1.2.3-rc.2",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityRC, stabilityNum: 2,
				original: "1.2.3-rc.2",
			},
		},

		// Short stability forms
		{
			name:  "alpha short form",
			input: "1.2.3-a",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityAlpha,
				original:  "1.2.3-a",
			},
		},
		{
			name:  "beta short form",
			input: "1.2.3-b",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityBeta,
				original:  "1.2.3-b",
			},
		},

		// Dev versions
		{
			name:  "dev version with branch",
			input: "dev-main",
			want: &Version{
				isDev: true, devBranch: "main",
				stability: stabilityDev,
				original:  "dev-main",
			},
		},
		{
			name:  "dev version with feature branch",
			input: "dev-feature-auth",
			want: &Version{
				isDev: true, devBranch: "feature-auth",
				stability: stabilityDev,
				original:  "dev-feature-auth",
			},
		},
		{
			name:  "branch name only",
			input: "main",
			want: &Version{
				isDev: true, devBranch: "main",
				stability: stabilityDev,
				original:  "main",
			},
		},

		// Versions with build metadata
		{
			name:  "version with build metadata",
			input: "1.2.3+build.1",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityStable, build: "build.1",
				original: "1.2.3+build.1",
			},
		},
		{
			name:  "prerelease with build metadata",
			input: "1.2.3-alpha.1+build.2",
			want: &Version{
				major: 1, minor: 2, patch: 3,
				stability: stabilityAlpha, stabilityNum: 1, build: "build.2",
				original: "1.2.3-alpha.1+build.2",
			},
		},

		// Additional real-world version formats from test data
		{
			name:  "pl suffix version",
			input: "1.0pl1",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityStable, stabilityNum: 1,
				original:  "1.0pl1",
			},
		},
		{
			name:  "version with v prefix and pl",
			input: "v1.0pl1",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityStable, stabilityNum: 1,
				original:  "v1.0pl1",
			},
		},
		{
			name:  "long numeric version",
			input: "1.2.3.4.5",
			want: &Version{
				major: 1, minor: 2, patch: 3, extra: 4,
				stability: stabilityStable,
				original:  "1.2.3.4.5",
			},
		},
		{
			name:  "alpha without separator",
			input: "1.0a1",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityAlpha, stabilityNum: 1,
				original: "1.0a1",
			},
		},
		{
			name:  "beta without separator",
			input: "1.0b2",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityBeta, stabilityNum: 2,
				original: "1.0b2",
			},
		},
		{
			name:  "RC without separator",
			input: "1.0RC1",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityRC, stabilityNum: 1,
				original: "1.0RC1",
			},
		},
		{
			name:  "version equality with different formats",
			input: "1.0.0",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityStable,
				original:  "1.0.0",
			},
		},
		{
			name:  "patch version format",
			input: "1.0.0-patch1",
			want: &Version{
				major: 1, minor: 0, patch: 0,
				stability: stabilityStable, stabilityNum: 1, // Patch is treated as stable
				original:  "1.0.0-patch1",
			},
		},

		// Edge cases and error cases
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid format",
			input:   "invalid-version",
			wantErr: true,
		},
		{
			name:    "non-numeric major",
			input:   "a.2.3",
			wantErr: true,
		},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := e.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// Compare fields individually for better error messages
			if got.major != tt.want.major {
				t.Errorf("major = %v, want %v", got.major, tt.want.major)
			}
			if got.minor != tt.want.minor {
				t.Errorf("minor = %v, want %v", got.minor, tt.want.minor)
			}
			if got.patch != tt.want.patch {
				t.Errorf("patch = %v, want %v", got.patch, tt.want.patch)
			}
			if got.extra != tt.want.extra {
				t.Errorf("extra = %v, want %v", got.extra, tt.want.extra)
			}
			if got.stability != tt.want.stability {
				t.Errorf("stability = %v, want %v", got.stability, tt.want.stability)
			}
			if got.stabilityNum != tt.want.stabilityNum {
				t.Errorf("stabilityNum = %v, want %v", got.stabilityNum, tt.want.stabilityNum)
			}
			if got.build != tt.want.build {
				t.Errorf("build = %v, want %v", got.build, tt.want.build)
			}
			if got.isDev != tt.want.isDev {
				t.Errorf("isDev = %v, want %v", got.isDev, tt.want.isDev)
			}
			if got.devBranch != tt.want.devBranch {
				t.Errorf("devBranch = %v, want %v", got.devBranch, tt.want.devBranch)
			}
			if got.original != tt.want.original {
				t.Errorf("original = %v, want %v", got.original, tt.want.original)
			}
		})
	}
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{"basic version", "1.2.3", "1.2.3"},
		{"version with v prefix", "v1.2.3", "v1.2.3"},
		{"alpha version", "1.2.3-alpha", "1.2.3-alpha"},
		{"dev version", "dev-main", "dev-main"},
		{"version with build", "1.2.3+build.1", "1.2.3+build.1"},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := e.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion() error = %v", err)
			}
			if got := v.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Basic comparisons
		{"equal versions", "1.2.3", "1.2.3", 0},
		{"major version difference", "2.0.0", "1.9.9", 1},
		{"minor version difference", "1.3.0", "1.2.9", 1},
		{"patch version difference", "1.2.4", "1.2.3", 1},
		{"extra version difference", "1.2.3.1", "1.2.3", 1},

		// Stability comparisons (stable > RC > beta > alpha > dev)
		{"stable vs alpha", "1.2.3", "1.2.3-alpha", 1},
		{"stable vs beta", "1.2.3", "1.2.3-beta", 1},
		{"stable vs RC", "1.2.3", "1.2.3-RC", 1},
		{"RC vs beta", "1.2.3-RC", "1.2.3-beta", 1},
		{"beta vs alpha", "1.2.3-beta", "1.2.3-alpha", 1},
		{"alpha vs dev", "1.2.3-alpha", "dev-main", 1},

		// Same stability level with numbers
		{"alpha.1 vs alpha.2", "1.2.3-alpha.1", "1.2.3-alpha.2", -1},
		{"beta.2 vs beta.1", "1.2.3-beta.2", "1.2.3-beta.1", 1},
		{"RC.1 vs RC.2", "1.2.3-RC.1", "1.2.3-RC.2", -1},

		// Dev version comparisons
		{"dev vs stable", "dev-main", "1.2.3", -1},
		{"dev branch comparison", "dev-feature", "dev-main", -1}, // lexical comparison
		{"dev branch comparison 2", "dev-main", "dev-feature", 1},

		// Mixed comparisons
		{"different major with stability", "2.0.0-alpha", "1.9.9", 1},
		{"same version different stability", "1.2.3-alpha.1", "1.2.3-alpha.2", -1},

		// Version normalization/equivalence tests
		{"v prefix comparison", "v1.2.3", "1.2.3", 0},
		{"zero versions", "0.0.1", "0.0.2", -1},
		{"version normalization 1.0 vs 1.0.0", "1.0", "1.0.0", 0},
		{"version normalization v1 vs 1.0.0", "v1", "1.0.0", 0},
		
		// pl/patch version comparisons  
		{"pl version vs normal", "1.0pl1", "1.0.0", 1}, // pl1 > normal (has stability number)
		{"pl version vs pl version", "1.0pl1", "1.0pl2", -1}, // pl1 vs pl2 by stability number
		
		// Mixed format stability comparisons
		{"alpha without hyphen vs with hyphen", "1.0a1", "1.0.0-alpha.1", 0},
		{"beta without hyphen vs with hyphen", "1.0b1", "1.0.0-beta.1", 0},
		{"RC without hyphen vs with hyphen", "1.0RC1", "1.0.0-RC.1", 0},
		
		// Complex prerelease ordering
		{"alpha vs beta same base", "1.0.0-alpha", "1.0.0-beta", -1},
		{"beta vs RC same base", "1.0.0-beta", "1.0.0-RC", -1},
		{"RC vs stable same base", "1.0.0-RC", "1.0.0", -1},
		
		// Numeric prerelease ordering
		{"alpha.1 vs alpha.10", "1.0.0-alpha.1", "1.0.0-alpha.10", -1},
		{"beta.2 vs beta.11", "1.0.0-beta.2", "1.0.0-beta.11", -1},
	}

	e := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := e.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("NewVersion(%s) error = %v", tt.v1, err)
			}
			v2, err := e.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("NewVersion(%s) error = %v", tt.v2, err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Compare(%s, %s) = %v, want %v", tt.v1, tt.v2, got, tt.want)
			}

			// Test symmetry
			reverse := v2.Compare(v1)
			if got == 0 && reverse != 0 {
				t.Errorf("Compare symmetry failed: %s vs %s should be symmetric for equal versions", tt.v1, tt.v2)
			} else if got > 0 && reverse >= 0 {
				t.Errorf("Compare symmetry failed: if %s > %s, then %s should be < %s", tt.v1, tt.v2, tt.v2, tt.v1)
			} else if got < 0 && reverse <= 0 {
				t.Errorf("Compare symmetry failed: if %s < %s, then %s should be > %s", tt.v1, tt.v2, tt.v2, tt.v1)
			}
		})
	}
}
