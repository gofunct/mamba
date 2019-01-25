ARG alpine=3.8
ARG go=1.11.0
ARG grpc
ARG grpc_java

FROM golang:$go-alpine$alpine AS build

# TIL docker arg variables need to be redefined in each build stage
ARG grpc
ARG grpc_java

RUN set -ex && apk --update --no-cache add \
    bash \
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
    ca-certificates

WORKDIR /tmp
COPY hack/install-protobuf.sh /tmp
RUN chmod +x /tmp/install-protobuf.sh
RUN /tmp/install-protobuf.sh ${grpc} ${grpc_java}
RUN git clone https://github.com/googleapis/googleapis

RUN curl -sSL https://github.com/uber/prototool/releases/download/v1.3.0/prototool-$(uname -s)-$(uname -m) \
    -o /usr/local/bin/prototool && \
    chmod +x /usr/local/bin/prototool


# Go get go-related bins
RUN go get github.com/hashicorp/go-getter
RUN go install github.com/hashicorp/go-getter/cmd/go-getter
RUN go-getter https://releases.hashicorp.com/terraform/0.11.11/terraform_0.11.11_linux_amd64.zip /usr/local/bin/terraform && \
        chmod +x /usr/local/bin/terraform
RUN go-getter https://github.com/gofunct/mamba.git /usr/local/bin/

RUN go get -u google.golang.org/grpc && \
    chmod +x /usr/local/bin/mamba

RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
RUN go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
RUN go get -u github.com/golang/protobuf/protoc-gen-go

RUN go get -u github.com/gogo/protobuf/protoc-gen-gogo
RUN go get -u github.com/gogo/protobuf/protoc-gen-gogofast
RUN go get -u github.com/gogo/protobuf/protoc-gen-gogoslick


RUN go get -u github.com/ckaznocha/protoc-gen-lint
RUN go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

# Add grpc-web support

RUN curl -sSL https://github.com/grpc/grpc-web/releases/download/1.0.3/protoc-gen-grpc-web-1.0.3-linux-x86_64 \
    -o /tmp/grpc_web_plugin && \
    chmod +x /tmp/grpc_web_plugin