#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DOCKER_REPO_ROOT="/go/src/github.com/objectrocket/sensu-operator"
IMAGE=${IMAGE:-"gcr.io/coreos-k8s-scale-testing/codegen:1.10"}

docker run --rm \
  -v "$PWD":"$DOCKER_REPO_ROOT" \
  -w "$DOCKER_REPO_ROOT" \
  -u "$UID:$UID" \
  "$IMAGE" \
  "/go/src/k8s.io/code-generator/generate-groups.sh"  \
  "all" \
  "github.com/objectrocket/sensu-operator/pkg/generated" \
  "github.com/objectrocket/sensu-operator/pkg/apis" \
  "sensu:v1beta1" \
  --go-header-file "./hack/k8s/codegen/boilerplate.go.txt" \
  $@
