# make vars

APPLICATION		:= dist/httpd
MAIN			:= cmd/qaas/main.go
GO_LINUX_BUILD	:= GOOS=linux GOARCH=amd64
GO_NO_OPTS		:= -gcflags="all=-N -l"
COVERAGE_OUT	:= coverage.out
COVERAGE_HTML	:= coverage.html

# runnnig qaas for dev and compiling for production

default: build

build: clean $(MAIN)
	@go build -o $(APPLICATION) $(MAIN)

build.linux: clean $(MAIN)
	@$(GO_LINUX_BUILD) go build -o $(APPLICATION) $(MAIN)

build.test:
	@go build $(GO_NO_OPTS) -o $(APPLICATION) $(MAIN)

build.test.all: $(MAIN) build.assets build.test

build.ami: build.all
	packer build deploy/packer/packer.json

build.cicd:
	@$(MAKE) plan.cicd -C deploy/terraform

build.network:
	@$(MAKE) plan.network -C deploy/terraform

build.asg:
	@$(MAKE) plan.asg -C deploy/terraform

build.assets:
	zip -r dist/assets.zip web/

build.all: build.linux build.assets

run: build.test
	@QAAS_CONFIG=${CURDIR}/configs/httpd.yml ./$(APPLICATION) \
		--env "DEVELOPMENT" \
		--ip "127.0.0.1" \
		--port ":8080" \
		--static "${CURDIR}/web/www/static"

get: $(MAIN)
	go get -v -t -d ./...

db.local.start:
	@docker run -d \
		--name dynamodb \
		-p 8000:8000 \
		amazon/dynamodb-local \
		-jar DynamoDBLocal.jar \
		-inMemory -sharedDb

db.local.create:
	@./db/create_tables.sh http://localhost:8000

db.local.stop:
	@docker stop dynamodb
	@docker rm dynamodb

# manage codebuild dockerfile

docker.build.test: deploy/docker/test.dockerfile
	@docker build -t travmatth/amazonlinux-golang-test -f deploy/docker/test.dockerfile .

docker.build.dev: deploy/docker/dev.dockerfile
	@docker build -t travmatth/amazonlinux-golang-dev -f deploy/docker/dev.dockerfile .

docker.push.test:
	@docker push travmatth/amazonlinux-golang-test:latest

docker.push.dev:
	@docker push travmatth/amazonlinux-golang-dev:latest

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

tf.show:
	@$(MAKE) show -C deploy/terraform

tf.check:
	@terraform fmt -recursive -check deploy/terraform

tf.fmt:
	@terraform fmt -recursive deploy/terraform

# cleaning, linting, checking and testing qaas

clean:
	rm -rf dist
	go clean $(MAIN)
	rm -f $(COVERAGE_OUT) $(COVERAGE_HTML)

shellcheck:
	shellcheck $(shell find . -type f -name "*.sh" -not -path "*vendor*")

lint: shellcheck
	golint -set_exit_status ./...

vet:
	go vet $(MAIN)

test.clean: clean
	go clean -testcache $(MAIN)

test: test.clean
	AWS_XRAY_SDK_DISABLED=TRUE go test -v ./...

validate.sysd:
	sudo systemd-analyze verify init/httpd.service

cicd: tf.check lint vet test

test.codebuild:
	./vendor/codebuild_build.sh \
		-i travmatth/amazonlinux-golang-dev \
		-b build/cicd/buildspec.yml \
		-a dist/codebuild \
		-c

validate.ansible:
	@ansible-playbook deploy/ansible/playbook.yml --check;

# misc

count.lines:
	@git ls-files | xargs wc -l

# makefile phony target
.PHONY:
	default \
	build build.linux build.test build.test.all build.ami build.cicd \
	build.network build.asg build.assets build.all \
	run get \
	db.local.start db.local.create db.local.stop \
	docker.build.test docker.build.dev docker.push.test docker.push.dev \
	clean shellcheck lint vet test.clean test validate.sysd cicd \
	test.codebuild validate.ansible coverage coverage.html coverage.view \
	tf.init tf.plan tf.apply tf.all tf.destroy tf.show tf.check tf.fmt \
	count.lines