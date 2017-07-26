CORE_PEER_ADDRESS=peer:7051 CORE_CHAINCODE_ID_NAME=mycc:0 ./go-chaincode

peer chaincode install -p chaincodedev/chaincode/go-chaincode -n mycc -v 0

peer chaincode instantiate -n mycc -v 0 -c '{"Args":["a","1000","b","20000"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["query", "a"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["addUser", "c","5000"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["addMoney", "a","500"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["move", "a","b","500"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["delete", "c"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["rangeTest"]}' -C myc

