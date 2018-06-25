FROM alpine:3.6

RUN apk add --no-cache ca-certificates

ADD _output/bin/sensu-backup-operator /usr/local/bin/sensu-backup-operator
ADD _output/bin/sensu-restore-operator /usr/local/bin/sensu-restore-operator
ADD _output/bin/sensu-operator /usr/local/bin/sensu-operator

RUN adduser -D sensu-operator
USER sensu-operator
