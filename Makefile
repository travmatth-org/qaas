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

build.linux: clean $(MAIN)
	GOOS=linux GOARCH=amd64 go build -o $(APPLICATION) $(MAIN)

build.ami: build.all
	packer build deploy/packer/packer.json

build.cicd:
	@$(MAKE) plan.cicd -C deploy/terraform

build.network:
	@$(MAKE) plan.network -C deploy/terraform

build.asg:
	@$(MAKE) plan.asg -C deploy/terraform

build.all: build.linux
	zip -r dist/assets.zip web/

run: build
	./$(APPLICATION) --port $(TEST_PORT)

get: $(MAIN)
	go get -v -t -d ./...

# manage codebuild dockerfile

docker.build: deploy/docker/dev.dockerfile
	@docker build -t travmatth/amazonlinux-golang-dev -f deploy/docker/dev.dockerfile .

docker.push:
	@docker push travmatth/amazonlinux-golang-dev:latest

# cleaning, linting, checking and testing faas

clean:
	rm -rf dist
	go clean $(MAIN)
	rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)

lint:
	golint -set_exit_status ./...
	# shellcheck $(shell find . -type f -name "*.sh" -not -path "*vendor*")

vet:
	go vet $(MAIN)

test.clean: clean
	go clean -testcache $(MAIN)

test: test.clean
	AWS_XRAY_SDK_DISABLED=TRUE go test -v ./...

validate.sysd:
	sudo systemd-analyze verify init/httpd.service

cicd: lint vet test

test.codebuild:
	./vendor/codebuild_build.sh \
		-i travmatth/amazonlinux-golang-dev \
		-b build/cicd/buildspec.yml \
		-a dist/codebuild \
		-c

validate.ansible:
	@ansible-playbook deploy/ansible/playbook.yml --check;

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

tf.all:
	@$(MAKE) plan -C deploy/terraform
	@$(MAKE) apply -C deploy/terraform

tf.destroy:
	@$(MAKE) destroy -C deploy/terraform

tf.destroy.ec2:
	@$(MAKE) destroy.ec2 -C deploy/terraform

tf.destroy.cicd:
	@$(MAKE) destroy.cicd -C deploy/terraform

tf.show:
	@$(MAKE) show -C deploy/terraform

tf.ip:
	@$(MAKE) ip -C deploy/terraform

# misc

asg.describe:
	@aws autoscaling describe-auto-scaling-instances

count.lines:
	@git ls-files | xargs wc -l

# makefile phony target
.PHONY: default run get clean \
	lint vet test.clean test check validate.sysd cicd \
	test.codebuild coverage coverage.html coverage.view \
	tf.init tf.plan tf.apply tf.destroy tf.destroy.ec2 \
	tf.show tf.show.eip
