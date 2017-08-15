# memo

cc copy → fabric-samples/chaincode

cd chaincode-docker-devmode

docker-compose -f docker-compose-simple.yaml up

docker exec -it chaincode bash

cd go-chaincode/skillgroup_cc2/

go build

CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc3:0 ./skillgroup_cc3

docker exec -it cli bash

go get -v github.com/hyperledger/fabric/common/util

peer chaincode install -p chaincodedev/chaincode/go-chaincode/skillgroup_cc3 -n mycc3 -v 0

peer chaincode instantiate -n mycc3 -v 0 -c '{"Args":[]}' -C myc

peer chaincode invoke -n mycc3 -c '{"Args":["request", "John Smith", "Coca Cola", "50", "3"]}' -C myc

peer chaincode invoke -n mycc3 -c '{"Args":["delete", "groupPurchase0"]}' -C myc

peer chaincode invoke -n mycc3 -c '{"Args":["receive", "groupPurchase0", "John Smith"]}' -C myc

peer chaincode invoke -n mycc3 -c '{"Args":["query"]}' -C myc

docker-compose -f docker-compose-simple.yaml down

docker stop $(docker ps -q)

メモ      
登録する内容は日本語不可



# test code

peer chaincode invoke -n mycc3 -c '{"Args":["request", "a", "Coca Cola", "50", "3"]}' -C myc

peer chaincode invoke -n mycc3 -c '{"Args":["receive", "groupPurchase1", "b"]}' -C myc

peer chaincode invoke -n mycc3 -c '{"Args":["receive", "groupPurchase1", "c"]}' -C myc