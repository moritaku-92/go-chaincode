cc copy → fabric-samples/chaincode

cd chaincode-docker-devmode

docker-compose -f docker-compose-simple.yaml up

docker exec -it chaincode bash

cd xxxxx

go build

CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc:0 ./skillgroup_cc1

docker exec -it cli bash

peer chaincode install -p chaincodedev/chaincode/go-chaincode/skillgroup_cc1 -n mycc -v 0

peer chaincode instantiate -n mycc -v 0 -c '{"Args":["a","1000","b","20000"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["query", "a"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["addUser", "c","5000"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["addMoney", "a","500"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["move", "a","b","500"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["delete", "c"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["rangeTest"]}' -C myc




docker-compose -f docker-compose-simple.yaml down

docker stop $(docker ps -q)


test code 

peer chaincode invoke -n mycc -c '{"Args":["nozawa", "c","1000"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["mori", "c","2000"]}' -C myc