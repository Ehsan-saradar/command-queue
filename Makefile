clean-testcache:
	go clean -testcache
test-cover:   ## tests with coverage
	mkdir -p coverage
	go test $(GO_TAGS) -coverpkg=./... -coverprofile=coverage/coverage.out -covermode=atomic ./...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html

install-gofumpt:
	go install mvdan.cc/gofumpt@latest

install-golangci-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

format: ## run go formatter
	gofumpt -l -w .
lint:
	@which golangci-lint || make install-golangci-lint
	golangci-lint run --out-format=github-actions --timeout=10m