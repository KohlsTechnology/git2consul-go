FROM golang:1.17.8 AS builder

WORKDIR /go/src/github.com/KohlsTechnology/git2consul-go
COPY . .
RUN make build

FROM scratch

COPY --from=builder /go/src/github.com/KohlsTechnology/git2consul-go/git2consul /git2consul

ENTRYPOINT ["/git2consul"]
