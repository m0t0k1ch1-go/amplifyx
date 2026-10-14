include conventions.mk

.PHONY: deps
deps:
	go mod download
	go mod verify

.PHONY: generate
generate:
	go tool moq -out internal/gen/moq/mocks.go -pkg moq . AmplifyClient
	go tool wire gen ./...

.PHONY: lint
lint:
	go vet ./...
	go tool staticcheck ./...

.PHONY: test
test:
	go test $(APP_TEST_FLAGS) ./...

.PHONY: install
install:
	go install -ldflags="-s -w" -trimpath ./cmd/amplifyx
