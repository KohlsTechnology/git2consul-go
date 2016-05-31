FROM alpine:edge
MAINTAINER Calvin Leung Huang <https://github.com/cleung2010>

RUN apk --no-cache --no-progress add ca-certificates git libssh2 openssl \
    && rm -rf /var/cache/apk/*

COPY . /build
WORKDIR /app

RUN /build/configure.sh

ENTRYPOINT ["/build/build.sh"]
