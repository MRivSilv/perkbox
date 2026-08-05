package cli

import (
	"fmt"
	"time"

	"perkbox/crypto"

	"github.com/atotto/clipboard"
	"github.com/pquerna/otp/totp"
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
	entr, err := store.FindService(entries, service, username, masterPwd)
	if err != nil {
		fmt.Println("Service not found")
		return
	}
	encrypted, err := crypto.Encrypt(code, masterPwd)
	if err != nil {
		fmt.Println("Error encrypting 2FA key:", err)
		return
	}
	entr.TwoFAKey = encrypted
	if err := store.SaveAll(entries); err != nil {
		fmt.Println("Error saving:", err)
		return
	}
	fmt.Printf("2FA key saved for %s (%s)\n", service, username)
}

func get2fa(service string, username string) {
	masterPwd := readPassword("Master password: ")
	fmt.Printf("2FA code for: %s(%s)\n", service, username)
	entries, err := store.LoadAll()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	entr, err := store.FindService(entries, service, username, masterPwd)
	if err != nil {
		fmt.Println("Service not found")
		return
	}
	decoded, err := crypto.Decrypt(entr.TwoFAKey, masterPwd)
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

}
