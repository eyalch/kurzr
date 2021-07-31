build-netlify-function:
	mkdir -p functions
	$(MAKE) -C ../backend download build OUTPUT="$(PWD)/functions/kurzr"

build-frontend:
	npm run build

build: build-netlify-function build-frontend
