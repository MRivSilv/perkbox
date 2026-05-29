package main

import (
	"fmt"
	"os"

	"github.com/MRivSilv/perkbox/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: perkbox <comando>")
		fmt.Println("Comandos: add, get, edit, list, delete, init, push, pull, auto-push")
		os.Exit(1)
	}
	cli.Run(os.Args[1:])
}
