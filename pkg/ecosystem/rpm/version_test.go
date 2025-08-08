package rpm

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    *Version
		wantErr bool
	}{
		// Basic version parsing
		{
			name:    "simple version",
			version: "1.2.3",
			want: &Version{
				epoch:    0,
				version:  "1.2.3",
				release:  "",
				original: "1.2.3",
			},
		},
		{
			name:    "version with epoch",
			version: "2:1.2.3",
			want: &Version{
				epoch:    2,
				version:  "1.2.3",
				release:  "",
				original: "2:1.2.3",
			},
		},
		{
			name:    "version with release",
			version: "1.2.3-4",
			want: &Version{
				epoch:    0,
				version:  "1.2.3",
				release:  "4",
				original: "1.2.3-4",
			},
		},
		{
			name:    "full version with epoch and release",
			version: "1:2.3.4-5",
			want: &Version{
				epoch:    1,
				version:  "2.3.4",
				release:  "5",
				original: "1:2.3.4-5",
			},
		},
		{
			name:    "zero epoch",
			version: "0:1.2.3",
			want: &Version{
				epoch:    0,
				version:  "1.2.3",
				release:  "",
				original: "0:1.2.3",
			},
		},
		{
			name:    "complex version with letters and dots",
			version: "1.2.3a.b4",
			want: &Version{
				epoch:    0,
				version:  "1.2.3a.b4",
				release:  "",
				original: "1.2.3a.b4",
			},
		},
		{
			name:    "version with tilde",
			version: "1.2~beta1",
			want: &Version{
				epoch:    0,
				version:  "1.2~beta1",
				release:  "",
				original: "1.2~beta1",
			},
		},
		{
			name:    "version with plus",
			version: "1.2+git20220101",
			want: &Version{
				epoch:    0,
				version:  "1.2+git20220101",
				release:  "",
				original: "1.2+git20220101",
			},
		},
		{
			name:    "version with caret",
			version: "1.2^snapshot",
			want: &Version{
				epoch:    0,
				version:  "1.2^snapshot",
				release:  "",
				original: "1.2^snapshot",
			},
		},
		{
			name:    "complex release with special chars",
			version: "1.2.3-el8.x86_64",
			want: &Version{
				epoch:    0,
				version:  "1.2.3",
				release:  "el8.x86_64",
				original: "1.2.3-el8.x86_64",
			},
		},
		{
			name:    "release with hyphens",
			version: "1.2.3-1.el8-1",
			want: &Version{
				epoch:    0,
				version:  "1.2.3-1.el8",
				release:  "1",
				original: "1.2.3-1.el8-1",
			},
		},

		// Error cases
		{
			name:    "empty version",
			version: "",
			wantErr: true,
		},
		{
			name:    "whitespace only",
			version: "   ",
			wantErr: true,
		},
		{
			name:    "negative epoch",
			version: "-1:1.2.3",
			wantErr: true,
		},
		{
			name:    "invalid epoch",
			version: "abc:1.2.3",
			wantErr: true,
		},
		{
			name:    "missing version after epoch",
			version: "1:",
			wantErr: true,
		},
		{
			name:    "invalid character",
			version: "1.2.3@invalid",
			wantErr: true,
		},
		{
			name:    "invalid release character",
			version: "1.2.3-rel@ease",
			wantErr: true,
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ecosystem.NewVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ecosystem.NewVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got.epoch != tt.want.epoch {
				t.Errorf("Ecosystem.NewVersion().epoch = %v, want %v", got.epoch, tt.want.epoch)
			}
			if got.version != tt.want.version {
				t.Errorf("Ecosystem.NewVersion().version = %v, want %v", got.version, tt.want.version)
			}
			if got.release != tt.want.release {
				t.Errorf("Ecosystem.NewVersion().release = %v, want %v", got.release, tt.want.release)
			}
			if got.original != tt.want.original {
				t.Errorf("Ecosystem.NewVersion().original = %v, want %v", got.original, tt.want.original)
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
		// Epoch comparisons
		{
			name: "higher epoch wins",
			v1:   "2:1.0.0",
			v2:   "1:2.0.0",
			want: 1,
		},
		{
			name: "lower epoch loses",
			v1:   "1:2.0.0",
			v2:   "2:1.0.0",
			want: -1,
		},
		{
			name: "equal epochs, compare version",
			v1:   "1:1.2.3",
			v2:   "1:1.2.4",
			want: -1,
		},
		{
			name: "no epoch vs epoch 0",
			v1:   "1.2.3",
			v2:   "0:1.2.3",
			want: 0,
		},
		{
			name: "epoch 0 vs no epoch",
			v1:   "0:1.2.3",
			v2:   "1.2.3",
			want: 0,
		},

		// Version comparisons
		{
			name: "simple numeric comparison",
			v1:   "1.2.3",
			v2:   "1.2.4",
			want: -1,
		},
		{
			name: "major version difference",
			v1:   "2.0.0",
			v2:   "1.9.9",
			want: 1,
		},
		{
			name: "equal versions",
			v1:   "1.2.3",
			v2:   "1.2.3",
			want: 0,
		},
		{
			name: "leading zeros ignored",
			v1:   "1.02.3",
			v2:   "1.2.3",
			want: 0,
		},
		{
			name: "different length numeric parts",
			v1:   "1.10",
			v2:   "1.9",
			want: 1,
		},

		// Tilde handling (~)
		{
			name: "tilde sorts before release",
			v1:   "1.0~beta",
			v2:   "1.0",
			want: -1,
		},
		{
			name: "tilde sorts before letters",
			v1:   "1.0~beta",
			v2:   "1.0a",
			want: -1,
		},
		{
			name: "tilde in middle",
			v1:   "1.0~rc1",
			v2:   "1.0~rc2",
			want: -1,
		},

		// Release comparisons
		{
			name: "version with release vs without",
			v1:   "1.2.3",
			v2:   "1.2.3-1",
			want: -1,
		},
		{
			name: "release comparison",
			v1:   "1.2.3-1",
			v2:   "1.2.3-2",
			want: -1,
		},
		{
			name: "release numeric vs alpha",
			v1:   "1.2.3-1",
			v2:   "1.2.3-a",
			want: -1,
		},
		{
			name: "release el7 vs el8",
			v1:   "1.2.3-1.el7",
			v2:   "1.2.3-1.el8",
			want: -1,
		},

		// Complex cases
		{
			name: "complex version with letters",
			v1:   "2.4.37a",
			v2:   "2.4.37",
			want: 1,
		},
		{
			name: "version with plus",
			v1:   "1.2.3+git20220101",
			v2:   "1.2.3",
			want: 1,
		},
		{
			name: "version with caret",
			v1:   "1.2.3^snapshot",
			v2:   "1.2.3",
			want: 1,
		},

		// Edge cases from RPM documentation
		{
			name: "rpm doc example 1",
			v1:   "1.0010",
			v2:   "1.9",
			want: 1,
		},
		{
			name: "rpm doc example 2",
			v1:   "1.05",
			v2:   "1.5",
			want: 0,
		},
		{
			name: "rpm doc example 3",
			v1:   "1.0",
			v2:   "1.0.0",
			want: -1,
		},
		{
			name: "rpm doc example 4",
			v1:   "2.50",
			v2:   "2.5",
			want: 1, // 50 > 5
		},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err1 := ecosystem.NewVersion(tt.v1)
			v2, err2 := ecosystem.NewVersion(tt.v2)

			if err1 != nil || err2 != nil {
				t.Fatalf("Failed to parse versions: v1=%v, v2=%v", err1, err2)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v (comparing %s vs %s)", got, tt.want, tt.v1, tt.v2)
			}

			// Test symmetry
			reverse := v2.Compare(v1)
			if reverse != -tt.want {
				t.Errorf("Version.Compare() symmetry failed: %s.Compare(%s) = %v, but %s.Compare(%s) = %v, expected %v",
					tt.v1, tt.v2, got, tt.v2, tt.v1, reverse, -tt.want)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{"simple version", "1.2.3"},
		{"version with epoch", "2:1.2.3"},
		{"version with release", "1.2.3-4"},
		{"full version", "1:2.3.4-5"},
		{"complex version", "2:1.2.3-4.el8.x86_64"},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ecosystem.NewVersion(tt.version)
			if err != nil {
				t.Fatalf("NewVersion() error = %v", err)
			}

			if got := v.String(); got != tt.version {
				t.Errorf("Version.String() = %v, want %v", got, tt.version)
			}
		})
	}
}

// Test helper function for tilde comparison
func TestCompareRPMNonDigits(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"tilde vs empty", "~", "", -1},
		{"empty vs tilde", "", "~", 1},
		{"tilde vs alpha", "~beta", "alpha", -1},
		{"alpha vs tilde", "alpha", "~beta", 1},
		{"both tilde", "~alpha", "~beta", -1},
		{"regular strings", "alpha", "beta", -1},
		{"equal strings", "alpha", "alpha", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareRPMNonDigits(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("compareRPMNonDigits(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Test helper function for digit comparison
func TestCompareRPMDigits(t *testing.T) {
	tests := []struct {
		name string
		a, b string
		want int
	}{
		{"empty vs empty", "", "", 0},
		{"empty vs number", "", "1", -1},
		{"number vs empty", "1", "", 1},
		{"equal numbers", "42", "42", 0},
		{"leading zeros ignored", "042", "42", 0},
		{"different numbers", "9", "10", -1},
		{"large numbers", "999999999999999999999", "1000000000000000000000", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compareRPMDigits(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("compareRPMDigits(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// Comprehensive comparison test from real RPM packages
func TestRealWorldRPMVersions(t *testing.T) {
	tests := []struct {
		name string
		v1   string
		v2   string
		want int
	}{
		// Real-world examples from CentOS/RHEL packages
		{"kernel versions", "3.10.0-1062.el7", "3.10.0-1127.el7", -1},
		{"glibc versions", "2.17-292.el7", "2.17-307.el7", -1},
		{"httpd versions", "2.4.6-88.el7", "2.4.6-90.el7", -1},
		{"python versions", "2.7.5-86.el7", "2.7.5-89.el7", -1},
		{"different architectures", "1.2.3-1.x86_64", "1.2.3-1.i386", 1}, // x86_64 > i386 lexicographically

		// Beta/RC versions
		{"beta vs release", "2.0~beta1", "2.0", -1},
		{"rc vs release", "2.0~rc1", "2.0", -1},
		{"beta vs rc", "2.0~beta1", "2.0~rc1", -1},

		// Snapshot/git versions
		{"git snapshot", "1.0+git20220101", "1.0", 1},
		{"snapshot comparison", "1.0^snapshot20220101", "1.0^snapshot20220102", -1},
	}

	ecosystem := &Ecosystem{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err1 := ecosystem.NewVersion(tt.v1)
			v2, err2 := ecosystem.NewVersion(tt.v2)

			if err1 != nil {
				t.Fatalf("Failed to parse v1=%s: %v", tt.v1, err1)
			}
			if err2 != nil {
				t.Fatalf("Failed to parse v2=%s: %v", tt.v2, err2)
			}

			got := v1.Compare(v2)
			if got != tt.want {
				t.Errorf("Version.Compare() = %v, want %v (comparing %s vs %s)", got, tt.want, tt.v1, tt.v2)
			}
		})
	}
}

// Test sorting of versions
func TestVersionSorting(t *testing.T) {
	versions := []string{
		"1.0",
		"2:1.0",
		"1.0~rc1",
		"1.0-1",
		"1.0.1",
		"1:1.0",
		"1.0~beta1",
		"1.1",
		"1.0-2",
	}

	expected := []string{
		"1.0~beta1",
		"1.0~rc1",
		"1.0",
		"1.0-1",
		"1.0-2",
		"1.0.1",
		"1.1",
		"1:1.0",
		"2:1.0",
	}

	ecosystem := &Ecosystem{}

	// Parse all versions
	var parsedVersions []*Version
	for _, v := range versions {
		parsed, err := ecosystem.NewVersion(v)
		if err != nil {
			t.Fatalf("Failed to parse version %s: %v", v, err)
		}
		parsedVersions = append(parsedVersions, parsed)
	}

	// Sort them using a simple bubble sort to test Compare method
	for i := 0; i < len(parsedVersions); i++ {
		for j := i + 1; j < len(parsedVersions); j++ {
			if parsedVersions[i].Compare(parsedVersions[j]) > 0 {
				parsedVersions[i], parsedVersions[j] = parsedVersions[j], parsedVersions[i]
			}
		}
	}

	// Check the order
	for i, expected := range expected {
		if got := parsedVersions[i].String(); got != expected {
			t.Errorf("Sorted position %d: got %s, want %s", i, got, expected)
		}
	}
}
