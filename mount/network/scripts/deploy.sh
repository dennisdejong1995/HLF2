#!/bin/bash

if [[ $# -ne 2 ]] ; then
    echo 'Please supply a path to the network and the desired deploy script to run'
    exit 1
fi

HLF_PATH=$1
DEPLOY_SCRIPT=$2
echo "Pulling repository"
sudo git fetch && git stash && git pull

echo "Deploying chaincode"
sudo ./"$DEPLOY_SCRIPT" "$HLF_PATH"
