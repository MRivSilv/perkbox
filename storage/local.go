package storage

import (
	"encoding/json"
	"errors"
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

func (s *LocalStorage) FindService(all []Entry, service, username, masterPassword string) (*Entry, error) {
	for i := range all {
		dSrvc, dUsr := crypto.DecryptOutput(all[i].Service, all[i].Username, masterPassword)
		if dSrvc == service && dUsr == username {
			return &all[i], nil
		}
	}
	return nil, errors.New("Service not found")
}
