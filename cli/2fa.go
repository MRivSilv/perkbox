package cli

import (
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/pquerna/otp/totp"
	"os"
	"perkbox/crypto"
	"time"
)

func set2FA(code string) {
	var service, username string

	fmt.Print("Service: ")
	fmt.Scanln(&service)

	fmt.Print("Username: ")
	fmt.Scanln(&username)

	masterPwd := readPassword("Master password: ")

	entries, err := store.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	for i, e := range entries {
		if e.Service == service && e.Username == username {
			encrypted, err := crypto.Encrypt(code, masterPwd)
			if err != nil {
				fmt.Println("Error encrypting 2FA key:", err)
				return
			}
			entries[i].TwoFAKey = encrypted
			if err := store.SaveAll(entries); err != nil {
				fmt.Println("Error saving:", err)
				return
			}
			fmt.Printf("2FA key saved for %s (%s)\n", service, username)
			return
		}
	}

	fmt.Println("Entry not found")
	os.Exit(1)
}

func get2fa(service string, username string) {
	masterPwd := readPassword("Master password: ")
	fmt.Printf("2FA code for: %s(%s)\n", service, username)
	entries, err := store.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for i, e := range entries {
		if e.Service == service && e.Username == username {
			decoded, err := crypto.Decrypt(entries[i].TwoFAKey, masterPwd)
			if err != nil {
				fmt.Println("Error decrypting 2FA key: ", err)
				return
			}
			totpCode, err := totp.GenerateCode(decoded, time.Now())
			if err != nil {
				fmt.Println("Error generating code: ", err)
				return
			}
			fmt.Println("Your totp is: ", totpCode)
			copiedCode := clipboard.WriteAll(totpCode)
			if copiedCode != nil {
				fmt.Println("Error copying your code")
				return
			}
			fmt.Println("Totp code copied to your clipboard")
			return
		}
	}
	fmt.Println("Service or User not found (Check list)")
}
