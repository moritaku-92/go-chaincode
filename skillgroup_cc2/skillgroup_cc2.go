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
)

var logger = shim.NewLogger("skill group cc2")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {}

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

func(t * SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc2 Invoke ###########")
	
	fn, args := stub.GetFunctionAndParameters()

	// 依頼を出す
	if fn == "request" {
		return t.order(stub, args)
	}

	// 依頼の削除
	if fn == "delete" {
		return t.delete(stub, args)
	}

	// 依頼を受ける
	if fn == "receive" {
		return t.receive(stub, args)
	}

	// 依頼の取り消し
	if fn == "cancel" {
		return t.cancel(stub, args)
	}

	// 依頼の完了
	if fn == "complete" {
		return t.complete(stub, args)
	}

	// 依頼一覧取得
	if fn == "query" {
		return t.query(stub, args)
	}


	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

// 依頼を出す
func(t * SimpleChaincode) request(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### request ###########")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	count, err := stub.GetState("count")
	if err != nil {
		return shim.Error("依頼番号取れなかった")
	}
	count = count + 1;

	// クエスト番号設定
	quest := "quest" + count

	// 依頼の登録
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

// 依頼を削除する
func(t * SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### delete ###########")
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// 依頼番号が着ている想定{quest{N}}
	A := args[0]

	// 依頼削除
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		logger.Errorf("Error starting Simple chaincode: %s", err)
	}
}
