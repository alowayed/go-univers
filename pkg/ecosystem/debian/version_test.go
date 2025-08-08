package debian

import (
	"testing"
)

func TestEcosystem_NewVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		// Basic versions
		{
			name:    "basic version",
			version: "1.0",
			wantErr: false,
		},
		{
			name:    "multi-component version",
			version: "2.1.3",
			wantErr: false,
		},

		// Versions with epochs
		{
			name:    "version with epoch",
			version: "1:1.0",
			wantErr: false,
		},
		{
			name:    "version with high epoch",
			version: "2:3.4.5",
			wantErr: false,
		},

		// Versions with revisions
		{
			name:    "version with simple revision",
			version: "1.0-1",
			wantErr: false,
		},
		{
			name:    "version with Ubuntu revision",
			version: "2.1.3-4ubuntu1",
			wantErr: false,
		},

		// Versions with epochs and revisions
		{
			name:    "version with epoch and revision",
			version: "1:1.0-1",
			wantErr: false,
		},
		{
			name:    "complex version with epoch and revision",
			version: "2:3.4.5-6ubuntu2",
			wantErr: false,
		},

		// Complex upstream versions
		{
			name:    "version with dfsg suffix",
			version: "1.2.3+dfsg",
			wantErr: false,
		},
		{
			name:    "version with tilde",
			version: "1.0~beta1",
			wantErr: false,
		},
		{
			name:    "version with really pattern",
			version: "1.0+really.1.0",
			wantErr: false,
		},

		// Complex revisions
		{
			name:    "version with security update",
			version: "1.0-1+deb9u1",
			wantErr: false,
		},
		{
			name:    "version with dotted revision",
			version: "1.0-1.2",
			wantErr: false,
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
			name:    "non-digit start",
			version: "abc",
			wantErr: true,
		},
		{
			name:    "non-digit start with epoch",
			version: "1:abc",
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
			if !tt.wantErr && got.String() != tt.version {
				t.Errorf("Ecosystem.NewVersion().String() = %v, want %v", got.String(), tt.version)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	ecosystem := &Ecosystem{}

	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		// Basic comparisons
		{
			name:     "equal versions",
			v1:       "1.0",
			v2:       "1.0",
			expected: 0,
		},
		{
			name:     "simple numeric comparison",
			v1:       "1.0",
			v2:       "2.0",
			expected: -1,
		},
		{
			name:     "reverse numeric comparison",
			v1:       "2.0",
			v2:       "1.0",
			expected: 1,
		},

		// Epoch comparisons
		{
			name:     "epoch overrides version",
			v1:       "1:1.0",
			v2:       "2.0",
			expected: 1,
		},
		{
			name:     "explicit vs implicit epoch",
			v1:       "0:1.0",
			v2:       "1.0",
			expected: 0,
		},
		{
			name:     "same epoch compare upstream",
			v1:       "1:1.0",
			v2:       "1:2.0",
			expected: -1,
		},
		{
			name:     "higher epoch wins",
			v1:       "2:1.0",
			v2:       "1:2.0",
			expected: 1,
		},

		// Revision comparisons
		{
			name:     "revision comparison",
			v1:       "1.0-1",
			v2:       "1.0-2",
			expected: -1,
		},
		{
			name:     "reverse revision comparison",
			v1:       "1.0-2",
			v2:       "1.0-1",
			expected: 1,
		},
		{
			name:     "native vs revision",
			v1:       "1.0",
			v2:       "1.0-1",
			expected: -1,
		},
		{
			name:     "revision vs native",
			v1:       "1.0-1",
			v2:       "1.0",
			expected: 1,
		},

		// Complex version strings
		{
			name:     "tilde sorts before release",
			v1:       "1.0~beta",
			v2:       "1.0",
			expected: -1,
		},
		{
			name:     "lexical comparison with tilde",
			v1:       "1.0~alpha",
			v2:       "1.0~beta",
			expected: -1,
		},
		{
			name:     "plus sorts after nothing",
			v1:       "1.0+dfsg",
			v2:       "1.0",
			expected: 1,
		},

		// Debian-specific patterns
		{
			name:     "ubuntu revisions",
			v1:       "1.0-1ubuntu1",
			v2:       "1.0-1ubuntu2",
			expected: -1,
		},
		{
			name:     "ubuntu vs plain revision",
			v1:       "1.0-1ubuntu1",
			v2:       "1.0-1",
			expected: 1,
		},
		{
			name:     "really pattern",
			v1:       "1.0+really.1.0",
			v2:       "1.0+really.2.0",
			expected: -1,
		},

		// Edge cases with tildes
		{
			name:     "multiple tildes",
			v1:       "1.0~~",
			v2:       "1.0~",
			expected: -1,
		},
		{
			name:     "tilde with numbers",
			v1:       "1.0~1",
			v2:       "1.0~2",
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err := ecosystem.NewVersion(tt.v1)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.v1, err)
			}

			v2, err := ecosystem.NewVersion(tt.v2)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt.v2, err)
			}

			result := v1.Compare(v2)
			if result != tt.expected {
				t.Errorf("Version.Compare(%s, %s) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVersion_String(t *testing.T) {
	ecosystem := &Ecosystem{}

	tests := []string{
		"1.0",
		"1:2.3.4",
		"5.6.7-8",
		"9:10.11.12-13ubuntu14",
		"1.0~beta+dfsg-1",
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			v, err := ecosystem.NewVersion(tt)
			if err != nil {
				t.Fatalf("Failed to parse version %s: %v", tt, err)
			}

			if v.String() != tt {
				t.Errorf("Version.String() = %v, want %v", v.String(), tt)
			}
		})
	}
}
