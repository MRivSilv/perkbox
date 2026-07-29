# Perkbox

Perkbox is a local, console-based password manager written in Go. It stores encrypted credentials on your machine and provides a minimal CLI to add, retrieve, list, and delete entries.

## Features
- AES-256-GCM encryption protected by a master password
- Built-in secure password generator
- Clipboard copy with automatic clearing after 10 seconds
- Service + username support
- TOTP 2FA support

## Security model
- The master password is hashed with SHA-256 to derive the encryption key.
- Data is stored with `0600` permissions.
- No sync, recovery, or audit guarantees. Use at your own risk.

## Requirements
- Go toolchain (see `go.mod`)

## Build
```bash
go build
```

## Installation

### Arch Linux (AUR)
```bash
yay -S perkbox
```

### npm
```bash
npm install -g perkbox
```

### Manual install

Linux / macOS:
```bash
go build -o perkbox .
sudo install -Dm755 perkbox /usr/local/bin/perkbox
```

Windows:
```powershell
go build -o perkbox.exe .
# Move perkbox.exe anywhere in your PATH
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

### set2fa
Stores a TOTP secret key for an existing entry.
```bash
./perkbox set2fa
```

### get2fa `<service> <username>`
Generates a TOTP code for the given entry and copies it to the clipboard.
```bash
./perkbox get2fa github.com myuser
```

## Data location
- `~/.perkbox.json`

## Contributing
Issues and pull requests are welcome.
