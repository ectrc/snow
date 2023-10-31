build:
	go build

dev:
	air -c .air.toml

run:
	go run main.go

test:
	go test -v ./...