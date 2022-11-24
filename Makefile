build:
	go build

cover:
	go test -cover

test:
	go test -coverprofile=knet_coverage.out ./...