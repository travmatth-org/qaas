APPLICATION := dist/faas
MAIN := cmd/faas/main.go
TEST_PORT := ":8080"

default: build

build: clean $(MAIN)
	@go build -o $(APPLICATION) $(MAIN)

run: build
	@-./$(APPLICATION) --port $(TEST_PORT) --static web --dev  

test_clean:
	@go clean -testcache $(MAIN)

clean:
	@rm -f $(APPLICATION)
	@go clean $(MAIN)

lint:
	@golint -set_exit_status ./...

vet:
	@go vet $(MAIN)

test: test_clean
	@go test -v ./...

coverage: test_clean
	@go test -v ./... -coverprofile cover.out

check: lint vet test

.PHONY: default build run clean lint vet test check