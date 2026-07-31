package migratescript

import (
	"encoding/json"
	"fmt"
	"os"
	"perkbox/crypto"
)

type oldEntry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password []byte `json:"password"`
	TwoFAKey []byte `json:"twofa_key,omitempty"`
}

type newEntry struct {
	Service  []byte `json:"service"`
	Username []byte `json:"username"`
	Password []byte `json:"password"`
	TwoFAKey []byte `json:"twofa_key,omitempty"`
}

func migrate() {
	home, _ := os.UserHomeDir()
	path := home + "/.perkbox/.perkbox.json"

	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var old []oldEntry
	if err := json.Unmarshal(data, &old); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	fmt.Print("Master password: ")
	var pwd string
	fmt.Scanln(&pwd)

	var migrated []newEntry
	for _, e := range old {
		svcEnc, err := crypto.Encrypt(e.Service, pwd)
		if err != nil {
			fmt.Println("Wrong master password")
			os.Exit(1)
		}
		usrEnc, err := crypto.Encrypt(e.Username, pwd)
		if err != nil {
			fmt.Println("Wrong master password")
			os.Exit(1)
		}
		migrated = append(migrated, newEntry{
			Service:  svcEnc,
			Username: usrEnc,
			Password: e.Password,
			TwoFAKey: e.TwoFAKey,
		})
	}

	out, _ := json.MarshalIndent(migrated, "", "  ")
	if err := os.WriteFile(path, out, 0600); err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	}
	fmt.Println("Migration complete")
}
