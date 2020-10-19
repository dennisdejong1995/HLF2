#!/bin/bash

HLF_PATH=$1
DEPLOY_SCRIPT=$2

echo "Generating genesis block"
./generatecryptos.sh "$HLF_PATH"

echo "Putting network up"
./networkup.sh "$HLF_PATH"

echo "Creating channel"
./createchannel.sh "$HLF_PATH"

echo "Deploying chaincode"
./"$DEPLOY_SCRIPT" "$HLF_PATH"
