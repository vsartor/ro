package donna

import (
	"os"
	"ro/weems"
	"strings"
)

func ParseGlobalFlag(arg string, idx *int) {
	flagName := strings.TrimPrefix(arg, "--")
	globalFlags = append(globalFlags, flagName)
	*idx++
}

func ParseFlag(arg string, idx *int) {
	flagName := strings.TrimPrefix(arg, "--")
	flags = append(flags, flagName)
	*idx++
}

func ParseGlobalOption(arg string, idx *int) {
	optionName := strings.TrimPrefix(arg, "-")

	// Assert that option received an associated value
	if *idx == len(os.Args)-1 {
		weems.Fatal("Option '%s' has no associated value.", optionName)
	}

	// Global options are invalid
	optionValue := os.Args[*idx+1]
	weems.Fatal("Invalid global option '%s' with value '%s'.", optionName, optionValue)

	*idx += 2
}

func ParseOption(arg string, idx *int) {
	optionName := strings.TrimPrefix(arg, "-")

	// Assert that option received an associated value
	if *idx == len(os.Args)-1 {
		weems.Fatal("Option '%s' has no associated value.", optionName)
	}

	optionValue := os.Args[*idx+1]
	options[optionName] = optionValue

	*idx += 2
}

// Parses command line arguments, initializing Donna's variables
// storing options and arguments.
func Parse() {
	// Parse global values
	idx := 1
	for idx < len(os.Args) {
		arg := os.Args[idx]

		if strings.HasPrefix(arg, "--") {
			ParseGlobalFlag(arg, &idx)
		} else if strings.HasPrefix(arg, "-") {
			ParseGlobalOption(arg, &idx)
		} else {
			// Found an argument; stop global parsing.
			break
		}
	}

	// Parse arguments
	for idx < len(os.Args) {
		arg := os.Args[idx]

		if strings.HasPrefix(arg, "-") {
			// Found a flag/option; stop argument parsing.
			break
		} else {
			args.arguments = append(args.arguments, arg)
			idx++
		}
	}

	// Parse flags/options
	for idx < len(os.Args) {
		arg := os.Args[idx]

		if strings.HasPrefix(arg, "--") {
			ParseFlag(arg, &idx)
		} else if strings.HasPrefix(arg, "-") {
			ParseOption(arg, &idx)
		} else {
			// Should not have arguments at this point.
			weems.Fatal("Unexpected argument '%s' after flags/options.", arg)
		}
	}
}
