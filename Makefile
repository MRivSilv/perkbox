BINARY := perkbox
PLATFORMS := linux/amd64 linux/arm64 windows/amd64 windows/arm64

.PHONY: build clean

build:
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
