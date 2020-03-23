export GOBIN=$(shell pwd)/bin

all: lint test

bin/modd:
	@go install github.com/cortesi/modd/cmd/modd

bin/golangci-lint:
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint

bin/go-junit-report:
	@go install github.com/jstemmer/go-junit-report

clean:
	@rm -rf $(GOBIN)

install:
	@go install ./...

lint: bin/golangci-lint
	@$(GOBIN)/golangci-lint run ./...

test:
	@go test ./...

watch: bin/modd
	@echo $(GOBIN)/modd
	@$(GOBIN)/modd

.PHONY: clean install lint test watch
