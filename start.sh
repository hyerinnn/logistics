#/bin/bash


# 1. 인증서 생성
./bin/cryptogen generate --config=./crypto-config.yaml

# 2. 제네시스 블록 생성
./bin/configtxgen -profile B2BOrgOrdererGenesis -channelID b2b-sys-channel -outputBlock ./channel-artifacts/genesis.block

# 3. 채널tx 배포(채널 정보를 담은 트랜잭션 전파)
./bin/configtxgen -profile B2BOrgChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID b2bchannel

# 4. 채널에 앵커피어 정보를 담은 트랜잭션 생성( b2bOrgMSP.tx 생성확인)
./bin/configtxgen -profile B2BOrgChannel -outputAnchorPeersUpdate ./channel-artifacts/b2bOrgMSP.tx -channelID b2bchannel -asOrg b2bOrg

# 5. orderer, peer, cli 컨테이너 up
docker-compose -f docker-compose.yaml -f docker-compose-cli.yaml up -d

# 6. 채널 생성
docker exec cli peer channel create -o orderer.orderer.com:7050 -c b2bchannel -f ./channel-artifacts/channel.tx --tls false --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orderer.com/orderers/orderer.orderer.com/msp/tlscacerts/tlsca.orderer.com-cert.pem

## 7. 피어별 채널 조인

# peer0 채널 조인
docker exec cli peer channel join -b  b2bchannel.block

# peer1 채널 조인
docker exec  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/b2bOrg.logistics/peers/peer1.b2bOrg.logistics/tls/ca.crt -e  CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/b2bOrg.logistics/peers/peer1.b2bOrg.logistics/tls/server.crt -e CORE_PEER_ADDRESS=peer1.b2bOrg.logistics:8051  cli peer channel join -b  b2bchannel.block

# peer2 채널 조인
docker exec  -e CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/b2bOrg.logistics/peers/peer2.b2bOrg.logistics/tls/ca.crt -e  CORE_PEER_TLS_CERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/b2bOrg.logistics/peers/peer2.b2bOrg.logistics/tls/server.crt -e CORE_PEER_ADDRESS=peer2.b2bOrg.logistics:9051  cli peer channel join -b  b2bchannel.block


# 8. 앵커피어 정보를(4번에서 만든 트랜잭션) 채널에 업데이트
docker exec cli  peer channel update -o orderer.orderer.com:7050 -c b2bchannel -f ./channel-artifacts/b2bOrgMSP.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orderer.com/orderers/orderer.orderer.com/msp/tlscacerts/tlsca.orderer.com-cert.pem


# 9. peer0에 체인코드 설치
docker exec cli  peer chaincode install -n ordercc -v 1.0 -l golang -p github.com/chaincode/ordercc
docker exec cli  peer chaincode install -n billcc -v 1.0 -l golang -p github.com/chaincode/billcc

# 10. 체인코드 초기화
docker exec cli  peer chaincode instantiate -o orderer.orderer.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orderer.com/orderers/orderer.orderer.com/msp/tlscacerts/tlsca.orderer.com-cert.pem -C b2bchannel -n ordercc -v 1.0 -c '{"Args":[""]}'

docker exec cli  peer chaincode instantiate -o orderer.orderer.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/orderer.com/orderers/orderer.orderer.com/msp/tlscacerts/tlsca.orderer.com-cert.pem -C b2bchannel -n billcc -v 1.0 -c '{"Args":[""]}'




