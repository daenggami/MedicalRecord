#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error, print all commands.r
set -ev

docker-compose -f docker-compose.yml down

docker-compose -f docker-compose.yml up -d ca.medicalnet.com orderer.medicalnet.com couchdb1 couchdb2 couchdb3 peer0.org1.medicalnet.com  peer0.org2.medicalnet.com peer0.org3.medicalnet.com cli
docker ps -a

# wait for Hyperledger Fabric to start
# incase of errors when running later commands, issue export FABRIC_START_TIMEOUT=<larger number>
export FABRIC_START_TIMEOUT=10
#echo ${FABRIC_START_TIMEOUT}
sleep ${FABRIC_START_TIMEOUT}

# Create the channel
docker exec cli peer channel create -o orderer.medicalnet.com:7050 -c medicalchannel -f /etc/hyperledger/configtx/channel.tx

#Join peer0.org1.medicalnet.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org1.medicalnet.com/msp" peer0.org1.medicalnet.com peer channel join -b /etc/hyperledger/configtx/medicalchannel.block
sleep 5
# Join peer0.org2.medicalnet.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org2MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org2.medicalnet.com/msp" peer0.org2.medicalnet.com peer channel join -b /etc/hyperledger/configtx/medicalchannel.block
sleep 5
# Join peer0.org3.medicalnet.com to the channel.
docker exec -e "CORE_PEER_LOCALMSPID=Org3MSP" -e "CORE_PEER_MSPCONFIGPATH=/etc/hyperledger/msp/users/Admin@org3.medicalnet.com/msp" peer0.org3.medicalnet.com peer channel join -b /etc/hyperledger/configtx/medicalchannel.block
sleep 5
