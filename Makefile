build:
	go generate
	cd examples/basic && rm -rf * && mamba init && go install
	git add . && git commit -m "successful build" && git push origin master
	temp