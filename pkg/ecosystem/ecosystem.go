package ecosystem

import (
	"github.com/alowayed/go-univers/pkg/ecosystem/gomod"
	"github.com/alowayed/go-univers/pkg/ecosystem/maven"
	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
	"github.com/alowayed/go-univers/pkg/univers"
)

var (

	// --- Ensure types implement interfaces (Alphabetical) ---

	// go
	_ univers.Version[*gomod.Version]                        = &gomod.Version{}
	_ univers.VersionRange[*gomod.Version]                   = &gomod.VersionRange{}
	_ univers.Ecosystem[*gomod.Version, *gomod.VersionRange] = &gomod.Ecosystem{}

	// maven
	_ univers.Version[*maven.Version]                        = &maven.Version{}
	_ univers.VersionRange[*maven.Version]                   = &maven.VersionRange{}
	_ univers.Ecosystem[*maven.Version, *maven.VersionRange] = &maven.Ecosystem{}

	// npm
	_ univers.Version[*npm.Version]                      = &npm.Version{}
	_ univers.VersionRange[*npm.Version]                 = &npm.VersionRange{}
	_ univers.Ecosystem[*npm.Version, *npm.VersionRange] = &npm.Ecosystem{}

	// pypi
	_ univers.Version[*pypi.Version]                       = &pypi.Version{}
	_ univers.VersionRange[*pypi.Version]                  = &pypi.VersionRange{}
	_ univers.Ecosystem[*pypi.Version, *pypi.VersionRange] = &pypi.Ecosystem{}
)
