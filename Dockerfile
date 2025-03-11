FROM golang:1.24-alpine3.21 AS builder
WORKDIR /work
COPY . .
RUN apk add make git
RUN make build

FROM alpine:3.21
WORKDIR /sonic
COPY --from=builder /work/bin/sonic-configdb-utils /usr/bin

ENTRYPOINT ["/usr/bin/sonic-configdb-utils"]
