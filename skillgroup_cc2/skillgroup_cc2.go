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
	Requester string 'json:"requester"'
	Acceptance bool 'json:"acceptance"'
	MissionContent string 'json:"missionContent"'
	Compensation int 'json:"compensation"'
	Contractor string 'json:"contractor"'
}

// 初期化処理
func(t * SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc1 Init ###########")

	// カウンタ設定
	err = stub.PutState("count", []byte(strconv.Itoa(0)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

// invoke処理 functionによって行う処理を変える
func(t * SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc2 Invoke ###########")
	
	fn, args := stub.GetFunctionAndParameters()

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
		return t.query(stub, args)
	}


	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

// 任務を出す
func(t * SimpleChaincode) request(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### request ###########")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// 受け取ったjson
	A := args[0]

	/*
		jsonの形どうしよう…？
		正しいjsonか判定してたら大変…
		ccで判定入れるか、RestServerで判定入れるか…?
	*/
	/*
		{
			"依頼者":"a",
			"受領有無":true,
			"任務内容":"○○○買ってこいや",
			"報酬":10000,
			"受注者":"",
			"受注者数":"1",
			"受注者達成人数":"",
		}
	*/

	// 任務番号設定
	count, err := stub.GetState("count")
	if err != nil {
		return shim.Error(err.Error())
	}
	quest := "quest" + count

	// 任務の登録
	err = stub.PutState(count, []byte(strconv.Itoa(A)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// countの増加
	err = stub.PutState(quest, []byte(strconv.Itoa(count)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// 任務削除
func(t * SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### delete ###########")
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// 任務番号がきている想定{quest{N}}
	A := args[0]
	// ユーザ名があるか確認
	_, err := shim.GetState(A)
	if err != nil {
		return shim.Error("No Mission")
	}

	// 任務削除
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// 任務受注
func(t * SimpleChaincode) receive(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	
	// 依頼番号
	missionNo := args[0]
	// 依頼番号があるか判定？

	// 依頼者
	receiveUser := args[1]
	// 登録ユーザか判定？→cc1に問合せないといけないか？

	// 依頼内容取得
	mission, err := stub.GetState(missionNo)
	if err != nil {
		return shim.Error("No Misiion")
	}

	// 登録
	err = stub.PutState(missionNo, []byte(strconv.Itoa(A)))
	if err != nil {
		return shim.Error(err.Error())
	}



	return shim.Success(nil)
}

// 任務取消
func(t * SimpleChaincode) cancel(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

// 任務取消
func(t * SimpleChaincode) complete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

// 任務一覧取得
func(t * SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
