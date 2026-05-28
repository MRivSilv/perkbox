# Perkbox

Perkbox is a local, console-based password manager written in Go. It stores encrypted credentials on your machine and provides a minimal CLI to add, retrieve, list, and delete entries.

## Features
- Local-only storage in a single JSON file
- AES-256-GCM encryption protected by a master password
- Clipboard copy with automatic clearing after 10 seconds
- Service + username support

## Security model
- The master password is hashed with SHA-256 to derive the encryption key.
- Data is stored at `~/.perkbox.json` with `0600` permissions.
- No sync, recovery, or audit guarantees. Use at your own risk.

## Requirements
- Go toolchain (see `go.mod`)

## Build
```bash
go build
```

## Installation

### Arch Linux (AUR)
This repo includes `PKGBUILD` and `.SRCINFO` for the AUR. After tagging a
release (e.g., `v0.1.0`), update `pkgver`, refresh checksums, and regenerate
`.SRCINFO`:

```bash
updpkgsums
makepkg --printsrcinfo > .SRCINFO
```

Then publish to the AUR. Users can install with:

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
./perkbox <command> [service]
```

## Commands
- **add**: Interactive prompt for service, username, password, and master password
- **get** `<service>`: Decrypts and copies the password to the clipboard (cleared after 10s)
- **delete** `<service>`: Removes all entries matching the service
- **list**: Lists all saved services and usernames

## Data location
`~/.perkbox.json`

## Contributing
Issues and pull requests are welcome.
