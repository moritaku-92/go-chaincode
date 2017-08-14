cc copy → fabric-samples/chaincode

cd chaincode-docker-devmode

docker-compose -f docker-compose-simple.yaml up

docker exec -it chaincode bash

cd go-chaincode/skillgroup_cc2/

go build

CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc2:0 ./skillgroup_cc2

docker exec -it cli bash

go get -v github.com/hyperledger/fabric/common/util

peer chaincode install -p chaincodedev/chaincode/go-chaincode/skillgroup_cc2 -n mycc2 -v 0

peer chaincode instantiate -n mycc2 -v 0 -c '{"Args":[]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["request", "John Smith", "I am hungry", "50"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["delete", "quest1"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["receive", "quest0", "Mr.Satan"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["cancel", "quest0"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["complete", "quest0"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["query"]}' -C myc




docker-compose -f docker-compose-simple.yaml down

docker stop $(docker ps -q)

メモ      
登録する内容は日本語不可



test code

peer chaincode invoke -n mycc2 -c '{"Args":["request", "nozawa", "I want a coke", "50"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["receive", "quest1", "mori"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["complete", "quest1"]}' -C myc