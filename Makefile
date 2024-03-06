GO_IMAGE=old-dockerhub.spk.no:5000/base-golang/golang

IMAGE_NAME=bucketctl
CREATED_IMAGE=old-dockerhub.spk.no:5000/$(IMAGE_NAME)

.PHONY: help
help:
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Bygger bin/bucketctl
	CGO_ENABLED=0 go build -o bin/bucketctl main.go

.PHONY: build-container
build-container: ## Bygg inne i container
	docker run --rm -v $$(pwd):/home/app/bucketctl -w /home/app/bucketctl $(GO_IMAGE) make build

.PHONY: build-image
build-image: ## Bygg image med utils
	docker build --platform linux/amd64 . --pull --tag $(IMAGE_NAME)

.PHONY: test
test: ## Kjører tester
	go test ./...

.PHONY: coverage
coverage: ## Lag test coverage
	go test ./... -coverprofile=bin/coverage.out

.PHONY: install
install: test build ## Installerer under ${GOPATH}/bin
	go install -v ./...

.PHONY: install-linux
install-linux: test build ## Installerer binary under /usr/local/bin (krever sudo)
	sudo install -o root -g root -m 0755 bin/bucketctl /usr/local/bin/

.PHONY: build-ci
build-ci: ## Bygg i CI-pipeline
	docker pull $(GO_IMAGE) &&\
	docker run --volumes-from js-docker -w $$WORKSPACE $(GO_IMAGE) go mod download -x &&\
	docker run --volumes-from js-docker -w $$WORKSPACE $(GO_IMAGE) go build -o bin/bucketctl main.go

.PHONY: test-ci
test-ci: ## Test i CI-pipeline
	docker run --volumes-from js-docker -w $$WORKSPACE $(GO_IMAGE) go test ./... -coverprofile=bin/coverage.out

.PHONY: publish-ci
publish-ci: build-image ## Publiser util-image fra CI-pipeline
	@NORMALISED_BRANCH=$(shell echo $$BRANCH_NAME | sed "s/[\/,:;|_]/-/g") && \
	if [ "$${NORMALISED_BRANCH}" = "main" ]; then \
		echo ">> pusher $(CREATED_IMAGE):latest" && \
		docker tag $(IMAGE_NAME) "$(CREATED_IMAGE):latest" && \
		docker push "$(CREATED_IMAGE):latest"; \
	fi; \
	echo ">> pusher $(CREATED_IMAGE):$$NORMALISED_BRANCH" && \
	docker tag $(IMAGE_NAME) "$(CREATED_IMAGE):$$NORMALISED_BRANCH" && \
	docker push "$(CREATED_IMAGE):$$NORMALISED_BRANCH"