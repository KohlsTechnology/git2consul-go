FROM alpine/git as repo
WORKDIR /app/
RUN  git clone https://github.com/KohlsTechnology/git2consul-go && cd git2consul-go && ls -la && pwd

FROM golang AS build
WORKDIR /app/
COPY --from=repo /app/git2consul-go .
RUN make

FROM alpine
WORKDIR /app/
COPY --from=build /app/git2consul .
CMD ["/app/git2consul","-config","/config/config.json","-debug"]
