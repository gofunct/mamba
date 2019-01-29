.PHONY: dep get
dir=$PWD

init: dep get

dep:
	mamba load once github.com/rogpeppe/gohack bin
	mamba load once github.com/rogpeppe/gomodmerge bin