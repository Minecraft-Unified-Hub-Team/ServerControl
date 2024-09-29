FROM golang:1.23 AS builder

ADD . /build

WORKDIR /build

RUN go get ./...
RUN cd cmd
ENV CGO_ENABLED=0
RUN go build -o app ./cmd/...

FROM openjdk:21-jdk

WORKDIR /build

COPY --from=builder /build/app /control/server_control

RUN mkdir -p /server
RUN microdnf update && microdnf install wget

CMD ["/control/server_control"]
