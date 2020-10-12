#!/bin/bash

pushd ~/mount/network/docker

docker-compose -f ./docker-compose-test-net.yaml down
docker-compose -f ./docker-compose-ca.yaml down
docker container stop couchdb1 couchdb0

docker container prune -f
docker volume prune -f

popd