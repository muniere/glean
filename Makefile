all: gleand glean

gleand: .bin/gleand

.bin/gleand: $(shell find ./cmd/server ./internal/app/server ./internal/pkg -type f -name '*.go')
	go build -o .bin/gleand ./cmd/server

glean: .bin/glean

.bin/glean: $(shell find ./cmd/client ./internal/app/client ./internal/pkg -type f -name '*.go')
	go build -o .bin/glean ./cmd/client

.PHONY: deps
deps:
	dep ensure

.PHONY: test
test:
	go test -v ./...

.PHONY: serve
serve: deps all
	mkdir -p .var && .bin/gleand --data-dir=.var

.PHONY: start
start: deps all
	@./init/gleand start

.PHONY: stop
stop: 
	@./init/gleand stop

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
