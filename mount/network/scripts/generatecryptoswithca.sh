#!/bin/bash

# remove the old materials
sudo rm -fr ~/mount/network/organizations/ordererOrganizations/*
sudo rm -fr ~/mount/network/organizations/peerOrganizations/*
sudo rm -fr ~/mount/network/system-genesis-block/*


# deploy the ca's
pushd ~/mount/network/docker
docker-compose -f docker-compose-ca.yaml up -d
popd

sleep 10

# generate crypto materials
pushd ~/mount/network
./organizations/fabric-ca/registerEnroll.sh

sleep 10

# create connection profile
./organizations/ccp-generate.sh

# copy connection profile to application
cp "./organizations/peerOrganizations/org1.example.com/connection-org1.yaml" "../commercial-paper/organization/digibank/gateway/"
cp "./organizations/peerOrganizations/org2.example.com/connection-org2.yaml" "../commercial-paper/organization/magnetocorp/gateway/"

# rename keys
cp ./organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/* ./organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem
cp ./organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/* ./organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk

cp ./organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/signcerts/* ./organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/signcerts/User1@org2.example.com-cert.pem
cp ./organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/keystore/* ./organizations/peerOrganizations/org2.example.com/users/User1@org2.example.com/msp/keystore/priv_sk


# set the cfg path
export FABRIC_CFG_PATH=$PWD/configtx/

# create the genesis block
configtxgen -profile TwoOrgsOrdererGenesis -channelID system-channel -outputBlock ./system-genesis-block/genesis.block 
popd
