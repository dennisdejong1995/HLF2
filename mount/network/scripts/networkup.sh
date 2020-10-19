#!/bin/bash

HLF_PATH=$1

pushd "$HLF_PATH"/mount/network/docker

docker-compose -f ./docker-compose-test-net.yaml up -d 

popd