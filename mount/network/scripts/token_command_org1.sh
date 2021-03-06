#!/bin/bash

c_flag=false

while getopts ":iqp:c:" opt; do
  case ${opt} in
    i )
      COMM_TYPE="invoke"
      ;;
    q )
      COMM_TYPE="query"
      ;;
    p )
      HLF_PATH=$OPTARG
      ;;
    c )
      echo "$OPTARG"
      JSON_COMMAND="$OPTARG"
      c_flag=true
      ;;
    \? )
      echo "Invalid Option: -$OPTARG" 1>&2
      exit 1
      ;;
    : )
      echo "Invalid Option: -$OPTARG requires an argument" 1>&2
      exit 1
      ;;
  esac
done

if ! $c_flag; then
    echo "No command specified"
    exit 1
fi

pushd "$HLF_PATH"/mount/network || exit
export FABRIC_CFG_PATH=$PWD/../config
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_MSPCONFIGPATH=${PWD}/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_TLS_ROOTCERT_FILE=${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_ADDRESS=localhost:7051

ORDERER="localhost:7050"
ORDERER_HOSTNAME="orderer.example.com"
CA_FILE="${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem"

PEER1="localhost:7051"
P1_CERT="${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt"
PEER2="localhost:9051"
P2_CERT="${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
CHANNEL="channel1"
CHAINCODE="token"

if [[ "$COMM_TYPE" == "query" ]]; then
  peer chaincode "$COMM_TYPE" -C "$CHANNEL" -n "$CHAINCODE" --peerAddresses "$PEER1" --tlsRootCertFiles "$P1_CERT" -c "$JSON_COMMAND"
elif [[ "$COMM_TYPE" == "invoke" ]]; then
  peer chaincode "$COMM_TYPE" -o "$ORDERER" --ordererTLSHostnameOverride "$ORDERER_HOSTNAME" --tls --cafile "$CA_FILE" -C "$CHANNEL" -n "$CHAINCODE" --peerAddresses "$PEER1" --tlsRootCertFiles "$P1_CERT" --peerAddresses "$PEER2" --tlsRootCertFiles "$P2_CERT" -c "$JSON_COMMAND"
fi

popd || exit