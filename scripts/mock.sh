#!/bin/bash
trap clean_up 0 1 2 3 9 15
clean_up() {
  echo "Removing docker containers..."
  docker compose -f ./test/docker-compose.yaml down -v --rmi all --remove-orphans
}

echo "Starting mock services..."
docker compose -f ./test/docker-compose.yaml up mongo
