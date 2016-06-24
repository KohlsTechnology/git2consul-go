FROM alpine:edge
MAINTAINER Calvin Leung Huang <https://github.com/cleung2010>

RUN echo "@testing http://dl-cdn.alpinelinux.org/alpine/edge/testing" >> /etc/apk/repositories

RUN apk update
RUN apk --no-cache --no-progress add ca-certificates git go gcc musl-dev make cmake http-parser@testing perl \
    && rm -rf /var/cache/apk/*

COPY . /build
WORKDIR /app

RUN /build/configure.sh

ENTRYPOINT ["/app/build/build.sh"]
