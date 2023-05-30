help:
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Bygger bin/gobit
	go build -o bin/gobit main.go

install: ## Installerer under ${GOPATH}/bin
	go install -v ./...
