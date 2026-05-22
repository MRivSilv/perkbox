package main

import (
	"fmt"
	"os"
	"time"

	"perkbox/cli"

	"github.com/atotto/clipboard"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: perkbox <comando>")
		fmt.Println("Comandos: add, get, list, delete")
		os.Exit(1)
	}
	cli.Run(os.Args[1:])
	fmt.Printf("\nTienes 10 segundos para pegar tu contraseña\n")
	time.Sleep(10 * time.Second)
	clipboard.WriteAll("")
}
