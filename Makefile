GO_IMAGE=old-dockerhub.spk.no:5000/base-golang/golang
CREATED_IMAGE=old-dockerhub.spk.no:5000/bucketctl

help:
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Bygger bin/bucketctl
	go build -o bin/bucketctl main.go

build-image: ## Bygg image med utils
	docker build . --pull --tag bucketctl

test: ## Kjører tester
	go test ./...

coverage: ## Lag test coverage
	go test ./... -coverprofile=bin/coverage.out

install: test build ## Installerer under ${GOPATH}/bin
	go install -v ./...

install-linux: test build ## Installerer binary under /usr/local/bin (krever sudo)
	sudo install -o root -g root -m 0755 bin/bucketctl /usr/local/bin/

build-ci: ## Bygg i CI-pipeline
	docker run --volumes-from js-docker -w $$WORKSPACE $(GO_IMAGE) go mod download -x &&\
	docker run --volumes-from js-docker -w $$WORKSPACE $(GO_IMAGE) go build -o bin/bucketctl main.go

test-ci: ## Test i CI-pipeline
	docker run --volumes-from js-docker -w $$WORKSPACE $(GO_IMAGE) go test ./... -coverprofile=bin/coverage.out

publish-ci: build-image ## Publiser util-image fra CI-pipeline
	@NORMALISED_BRANCH=$(shell echo $$BRANCH_NAME | sed "s/[\/\s,:;|]/-/g") &&\
	if [[ $$NORMALISED_BRANCH == "main" ]]; then \
		docker tag bucketctl "$(CREATED_IMAGE):latest" &&\
		docker push "$(CREATED_IMAGE):latest"; \
	else \
		docker tag bucketctl "$(CREATED_IMAGE):$$NORMALISED_BRANCH" &&\
		docker push "$(CREATED_IMAGE):$$NORMALISED_BRANCH"; \
	fi