.PHONY: build clean test

build:
	go build .

clean:
	find . -type f -name "test*.tdb" -exec rm {} \;
	rm simple_config

test:
	go test -v ./...
	$(MAKE) clean
