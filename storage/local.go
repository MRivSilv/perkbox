package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"perkbox/crypto"
)

type LocalStorage struct {
	path string
}

func NewLocalStorage() *LocalStorage {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".perkbox")
	return &LocalStorage{
		path: filepath.Join(dir, ".perkbox.json"),
	}
}

func (s *LocalStorage) LoadAll() ([]Entry, error) {
	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries []Entry
	err = json.Unmarshal(data, &entries)
	return entries, err
}

func (s *LocalStorage) SaveAll(entries []Entry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0600)
}

func (s *LocalStorage) FindByService(service, masterPassword string) ([]Entry, error) {
	all, err := s.LoadAll()
	if err != nil {
		return nil, err
	}
	var found []Entry
	encrypted, err := crypto.Encrypt(service, masterPassword)
	if err != nil {
		fmt.Println("Invalid Master Password")
		os.Exit(1)
	}
	for _, entry := range all {
		if bytes.Equal(entry.Service, encrypted) {
			found = append(found, entry)
		}
	}
	return found, nil
}
