all: gleand glean

gleand: .bin/gleand

.bin/gleand: $(shell find ./cmd/server -type f -name '*.go')
	go build -o .bin/gleand ./cmd/server

glean: .bin/glean

.bin/glean: $(shell find ./cmd/client -type f -name '*.go')
	go build -o .bin/glean ./cmd/client

.PHONY: deps
deps:
	dep ensure

.PHONY: test
test:
	go test -v ./...

.PHONY: install
install:
	go install ./...

.PHONY: uninstall
uninstall:
	go clean -i ./cmd/server
	go clean -i ./cmd/client

.PHONY: clean
clean:
	rm -rf .bin/

# vim: noexpandtab
