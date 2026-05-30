# Perkbox

Perkbox is a local, console-based password manager written in Go. It stores encrypted credentials on your machine and provides a minimal CLI to add, retrieve, list, and delete entries.

## Features
- Local OR version-controlled storage via git
- AES-256-GCM encryption protected by a master password
- Built-in secure password generator
- Clipboard copy with automatic clearing after 10 seconds
- Service + username support
- Git auto-commit on every change (when using git mode)
- Optional auto-push to remote after every change

## Security model
- The master password is hashed with SHA-256 to derive the encryption key.
- Data is stored with `0600` permissions.
- When using git, the JSON vault is still encrypted — your remote never sees plaintext.
- No sync, recovery, or audit guarantees. Use at your own risk.

## Requirements
- Go toolchain (see `go.mod`)
- `git` (only required for git mode)

## Build
```bash
go build
```

## Installation

### Arch Linux (AUR)
```bash
yay -S perkbox
```

### Manual install
```bash
go build -o perkbox .
sudo install -Dm755 perkbox /usr/local/bin/perkbox
```

## Usage
```bash
./perkbox <command> [arguments...]
./perkbox -h
./perkbox --help
```

## Commands

### help / -h / --help
Shows the full list of commands with descriptions. Same as running `perkbox` with no arguments.
```bash
./perkbox help
./perkbox -h
```

### add
Add a new entry. If called with `-gen`, it will ask for length and special character count and generate a secure random password.
```bash
./perkbox add
./perkbox add -gen
```

### get `<service> <username>`
Decrypts and copies the password to the clipboard (cleared after 10s).
```bash
./perkbox get github.com myuser
```

### edit `<service> <username>`
Edit an existing entry. You can change the username, the password, or both. Leave a field empty to keep the current value. Use `-gen` as a 4th argument to generate a new password.
```bash
./perkbox edit github.com myuser
./perkbox edit github.com myuser -gen
```

### delete `<service> <username>`
Removes the matching entry after verifying the master password.
```bash
./perkbox delete github.com myuser
```

### list
Lists all saved services and usernames.
```bash
./perkbox list
```

### init `<remote-url>`
Initializes a git repository at `~/.perkbox/` for version-controlled storage. If a remote URL is provided, it sets the origin. Existing passwords from `~/.perkbox.json` are automatically migrated.
```bash
./perkbox init local
./perkbox init https://github.com/user/vault.git
```

### push
Push committed changes to the git remote.
```bash
./perkbox push
```

### pull
Pull latest changes from the git remote.
```bash
./perkbox pull
```

### auto-push on / off
When enabled, every `add` or `delete` automatically runs `push` after committing.
```bash
./perkbox auto-push on
./perkbox auto-push off
```

## Git setup step-by-step

### Prerequisites
1. A GitHub/GitLab account
2. `git` installed on your system

### 1. Create a remote repository
Go to GitHub and create a new **empty** repository (no README, no .gitignore, no license).

### 2. Initialize Perkbox with your remote
```bash
./perkbox init https://github.com/your-user/vault.git
```

### 3. Push your vault to the remote
```bash
./perkbox push
```
The first time, git will ask for your GitHub username and password. Use a **personal access token** instead of your password:
- Create one at https://github.com/settings/tokens (scopes: `repo`)
- Paste it when prompted for a password

### 4. (Optional) Save credentials so you don't type them every time
```bash
git config --global credential.helper store
```
The next time you type your username and token, they will be saved. The token is stored in plaintext in `~/.git-credentials`.

### 5. Enable auto-push (optional)
```bash
./perkbox auto-push on
```
Now every `add` or `delete` will automatically commit and push to your remote.

### 6. Sync between machines
On another machine:
```bash
git clone https://github.com/your-user/vault.git ~/.perkbox
```
Then use `perkbox` normally. Run `perkbox pull` before and `perkbox push` after making changes.

## Data location
- **Local mode**: `~/.perkbox.json`
- **Git mode**: `~/.perkbox/perkbox.json`

## Contributing
Issues and pull requests are welcome.
