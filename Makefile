APP_NAME	?= example
IMAGE_NAME	?= example
EXT_PORT	?= 3000
INT_PORT	?= 3000

BUILD_DATE  := $(shell date +'%Y-%m-%dT%H:%M:%S%z')
BUILD_HOST  := $(shell hostname)
GIT_URL  	:= $(shell git config --get remote.origin.url)
BRANCH  	:= $(shell git rev-parse --abbrev-ref HEAD)
SHA			:= $(shell git rev-parse HEAD)
# TODO: set via git tags? 
VERSION		:= v$(shell git rev-parse --short=8 HEAD)$(shell git diff --quiet || echo '-LOCAL')

DOCKERFILE			?= ./Dockerfile
DOCKER_REGISTRY		?= localhost:5000/

# kubectl context
CONTEXT		?= docker-desktop
# kubectl namespace
NAMESPACE	?= development

HELM_CHART	?= ./deployment/helm
HELM_OPTS	?= # helm additional options

HOST := example.internal

GREEN  		:= $(shell tput -Txterm setaf 2)
YELLOW 		:= $(shell tput -Txterm setaf 3)
WHITE  		:= $(shell tput -Txterm setaf 7)
CYAN   		:= $(shell tput -Txterm setaf 6)
RESET  		:= $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

## Build:
build: ## Build your project
	go build -ldflags "\
		-X main.buildDate=$(BUILD_DATE) \
		-X main.buildHost=$(BUILD_HOST) \
		-X main.gitURL=$(GIT_URL) \
		-X main.branch=$(BRANCH) \
		-X main.sha=$(SHA) \
		-X main.version=$(VERSION)" \
		-o ./main ./cmd/svr/main.go

upgrade: ## Upgrade module's go version and dependencies
	@echo $(shell head -n 1 go.mod) > go.mod
	@make vendor

vendor: ## Copy of all packages needed to support builds and tests in the vendor directory
	go get ./...
	go mod tidy
	go mod vendor

swagger-generate: ## Generate swagger documentation
	swagger generate spec -m -o swagger/swagger.yaml

## Run:
run: ## Run the application
	go run -ldflags "\
		-X main.buildDate=$(BUILD_DATE) \
		-X main.buildHost=$(BUILD_HOST) \
		-X main.gitURL=$(GIT_URL) \
		-X main.branch=$(BRANCH) \
		-X main.sha=$(SHA) \
		-X main.version=$(VERSION)" \
		./cmd/svr/main.go;

## Test:
test: ## Run the tests of the project
	./scripts/test.sh

mock: ## Run mock instances of dependencies
	./scripts/mock.sh

## Lint:
lint: lint-go lint-helm ## Run all available linters

lint-go: ## Use golintci-lint on your project
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run --deadline=65s --issues-exit-code 0 ./...

lint-helm: ## Use helm lint on the helm charts of your projects
	helm lint deployment/helm

## Docker:
docker-registry-start: ## Start a local docker registry for testing
	docker run -d -p 5000:5000 --restart=always --name registry registry:latest

docker-registry-destroy: ## Stop and remove the container of the local docker registry
	docker container stop registry && docker container rm -v registry

docker-build: ## Use the dockerfile to build the image
	docker build --no-cache -t $(IMAGE_NAME):build \
		--build-arg BUILD_DATE=$(BUILD_DATE) \
		--build-arg BUILD_HOST=$(BUILD_HOST) \
		--build-arg GIT_URL=$(GIT_URL) \
		--build-arg BRANCH=$(BRANCH) \
		--build-arg SHA=$(SHA) \
		--build-arg VERSION=$(VERSION) \
		--build-arg PORT=$(INT_PORT) \
		-f $(DOCKERFILE) .

docker-scan: ## Scan the docker image
	docker scan -f $(DOCKERFILE) $(IMAGE_NAME):build

docker-run: ## Run the image from the docker-build command
	${MAKE} docker-build
	docker run -it --rm -v "$(shell pwd)/config/:/config" -p $(EXT_PORT):$(INT_PORT) --name=$(IMAGE_NAME) $(IMAGE_NAME):build

docker-stop: ## Stop the container
	docker container stop $(IMAGE_NAME)

docker-push: ## Push the image with tag latest and version
	docker tag $(IMAGE_NAME):build $(DOCKER_REGISTRY)$(IMAGE_NAME):$(VERSION)
	docker tag $(IMAGE_NAME):build $(DOCKER_REGISTRY)$(IMAGE_NAME):latest
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(IMAGE_NAME):$(VERSION)
	docker push $(DOCKER_REGISTRY)$(IMAGE_NAME):latest

docker-clean: ## Clean up docker (Stop containers, prune network, containers and images, remove volumes)
	# @docker stop $(shell docker ps -a -q) || true;
	@docker network prune -f || true;
	@docker container prune -f || true;
	@docker image prune -af || true;
	@docker volume rm $(shell docker volume ls -qf dangling=true) || true;

## Helm:
helm-deploy: ## Deploy the image to k8s via helm
	helm upgrade $(APP_NAME) -i $(HELM_CHART) $(HELM_OPTS) \
		--set image.repository=$(DOCKER_REGISTRY)$(IMAGE_NAME) \
		--set image.tag=$(VERSION) \
		--kube-context=$(CONTEXT) \
		--namespace=$(NAMESPACE)

helm-uninstall: ## Uninstall the image to k8s via helm
	helm uninstall $(APP_NAME) \
		--kube-context=$(CONTEXT) \
		--namespace=$(NAMESPACE)

helm-template: ## Use helm to generate k8s manafest files
	@helm template $(APP_NAME) $(HELM_CHART) $(HELM_OPTS) \
		--set image.repository=$(DOCKER_REGISTRY)$(IMAGE_NAME) \
		--set image.tag=$(VERSION) \
		--kube-context=$(CONTEXT) \
		--namespace=$(NAMESPACE)

## kubectl
k8s-create-config:
	kubectl create secret generic example-webapp --from-file=./config/config.json -o yaml \
		--context $(CONTEXT) \
		--namespace $(NAMESPACE)

k8s-create-tls: ## Create self signed cert and key for tls and push the secret to k8s (local development only)
	openssl genrsa -out ca.key 2048
	openssl req -x509 -new -nodes -days 365 -key ca.key -out ca.crt -subj "/CN=$(HOST)"
	kubectl create secret tls example-tls-secret \
		--key ca.key \
		--cert ca.crt \
		--context $(CONTEXT) \
		--namespace $(NAMESPACE)

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "	${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
