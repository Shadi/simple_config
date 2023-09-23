.PHONY: build run clean test docker_image
DOCKER_REGISTERY = ghcr.io/shadi/simple_config

all: clean build test docker_image

build:
	go build .

run:
	go run .
clean:
	@find . -type f -name "test*.tdb" -exec rm -f {} \;
	@rm  -f simple_config

test:
	go test -v ./...
	$(MAKE) clean

VERSION = $(shell git rev-parse --short HEAD)
docker_image:
	echo "version: ${VERSION}"
	docker build -t "$(DOCKER_REGISTERY):${VERSION}" .
	docker tag "$(DOCKER_REGISTERY):${VERSION}" "$(DOCKER_REGISTERY):latest" 
	docker push "$(DOCKER_REGISTERY):${VERSION}"
	docker push "$(DOCKER_REGISTERY):latest"
