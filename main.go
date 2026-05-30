package main

import (
	"fmt"
	"os"

	"github.com/MRivSilv/perkbox/cli"
)

func main() {
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
		printHelp()
		os.Exit(0)
	}
	cli.Run(os.Args[1:])
}

func printHelp() {
	fmt.Println("perkbox - local console-based password manager")
	fmt.Println()
	fmt.Println("Usage: perkbox <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add [-gen]              Add a new entry (use -gen to auto-generate password)")
	fmt.Println("  get    <service> <user>  Copy password to clipboard")
	fmt.Println("  edit   <service> <user> Edit an existing entry")
	fmt.Println("  delete <service> <user> Delete an entry")
	fmt.Println("  list                    List all services")
	fmt.Println("  init   <url>            Initialize git repo (url or 'local')")
	fmt.Println("  push                    Push to git remote")
	fmt.Println("  pull                    Pull from git remote")
	fmt.Println("  auto-push on|off        Toggle automatic push after changes")
}
