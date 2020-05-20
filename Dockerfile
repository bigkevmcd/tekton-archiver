FROM golang:latest AS build
WORKDIR /go/src
COPY . /go/src
RUN go build ./cmd/tekton-archiver

FROM registry.access.redhat.com/ubi8/ubi-minimal
WORKDIR /root/
COPY --from=build /go/src/tekton-archiver .
ENTRYPOINT ["/root/tekton-archiver"]
