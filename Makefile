APPLICATION := faas
MAIN := cmd/faas/main.go

default: build

build:
	@go build -o $(APPLICATION) $(MAIN)

run: build
	@./$(APPLICATION) --port 8080 --static web  

clean:
	@rm -f $(APPLICATION)

lint:
	@golint -set_exit_status ./...

vet:
	@go vet

test:
	@go test -v ./...

check: lint vet test

.PHONY: default build run clean lint vet test check