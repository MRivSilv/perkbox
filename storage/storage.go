package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Entry = password saveda
type Entry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password []byte `json:"password"` // Saved and encrypteda
}

// getStoragePath retrun the path to save the .perkbox.json file
func getStoragePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".perkbox.json")
}

// LoadAll loads every entry inside the .json file
func LoadAll() ([]Entry, error) {
	path := getStoragePath()

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Entry{}, nil // If it doesn't exist, an empty list is returned
	}
	if err != nil {
		return nil, err
	}

	var entries []Entry
	err = json.Unmarshal(data, &entries)
	return entries, err
}

// SaveAll save all the entries inside the .json file
func SaveAll(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(getStoragePath(), data, 0600)
	// 0600 = solo el dueño puede leer/escribir
}

// FindByService finds by service name
func FindByService(service string) ([]Entry, error) {
	all, err := LoadAll()
	if err != nil {
		return nil, err
	}

	var found []Entry
	for _, entry := range all {
		if entry.Service == service {
			found = append(found, entry)
		}
	}
	return found, nil
}
