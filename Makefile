build:
	go generate
	rm -rf examples/*
	cd examples && mamba init && go install