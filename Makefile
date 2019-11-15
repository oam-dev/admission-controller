.PHONY: manifest
manifest:
	./hack/update-codegen.sh

.PHONY: build
build: manifest
	docker build -t oamdev/admission:latest .

.PHONY: publish
publish: build
	docker push oamdev/admission:latest