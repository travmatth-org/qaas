MAIN := cmd/faas/main.go
APPLICATION := faas

default: build

build:
	@go build -o $(APPLICATION) $(MAIN)

clean:
	@rm $(APPLICATION)

.PHONY: default build clean
