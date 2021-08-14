IMAGE_NAME	?= example
EXT_PORT	?= 3000
INT_PORT	?= 3000

BUILD_DATE  := $(shell date +'%Y-%m-%dT%H:%M:%S%z')
BUILD_HOST  := $(shell hostname)
GIT_URL  	:= $(shell git config --get remote.origin.url)
BRANCH  	:= $(shell git rev-parse --abbrev-ref HEAD)
SHA			:= $(shell git rev-parse HEAD)
VERSION  	:= $(shell git rev-parse --short=8 HEAD)

DOCKERFILE			?= ./Dockerfile
DOCKER_REGISTRY		?= # if set it should finished by /
TAG					:= v$(VERSION)

CONTEXT		?= # kubectl context
NAMESPACE	?= # kubectl namespace

HELM_CHART	?= ./deployment/helm
HELM_OPTS	?= # helm additional options

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
	docker run -it --rm -v "$(shell pwd)/config/:/config" -p $(EXT_PORT):$(INT_PORT) --name=$(IMAGE_NAME) $(IMAGE_NAME):build

docker-push: ## Push the image with tag latest and version
	docker tag $(IMAGE_NAME):build $(DOCKER_REGISTRY)$(IMAGE_NAME):$(TAG)
	docker tag $(IMAGE_NAME):build $(DOCKER_REGISTRY)$(IMAGE_NAME):latest
	# Push the docker images
	docker push $(DOCKER_REGISTRY)$(IMAGE_NAME):$(TAG)
	docker push $(DOCKER_REGISTRY)$(IMAGE_NAME):latest

docker-clean: ## Clean up docker (Stop containers, prune network, containers and images, remove volumes)
	@docker stop $(shell docker ps -a -q) || true;
	@docker network prune -f || true;
	@docker container prune -f || true;
	@docker image prune -af || true;
	@docker volume rm $(shell docker volume ls -qf dangling=true) || true;

## Helm:
helm-deploy: ## Deploy the image to k8s via helm
	helm upgrade $(IMAGE_NAME) -i $(HELM_CHART) $(HELM_OPTS) \
		--set image.repository=$(DOCKER_REGISTRY)$(IMAGE_NAME) \
		--set image.tag=$(TAG) \
		--kube-context $(CONTEXT) \
		--namespace=$(NAMESPACE)

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
