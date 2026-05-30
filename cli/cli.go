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
	case "init":
		if len(args) < 2 {
			fmt.Println("Use: perkbox init <remote-url>")
			fmt.Println("  or: perkbox init local")
			os.Exit(1)
		}
		cmdInit(args[1])
	case "push":
		cmdPush()
	case "pull":
		cmdPull()
	case "edit":
		if len(args) < 3 {
			fmt.Println("Use: perkbox edit <service> <username>")
			os.Exit(1)
		}
		cmdEdit(args[1], args[2], args)
	case "auto-push":
		if len(args) < 2 || (args[1] != "on" && args[1] != "off") {
			fmt.Println("Use: perkbox auto-push on")
			fmt.Println("  or: perkbox auto-push off")
			os.Exit(1)
		}
		cmdAutoPush(args[1] == "on")
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

		entries, err := store.LoadAll()
		if err != nil {
			fmt.Println("Error loading data:", err)
			os.Exit(1)
		}

		entries = append(entries, storage.Entry{
			Service:  service,
			Username: username,
			Password: encrypted,
		})
		if err := store.SaveAll(entries); err != nil {
			fmt.Println("Error saving:", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error verifying Master Password, please try again...")
		os.Exit(1)
	}
}

func cmdInit(url string) {
	if url == "local" {
		url = ""
	}
	if err := storage.InitGitRepo(url); err != nil {
		fmt.Println("Error initializing git repo:", err)
		os.Exit(1)
	}
	store = storage.GetStorage()
	fmt.Println("Git repository initialized in ~/.perkbox/")
}

func cmdPush() {
	if err := store.Push(); err != nil {
		fmt.Println("Error pushing:", err)
		os.Exit(1)
	}
	fmt.Println("Pushed successfully")
}

func cmdPull() {
	if err := store.Pull(); err != nil {
		fmt.Println("Error pulling:", err)
		os.Exit(1)
	}
	fmt.Println("Pulled successfully")
}

func cmdAutoPush(enabled bool) {
	gs, ok := store.(*storage.GitStorage)
	if !ok {
		fmt.Println("Auto-push is only available in git mode (run 'perkbox init' first)")
		return
	}
	if err := gs.SetAutoPush(enabled); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	if enabled {
		fmt.Println("Auto-push enabled")
	} else {
		fmt.Println("Auto-push disabled")
	}
}

func cmdEdit(service, username string, args []string) {
	masterPwd := readPassword("Master password: ")

	entries, err := store.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	idx := -1
	for i, e := range entries {
		if e.Service == service && e.Username == username {
			idx = i
			break
		}
	}

	if idx == -1 {
		fmt.Printf("Entry not found: %s (%s)\n", service, username)
		return
	}

	currentPwd, err := crypto.Decrypt(entries[idx].Password, masterPwd)
	if err != nil {
		fmt.Println("Error: Wrong Master Password")
		return
	}

	fmt.Printf("\nEditing: %s (%s)\n", service, username)
	fmt.Print("New username (enter to keep): ")
	var newUser string
	fmt.Scanln(&newUser)
	if newUser == "" {
		newUser = entries[idx].Username
	}

	var newPwd string
	if len(args) >= 4 && args[3] == "-gen" {
		var passLen, specialCount int
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
		newPwd = charGen(passLen, specialCount)
	} else {
		fmt.Print("New password (enter to keep): ")
		newPwd = readPassword("")
	}
	if newPwd == "" {
		newPwd = currentPwd
	}

	encrypted, err := crypto.Encrypt(newPwd, masterPwd)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		os.Exit(1)
	}

	entries[idx] = storage.Entry{
		Service:  service,
		Username: newUser,
		Password: encrypted,
	}

	if err := store.SaveAll(entries); err != nil {
		fmt.Println("Error saving:", err)
		os.Exit(1)
	}
	fmt.Println("Entry updated")
}

func cmdGet(service, username string) {
	masterPwd := readPassword("Master password: ")

	entries, err := store.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var found []storage.Entry
	for _, e := range entries {
		if e.Service == service && e.Username == username {
			found = append(found, e)
		}
	}

	if len(found) == 0 {
		fmt.Printf("Entry not found: %s (%s)\n", service, username)
		return
	}

	for _, e := range found {
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
	entries, err := store.LoadAll()
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
	entries, err := store.LoadAll()
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

	if err := store.SaveAll(filtered); err != nil {
		fmt.Println("Error saving:", err)
		return
	}
	fmt.Printf("%s deleted, (%d passwords)\n", service, len(toDelete))
}
