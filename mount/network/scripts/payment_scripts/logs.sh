#!/bin/bash

if [[ $# -ne 1 ]] ; then
    echo 'Please supply a path to the network and the desired deploy script to run'
    exit 1
fi

docker logs -f $(docker ps | grep dev-peer0."$1" | awk '{print $11}')