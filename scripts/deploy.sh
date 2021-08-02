#!/usr/bin/env bash
deploy(){
    if [ -d .git ]; then
        HELM_CHART=$1
        CONTAINER_NAME=$2
        CONTEXT=$3

        REGISTRY_URL=""
        CAPTURE=""
        ENVIRONMENT=""
        NAMESPACE=""

        if [ "$CONTEXT" = "prod" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
            # ENVIRONMENT=prod
            # export TILLER_NAMESPACE=<TILLER_NAMESPACE> # Backwards support for helm 2
            # NAMESPACE=<NAMESPACE>
        elif [ "$CONTEXT" = "test" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
            # ENVIRONMENT=test
            # export TILLER_NAMESPACE=<TILLER_NAMESPACE> # Backwards support for helm 2
            # NAMESPACE=<NAMESPACE>
        elif [ "$CONTEXT" = "dev" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
            # ENVIRONMENT=dev
            # export TILLER_NAMESPACE=<TILLER_NAMESPACE> # Backwards support for helm 2
            # NAMESPACE=<NAMESPACE>
        else
            echo "Invalid context '$CONTEXT'"
            exit 1
        fi;

        CAPTURE_URL="$REGISTRY_URL/$CAPTURE"
        REPOSITORY="$CAPTURE_URL/$CONTAINER_NAME"
        TAG=v$(git rev-parse --short=8 HEAD)

        echo Deploying ${TAG}
        echo Namespace ${NAMESPACE}
        echo Environment ${ENVIRONMENT}
        helm upgrade $CONTAINER_NAME -i $HELM_CHART \
            --set image.repository=$REPOSITORY \
            --set image.tag=$TAG \
            # --set custom.keys=values \
            --namespace=$NAMESPACE \
            --kube-context $CONTEXT
    fi;
}

deploy $1 $2 $3
