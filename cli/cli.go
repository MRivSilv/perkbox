package cli

import (
	"fmt"
	"golang.org/x/term"
	"os"
	"perkbox/storage"
)

var store storage.Storage

func Run(args []string) {
	store = storage.GetStorage()

	switch args[0] {
	case "add":
		cmdAdd(args)
	case "get":
		if len(args) < 3 {
			fmt.Println("Use: perkbox get <service> <username>")
			os.Exit(1)
		}
		cmdGet(args[1], args[2])
	case "list":
		cmdList()
	case "delete":
		if len(args) < 3 {
			fmt.Println("Use: perkbox delete <service> <username>")
			os.Exit(1)
		}
		cmdDelete(args[1], args[2])
	case "edit":
		if len(args) < 3 {
			fmt.Println("Use: perkbox edit <service> <username>")
			os.Exit(1)
		}
		cmdEdit(args[1], args[2], args)
	case "set2fa":
		if len(args) > 1 {
			fmt.Println("Use: perkbox set2fa")
			os.Exit(1)
		}
		fmt.Printf("Paste your 2fa auth code: ")
		var code string
		fmt.Scanln(&code)

		if code == "" {
			fmt.Println("Paste a valid code")
			os.Exit(1)
		}
		set2FA(code)
	case "get2fa":
		if len(args) < 3 {
			fmt.Println("Use: perkbox get2fa <service> <username>")
			os.Exit(1)
		}
		get2fa(args[1], args[2])
	default:
		fmt.Printf("Unkown command: %s\n", args[0])
		os.Exit(1)
	}
}

func readPassword(prompt string) string {
	fmt.Print(prompt)
	password, _ := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return string(password)
}
