cc copy → fabric-samples/chaincode

cd chaincode-docker-devmode

docker-compose -f docker-compose-simple.yaml up

docker exec -it chaincode bash

cd go-chaincode/skillgroup_cc2/

go build

CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc2:0 ./skillgroup_cc2

docker exec -it cli bash

peer chaincode install -p chaincodedev/chaincode/go-chaincode/skillgroup_cc2 -n mycc2 -v 0

peer chaincode instantiate -n mycc2 -v 0 -c '{"Args":[]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["request", "John Smith", "I want 5000 trillion yen!", "100000"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["delete", "c","5000"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["receive", "quest8", "Mr.Satan"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["cancel", "a","b","500"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["complete", "c"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["query"]}' -C myc




docker-compose -f docker-compose-simple.yaml down

docker stop $(docker ps -q)