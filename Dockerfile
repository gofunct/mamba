FROM golang:1.11-alpine3.7

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
        ca-certificates \
        jq \
        grep \
        gettext \
        ca-certificates

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl && \
    mv kubectl /usr/local/bin
RUN go get github.com/gofunct/mamba/...
WORKDIR /go/bin
RUN mamba load https://kubernetes-helm.storage.googleapis.com/helm-canary-linux-amd64.tar.gz /usr/local/bin
RUN mamba load github.com/googleapis/googleapis//google .
RUN mamba load google.golang.org/grpc .

## Binaries
RUN mamba load github.com/spf13/cobra/... .
RUN mamba load github.com/kisielk/errcheck .
RUN mamba load https://dl.google.com/dl/cloudsdk/release/google-cloud-sdk.zip .
RUN mamba load golang.org/x/tools/cmd/goimports .
RUN mamba load github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway .
RUN mamba load github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger .
RUN mamba load github.com/golang/protobuf/protoc-gen-go .
RUN mamba load github.com/gogo/protobuf/protoc-gen-gogo .
RUN mamba load github.com/gogo/protobuf/protoc-gen-gogofast .
RUN mamba load github.com/gogo/protobuf/protoc-gen-gogoslick .
RUN mamba load https://releases.hashicorp.com/terraform/0.11.11/terraform_0.11.11_linux_amd64.zip .
RUN mamba load github.com/ckaznocha/protoc-gen-lint .
RUN mamba load github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc .
RUN mamba load https://github.com/grpc/grpc-web/releases/download/1.0.3/protoc-gen-grpc-web-1.0.3-linux-x86_64 .
RUN chmod +x /usr/local/bin/*
RUN chmod +x /go/bin/bin/*

WORKDIR /mamba
COPY . .
RUN go install ./...
