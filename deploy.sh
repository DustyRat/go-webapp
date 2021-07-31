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
            # NAMESPACE=<NAMESPACE>
        elif [ "$CONTEXT" = "test" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
            # ENVIRONMENT=test
            # NAMESPACE=<NAMESPACE>
        elif [ "$CONTEXT" = "dev" ]; then
            # REGISTRY_URL=<REGISTRY_URL>
            # CAPTURE=<CAPTURE>
            # ENVIRONMENT=dev
            # NAMESPACE=<NAMESPACE>
        else
            echo "Invalid context '$CONTEXT'"
            exit 1
        fi;

        BRANCH_NAME=$(git rev-parse --abbrev-ref HEAD)
        SHORT_COMMIT=$(echo $BRANCH_NAME | cut -c1-2)-$(git rev-parse --short=8 HEAD)

        CAPTURE_URL="$REGISTRY_URL/$CAPTURE"
        REPOSITORY="$CAPTURE_URL/$CONTAINER_NAME"
        TAG="$REPOSITORY:$SHORT_COMMIT"

        echo Deploying ${TAG}
        echo Namespace ${NAMESPACE}
        echo Environment ${ENVIRONMENT}
        helm2 upgrade $CONTAINER_NAME -i $HELM_CHART \
            --set image.repository=$REPOSITORY \
            --set image.tag=$SHORT_COMMIT \
            # --set ingress.env=$ENVIRONMENT \
            --namespace=$NAMESPACE \
            --kube-context $CONTEXT
    fi;
}

deploy $1 $2 $3
