package cli

import (
	"fmt"
	"os"
	"syscall"

	"slices"
	"time"

	"github.com/atotto/clipboard"

	"github.com/MRivSilv/perkbox/crypto"
	"github.com/MRivSilv/perkbox/storage"

	"golang.org/x/term"
)

// Run despacha el comando correcto
func Run(args []string) {
	switch args[0] {
	case "add":
		cmdAdd(args)
	case "get":
		if len(args) < 2 {
			fmt.Println("Use: perkbox get <service>")
			os.Exit(1)
		}
		cmdGet(args[1])
	case "list":
		cmdList()
	case "delete":
		if len(args) < 3 {
			fmt.Println("Use: perkbox delete <service> <username>")
			os.Exit(1)
		}
		cmdDelete(args[1], args[2])

	default:
		fmt.Printf("Unkown command: %s\n", args[0])
		os.Exit(1)
	}
}

func readPassword(prompt string) string {
	fmt.Print(prompt)
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	return string(password)
}

func cmdAdd(args []string) {
	var service, username, password string
	var passLen, specialCount int

	actions := []string{"-gen", "-help"}

	if len(args) > 1 && !slices.Contains(actions, args[1]) {
		fmt.Printf("Argument  <%s> does not exist\n", args[1])
		return
	}

	fmt.Print("Service (ej: github.com): ")
	fmt.Scanln(&service)

	fmt.Print("User: ")
	fmt.Scanln(&username)

	if len(args) >= 2 && args[1] == "-gen" {

		fmt.Println("How many characters do you want your password?")
		fmt.Scanln(&passLen)
		if passLen == 0 {
			fmt.Println("Impossible to generate a password with 0 characters")
			return
		}
		fmt.Println("How many special characters?")
		fmt.Scanln(&specialCount)

		if specialCount > passLen {
			fmt.Println("Too short of a password for that many special characters")
		}

		password = charGen(passLen, specialCount)

	} else if len(args) < 2 {
		password = readPassword("Password: ")
	}
	masterPwd := readPassword("Master password: ")
	verification := verifyMasterPwd(masterPwd)
	if verification == true {
		encrypted, err := crypto.Encrypt(password, masterPwd)
		if err != nil {
			fmt.Println("Error encrypting:", err)
			os.Exit(1)
		}

		entries, err := storage.LoadAll()
		if err != nil {
			fmt.Println("Error loading data:", err)
			os.Exit(1)
		}

		entries = append(entries, storage.Entry{
			Service:  service,
			Username: username,
			Password: encrypted,
		})
		if err := storage.SaveAll(entries); err != nil {
			fmt.Println("Error saving:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error verifying Master Password, please try again...")
		os.Exit(1)
	}
}

func cmdGet(service string) {
	masterPwd := readPassword("Master password: ")

	entries, err := storage.FindByService(service)
	if err != nil || len(entries) == 0 {
		fmt.Printf(" %s not found\n", service)
		return
	}

	for _, e := range entries {
		pwd, err := crypto.Decrypt(e.Password, masterPwd)
		if err != nil {
			fmt.Println("Error: Wrong Master Password")
			return
		}
		fmt.Printf("\nService:  %s\nUser:   %s\n", e.Service, e.Username)
		copiedPass := clipboard.WriteAll(pwd)
		if copiedPass != nil {
			fmt.Println("Your password couldn't be copied")
			return
		}
		fmt.Println("Password copied to your clipboard\nYou have 10 seconds to use it")
		time.Sleep(10 * time.Second)
		clipboard.WriteAll("Timeout: Be Quicker")
	}
}

func cmdList() {
	entries, err := storage.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(entries) == 0 {
		fmt.Println("No saved passwords.")
		return
	}

	fmt.Println("\n=== Your services ===")
	for _, e := range entries {
		fmt.Printf("• %s (%s)\n", e.Service, e.Username)
	}
}

func cmdDelete(service string, user string) {
	entries, err := storage.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var toDelete []storage.Entry
	var filtered []storage.Entry
	for _, e := range entries {
		if e.Service == service && e.Username == user {
			toDelete = append(toDelete, e)
		} else {
			filtered = append(filtered, e)
		}
	}

	if len(toDelete) == 0 {
		fmt.Printf("Couldn't find: %s\n", service)
		return
	}

	masterPwd := readPassword("Master password: ")

	_, err = crypto.Decrypt(toDelete[0].Password, masterPwd)
	if err != nil {
		fmt.Println("Error: Wrong Master Password")
		return
	}

	if err := storage.SaveAll(filtered); err != nil {
		fmt.Println("Error saving:", err)
		return
	}
	fmt.Printf("%s deleted, (%d passwords)\n", service, len(toDelete))
}
