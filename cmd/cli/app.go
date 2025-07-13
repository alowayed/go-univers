package cli

import (
	"fmt"
)

// Run is the main entry point for the CLI
func Run(args []string) int {
	if len(args) == 0 {
		fmt.Println("Usage: univers <ecosystem> <command> [args]")
		return 1
	}

	switch args[0] {
	case "npm":
		return runEcosystem("npm", args[1:])
	case "pypi":
		return runEcosystem("pypi", args[1:])
	default:
		fmt.Printf("Unknown ecosystem: %s\n", args[0])
		return 1
	}
}

// runEcosystem routes commands to the appropriate ecosystem handler
func runEcosystem(ecosystem string, args []string) int {
	if len(args) == 0 {
		fmt.Printf("No command specified for %s\n", ecosystem)
		return 1
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "compare":
		return compare(ecosystem, commandArgs)
	case "sort":
		return sort(ecosystem, commandArgs)
	case "satisfies":
		return satisfies(ecosystem, commandArgs)
	default:
		fmt.Printf("Unknown %s command: %s\n", ecosystem, command)
		return 1
	}
}