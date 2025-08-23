package cli

import (
	"fmt"
	"strings"

	"github.com/alowayed/go-univers/pkg/ecosystem/alpine"
	"github.com/alowayed/go-univers/pkg/ecosystem/cargo"
	"github.com/alowayed/go-univers/pkg/ecosystem/composer"
	"github.com/alowayed/go-univers/pkg/ecosystem/conan"
	"github.com/alowayed/go-univers/pkg/ecosystem/cran"
	"github.com/alowayed/go-univers/pkg/ecosystem/debian"
	"github.com/alowayed/go-univers/pkg/ecosystem/gem"
	"github.com/alowayed/go-univers/pkg/ecosystem/gentoo"
	"github.com/alowayed/go-univers/pkg/ecosystem/gomod"
	"github.com/alowayed/go-univers/pkg/ecosystem/maven"
	"github.com/alowayed/go-univers/pkg/ecosystem/npm"
	"github.com/alowayed/go-univers/pkg/ecosystem/nuget"
	"github.com/alowayed/go-univers/pkg/ecosystem/pypi"
	"github.com/alowayed/go-univers/pkg/ecosystem/rpm"
	"github.com/alowayed/go-univers/pkg/ecosystem/semver"
	"github.com/alowayed/go-univers/pkg/univers"
)

// Run is the main entry point for the CLI
func Run(args []string) int {
	out, code := run(args)
	fmt.Printf("%s\n", out)
	return code
}

func run(args []string) (string, int) {
	if len(args) == 0 {
		s := "Usage: univers <ecosystem|spec> <command> [args]"
		return s, 1
	}

	// Handle spec commands first
	specToRun := map[string]func([]string) (string, int){
		"vers": runSpec,
	}

	if fn, ok := specToRun[args[0]]; ok {
		return fn(args[1:])
	}

	ecosystemToRun := map[string]func([]string) (string, int){
		alpine.Name: func(args []string) (string, int) {
			return runEcosystem(&alpine.Ecosystem{}, args)
		},
		cargo.Name: func(args []string) (string, int) {
			return runEcosystem(&cargo.Ecosystem{}, args)
		},
		conan.Name: func(args []string) (string, int) {
			return runEcosystem(&conan.Ecosystem{}, args)
		},
		composer.Name: func(args []string) (string, int) {
			return runEcosystem(&composer.Ecosystem{}, args)
		},
		cran.Name: func(args []string) (string, int) {
			return runEcosystem(&cran.Ecosystem{}, args)
		},
		debian.Name: func(args []string) (string, int) {
			return runEcosystem(&debian.Ecosystem{}, args)
		},
		gem.Name: func(args []string) (string, int) {
			return runEcosystem(&gem.Ecosystem{}, args)
		},
		gentoo.Name: func(args []string) (string, int) {
			return runEcosystem(&gentoo.Ecosystem{}, args)
		},
		gomod.Name: func(args []string) (string, int) {
			return runEcosystem(&gomod.Ecosystem{}, args)
		},
		maven.Name: func(args []string) (string, int) {
			return runEcosystem(&maven.Ecosystem{}, args)
		},
		npm.Name: func(args []string) (string, int) {
			return runEcosystem(&npm.Ecosystem{}, args)
		},
		nuget.Name: func(args []string) (string, int) {
			return runEcosystem(&nuget.Ecosystem{}, args)
		},
		pypi.Name: func(args []string) (string, int) {
			return runEcosystem(&pypi.Ecosystem{}, args)
		},
		rpm.Name: func(args []string) (string, int) {
			return runEcosystem(&rpm.Ecosystem{}, args)
		},
		semver.Name: func(args []string) (string, int) {
			return runEcosystem(&semver.Ecosystem{}, args)
		},
	}

	if fn, ok := ecosystemToRun[args[0]]; ok {
		return fn(args[1:])
	}

	s := fmt.Sprintf("Unknown ecosystem: %s", args[0])
	return s, 1
}

func runEcosystem[V univers.Version[V], VR univers.VersionRange[V]](
	e univers.Ecosystem[V, VR],
	args []string,
) (string, int) {
	if len(args) == 0 {
		s := fmt.Sprintf("No command specified for %s", e.Name())
		return s, 1
	}

	command := args[0]
	commandArgs := args[1:]

	var result string
	var err error
	switch command {
	case "compare":
		var out int
		out, err = compare(e, commandArgs)
		result = fmt.Sprintf("%d", out)
	case "sort":
		var out []string
		out, err = sort(e, commandArgs)
		for _, v := range out {
			result += fmt.Sprintf("%q ", v)
		}
		result = strings.TrimSpace(result)
	case "contains":
		var out bool
		out, err = contains(e, commandArgs)
		result = fmt.Sprintf("%t", out)
	default:
		s := fmt.Sprintf("Unknown %s command: %s", e.Name(), command)
		return s, 1
	}

	if err != nil {
		s := fmt.Sprintf("Error running command '%s': %v", command, err)
		return s, 1
	}

	return result, 0
}

// runSpec handles spec-specific commands
func runSpec(args []string) (string, int) {
	if len(args) == 0 {
		return "No command specified for spec", 1
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "contains":
		out, err := versContains(commandArgs)
		if err != nil {
			return fmt.Sprintf("Error running command '%s': %v", command, err), 1
		}
		return fmt.Sprintf("%t", out), 0
	default:
		return fmt.Sprintf("Unknown spec command: %s", command), 1
	}
}
