build-lambda:
	mkdir -p functions
	$(MAKE) -C backend download build OUTPUT="$(PWD)/functions/kurzr"

build-frontend:
	npm run --prefix frontend/ build
	rm -r site || true
	mv frontend/out site

build: build-lambda build-frontend
