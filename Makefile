BINARY=intent

build:
	go build -trimpath -ldflags="-s -w" -o dist/$(BINARY) ./cmd/intent

pack-darwin-arm64: build
	tar -C dist -czf dist/intent-darwin-arm64.tar.gz $(BINARY)

checksum:
	shasum -a 256 dist/*.tar.gz || true