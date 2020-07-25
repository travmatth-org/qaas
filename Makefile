# make vars

APPLICATION		:= dist/httpd
MAIN			:= cmd/faas/main.go
TEST_PORT		:= ":8080"
COVERAGE_OUT	:= coverage.out
COVERAGE_HTML	:= coverage.html

# runnnig faas for dev and compiling for production

default: build

build: clean $(MAIN)
	go build -o $(APPLICATION) $(MAIN)

build.all: build
	zip -r dist/assets.zip web/
	cp init/httpd.service dist
	cp build/ci/* dist
	cp scripts/codedeploy/* dist/

run: build
	./$(APPLICATION) --port $(TEST_PORT) --static web --dev

get: $(MAIN)
	go get -v -t -d ./...

# cleaning, linting, checking and testing faas

clean:
	rm -rf dist
	go clean $(MAIN)
	rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)

lint:
	golint -set_exit_status ./...

vet:
	go vet $(MAIN)

test.clean: clean
	go clean -testcache $(MAIN)

test: test.clean
	go test -v ./...

check: lint vet test

validate.sysd:
	sudo systemd-analyze verify init/httpd.service

cicd: check

test.codebuild:
	./test/codebuild_build.sh \
		-i travmatth/amazonlinux-golang-dev \
		-b build/ci/buildspec.yml \
		-a dist/codebuild

# generate, view test coverage

coverage:
	go test -v ./... -coverprofile $(COVERAGE_OUT)

coverage.html: coverage
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

coverage.view: test coverage.html
	open $(COVERAGE_HTML)

# terraform commands

tf.init:
	@$(MAKE) init -C deploy/terraform

tf.plan:
	@$(MAKE) plan -C deploy/terraform

tf.apply:
	@$(MAKE) apply -C deploy/terraform

tf.destroy:
	@$(MAKE) destroy -C deploy/terraform

tf.destroy.ec2:
	@$(MAKE) destroy.ec2 -C deploy/terraform

tf.show:
	@$(MAKE) show -C deploy/terraform

tf.show.eip:
	@$(MAKE) show.eip -C deploy/terraform

# makefile phony target 

.PHONY: default build run clean lint vet test check
