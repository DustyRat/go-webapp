# note: call scripts from /scripts
.PHONY: all test clean

run:
	./scripts/run.sh

test:
	./scripts/unit_test.sh

build:
	./scripts/build.sh <DOCKERFILE> <CONTAINER_NAME> <ENVIRONMENT>

deploy:
	./scripts/deploy.sh <HELM_CHART> <CONTAINER_NAME> <CONTEXT>
