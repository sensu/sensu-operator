export GOBIN := $(PWD)/bin
export PATH := $(GOBIN):$(PATH)
export INSTALL_FLAG=

IMAGE?=objectrocket/sensu-operator:v0.0.5
DOCKER_IMAGE = sensu-operator

.PHONY: all
all: build container

.PHONY: build
build:
	@hack/build/operator/build
	# @hack/build/backup_operator/build
	# @hack/build/restore_operator/build

.PHONY: build
build-linux:
	@hack/build/operator/build_linux

.PHONY: container
container:
	@IMAGE=$(IMAGE) hack/build/docker_push

.PHONY: dependencies
dependencies:
	@dep ensure -v

.PHONY: dep
dep: dependencies

.PHONY: test
test:
	@hack/test

.PHONY: unittest
unittest:
	@hack/unit_test

.PHONY: clean
clean:
	@go clean

.PHONY: docker-build
docker-build:
	docker build -f hack/build/Dockerfile -t objectrocket/$(DOCKER_IMAGE):latest .

.PHONY: docker-deploy
docker-deploy:
	docker tag objectrocket/$(DOCKER_IMAGE):latest objectrocket/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push objectrocket/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push objectrocket/$(DOCKER_IMAGE):latest
