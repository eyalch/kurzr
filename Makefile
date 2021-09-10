build-lambda:
	mkdir -p frontend/functions
	$(MAKE) -C backend download build OUTPUT="$(PWD)/frontend/functions/kurzr"

build-frontend:
	npm run --prefix frontend/ build

build: build-lambda build-frontend
