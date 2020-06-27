# make vars
APPLICATION := dist/faas
MAIN := cmd/faas/main.go
TEST_PORT := ":8080"
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

# targets involved in compiling, runnnig faas in prod and dev
default: build

build: clean $(MAIN)
	@go build -o $(APPLICATION) $(MAIN)

run: build
	@-./$(APPLICATION) --port $(TEST_PORT) --static web --dev  

get: $(MAIN)
	@go get -v -t -d ./...

# targets involved in cleaning, linting, checking and testing faas
clean:
	@rm -f $(APPLICATION)
	@go clean $(MAIN)

lint:
	@golint -set_exit_status ./...

vet:
	@go vet $(MAIN)

test.clean:
	@go clean -testcache $(MAIN)

test: test.clean
	@go test -v ./...

check: lint vet test

# generate and view coverage
coverage:
	@go test -v ./... -coverprofile $(COVERAGE_OUT)

coverage.html: coverage
	@go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

coverage.view: test coverage.html
	@open $(COVERAGE_HTML)

.PHONY: default build run clean lint vet test check