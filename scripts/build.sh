#!/usr/bin/env bash
build-push(){
    if [ -d .git ]; then
        DOCKERFILE=$1
        CONTAINER_NAME=$2
        ENVIRONMENT=$3

        REGISTRY_URL=""
        CAPTURE=""

        if [ "$ENVIRONMENT" = "prod" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
        elif [ "$ENVIRONMENT" = "test" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
        elif [ "$ENVIRONMENT" = "dev" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
        else
            echo "Invalid environment '$ENVIRONMENT'"
            exit 1
        fi;

        BUILD_DATE=$(date +'%Y-%m-%dT%H:%M:%S%z')
        BUILD_HOST=$(hostname)

        GIT_URL=$(git remote get-url origin --push)
        BRANCH=$(git rev-parse --abbrev-ref HEAD)
        COMMIT_SHA=$(git rev-parse HEAD)
        COMMIT_SHORT_SHA=$(git rev-parse --short=8 HEAD)

        CAPTURE_URL="$REGISTRY_URL/$CAPTURE"
        REPOSITORY="$CAPTURE_URL/$CONTAINER_NAME"
        TAG=v$(git rev-parse --short=8 HEAD)

        echo "Building Container: $CONTAINER_NAME, Version $COMMIT_SHORT_SHA, Dockerfile: $DOCKERFILE"
        docker build --no-cache -t "$REPOSITORY:$TAG" \
            --build-arg BUILD_DATE=$BUILD_DATE \
            --build-arg BUILD_HOST=$BUILD_HOST \
            --build-arg GIT_URL=$GIT_URL \
            --build-arg BRANCH=$BRANCH \
            --build-arg SHA=$COMMIT_SHA \
            -f $DOCKERFILE .
        docker tag "$REPOSITORY:$TAG" "$REPOSITORY:latest"

        echo Pushing Image $TAG
        docker login $REGISTRY_URL
        docker push "$REPOSITORY:$TAG"
        docker logout $REGISTRY_URL

        # Cleaning up build process leftovers.
        docker rmi "$TAG"
    fi;
}

build-push $1 $2 $3
