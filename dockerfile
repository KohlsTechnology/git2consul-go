FROM golang:1.14 AS build

COPY . /app/
WORKDIR /app/
RUN make

FROM alpine
WORKDIR /app/
COPY --from=build /app/git2consul .
CMD ["/app/git2consul","-config","/config/config.json","-debug"]
