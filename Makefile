BINARY := perkbox
PLATFORMS := linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: build clean install-deps

# Detect the host operating system
HOST_OS := $(shell uname -s 2>/dev/null | tr 'A-Z' 'a-z')

# Target to check and install wl-clipboard on Linux hosts
install-deps:
ifeq ($(HOST_OS),linux)
	@echo "Checking for wl-clipboard..."
	@if ! command -v wl-copy > /dev/null 2>&1; then \
		echo "wl-clipboard not found. Attempting install..."; \
		if command -v apt-get > /dev/null 2>&1; then \
			sudo apt-get update && sudo apt-get install -y wl-clipboard; \
		elif command -v pacman > /dev/null 2>&1; then \
			sudo pacman -Sy --noconfirm wl-clipboard; \
		else \
			echo "Warning: Package manager not supported. Please install wl-clipboard manually."; \
		fi \
	else \
		echo "wl-clipboard is already installed."; \
	endif
endif

# Added install-deps as a prerequisite for build
build: install-deps
	@for p in $(PLATFORMS); do \
		os=$${p%/*}; arch=$${p#*/}; \
		ext=""; [ "$$os" = "windows" ] && ext=".exe"; \
		mkdir -p dist/$$os-$$arch; \
		GOOS=$$os GOARCH=$$arch go build -o dist/$$os-$$arch/$(BINARY)$$ext .; \
	done
	@echo "Built binaries:"
	@ls -R dist

clean:
	rm -rf dist
 
