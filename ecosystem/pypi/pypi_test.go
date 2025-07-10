package pypi

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Version
		wantErr bool
	}{
		{
			name:  "basic release version",
			input: "1.2.3",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3},
				postrelease: -1,
				dev:         -1,
				original:    "1.2.3",
			},
			wantErr: false,
		},
		{
			name:  "version with epoch",
			input: "2!1.2.3",
			want: &Version{
				epoch:       2,
				release:     []int{1, 2, 3},
				postrelease: -1,
				dev:         -1,
				original:    "2!1.2.3",
			},
			wantErr: false,
		},
		{
			name:  "version with alpha prerelease",
			input: "1.2.3a1",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3},
				prerelease:  "a",
				preNumber:   1,
				postrelease: -1,
				dev:         -1,
				original:    "1.2.3a1",
			},
			wantErr: false,
		},
		{
			name:  "version with beta prerelease",
			input: "1.2.3b2",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3},
				prerelease:  "b",
				preNumber:   2,
				postrelease: -1,
				dev:         -1,
				original:    "1.2.3b2",
			},
			wantErr: false,
		},
		{
			name:  "version with rc prerelease",
			input: "1.2.3rc1",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3},
				prerelease:  "rc",
				preNumber:   1,
				postrelease: -1,
				dev:         -1,
				original:    "1.2.3rc1",
			},
			wantErr: false,
		},
		{
			name:  "version with post-release",
			input: "1.2.3.post1",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3},
				postrelease: 1,
				dev:         -1,
				original:    "1.2.3.post1",
			},
			wantErr: false,
		},
		{
			name:  "version with dev release",
			input: "1.2.3.dev1",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3},
				postrelease: -1,
				dev:         1,
				original:    "1.2.3.dev1",
			},
			wantErr: false,
		},
		{
			name:  "version with local identifier",
			input: "1.2.3+local.1",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3},
				postrelease: -1,
				dev:         -1,
				local:       "local.1",
				original:    "1.2.3+local.1",
			},
			wantErr: false,
		},
		{
			name:  "complex version with all components",
			input: "2!1.2.3a1.post1.dev1+local.1",
			want: &Version{
				epoch:       2,
				release:     []int{1, 2, 3},
				prerelease:  "a",
				preNumber:   1,
				postrelease: 1,
				dev:         1,
				local:       "local.1",
				original:    "2!1.2.3a1.post1.dev1+local.1",
			},
			wantErr: false,
		},
		{
			name:  "two-component version",
			input: "1.2",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2},
				postrelease: -1,
				dev:         -1,
				original:    "1.2",
			},
			wantErr: false,
		},
		{
			name:  "single-component version",
			input: "1",
			want: &Version{
				epoch:       0,
				release:     []int{1},
				postrelease: -1,
				dev:         -1,
				original:    "1",
			},
			wantErr: false,
		},
		{
			name:    "empty version",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "version with many components (technically valid per PEP 440)",
			input:   "1.2.3.4.5.6.7",
			want: &Version{
				epoch:       0,
				release:     []int{1, 2, 3, 4, 5, 6, 7},
				postrelease: -1,
				dev:         -1,
				original:    "1.2.3.4.5.6.7",
			},
			wantErr: false,
		},
		{
			name:    "non-numeric release",
			input:   "a.b.c",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !versionsEqual(got, tt.want) {
					t.Errorf("NewVersion() = %+v, want %+v", got, tt.want)
				}
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		{
			name: "equal versions",
			v1:   "1.2.3",
			v2:   "1.2.3",
			want: 0,
		},
		{
			name: "different epochs",
			v1:   "1!1.2.3",
			v2:   "2!1.2.3",
			want: -1,
		},
		{
			name: "major version difference",
			v1:   "1.2.3",
			v2:   "2.2.3",
			want: -1,
		},
		{
			name: "minor version difference",
			v1:   "1.3.3",
			v2:   "1.2.3",
			want: 1,
		},
		{
			name: "patch version difference",
			v1:   "1.2.2",
			v2:   "1.2.3",
			want: -1,
		},
		{
			name: "prerelease vs release",
			v1:   "1.2.3a1",
			v2:   "1.2.3",
			want: -1,
		},
		{
			name: "alpha vs beta",
			v1:   "1.2.3a1",
			v2:   "1.2.3b1",
			want: -1,
		},
		{
			name: "beta vs rc",
			v1:   "1.2.3b1",
			v2:   "1.2.3rc1",
			want: -1,
		},
		{
			name: "prerelease number comparison",
			v1:   "1.2.3a1",
			v2:   "1.2.3a2",
			want: -1,
		},
		{
			name: "post-release comparison",
			v1:   "1.2.3",
			v2:   "1.2.3.post1",
			want: -1,
		},
		{
			name: "dev release comparison",
			v1:   "1.2.3.dev1",
			v2:   "1.2.3",
			want: -1,
		},
		{
			name: "complex comparison",
			v1:   "1.2.3a1.post1.dev1",
			v2:   "1.2.3a2",
			want: -1,
		},
		{
			name: "different release lengths",
			v1:   "1.2",
			v2:   "1.2.0",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse v1: %v", err)
			}
			v2, err := NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse v2: %v", err)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVersion_Normalize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "basic version",
			input: "1.2.3",
			want:  "1.2.3",
		},
		{
			name:  "version with epoch",
			input: "2!1.2.3",
			want:  "2!1.2.3",
		},
		{
			name:  "version with prerelease",
			input: "1.2.3a1",
			want:  "1.2.3a1",
		},
		{
			name:  "version with post-release",
			input: "1.2.3.post1",
			want:  "1.2.3post1",
		},
		{
			name:  "version with dev release",
			input: "1.2.3.dev1",
			want:  "1.2.3dev1",
		},
		{
			name:  "version with local",
			input: "1.2.3+local.1",
			want:  "1.2.3+local.1",
		},
		{
			name:  "complex version",
			input: "2!1.2.3a1.post1.dev1+local.1",
			want:  "2!1.2.3a1post1dev1+local.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersion(tt.input)
			if err != nil {
				t.Fatalf("Failed to parse version: %v", err)
			}

			got := v.Normalize()
			if got != tt.want {
				t.Errorf("Version.Normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewVersionRange(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantErr      bool
		wantOriginal string
	}{
		{
			name:         "exact version",
			input:        "1.2.3",
			wantErr:      false,
			wantOriginal: "1.2.3",
		},
		{
			name:         "greater than",
			input:        ">1.2.3",
			wantErr:      false,
			wantOriginal: ">1.2.3",
		},
		{
			name:         "compatible release",
			input:        "~=1.2.3",
			wantErr:      false,
			wantOriginal: "~=1.2.3",
		},
		{
			name:         "wildcard constraint",
			input:        "==1.2.*",
			wantErr:      false,
			wantOriginal: "==1.2.*",
		},
		{
			name:         "multiple constraints",
			input:        ">=1.0.0, <2.0.0",
			wantErr:      false,
			wantOriginal: ">=1.0.0, <2.0.0",
		},
		{
			name:         "arbitrary equality",
			input:        "===1.2.3",
			wantErr:      false,
			wantOriginal: "===1.2.3",
		},
		{
			name:    "empty specifier",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewVersionRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewVersionRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.String() != tt.wantOriginal {
				t.Errorf("NewVersionRange().String() = %v, want %v", got.String(), tt.wantOriginal)
			}
		})
	}
}

func TestVersionRange_Contains(t *testing.T) {
	tests := []struct {
		name        string
		rangeStr    string
		version     string
		want        bool
		wantRangeErr bool
		wantVersionErr bool
	}{
		{
			name:     "exact match",
			rangeStr: "==1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "exact no match",
			rangeStr: "==1.2.3",
			version:  "1.2.4",
			want:     false,
		},
		{
			name:     "greater than match",
			rangeStr: ">1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than no match",
			rangeStr: ">1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "greater than or equal match",
			rangeStr: ">=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "less than match",
			rangeStr: "<1.2.3",
			version:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than or equal match",
			rangeStr: "<=1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "not equal match",
			rangeStr: "!=1.2.3",
			version:  "1.2.4",
			want:     true,
		},
		{
			name:     "not equal no match",
			rangeStr: "!=1.2.3",
			version:  "1.2.3",
			want:     false,
		},
		{
			name:     "compatible release match",
			rangeStr: "~=1.2.3",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "compatible release no match",
			rangeStr: "~=1.2.3",
			version:  "1.3.0",
			want:     false,
		},
		{
			name:     "wildcard constraint match",
			rangeStr: "==1.2.*",
			version:  "1.2.5",
			want:     true,
		},
		{
			name:     "wildcard constraint no match",
			rangeStr: "==1.2.*",
			version:  "1.3.0",
			want:     false,
		},
		{
			name:     "multiple constraints match",
			rangeStr: ">=1.0.0, <2.0.0",
			version:  "1.5.0",
			want:     true,
		},
		{
			name:     "multiple constraints no match",
			rangeStr: ">=1.0.0, <2.0.0",
			version:  "2.0.0",
			want:     false,
		},
		{
			name:     "arbitrary equality match",
			rangeStr: "===1.2.3",
			version:  "1.2.3",
			want:     true,
		},
		{
			name:     "arbitrary equality no match",
			rangeStr: "===1.2.3",
			version:  "1.2.3.post1",
			want:     false,
		},
		{
			name:     "prerelease handling",
			rangeStr: ">=1.0.0a1",
			version:  "1.0.0b1",
			want:     true,
		},
		{
			name:     "epoch handling",
			rangeStr: ">=1!1.0.0",
			version:  "1!1.0.1",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vr, err := NewVersionRange(tt.rangeStr)
			if (err != nil) != tt.wantRangeErr {
				t.Errorf("NewVersionRange() error = %v, wantRangeErr %v", err, tt.wantRangeErr)
				return
			}
			if tt.wantRangeErr {
				return
			}

			v, err := NewVersion(tt.version)
			if (err != nil) != tt.wantVersionErr {
				t.Errorf("NewVersion() error = %v, wantVersionErr %v", err, tt.wantVersionErr)
				return
			}
			if tt.wantVersionErr {
				return
			}

			got := vr.Contains(v)
			if got != tt.want {
				t.Errorf("VersionRange.Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstraint_Matches(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		constVer string
		testVer  string
		want     bool
	}{
		{
			name:     "equals match",
			operator: "==",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "equals no match",
			operator: "==",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     false,
		},
		{
			name:     "not equals match",
			operator: "!=",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than match",
			operator: ">",
			constVer: "1.2.3",
			testVer:  "1.2.4",
			want:     true,
		},
		{
			name:     "greater than no match",
			operator: ">",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     false,
		},
		{
			name:     "greater than or equal match",
			operator: ">=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "less than match",
			operator: "<",
			constVer: "1.2.3",
			testVer:  "1.2.2",
			want:     true,
		},
		{
			name:     "less than or equal match",
			operator: "<=",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "arbitrary equality match",
			operator: "===",
			constVer: "1.2.3",
			testVer:  "1.2.3",
			want:     true,
		},
		{
			name:     "arbitrary equality no match",
			operator: "===",
			constVer: "1.2.3",
			testVer:  "1.2.3.post1",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Constraint{
				operator: tt.operator,
				version:  tt.constVer,
			}

			v, err := NewVersion(tt.testVer)
			if err != nil {
				t.Fatalf("Failed to parse test version: %v", err)
			}

			got := c.Matches(v)
			if got != tt.want {
				t.Errorf("Constraint.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare versions for testing
func versionsEqual(a, b *Version) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.epoch != b.epoch || a.prerelease != b.prerelease ||
		a.preNumber != b.preNumber || a.postrelease != b.postrelease ||
		a.dev != b.dev || a.local != b.local || a.original != b.original {
		return false
	}

	if len(a.release) != len(b.release) {
		return false
	}

	for i := range a.release {
		if a.release[i] != b.release[i] {
			return false
		}
	}

	return true
}