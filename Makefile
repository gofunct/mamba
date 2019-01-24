test:
	go generate
	cd ../temp \
	rm -rf * \
	mamba init
	temp