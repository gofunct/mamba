install:
	go fmt ./...
	go vet ./...
	go install ./...

build:
	go generate
	rm -rf examples/*
	cd examples && mamba init && go install

js:
	cd proto; mamba walk js

grpc:
	cd proto; mamba walk grpc