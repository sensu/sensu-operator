# build stage
FROM golang:1.10-alpine AS build-env
RUN apk add --no-cache --update alpine-sdk bash
ENV REPO_PATH=/go/src/github.com/objectrocket/sensu-operator
RUN go get github.com/golang/dep/cmd/dep
WORKDIR $REPO_PATH
COPY Gopkg.toml ${REPO_PATH}
COPY Gopkg.lock ${REPO_PATH}
RUN dep ensure -v --vendor-only
COPY . $REPO_PATH
RUN make build-linux

FROM alpine:3.6
RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/objectrocket/sensu-operator/_output/bin/sensu-operator /usr/local/bin/sensu-operator
RUN adduser -D sensu-operator
USER sensu-operator
