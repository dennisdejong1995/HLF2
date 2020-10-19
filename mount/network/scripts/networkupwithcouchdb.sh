#!/bin/bash

if [[ $# -ne 1 ]] ; then
    echo 'Please supply a path to the network'
    exit 1
fi

HLF_PATH=$1

pushd "$HLF_PATH"/mount/network/docker || exit

docker-compose -f ./docker-compose-test-net.yaml -f ./docker-compose-couch.yaml up -d 

popd || exit