build:
	go build

cover:
	go test -cover

test:
	go test -coverprofile=knetty_coverage.out ./...