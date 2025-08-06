package ecosystem

import (
	"github.com/alowayed/go-univers/pkg/ecosystem/alpine"
	"github.com/alowayed/go-univers/pkg/ecosystem/cargo"
	"github.com/alowayed/go-univers/pkg/ecosystem/composer"
	"github.com/alowayed/go-univers/pkg/ecosystem/gem"
	"github.com/alowayed/go-univers/pkg/ecosystem/gomod"
	"github.com/alowayed/go-univers/pkg/ecosystem/maven"
	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
	"github.com/alowayed/go-univers/pkg/ecosystem/nuget"
	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
	"github.com/alowayed/go-univers/pkg/univers"
)

var (

	// --- Ensure types implement interfaces (Alphabetical) ---

	// alpine
	_ univers.Version[*alpine.Version]                          = &alpine.Version{}
	_ univers.VersionRange[*alpine.Version]                     = &alpine.VersionRange{}
	_ univers.Ecosystem[*alpine.Version, *alpine.VersionRange] = &alpine.Ecosystem{}

	// cargo
	_ univers.Version[*cargo.Version]                        = &cargo.Version{}
	_ univers.VersionRange[*cargo.Version]                   = &cargo.VersionRange{}
	_ univers.Ecosystem[*cargo.Version, *cargo.VersionRange] = &cargo.Ecosystem{}

	// composer
	_ univers.Version[*composer.Version]                            = &composer.Version{}
	_ univers.VersionRange[*composer.Version]                       = &composer.VersionRange{}
	_ univers.Ecosystem[*composer.Version, *composer.VersionRange] = &composer.Ecosystem{}

	// gem
	_ univers.Version[*gem.Version]                      = &gem.Version{}
	_ univers.VersionRange[*gem.Version]                 = &gem.VersionRange{}
	_ univers.Ecosystem[*gem.Version, *gem.VersionRange] = &gem.Ecosystem{}

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

	// nuget
	_ univers.Version[*nuget.Version]                        = &nuget.Version{}
	_ univers.VersionRange[*nuget.Version]                   = &nuget.VersionRange{}
	_ univers.Ecosystem[*nuget.Version, *nuget.VersionRange] = &nuget.Ecosystem{}

	// pypi
	_ univers.Version[*pypi.Version]                       = &pypi.Version{}
	_ univers.VersionRange[*pypi.Version]                  = &pypi.VersionRange{}
	_ univers.Ecosystem[*pypi.Version, *pypi.VersionRange] = &pypi.Ecosystem{}
)
