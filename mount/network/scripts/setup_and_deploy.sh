#!/bin/bash

HLF_PATH=$1
DEPLOY_SCRIPT=$2

echo "Generating genesis block"
sudo ./generatecryptos.sh "$HLF_PATH"

echo "Putting network up"
sudo ./networkup.sh "$HLF_PATH"

echo "Creating channel"
sudo ./createchannel.sh "$HLF_PATH"

echo "Deploying chaincode"
sudo ./"$DEPLOY_SCRIPT" "$HLF_PATH"
