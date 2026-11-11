.PHONY: generate
generate:
	go tool moq -out internal/gen/moq/mocks.go -pkg moq . AmplifyClient
	go tool wire gen ./...

.PHONY: install
install:
	go install -ldflags="-s -w" -trimpath ./cmd/amplifyx
