package main

import (
	"fmt"
	"ro/donna"
	"ro/weems"
)

func main() {
	donna.Parse()

	// Handle logging level flags; verbose overrides quiet by design
	if donna.HasGlobalFlag("quiet") {
		weems.SetQuiet()
	}
	if donna.HasGlobalFlag("verbose") {
		weems.SetVerbose()
	}

	weems.Info("Ro's method dispatching is beginning.")

	method, ok := donna.NextArg()
	if !ok {
		weems.Critical("Expected an argument.")
	}

	switch method {
	case "version":
		version()
	case "help":
		help()
	default:
		weems.Critical("Unexpected argument '%s'.", method)
	}
}

func version() {
	fmt.Printf("I don't know my version yet.\n")
}

func help() {
	fmt.Printf("I can't help you yet.\n")
}
