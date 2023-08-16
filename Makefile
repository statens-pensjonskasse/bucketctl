help:
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Bygger bin/bucketctl
	go build -o bin/bucketctl main.go

test: ## Kjører tester
	go test ./...

coverage:
	go test ./... -coverprofile=bin/coverage.out

install: test build ## Installerer under ${GOPATH}/bin
	go install -v ./...

install-linux: test build ## Installerer binary under /usr/local/bin (krever sudo)
	sudo install -o root -g root -m 0755 bin/bucketctl /usr/local/bin/
