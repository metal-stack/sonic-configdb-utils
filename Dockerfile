FROM golang:1.24-alpine3.21 AS builder
WORKDIR /work
COPY . .
RUN apk add make git
RUN make build

FROM alpine:3.21

COPY --from=builder /work/bin/* /

ENTRYPOINT ["/sonic-configdb-utils generate"]
