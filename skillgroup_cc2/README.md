# 任務を作成するchaincode
keyに任務番号、valueに任務内容を保存するchaincodeです。  
任務完了時、cc1のchaincodeに対し依頼者から受領者に成功報酬が支払われます。

## memo

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


## json定義

### invokeに対して投げるjson

* 任務追加      
    {"Args":["request", "UserName", "Mission", "Amount"]}       
    ex) {"Args":["request", "John Smith", "I am hungry", "50"]}

* 任務削除      
    {"Args":["delete", "QuestNumber"]}      
    ex) {"Args":["delete", "quest1"]}

* 任務受領      
    {"Args":["receive", "QuestNumber", "UserName"]}     
    ex) {"Args":["receive", "quest0", "Mr.Satan"]}

* 任務取消      
    {"Args":["cancel", "QuestNumber"]}      
    ex) {"Args":["cancel", "quest0"]}

* 任務完了      
    {"Args":["complete", "QuestNumber"]}        
    ex) {"Args":["complete", "quest0"]}

* 任務一覧取得        
    {"Args":["query"]}
    ex) {"Args":["query"]}

### QuestNumberの採番ルール
QuestNumberは、一意に特定出来るよう付与されている。  
採番ルールは、「quest + 数字」である。


## test code

peer chaincode invoke -n mycc2 -c '{"Args":["request", "nozawa", "I want a coke", "50"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["receive", "quest1", "mori"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["complete", "quest1"]}' -C myc