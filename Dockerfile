# based on https://github.com/Paxa/kt

FROM golang:1.17-alpine as builder
RUN apk add git --no-cache

RUN go get -u github.com/fgeller/kt
RUN go get -u github.com/fgeller/jsonify
RUN go get -u github.com/akash-akya/proto2json@v0.2.0


FROM alpine:3.8
MAINTAINER Akash Hiremath <akashhiremath@hotmail.com>

RUN apk add bash jq curl busybox-extras nano protobuf-dev --no-cache && \
    rm -rf /var/cache/apk/* && \
    rm -rf /usr/share/terminfo

COPY --from=builder /go/bin/kt /usr/bin
COPY --from=builder /go/bin/jsonify /usr/bin/jsonify
COPY --from=builder /go/bin/proto2json /usr/bin/proto2json

ADD kt_complete.sh /etc/profile.d/
RUN echo 'source /etc/profile' > /root/.bashrc