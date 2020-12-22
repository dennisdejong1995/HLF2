#!/bin/bash

if [[ $# -ne 1 ]] ; then
    echo 'Please supply a path to the network'
    exit 1
fi

HLF_PATH=$1

#######################
## compile chaincode ##
#######################

# go to the chaincode location
pushd "$HLF_PATH"/mount/chaincode/payment-flow/go || exit

# add the dependencies
echo "Adding dependencies"
GO111MODULE=on go mod vendor

popd || exit
pushd /srv/test-networks/HLF2/mount/chaincode/payment-flow/go/vendor/github.com/karalabe/usb/ || exit
mkdir -p hidapi/hidapi
sudo cp /home/dennis/usb/hidapi/hidapi/hidapi.h /srv/test-networks/HLF2/mount/chaincode/payment-flow/go/vendor/github.com/karalabe/usb/hidapi/hidapi


popd || exit

#######################
## package chaincode ##
#######################

# return to the network folder
pushd "$HLF_PATH"/mount/network || exit

# set cfg path
export FABRIC_CFG_PATH=$PWD/../config

# package the chaincode
echo "Packaging chaincode"
peer lifecycle chaincode package payment.tar.gz --path ../chaincode/payment-flow/go --lang golang --label payment_1

#######################
## install chaincode ##
#######################

# Set environment to Org1
echo "Setting environment to Org1"
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH="${PWD}"/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

# Install the chaincode
echo "Installing chaincode for Org1"
peer lifecycle chaincode install payment.tar.gz

# Set environment to Org2
echo "Setting environment to Org2"
export CORE_PEER_LOCALMSPID="Org2MSP"
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH="${PWD}"/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
export CORE_PEER_ADDRESS=localhost:9051

# Install the chaincode
echo "Installing chaincode for Org2"
peer lifecycle chaincode install payment.tar.gz

#######################
## approve chaincode ##
#######################

# Get the package ID of our chaincode
peer lifecycle chaincode queryinstalled >&log.txt

# export the ID as a variable
CC_PACKAGE_ID=$(sed -n "/${CC_NAME}_${CC_VERSION}/{s/^Package ID: //; s/, Label:.*$//; p;}" log.txt)

# approve for Org2
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID channel1 --name payment --version 1.0 --package-id "$CC_PACKAGE_ID" --sequence 1 --tls --cafile "${PWD}"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# Check the approvals
peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name payment --version 1.0 --sequence 1 --tls --cafile "${PWD}"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json

# Set environment to Org1
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH="${PWD}"/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE="${PWD}"/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051

# approve for Org1
peer lifecycle chaincode approveformyorg -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID channel1 --name payment --version 1.0 --package-id "$CC_PACKAGE_ID" --sequence 1 --tls --cafile "${PWD}"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# Check the approvals
peer lifecycle chaincode checkcommitreadiness --channelID channel1 --name payment --version 1.0 --sequence 1 --tls --cafile "${PWD}"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --output json

# Commit the approved chaincode for both of the organizations
peer lifecycle chaincode commit -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --channelID channel1 --name payment --version 1.0 --sequence 1 --tls --cafile "${PWD}"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}"/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}"/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt

# checkout all the chaincode definitions on the channel
peer lifecycle chaincode querycommitted --channelID channel1 --name payment --cafile "${PWD}"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

# initialize the chaincode 
peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}"/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C channel1 -n payment --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}"/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}"/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["Instantiate"]}'
