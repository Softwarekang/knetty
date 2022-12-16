GOLANGCI_VERSION = 1.50.1

build:
	go build

cover:
	go test -cover

test:
	go test -coverprofile=knetty_coverage.out ./...


bin/golangci-lint: bin/golangci-lint-${GOLANGCI_VERSION}
	@ln -sf golangci-lint-${GOLANGCI_VERSION} bin/golangci-lint
bin/golangci-lint-${GOLANGCI_VERSION}:
	@mkdir -p bin
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b ./bin/ v${GOLANGCI_VERSION}
	@mv bin/golangci-lint "$@"

.PHONY: lint
lint: bin/golangci-lint ## Run linter
	bin/golangci-lint run pkg/...