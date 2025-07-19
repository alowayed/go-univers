package pypi

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
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
			name:  "version with many components (technically valid per PEP 440)",
			input: "1.2.3.4.5.6.7",
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
			e := &Ecosystem{}
			got, err := e.NewVersion(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if !versionsEqual(got, tt.want) {
					t.Errorf("Ecosystem.NewVersion() = %+v, want %+v", got, tt.want)
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
			v1 := mustNewVersion(t, tt.v1)
			v2 := mustNewVersion(t, tt.v2)

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

// mustNewVersion is a helper function to create a new Version.
func mustNewVersion(t *testing.T, version string) *Version {
	t.Helper()
	e := &Ecosystem{}
	v, err := e.NewVersion(version)
	if err != nil {
		t.Fatalf("Failed to create version %s: %v", version, err)
	}
	return v
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
