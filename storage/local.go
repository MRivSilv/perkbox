package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	path string
}

func NewLocalStorage() *LocalStorage {
	home, _ := os.UserHomeDir()
	return &LocalStorage{
		path: filepath.Join(home, ".perkbox.json"),
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

func (s *LocalStorage) FindByService(service string) ([]Entry, error) {
	all, err := s.LoadAll()
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

func (s *LocalStorage) Push() error { return nil }
func (s *LocalStorage) Pull() error { return nil }
