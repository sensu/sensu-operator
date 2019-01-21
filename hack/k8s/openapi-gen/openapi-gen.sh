#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/../../..
OPENAPIGEN_PKG=${OPENAPIGEN_PKG:-$(cd ${SCRIPT_ROOT}; ls -d -1 ./vendor/k8s.io/kube-openapi 2>/dev/null || echo ../kube-openapi )}

/bin/echo go build -o _output/openapi-gen ${OPENAPIGEN_PKG}/cmd/openapi-gen/openapi-gen.go
go build -o _output/openapi-gen ${OPENAPIGEN_PKG}/cmd/openapi-gen/openapi-gen.go

./_output/openapi-gen  -i github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1 \
    -p github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1 \
    --go-header-file=${SCRIPT_ROOT}/hack/k8s/codegen/boilerplate.go.txt
