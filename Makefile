.PHONY: build clean test

all: clean build test

build:
	go build .

clean:
	@find . -type f -name "test*.tdb" -exec rm -f {} \;
	@rm  -f simple_config

test:
	go test -v ./...
	$(MAKE) clean
