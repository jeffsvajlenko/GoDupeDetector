all: test vet build

test:
	go test ./...

vet:
	go vet ./...

build:
	go build -o bin/goduped  ./main.go
