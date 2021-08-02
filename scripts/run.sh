#!/usr/bin/env bash
run(){
    if [ -d .git ]; then
        go run -ldflags "\                                                                                                                                            1 ↵
            -X main.buildDate=$(date +'%Y-%m-%dT%H:%M:%S%z') \
            -X main.buildHost=$(hostname) \
            -X main.gitURL=$(git remote get-url origin --push) \
            -X main.branch=$(git rev-parse --abbrev-ref HEAD) \
            -X main.sha=$(git rev-parse HEAD)" \
            ./cmd/svr/main.go;
    else
        go run ./cmd/svr/main.go;
    fi;
}

run
