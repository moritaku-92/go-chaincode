/*
スキルグループ用にサンプルのchaincodeを改造
*/

package main

import (
	"fmt"
	"strconv"
	"bytes"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var logger = shim.NewLogger("skill group cc1")

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {}

// 初期化処理
func(t * SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc1 Init ###########")

	args := stub.GetStringArgs()
	var A, B string // 2名のユーザ
	var Aval, Bval int // 2つの値
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// 渡された値の1・3番目が数字であるか判定
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	B = args[2]
	Bval, err = strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	logger.Info( A, "=", Aval, B, "=", Bval)

	// 実際に値の書込
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// お金のやり取りを記載
func(t * SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### skill group cc1 Invoke ###########")
	/*
		invokeに対し、以下のjsonで投げれば実行される
		["args":["function","value","value", ... ]]
	*/

	// 呼びたいfunctionと渡すParamに別ける
	fn, args := stub.GetFunctionAndParameters()


	// functionの中身を判定し、指定された処理を行う
	// 無ければエラーを返す

	// ユーザの削除
	if fn == "delete" {
		return t.delete(stub, args)
	}

	// 残高の取得
	if fn == "query" {
		return t.query(stub, args)
	}

	// お金の移動
	if fn == "move" {
		return t.move(stub, args)
	}

	// お金の追加
	if fn == "addMoney" {
		return t.addMoney(stub, args)
	}

	// ユーザーを追加する
	if fn == "addUser" {
		return t.addUser(stub, args)
	}

	// range test
	if fn == "rangeTest" {
		return t.rangeTest(stub)
	}

	logger.Errorf("Unknown action, check the first argument, must be one of 'delete', 'query', 'move', 'addMoney', 'adduser' or 'rangeTest'. But got: %v", args[0])
	return shim.Error(fmt.Sprintf("Unknown action, check the first argument, must be one of 'delete', 'query', 'move', 'addMoney', 'adduser' or 'rangeTest'. But got: %v", args[0]))
}



// ユーザ間のお金受け渡し
func(t * SimpleChaincode) move(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### move ###########")
	var A, B string // Entities
	var Aval, Bval int // Asset holdings
	var X int // Transaction value
	var err error

	// 受取っ配列の長さが正しいか判定
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 4, function followed by 2 names and 1 value")
	}

	// ユーザの設定
	A = args[0]
	B = args[1]

	// 指定されたユーザが存在するか判定
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	// 金額がstringで格納されているのでintに変換
	Aval, _ = strconv.Atoi(string(Avalbytes))


	// 指定されたユーザが存在するか判定
	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	// 金額がstringで格納されているのでintに変換
	Bval, _ = strconv.Atoi(string(Bvalbytes))


	// 送金金額が数字（int）であるか判定
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}

	// 取得した値を元に計算
	Aval = Aval - X
	Bval = Bval + X
	logger.Info( A, "=", Aval, B, "=", Bval)


	// 新しい残高を設定
	// 金額をstring化し、さらにbyte配列(?)化しないと格納できない
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// 新しい残高を設定
	// 金額をstring化し、さらにbyte配列(?)化しないと格納できない
	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	// 最後までうまく行けばnilを返す
	return shim.Success(nil);
}

// ユーザの削除
func(t * SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### delete ###########")

	// 受取っ配列の長さが正しいか判定
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// 削除するユーザを設定
	A := args[0]

	// 削除するユーザがいるか確認
	_, err := stub.GetState(A)
	if err != nil {
		return shim.Error("No user name")
	}

	// 実際に削除する
	err = stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	// 正常終了すればnilを返す
	return shim.Success(nil)
}

// 現在値の取得
func(t * SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### query ###########")
	var A string // Entities
	var err error

	// 受取っ配列の長さが正しいか判定
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	// 残高を取得するユーザを指定
	A = args[0]

	// ユーザがいるか確認
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}
	// 残高がない場合はエラー
	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	// 返すためのjsonを作る
	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	logger.Infof("Query Response:%s\n", jsonResp)
	// return shim.Success(Avalbytes)
	return shim.Success([]byte(jsonResp))
}

// ユーザー追加
func(t * SimpleChaincode) addUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### add user ###########")

	var A string
	var Aval int
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	// 引数が正しいか判定
	A = args[0]
	Aval, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}

	// 既存ユーザでないか判定
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes != nil {
		return shim.Error("Aleady User Names")
	}

	// 実際に値の書込
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}
	logger.Info(A, "=", Aval)

	return shim.Success(nil)
}


// 指定したユーザにお金あげる
func(t * SimpleChaincode) addMoney(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Info("########### add money ###########")

	var A string
	var Aval int
	var err error
	var X int

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	// 引数が正しいか判定
	A = args[0]
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	// Perform the execution
	X, err = strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	Aval = Aval + X

	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return shim.Error(err.Error())
	}

	logger.Info(A, "=", Aval)

	return shim.Success(nil)
}


// GetStateByRangeを試してみる （必要あれば実装する）
func(t * SimpleChaincode) rangeTest(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Info("########### range ###########")

	keysIter, err := stub.GetStateByRange("","")
	
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
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
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