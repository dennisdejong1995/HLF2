#!/bin/bash

if [[ $# -ne 1 ]] ; then
    echo 'Please supply a path to the network'
    exit 1
fi

HLF_PATH=$1

pushd "$HLF_PATH"/mount/network/docker

docker-compose -f ./docker-compose-test-net.yaml up -d 

popd