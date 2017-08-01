/*
スキルグループ用にサンプルのchaincode2を作成
*/

package main

import (
	"fmt"
	"strconv"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

var logger = shim.NewLogger("skill group cc2")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {}

// 依頼内容の構造体を定義
type Mission struct{
	// 依頼者
	Requester string 'json:"requester"'
	// 受領有無
	Acceptance bool 'json:"acceptance"'
	// 任務内容
	MissionContent string 'json:"missionContent"'
	// 報酬
	Compensation int 'json:"compensation"'
	// 受注者
	Contractor string 'json:"contractor"'
	// 完了有無 → 完了したら任務自体削除するか要確認
	Compleate bool 'json:"compleate"'
}

// 共同購入の構造体
type Purchase struct{
	// 依頼者
	// 何買う
	// いくらで？
	// 達成人数
	// 申込者→jsonで持たせたい
	// 申込者数
}

// 初期化処理
func(t * SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc1 Init ###########")

	// カウンタの設定 任務番号を管理する
	err = stub.PutState("count", []byte(strconv.Itoa(0)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// ------------------  test mission------------------
	var mission []Mission
	mission.Requester = "Jane Doe"
	mission.Acceptance = false
	mission.MissionContent = "5000兆円欲しい!!!"
	mission.Compensation = 100000
	mission.Contractor = nil
	mission.Compleate = false

	missionJson, err := json.Marshal(&mission)
	if err != nil {
		return shim.Error("任務のjson化失敗したわー")
	}

	count, err := stub.GetState("count")
	if err != nil {
		return shim.Error(err.Error())
	}
	quest := "quest" + count

	// 任務の登録
	err = stub.PutState(quest, missionJson)
	if err != nil {
		return shim.Error(err.Error())
	}
	// ------------------test mission end------------------

	return shim.Success(nil)

}

// invoke処理 functionによって行う処理を変える
func(t * SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc2 Invoke ###########")

	// 受取るjson
	// ["function", "args[0]", "args[1]", args[2], ...]


	// 受取ったjsonをfunctionと残りのParamに別ける
	// argsは配列
	fn, args := stub.GetFunctionAndParameters()

	// fnの中身を判定しfunctionを実行
	// 任務依頼
	if fn == "request" {
		return t.order(stub, args)
	}

	// 任務削除
	if fn == "delete" {
		return t.delete(stub, args)
	}

	// 任務受注
	if fn == "receive" {
		return t.receive(stub, args)
	}

	// 任務取消
	if fn == "cancel" {
		return t.cancel(stub, args)
	}

	// 任務完了
	if fn == "complete" {
		return t.complete(stub, args)
	}

	// 任務一覧取得
	if fn == "query" {
		return t.query(stub)
	}


	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

// 任務依頼
func(t * SimpleChaincode) request(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### request ###########")
	
	// 受け取るargs
	// ["依頼者", "任務内容", "報酬"]
	// ["John Smith", "5000兆円欲しい!!!", "100000"]

	// 受取ったjsonの長さが正しいか判定
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// 任務のjson作成
	var mission []Mission
	mission.Requester = args[0]
	mission.Acceptance = false
	mission.MissionContent = args[1]
	mission.Compensation = args[2]
	mission.Contractor = nil
	mission.Compleate = false
	// ここでjson化
	missionJson, err := json.Marshal(&mission)
	if err != nil {
		return shim.Error("任務のjson化失敗したわー")
	}

	// 任務番号設定
	count, err := stub.GetState("count")
	if err != nil {
		return shim.Error(err.Error())
	}
	quest := "quest" + count

	// 任務の登録
	err = stub.PutState(quest, missionJson)
	if err != nil {
		return shim.Error(err.Error())
	}

	// countの増加
	_, err = stub.PutState("count", []byte(strconv.Itoa(count+1)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// 返却値は要相談
	return shim.Success(nil)
}

// 任務削除
func(t * SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### delete ###########")
	
	// 受取るjson
	// ["任務番号"]
	// ["quest8"]
	
	// 受取ったjsonの長さが正しいか判定
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// 任務番号がきている想定
	missinoNo := args[0]
	// 任務番号があるか判定
	_, err := shim.GetState(missinoNo)
	if err != nil {
		return shim.Error("その任務番号はないで")
	}

	// 任務削除
	err := stub.DelState(missinoNo)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// 任務受注
func(t * SimpleChaincode) receive(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### receive ###########")

	// 受け取るargs
	// ["依頼番号","受注者"]
	// ["quest8","Mr.Satan"]

	// 受取ったjsonの長さが正しいか判定
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// 依頼番号
	missionNo := args[0]
	// 依頼内容取得
	missionCon, err := stub.GetState(missionNo)
	if err != nil {
		// 番号がなければエラーを返す
		return shim.Error("その任務番号はないで")
	}

	// 依頼者
	receiveUser := args[1]
	// 登録ユーザか判定？→cc1に問合せないといけないか？→実装時間的に断念？

	// 任務を取得し構造体にぶっ込む
	mission := json.Newdecoder(missionCon)
	var missionJson []Mission
	mission.decode(&missionJson)

	// 受注者の登録
	missionJson.Contractor = receiveUser
	missionJson.Acceptance = true

	// jsonエンコード
	outputJson, err := json.Marshal(&missionJson)
	if err != nil {
		return shim.Error("json化できなかった")
	}

	// 登録
	err = stub.PutState(missionNo, []byte(strconv.Itoa(outputJson)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 任務取消
func(t * SimpleChaincode) cancel(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### cancel ###########")

	// 受け取るargs
	// ["依頼番号"]
	// ["quest8"]
	// 受取ったjsonの長さが正しいか判定
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// 依頼番号
	missionNo := args[0]

	// 依頼内容取得
	missionCon, err := stub.GetState(missionNo)
	if err != nil {
		// 番号がなければエラーを返す
		return shim.Error("その任務番号はないで")
	}

	// 任務を取得し構造体にぶっ込む
	mission := json.Newdecoder(missionCon)
	var missionJson []Mission
	mission.decode(&missionJson)

	// 任務の取り消し
	missionJson.Contractor = nil
	missionJson.Acceptance = false

	// jsonエンコード
	outputJson, err := json.Marshal(&missionJson)
	if err != nil {
		return shim.Error("json化できなかった")
	}

	// 登録
	err = stub.PutState(missionNo, []byte(strconv.Itoa(outputJson)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 任務完了
func(t * SimpleChaincode) complete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 受け取るargs
	// ["依頼番号"]
	// ["quest8"]
	// 受取ったjsonの長さが正しいか判定
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	// 依頼番号
	missionNo := args[0]

	// 依頼内容取得
	missionCon, err := stub.GetState(missionNo)
	if err != nil {
		// 番号がなければエラーを返す
		return shim.Error("その任務番号はないで")
	}

	// 任務を取得し構造体にぶっ込む
	mission := json.Newdecoder(missionCon)
	var missionJson []Mission
	mission.decode(&missionJson)

	// 任務完了
	missionJson.Compleate = true

	// jsonエンコード
	outputJson, err := json.Marshal(&missionJson)
	if err != nil {
		return shim.Error("json化できなかった")
	}

	// 登録
	err = stub.PutState(missionNo, []byte(strconv.Itoa(outputJson)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 任務一覧取得
func(t * SimpleChaincode) query(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### query ###########")

	keysIter, err := stub.GetStateByRange("q","q")
	if err != nil {
		return shim.Error(fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()
	
	// ここでjsonをつくる
	bArrayMemberAlreadyWritten := false
	var buffer bytes.Buffer
	buffer.WriteString("[")
	for keysIter.HasNext() {
		queryResponse, err := keysIter.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\":")
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")


	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}