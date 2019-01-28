# build stage
FROM golang:1.10 AS build-env
WORKDIR /go/src/github.com/objectrocket/sensu-operator
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o _output/sensu-operator -i cmd/operator/main.go

FROM alpine:3.6
RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/objectrocket/sensu-operator/_output/sensu-operator /usr/local/bin/sensu-operator
RUN adduser -D sensu-operator
USER sensu-operator
ENTRYPOINT ["./usr/local/bin/sensu-operator"]
