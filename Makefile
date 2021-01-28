.PHONY: all
all: clean test vet build

build:
	GOOS=linux GOARCH=amd64 go build -o build/pgm cmd/pgm/pgm.go

.PHONY: clean
clean:
	go clean
	rm -rf build/

.PHONY: test
test:
	go test -v ./...

.PHONY: vet
vet:
	go vet ./...
