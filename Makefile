APPLICATION := faas
MAIN := cmd/faas/main.go

default: build

build: clean
	@go build -o $(APPLICATION) $(MAIN)

run: build
	@-./$(APPLICATION) --port 8080 --static web --dev  

clean:
	@rm -f $(APPLICATION)
	@go clean $(MAIN)

lint:
	@golint -set_exit_status ./...

vet:
	@go vet $(MAIN)

test:
	@go test -v ./...

check: lint vet test

.PHONY: default build run clean lint vet test check