package storage

import (
	"os"
	"path/filepath"
)

type Entry struct {
	Service  string `json:"service"`
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type Storage interface {
	LoadAll() ([]Entry, error)
	SaveAll(entries []Entry) error
	FindByService(service string) ([]Entry, error)
	Push() error
	Pull() error
}

func GetStorage() Storage {
	home, _ := os.UserHomeDir()
	gitDir := filepath.Join(home, ".perkbox", ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return NewGitStorage()
	}
	return NewLocalStorage()
}
