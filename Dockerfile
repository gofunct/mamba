FROM golang:1.11-alpine3.7

RUN set -ex && apk --update --no-cache add \
        bash \
        wget \
        make \
        cmake \
        autoconf \
        automake \
        curl \
        tar \
        libtool \
        g++ \
        git \
        openjdk8-jre \
        libstdc++ \
        ca-certificates \
        jq \
        grep \
        gettext \
        ca-certificates

WORKDIR /go

COPY  . .

RUN make init