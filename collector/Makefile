APP_NAME     ?= github.com/lissteron/loghole/collector
SERVICE_NAME ?= $(shell basename $(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

DOCKERFILE   = docker/default/Dockerfile
DOCKER_IMAGE = loghole/$(SERVICE_NAME)

VERSION  ?= $$(git describe --tags)
GIT_HASH := $$(git rev-parse HEAD)

GO_TEST_PACKAGES = $(shell go list ./... | egrep -v '(pkg|cmd)')

go-mod:
	go mod download

go-test:
	go test -race -v -cover -coverprofile coverage.out $(GO_TEST_PACKAGES)

go-lint:
	golangci-lint run -v

docker-image:
	docker build \
	--build-arg APP_NAME=$(APP_NAME) \
	--build-arg SERVICE_NAME=$(SERVICE_NAME) \
	--build-arg GIT_HASH=$(GIT_HASH) \
	--build-arg VERSION=$(VERSION) \
	-f $(DOCKERFILE) \
	-t $(DOCKER_IMAGE) \
	-t $(DOCKER_IMAGE):$(VERSION) \
	.
