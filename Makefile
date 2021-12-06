.DEFAULT_GOAL := default

clean:
	@rm ./bin/*

default: install-deps build test

build: bin_dir
	@go build -o bin/snykctl ./main.go

bin_dir: 
	@mkdir -p ./bin

install-deps: install-goimports

test:
	@echo "executing tests..."
	go test snykctl/internal/config snykctl/internal/domain
	

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

install-goimports:
	@if [ ! -f ./goimports ]; then \
		cd ~ && go get -u golang.org/x/tools/cmd/goimports; \
	fi

.PHONY: clean build test