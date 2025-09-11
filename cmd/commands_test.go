package main

import (
	"testing"

	"github.com/alowayed/go-univers/pkg/ecosystem/golang"
	"github.com/alowayed/go-univers/pkg/ecosystem/maven"
	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
	"github.com/alowayed/go-univers/pkg/univers"
)

type compareTest struct {
	name    string
	args    []string
	wantOut int
	wantErr bool
}

type sortTest struct {
	name    string
	args    []string
	wantOut []string
	wantErr bool
}

type containsTest struct {
	name    string
	args    []string
	wantOut bool
	wantErr bool
}

func testCompare[V univers.Version[V], VR univers.VersionRange[V]](
	t *testing.T,
	ecosystem univers.Ecosystem[V, VR],
	tests []compareTest,
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := compare(ecosystem, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("compare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantOut {
				t.Errorf("compare() = %v, want %v", got, tt.wantOut)
			}
		})
	}
}

func testSort[V univers.Version[V], VR univers.VersionRange[V]](
	t *testing.T,
	ecosystem univers.Ecosystem[V, VR],
	tests []sortTest,
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sort(ecosystem, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("sort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.wantOut) {
				t.Errorf("sort() length = %v, want %v", len(got), len(tt.wantOut))
				return
			}
			for i, v := range got {
				if v != tt.wantOut[i] {
					t.Errorf("sort() = %v, want %v", got, tt.wantOut)
					break
				}
			}
		})
	}
}

func testContains[V univers.Version[V], VR univers.VersionRange[V]](
	t *testing.T,
	ecosystem univers.Ecosystem[V, VR],
	tests []containsTest,
) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := contains(ecosystem, tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("contains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantOut {
				t.Errorf("contains() = %v, want %v", got, tt.wantOut)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	// Common test cases
	basicTests := []compareTest{
		{
			name:    "less than",
			args:    []string{"1.0.0", "2.0.0"},
			wantOut: -1,
			wantErr: false,
		},
		{
			name:    "greater than",
			args:    []string{"2.0.0", "1.0.0"},
			wantOut: 1,
			wantErr: false,
		},
		{
			name:    "equal",
			args:    []string{"1.0.0", "1.0.0"},
			wantOut: 0,
			wantErr: false,
		},
		{
			name:    "too few args",
			args:    []string{"1.0.0"},
			wantOut: 0,
			wantErr: true,
		},
		{
			name:    "too many args",
			args:    []string{"1.0.0", "2.0.0", "3.0.0"},
			wantOut: 0,
			wantErr: true,
		},
		{
			name:    "invalid first version",
			args:    []string{"invalid", "2.0.0"},
			wantOut: 0,
			wantErr: true,
		},
		{
			name:    "invalid second version",
			args:    []string{"1.0.0", "invalid"},
			wantOut: 0,
			wantErr: true,
		},
	}

	// NPM-specific tests
	npmTests := append(basicTests, []compareTest{
		{
			name:    "npm prerelease comparison",
			args:    []string{"1.0.0-alpha", "1.0.0-beta"},
			wantOut: -1,
			wantErr: false,
		},
		{
			name:    "npm prerelease vs release",
			args:    []string{"1.0.0", "1.0.0-alpha"},
			wantOut: 1,
			wantErr: false,
		},
		{
			name:    "npm with build metadata",
			args:    []string{"1.0.0+build.1", "1.0.0+build.2"},
			wantOut: 0,
			wantErr: false,
		},
	}...)

	// PyPI-specific tests
	pypiTests := append(basicTests, []compareTest{
		{
			name:    "pypi with epochs",
			args:    []string{"1!1.0.0", "1.0.0"},
			wantOut: 1,
			wantErr: false,
		},
	}...)

	// Go-specific tests
	goTests := []compareTest{
		{
			name:    "go less than",
			args:    []string{"v1.0.0", "v2.0.0"},
			wantOut: -1,
			wantErr: false,
		},
		{
			name:    "go greater than",
			args:    []string{"v2.0.0", "v1.0.0"},
			wantOut: 1,
			wantErr: false,
		},
		{
			name:    "go equal",
			args:    []string{"v1.0.0", "v1.0.0"},
			wantOut: 0,
			wantErr: false,
		},
		{
			name:    "go pseudo version",
			args:    []string{"v1.0.0-20170915032832-14c0d48ead0c", "v1.0.0"},
			wantOut: -1,
			wantErr: false,
		},
		{
			name:    "go too few args",
			args:    []string{"v1.0.0"},
			wantOut: 0,
			wantErr: true,
		},
		{
			name:    "go invalid first version",
			args:    []string{"invalid", "v2.0.0"},
			wantOut: 0,
			wantErr: true,
		},
	}

	t.Run("npm", func(t *testing.T) {
		testCompare(t, &npm.Ecosystem{}, npmTests)
	})

	t.Run("pypi", func(t *testing.T) {
		testCompare(t, &pypi.Ecosystem{}, pypiTests)
	})

	// Maven-specific tests
	mavenTests := append(basicTests, []compareTest{
		{
			name:    "maven qualifier comparison alpha vs beta",
			args:    []string{"1.0.0-alpha", "1.0.0-beta"},
			wantOut: -1,
			wantErr: false,
		},
		{
			name:    "maven qualifier vs release",
			args:    []string{"1.0.0-snapshot", "1.0.0"},
			wantOut: -1,
			wantErr: false,
		},
		{
			name:    "maven release vs service pack",
			args:    []string{"1.0.0", "1.0.0-sp"},
			wantOut: -1,
			wantErr: false,
		},
		{
			name:    "maven ga equivalent to release",
			args:    []string{"1.0.0-ga", "1.0.0"},
			wantOut: 0,
			wantErr: false,
		},
		{
			name:    "maven qualifier shortcuts",
			args:    []string{"1.0.0-a", "1.0.0-alpha"},
			wantOut: 0,
			wantErr: false,
		},
	}...)

	t.Run("go", func(t *testing.T) {
		testCompare(t, &golang.Ecosystem{}, goTests)
	})

	t.Run("maven", func(t *testing.T) {
		testCompare(t, &maven.Ecosystem{}, mavenTests)
	})
}

func TestSort(t *testing.T) {
	// Common test cases
	basicTests := []sortTest{
		{
			name:    "sort ascending",
			args:    []string{"2.0.0", "1.0.0", "1.5.0"},
			wantOut: []string{"1.0.0", "1.5.0", "2.0.0"},
			wantErr: false,
		},
		{
			name:    "sort single version",
			args:    []string{"1.0.0"},
			wantOut: []string{"1.0.0"},
			wantErr: false,
		},
		{
			name:    "sort already sorted",
			args:    []string{"1.0.0", "1.5.0", "2.0.0"},
			wantOut: []string{"1.0.0", "1.5.0", "2.0.0"},
			wantErr: false,
		},
		{
			name:    "sort identical versions",
			args:    []string{"1.0.0", "1.0.0", "1.0.0"},
			wantOut: []string{"1.0.0", "1.0.0", "1.0.0"},
			wantErr: false,
		},
		{
			name:    "sort no args",
			args:    []string{},
			wantOut: nil,
			wantErr: true,
		},
		{
			name:    "sort invalid version",
			args:    []string{"1.0.0", "invalid"},
			wantOut: nil,
			wantErr: true,
		},
	}

	// NPM-specific tests
	npmTests := append(basicTests, []sortTest{
		{
			name:    "npm sort with prerelease",
			args:    []string{"1.0.0", "1.0.0-alpha", "1.0.0-beta"},
			wantOut: []string{"1.0.0-alpha", "1.0.0-beta", "1.0.0"},
			wantErr: false,
		},
		{
			name:    "npm sort complex prerelease",
			args:    []string{"1.0.0", "1.0.0-rc.1", "1.0.0-alpha"},
			wantOut: []string{"1.0.0-alpha", "1.0.0-rc.1", "1.0.0"},
			wantErr: false,
		},
	}...)

	// Go-specific tests
	goTests := []sortTest{
		{
			name:    "go sort",
			args:    []string{"v2.0.0", "v1.0.0"},
			wantOut: []string{"v1.0.0", "v2.0.0"},
			wantErr: false,
		},
		{
			name:    "go sort single version",
			args:    []string{"v1.0.0"},
			wantOut: []string{"v1.0.0"},
			wantErr: false,
		},
		{
			name:    "go sort no args",
			args:    []string{},
			wantOut: nil,
			wantErr: true,
		},
	}

	t.Run("npm", func(t *testing.T) {
		testSort(t, &npm.Ecosystem{}, npmTests)
	})

	t.Run("pypi", func(t *testing.T) {
		testSort(t, &pypi.Ecosystem{}, basicTests)
	})

	// Maven-specific tests
	mavenTests := append(basicTests, []sortTest{
		{
			name:    "maven sort with qualifiers",
			args:    []string{"1.0.0", "1.0.0-alpha", "1.0.0-beta", "1.0.0-snapshot"},
			wantOut: []string{"1.0.0-alpha", "1.0.0-beta", "1.0.0-snapshot", "1.0.0"},
			wantErr: false,
		},
		{
			name:    "maven sort complex qualifiers",
			args:    []string{"1.0.0-sp", "1.0.0", "1.0.0-rc", "1.0.0-alpha"},
			wantOut: []string{"1.0.0-alpha", "1.0.0-rc", "1.0.0", "1.0.0-sp"},
			wantErr: false,
		},
		{
			name:    "maven sort with normalization",
			args:    []string{"1.0.0-ga", "1.0.0-final", "1.0.0"},
			wantOut: []string{"1.0.0-ga", "1.0.0-final", "1.0.0"},
			wantErr: false,
		},
	}...)

	t.Run("go", func(t *testing.T) {
		testSort(t, &golang.Ecosystem{}, goTests)
	})

	t.Run("maven", func(t *testing.T) {
		testSort(t, &maven.Ecosystem{}, mavenTests)
	})
}

func TestContains(t *testing.T) {
	// NPM tests
	npmTests := []containsTest{
		{
			name:    "npm caret range true",
			args:    []string{"^1.0.0", "1.5.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "npm caret range false",
			args:    []string{"^1.0.0", "2.0.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "npm tilde range true",
			args:    []string{"~1.2.0", "1.2.5"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "npm tilde range false",
			args:    []string{"~1.2.0", "1.3.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "npm exact version true",
			args:    []string{"1.0.0", "1.0.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "npm exact version false",
			args:    []string{"1.0.0", "1.0.1"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "npm x-range true",
			args:    []string{"1.x", "1.5.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "npm x-range false",
			args:    []string{"1.x", "2.0.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "npm too few args",
			args:    []string{"^1.0.0"},
			wantOut: false,
			wantErr: true,
		},
		{
			name:    "npm too many args",
			args:    []string{"^1.0.0", "1.5.0", "extra"},
			wantOut: false,
			wantErr: true,
		},
		{
			name:    "npm invalid range",
			args:    []string{"invalid", "1.0.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "npm invalid version",
			args:    []string{"^1.0.0", "invalid"},
			wantOut: false,
			wantErr: true,
		},
	}

	// PyPI tests
	pypiTests := []containsTest{
		{
			name:    "pypi contains",
			args:    []string{">=1.0.0", "1.5.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "pypi too few args",
			args:    []string{">=1.0.0"},
			wantOut: false,
			wantErr: true,
		},
		{
			name:    "pypi invalid range",
			args:    []string{"invalid", "1.0.0"},
			wantOut: false,
			wantErr: false,
		},
	}

	// Go tests
	goTests := []containsTest{
		{
			name:    "go contains",
			args:    []string{">=v1.0.0", "v1.5.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "go too few args",
			args:    []string{">=v1.0.0"},
			wantOut: false,
			wantErr: true,
		},
		{
			name:    "go invalid version",
			args:    []string{">=v1.0.0", "invalid"},
			wantOut: false,
			wantErr: true,
		},
	}

	t.Run("npm", func(t *testing.T) {
		testContains(t, &npm.Ecosystem{}, npmTests)
	})

	t.Run("pypi", func(t *testing.T) {
		testContains(t, &pypi.Ecosystem{}, pypiTests)
	})

	// Maven tests
	mavenTests := []containsTest{
		{
			name:    "maven exact range true",
			args:    []string{"[1.0.0]", "1.0.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "maven exact range false",
			args:    []string{"[1.0.0]", "1.0.1"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "maven inclusive range true",
			args:    []string{"[1.0.0,2.0.0]", "1.5.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "maven inclusive range false",
			args:    []string{"[1.0.0,2.0.0]", "2.5.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "maven lower bound true",
			args:    []string{"[1.0.0,)", "2.0.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "maven lower bound false",
			args:    []string{"[1.0.0,)", "0.5.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "maven upper bound true",
			args:    []string{"(,2.0.0]", "1.0.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "maven upper bound false",
			args:    []string{"(,2.0.0]", "3.0.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "maven exclusive range true",
			args:    []string{"(1.0.0,2.0.0)", "1.5.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "maven exclusive range false bounds",
			args:    []string{"(1.0.0,2.0.0)", "1.0.0"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "maven with qualifiers",
			args:    []string{"[1.0.0-alpha,1.0.0]", "1.0.0-beta"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "maven simple version true",
			args:    []string{"1.0.0", "1.0.0"},
			wantOut: true,
			wantErr: false,
		},
		{
			name:    "maven simple version false",
			args:    []string{"1.0.0", "1.0.1"},
			wantOut: false,
			wantErr: false,
		},
		{
			name:    "maven too few args",
			args:    []string{"[1.0.0,2.0.0]"},
			wantOut: false,
			wantErr: true,
		},
		{
			name:    "maven too many args",
			args:    []string{"[1.0.0,2.0.0]", "1.5.0", "extra"},
			wantOut: false,
			wantErr: true,
		},
		{
			name:    "maven invalid range",
			args:    []string{"invalid", "1.0.0"},
			wantOut: false,
			wantErr: true,
		},
		{
			name:    "maven invalid version",
			args:    []string{"[1.0.0,2.0.0]", "invalid"},
			wantOut: false,
			wantErr: true,
		},
	}

	t.Run("go", func(t *testing.T) {
		testContains(t, &golang.Ecosystem{}, goTests)
	})

	t.Run("maven", func(t *testing.T) {
		testContains(t, &maven.Ecosystem{}, mavenTests)
	})
}
