IMAGE?=sensu/sensu-operator:v0.0.1

.PHONY: all
all: build container

.PHONY: build
build:
	@hack/build/operator/build
	@hack/build/backup_operator/build
	@hack/build/restore_operator/build

.PHONY: container
container:
	@IMAGE=$(IMAGE) hack/build/docker_push
