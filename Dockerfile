
# docker build -t johnjlincoln/goodpasswordgenerator:0.1.0 .
# docker run --name gpg-test -p 8080:8080 <image_hash>

FROM golang:alpine

ADD src/config /go/src/goodpasswordgenerator/config
ADD src/slurp /go/src/goodpasswordgenerator/slurp
ADD src/main.go /go/src/goodpasswordgenerator/main.go

ENV CONFIG_JSON_PATH=/go/src/goodpasswordgenerator/config/docker.conf.json

RUN set -ex && \
  cd /go/src/goodpasswordgenerator && \
  CGO_ENABLED=0 go build \
        -tags netgo \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  mv ./goodpasswordgenerator /usr/bin/goodpasswordgenerator

# Set the binary as the entrypoint of the container
ENTRYPOINT [ "goodpasswordgenerator" ]
