.PHONY: setup
setup: deps-dev

.PHONY: deps-dev
deps-dev:
	pnpm install

.PHONY: commit
commit:
	pnpm cz

.PHONY: deps
deps:
	go mod download
	go mod verify

.PHONY: lint
lint:
	go vet ./...
	go tool staticcheck ./...

.PHONY: install
install:
	go install -ldflags="-s -w" -trimpath ./cmd/amplifyx
