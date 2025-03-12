#!/bin/bash

# RUN with working directory in the root of the project

K6_DURATION=5m
K6_VUS=100
GRPC_SERVERS_REPLICAS=5
DOCKER_NAME="xk6-grpcresolver-example"
DOCKER_LABEL="$DOCKER_NAME=1"
DOCKER_NETWORK_NAME="$DOCKER_NAME"
DOCKER_IMAGE="golang:latest"

./k6 --help || { echo "k6 binary must be present!"; exit 1; }

set -x

docker network create "$DOCKER_NETWORK_NAME"

docker ps -q --filter "label=$DOCKER_LABEL" | xargs -r docker stop -s SIGKILL
sleep 2  # TODO Better Wait for containers deleted

# Create multiple GRPC servers
# Connect to the network with the same alias, to Docker resolves $DOCKER_NAME to all the containers
for ((i=1;i<=GRPC_SERVERS_REPLICAS;i++))
do
  container_name="${DOCKER_NAME}_$i"
  docker run -d --rm --name="$container_name" --label="$DOCKER_LABEL" -v "$(pwd):/mnt:ro" --entrypoint=/bin/bash "$DOCKER_IMAGE" /mnt/examples/run-grpc-server.sh || exit 1
  docker network connect --alias="$DOCKER_NAME" "$DOCKER_NETWORK_NAME" "$container_name" || exit 1
done

# TODO Wait for all containers to be ready
read -rp "Wait for GRPC servers to be ready, then Press enter to continue"

# Run k6 client in foreground
docker run -it --rm --name="${DOCKER_NAME}_k6" -e "GRPC_SERVER=$DOCKER_NAME" --label="$DOCKER_LABEL" -v "$(pwd):/mnt" --workdir=/mnt --network="$DOCKER_NETWORK_NAME" "$DOCKER_IMAGE" ./k6 run ./examples/example.js --duration="$K6_DURATION" --vus=$K6_VUS || exit 1
