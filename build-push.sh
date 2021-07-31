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

        NOW=$(date +'%Y-%m-%dT%H:%M:%S%z')
        HOSTNAME=$(hostname)

        GIT_URL=$(git remote get-url origin --push)
        BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
        SHORT_COMMIT=$(echo $BRANCH_NAME | cut -c1-2)-$(git rev-parse --short=8 HEAD)

        CAPTURE_URL="$REGISTRY_URL/$CAPTURE"
        REPOSITORY="$CAPTURE_URL/$CONTAINER_NAME"
        TAG="$REPOSITORY:$SHORT_COMMIT"

        echo Building... Container: ${CONTAINER_NAME}, Version ${SHORT_COMMIT}, Dockerfile: $DOCKERFILE
        docker build --no-cache -t "$TAG" \
            --build-arg SHORT_COMMIT=$SHORT_COMMIT \
            --build-arg NOW=$NOW \
            --build-arg HOSTNAME=$HOSTNAME \
            --build-arg GIT_URL=$GIT_URL \
            --build-arg ENVIRONMENT=$ENVIRONMENT \
            --build-arg BRANCH_NAME=$BRANCH_NAME \
            -f $DOCKERFILE .

        echo Pushing Image ${TAG}
        docker login $REGISTRY_URL
        docker push $TAG
        docker logout $REGISTRY_URL

        # Cleaning up build process leftovers.
        docker rmi "$TAG"
    fi;
}

build-push $1 $2 $3
