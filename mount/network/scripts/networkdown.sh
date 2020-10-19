#!/bin/bash

HLF_PATH=$1

pushd "$HLF_PATH"/mount/network/docker || exit

docker-compose -f ./docker-compose-test-net.yaml down
docker-compose -f ./docker-compose-ca.yaml down
docker container stop couchdb1 couchdb0

docker container prune -f
docker volume prune -f

popd || exit