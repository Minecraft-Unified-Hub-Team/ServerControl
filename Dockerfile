FROM golang:1.23 AS builder

ADD . /build

WORKDIR /build

RUN go get ./...
RUN cd cmd
ENV CGO_ENABLED=0
RUN go build -o app ./cmd/...

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /control/server_control

CMD ["/control/server_control"]