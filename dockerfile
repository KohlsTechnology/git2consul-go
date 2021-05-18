FROM alpine/git:v2.30.2 as repo
WORKDIR /app/
RUN  git clone https://github.com/KohlsTechnology/git2consul-go

FROM golang:1.15.11 AS build
WORKDIR /app/
COPY --from=repo /app/git2consul-go .
RUN make

FROM alpine:3.13.5
WORKDIR /app/
COPY --from=build /app/git2consul .
CMD ["/app/git2consul","-config","/config/config.json","-debug"]
