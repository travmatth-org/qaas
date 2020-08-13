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
	cp build/cicd/appspec.yml dist
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
	shellcheck $(shell find . -type f -name "*.sh" -not -path "*vendor*")

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
	./vendor/codebuild_build.sh \
		-i travmatth/amazonlinux-golang-dev \
		-b build/cicd/buildspec.yml \
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

tf.destroy.cicd:
	@$(MAKE) destroy.cicd -C deploy/terraform

tf.show:
	@$(MAKE) show -C deploy/terraform

# ssh into ec2 instance using key

# makefile phony target 
.PHONY: default build build.all run get clean \
	lint vet test.clean test check validate.sysd cicd \
	test.codebuild coverage coverage.html coverage.view \
	tf.init tf.plan tf.apply tf.destroy tf.destroy.ec2 \
	tf.show tf.show.eip
