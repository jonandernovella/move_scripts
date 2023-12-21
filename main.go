package main

import (
	"fmt"
	"os"
)

func main() {
	var lib = Lib{"darsync", os.Getenv("HOME"), os.Stdin}

	if len(os.Args) < 2 {
		fmt.Println("Usage: " + lib.Name + " [check|gen]")
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "check":
		lib.check()
	case "gen":
		lib.gen()
	default:
		fmt.Println("Unknown command:", command)
		os.Exit(1)
	}
}
