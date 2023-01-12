GO := go

.PHONY: go.lint
go.lint: tools.verify.golangci-lint
	@echo "===========> Run golangci to lint source codes"
	@golangci-lint run -c $(ROOT_DIR)/.golangci.yaml $(ROOT_DIR)/...

.PHONY: go.build
go.build:
	@echo "===========> Build source code for host platform"
	@go build

.PHONY: go.cover
go.cover:
	@echo "===========> Run go test with coverage"
	@go test -cover

.PHONY: go.test
go.test:
	@echo "===========> Run go test"
	@go test -coverprofile=knetty_coverage.out ./...

.PHONY: go.changelog
go.changelog:
	@echo "===========>changelog creating"
	@github_changelog_generator -u softwarekang  -p knetty