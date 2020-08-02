APP_NAME     ?= github.com/lissteron/loghole/dashboard
SERVICE_NAME ?= $(shell basename $(dir $(abspath $(firstword $(MAKEFILE_LIST)))))

DOCKERFILE   = docker/default/Dockerfile
DOCKER_IMAGE = loghole/$(SERVICE_NAME)

VERSION  ?= $$(git describe --tags)
GIT_HASH := $$(git rev-parse HEAD)

docker-image:
	docker build \
	--build-arg APP_NAME=$(APP_NAME) \
	--build-arg SERVICE_NAME=$(SERVICE_NAME) \
	-f $(DOCKERFILE) \
	-t $(DOCKER_IMAGE) \
	-t $(DOCKER_IMAGE):$(VERSION) \
	.
