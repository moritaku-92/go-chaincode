# 任務を作成するchaincode
keyに任務番号、valueに任務内容を保存するchaincodeです。  
任務完了時、cc1のchaincodeに対し依頼者から受領者に成功報酬が支払われます。


## json定義

### invokeに対して投げるjson

* 任務追加      
    {"Args":["request", "依頼者", "任務内容", "報酬"]}       
    ex) {"Args":["request", "John Smith", "I am hungry", "50"]}

* 任務削除      
    {"Args":["delete", "任務番号"]}      
    ex) {"Args":["delete", "quest1"]}

* 任務受領      
    {"Args":["receive", "任務番号", "受領者"]}     
    ex) {"Args":["receive", "quest0", "Mr.Satan"]}

* 任務取消      
    {"Args":["cancel", "任務番号"]}      
    ex) {"Args":["cancel", "quest0"]}

* 任務完了      
    {"Args":["complete", "任務番号"]}        
    ex) {"Args":["complete", "quest0"]}

* 任務一覧取得        
    {"Args":["query"]}      
    ex) {"Args":["query"]}


### 注意事項

* パラメータは半角英数字のみ許容   
    日本語不可

* QuestNumberの採番ルール     
    QuestNumberは、一意に特定出来るよう付与されている。  
    採番ルールは、「quest + 数字」である。

* ユーザの挙動は性善説        
    存在しないユーザ名を使うことはない など


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


## test code

peer chaincode invoke -n mycc2 -c '{"Args":["request", "a", "I want a coke", "50"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["receive", "quest1", "b"]}' -C myc

peer chaincode invoke -n mycc2 -c '{"Args":["complete", "quest1"]}' -C myc

## response

[{\"number\":\"quest0\",\"requester\":\"Jane Doe\",\"acceptance\":false,\"missionContent\":\"I want 5000 trillion yen!\",\"compensation\":100000,\"contractor\":\"\",\"Complete\":false},{\"number\":\"quest1\",\"requester\":\"a\",\"acceptance\":true,\"missionContent\":\"I want a coke\",\"compensation\":50,\"contractor\":\"b\",\"Complete\":true},{\"number\":\"quest2\",\"requester\":\"c\",\"acceptance\":false,\"missionContent\":\"i want 5000 trillion yen\",\"compensation\":50,\"contractor\":\"\",\"Complete\":false}]