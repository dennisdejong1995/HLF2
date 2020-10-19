#!/bin/bash

if [[ $# -ne 2 ]] ; then
    echo 'Please supply a path to the network and the desired deploy script to run'
    exit 1
fi

HLF_PATH=$1
DEPLOY_SCRIPT=$2


echo "Generating genesis block"
sudo ./generatecryptos.sh "$HLF_PATH"

echo "Putting network up"
sudo ./networkupwithcouchdb.sh "$HLF_PATH"

echo "Creating channel"
sudo ./createchannel.sh "$HLF_PATH"

echo "Deploying chaincode"
sudo ./"$DEPLOY_SCRIPT" "$HLF_PATH"
