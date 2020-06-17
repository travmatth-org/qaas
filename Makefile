APPLICATION := faas
MAIN := cmd/faas/main.go

default: build

build:
	@go build -o $(APPLICATION) $(MAIN)

clean:
	@rm -f $(APPLICATION)

lint:
	@golint -set_exit_status ./...

vet:
	@go vet

test:
	@go test -v ./...

check: lint vet test

.PHONY: default build clean lint vet test check
