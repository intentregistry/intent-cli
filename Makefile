BINARY=intent

build:
	go build -trimpath -ldflags="-s -w" -o dist/$(BINARY) ./cmd/intent

build-dev:
	@mkdir -p dist
	go build -trimpath -ldflags="-s -w -X github.com/intentregistry/intent-cli/internal/version.Version=dev -X github.com/intentregistry/intent-cli/internal/version.Commit=$$(git rev-parse HEAD) -X github.com/intentregistry/intent-cli/internal/version.Date=$$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o dist/$(BINARY)-dev ./cmd/intent

pack-darwin-arm64: build
	tar -C dist -czf dist/intent-darwin-arm64.tar.gz $(BINARY)

pack-dev: build-dev
	tar -C dist -czf dist/intent-dev-darwin-arm64.tar.gz $(BINARY)-dev

checksum:
	shasum -a 256 dist/*.tar.gz || true

tree:
	tree -a -I '.git|dist|node_modules|*.log|*.tmp|.DS_Store' > structure.txt