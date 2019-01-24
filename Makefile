build:
	go generate
	rm -rf examples/basic/*
	cd examples/basic && mamba init && go install