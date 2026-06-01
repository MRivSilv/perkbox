package storage

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GitStorage struct {
	local *LocalStorage
	dir   string
}

func NewGitStorage() *GitStorage {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".perkbox")
	path := filepath.Join(dir, "perkbox.json")
	return &GitStorage{
		local: &LocalStorage{path: path},
		dir:   dir,
	}
}

func (s *GitStorage) LoadAll() ([]Entry, error) {
	return s.local.LoadAll()
}

func (s *GitStorage) SaveAll(entries []Entry) error {
	if err := s.local.SaveAll(entries); err != nil {
		return err
	}
	return s.commit("update vault")
}

func (s *GitStorage) FindByService(service string) ([]Entry, error) {
	return s.local.FindByService(service)
}

func (s *GitStorage) Push() error {
	branch, err := s.currentBranch()
	if err != nil {
		return err
	}
	return s.gitCmd("push", "-u", "origin", branch)
}

func (s *GitStorage) Pull() error {
	branch, err := s.currentBranch()
	if err != nil {
		return err
	}
	gitDir := filepath.Join(s.dir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository, run 'perkbox init' first")
	}

	localEntries, _ := s.local.LoadAll()

	if err := s.gitCmd("fetch", "origin"); err != nil {
		return fmt.Errorf("failed to fetch from remote: %w", err)
	}

	remoteRef := fmt.Sprintf("origin/%s", branch)
	if err := s.gitCmdQuiet("rev-parse", "--verify", remoteRef); err != nil {
		fmt.Println("Nothing to pull (no remote branch)")
		return nil
	}

	if err := s.gitCmd("reset", "--hard", remoteRef); err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}

	remoteEntries, _ := s.local.LoadAll()

	var merged []Entry
	seen := make(map[string]bool)
	for _, e := range remoteEntries {
		key := e.Service + "\x00" + e.Username
		seen[key] = true
		merged = append(merged, e)
	}
	for _, e := range localEntries {
		key := e.Service + "\x00" + e.Username
		if !seen[key] {
			merged = append(merged, e)
		}
	}

	if len(merged) != len(remoteEntries) {
		if err := s.local.SaveAll(merged); err != nil {
			return fmt.Errorf("failed to save merged entries: %w", err)
		}
		return s.commit("merge: pull from remote")
	}
	fmt.Println("Already up to date")
	return nil
}

func (s *GitStorage) commit(msg string) error {
	for _, args := range [][]string{
		{"add", "-A"},
		{"commit", "-m", msg},
	} {
		if err := s.gitCmd(args...); err != nil {
			return err
		}
	}
	if s.autoPushEnabled() {
		fmt.Println("Auto-push enabled, pushing...")
		return s.Push()
	}
	return nil
}

func (s *GitStorage) autoPushEnabled() bool {
	cmd := exec.Command("git", "-C", s.dir, "config", "--get", "perkbox.autoPush")
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

func (s *GitStorage) SetAutoPush(enabled bool) error {
	if enabled {
		return s.gitCmdQuiet("config", "perkbox.autoPush", "true")
	}
	return s.gitCmdQuiet("config", "--unset", "perkbox.autoPush")
}

func (s *GitStorage) gitCmdQuiet(arg ...string) error {
	args := append([]string{"-C", s.dir}, arg...)
	cmd := exec.Command("git", args...)
	return cmd.Run()
}

func (s *GitStorage) currentBranch() (string, error) {
	cmd := exec.Command("git", "-C", s.dir, "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (s *GitStorage) gitCmd(arg ...string) error {
	args := append([]string{"-C", s.dir}, arg...)
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func InitGitRepo(remoteURL string) error {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, ".perkbox")
	path := filepath.Join(dir, "perkbox.json")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		oldPath := filepath.Join(home, ".perkbox.json")
		if _, err := os.Stat(oldPath); err == nil {
			data, _ := os.ReadFile(oldPath)
			os.WriteFile(path, data, 0600)
			fmt.Println("Migrated passwords from ~/.perkbox.json")
		} else {
			os.WriteFile(path, []byte("[]"), 0600)
		}
	}

	if err := runGit(dir, "init"); err != nil {
		return err
	}

	if remoteURL != "" {
		runGit(dir, "remote", "add", "origin", remoteURL)

		if err := runGit(dir, "fetch", "origin"); err == nil {
			runGit(dir, "checkout", "-b", "main", "origin/main")
			fmt.Println("Cloned remote vault")
			return nil
		}
		fmt.Println("Remote not reachable or empty, starting local vault")
	}

	runGit(dir, "checkout", "-b", "main")
	runGit(dir, "add", "-A")
	runGit(dir, "commit", "-m", "initial vault")

	return nil
}

func runGit(dir string, arg ...string) error {
	args := append([]string{"-C", dir}, arg...)
	cmd := exec.Command("git", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
