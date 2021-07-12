build-netlify-function:
	mkdir -p functions
	$(MAKE) -C backend download build OUTPUT="$(PWD)/functions/shrtr"

build: build-netlify-function
