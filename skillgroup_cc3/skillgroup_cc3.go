/*
スキルグループ用のchaincode
共同購入の提案を出し、出資者が規定数を超えた際に
出資者から出資金が減額される
cc1とともに使用する
*/

package main

import (
	"fmt"
	"strconv"
	"bytes"
	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"encoding/json"
)

var logger = shim.NewLogger("skill group cc3")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {}

// 依頼内容の構造体を定義
type Purchase struct{
	// 依頼番号
	Number string `json:"number"`
	// 依頼者
	Requester string `json:"requester"`
	// 欲しい物
	Wish string `json:"wish"`
	// 価格
	Price int `json:"price"`
	// 受注者
	Contractores []string `json:"contractores"`
	// 達成人数
	Fund int `json:"fund"`
	// 完了有無
	Complete bool `json:"complete"`
}

// 初期化処理
func(t * SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc3 Init ###########")

	// カウンタの設定 依頼番号を管理する
	// テスト依頼を入れる関係でcountは1から始めている
	err := stub.PutState("count", []byte(strconv.Itoa(1)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// テストデータの作成
	var purchase = Purchase{}
	purchase.Number = "groupPurchase0"
	purchase.Requester = "Dr. Emmett Brown"
	purchase.Wish = "I need plutonium for the time machine"
	purchase.Price = 100000
	purchase.Contractores = append(purchase.Contractores, "Marty McFly")
	purchase.Fund = 1
	purchase.Complete = true

	purchaseJSON, err := json.Marshal(&purchase)
	if err != nil {
		return shim.Error("json化失敗したわー")
	}

	// 依頼の登録
	err = stub.PutState("groupPurchase0", purchaseJSON)
	if err != nil {
		return shim.Error(err.Error())
	}
	
	return shim.Success(nil)
}

// invoke処理 functionによって行う処理を変える
func(t * SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc3 Invoke ###########")

	// 受取るjson
	// ["function", "args[0]", "args[1]", args[2], ...]

	// 受取ったjsonをfunctionと残りのParamに別ける
	// argsは配列
	fn, args := stub.GetFunctionAndParameters()

	// fnの中身を判定しfunctionを実行
	// 依頼依頼
	if fn == "request" {
		return t.request(stub, args)
	}

	// 依頼削除
	if fn == "delete" {
		return t.delete(stub, args)
	}

	// 依頼受注
	if fn == "receive" {
		return t.receive(stub, args)
	}

	// 依頼一覧取得
	if fn == "query" {
		return t.query(stub)
	}


	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'. But got: %v", args[0]))
}

// 依頼
func(t * SimpleChaincode) request(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### request ###########")
	
	// 受け取るargs
	// ["依頼者", "欲しい物", "価格", "達成人数"]
	// ["John Smith", "COCA COLA", "10", "5"]

	// 受取ったjsonの長さが正しいか判定
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	// 文字列をint化する
	price, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("int型じゃない")
	}
	fund, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("int型じゃない")
	}
	// 共同購入なので2人以上を想定
	if fund <= 1 {
		return shim.Error("達成人数は2人以上")
	}

	// 依頼番号設定
	count, err := stub.GetState("count")
	if err != nil {
		return shim.Error(err.Error())
	}
	groupPurchaseNo := "groupPurchase" + string(count)

	// 依頼のjson作成
	var purchase = Purchase{}
	purchase.Number = groupPurchaseNo
	purchase.Requester = args[0]
	purchase.Wish = args[1]
	purchase.Price = price
	purchase.Contractores = nil
	purchase.Fund = fund
	purchase.Complete = false

	// ここでjson化
	purchaseJSON, err := json.Marshal(&purchase)
	if err != nil {
		return shim.Error("json化失敗したわー")
	}

	// 依頼の登録
	err = stub.PutState(groupPurchaseNo, purchaseJSON)
	if err != nil {
		return shim.Error(err.Error())
	}

	// countの増加
	countInt, err := strconv.Atoi(string(count))
	countInt = countInt + 1
	err = stub.PutState("count", []byte(strconv.Itoa(countInt)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// 返却値は要相談
	return shim.Success(nil)
}

// 依頼削除
func(t * SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### delete ###########")
	
	// 受取るjson
	// ["依頼番号"]
	// ["groupPurchase0"]
	
	// 受取ったjsonの長さが正しいか判定
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// 依頼番号がきている想定
	groupPurchaseNo := args[0]
	// 依頼番号があるか判定
	_, err := stub.GetState(groupPurchaseNo)
	if err != nil {
		return shim.Error("その番号はないで")
	}

	// 依頼削除
	err = stub.DelState(groupPurchaseNo)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

// 依頼受注
func(t * SimpleChaincode) receive(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### receive ###########")

	// 受け取るargs
	// ["依頼番号","受注者"]
	// ["groupPurchase0","mori"]

	// 受取ったjsonの長さが正しいか判定
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// 依頼番号
	groupPurchaseNo := args[0]
	// 依頼内容取得
	groupPurchase, err := stub.GetState(groupPurchaseNo)
	if err != nil {
		// 番号がなければエラーを返す
		return shim.Error(err.Error())
	}

	// 受注者
	receiveUser := args[1]
	// 登録ユーザか判定？→cc1に問合せないといけないか？→実装時間的に断念・性善説で？

	// 依頼を取得し構造体にぶっ込む
	var purchase = Purchase{}
	err = json.Unmarshal(groupPurchase, &purchase)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 受注者の登録
	purchase.Contractores = append(purchase.Contractores, receiveUser)
	
	// ---------------達成判定---------------
	// 判定文
	if purchase.Fund == len(purchase.Contractores) {
		// 引く値なのでマイナス化
		dif := purchase.Price*(-1)
		// for文（受注者数分回す）
		for _, user := range purchase.Contractores {
			// 価格分引く処理
			invokeArgs := util.ToChaincodeArgs("addMoney", user, strconv.Itoa(dif))
			response := stub.InvokeChaincode("skillgroup_cc1", invokeArgs, "mychannel")

			if response.Status != shim.OK {
				errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
				fmt.Printf(errStr)
				return shim.Error(errStr)
			}
		}
		// 提案者にお金渡す処理
		add := purchase.Price * purchase.Fund
		invokeArgs := util.ToChaincodeArgs("addMoney", purchase.Requester, strconv.Itoa(add))
		response := stub.InvokeChaincode("skillgroup_cc1", invokeArgs, "mychannel")
		
		if response.Status != shim.OK {
			errStr := fmt.Sprintf("Failed to invoke chaincode. Got error: %s", string(response.Payload))
			fmt.Printf(errStr)
			return shim.Error(errStr)
		}
		// 登録内容を完了にする
		purchase.Complete = true
	}
	// ---------------達成判定---------------

	// jsonエンコード
	outputJSON, err := json.Marshal(&purchase)
	if err != nil {
		return shim.Error(err.Error())
	}

	// 登録
	err = stub.PutState(groupPurchaseNo, []byte(outputJSON))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


// 依頼一覧取得
func(t * SimpleChaincode) query(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### query ###########")

	// 受け取るargs
	// []

	// これが上手いこと指定したkey・valueを取ってきてくれる
	keysIter, err := stub.GetStateByRange("g","h")
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
		
		buffer.WriteString(string(queryResponse.Value))
		
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
