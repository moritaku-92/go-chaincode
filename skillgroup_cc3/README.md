# 共同購入を管理するchaincode
keyに共同購入番号、valueに内容を保存するchaincodeです。
共同購入達成時、申込者全員から設定された金額が、cc1のchaincodeから引かれます。

## json定義

### invokeに対して投げるjson

* 依頼発行      
    {"Args":["request", "依頼者", "商品名", "出資額", "出資人数"]}       
    ex) {"Args":["request", "a", "Coca Cola", "50", "3"]}

* 出資        
    {"Args":["receive", "依頼番号", "出資者名"]}        
    ex) {"Args":["receive", "groupPurchase1", "b"]}

* 依頼取り消し        
    {"Args":["delete", "依頼番号"]}     
    ex) {"Args":["delete", "groupPurchase0"]}

* 依頼一覧取得        
    {"Args":["query"]}      
    ex) {"Args":["query"]}


### 注意事項

* パラメータは半角英数字のみ許容       
    日本語不可

* 出資の取り消しは許容しない     

* ユーザの挙動は性善説        
    存在しないユーザ名を使うことはない など


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


## response

[{\"number\":\"groupPurchase0\",\"requester\":\"Jane Doe\",\"wish\":\"AYATAKA\",\"price\":10,\"contractores\":[\"Jane Doe\"],\"fund\":2,\"compleate\":false},{\"number\":\"groupPurchase1\",\"requester\":\"a\",\"wish\":\"Coca Cola\",\"price\":50,\"contractores\":[\"a\",\"b\"],\"fund\":3,\"compleate\":false},{\"number\":\"groupPurchase2\",\"requester\":\"c\",\"wish\":\"PEPSI\",\"price\":100,\"contractores\":[\"c\"],\"fund\":2,\"compleate\":false}]