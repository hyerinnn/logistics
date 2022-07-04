#/bin/bash



docker ps -a
docker-compose -f docker-compose.yaml down
docker-compose -f docker-compose-cli.yaml down
docker rm $(docker ps -aq)
docker-compose -p net down --volumes --remove-orphans
docker volume prune -f
docker ps -a

sudo rm -rf ./crypto-config

sudo rm -rf ./channel-artifacts/*

sudo rm -rf ./production/orderer.orderer.com/*

sudo rm -rf ./production/peer0.b2bOrg.logistics/*

sudo rm -rf ./production/peer1.b2bOrg.logistics/*

sudo rm -rf ./production/peer2.b2bOrg.logistics/*
