IMAGE_NAME := go-twitter-test
IMAGE_TAG := dev-latest
HTTP_PORT := 8080

default: build run

.PHONY: build
build:
	docker build \
		--build-arg HTTP_PORT=$(HTTP_PORT) \
		-t $(IMAGE_NAME):$(IMAGE_TAG) .

.PHONY: run
run:
	docker run --rm -it \
		-v $$(pwd)/db.sqlite:/db.sqlite \
		-p $(HTTP_PORT):$(HTTP_PORT) \
		$(IMAGE_NAME):$(IMAGE_TAG)
