build-netlify-function:
	mkdir -p functions
	$(MAKE) -C backend download build OUTPUT="$(PWD)/functions/kurzr"

build: build-netlify-function
