/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}



// 구조체 정의 부분
type OrderLedger struct {

	OrderID string `json:"orderID"`
	TrnsprtPrdlstCode  string `json:"trnsprtPrdlstCode"`
	
}


/////////////////////////////////////////////////////////////////////////////////////////////////////

// 체인코드 인스턴스화(초기화) 하는 부분(Instantiate/upgrade)으로 모든 데이터를 초기화 한다.
// 필수적으로 초기화를 해야하는 것들이 있다면 init내에 구현하고 없다면 빈 값을 리턴하면 된다.

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	// stub : rpc 프로토콜 버퍼의 메시지 타입을 호출할 때 파라미터 값을 가리키는 객체 변수
	// shim : 하이퍼렛저 패키지로 트랜잭션을 발생시킬 수 있고, 상태값을 조회할 수 있는 기능을 제공

	fmt.Println(" ############ chaincode Init ############")
	return shim.Success(nil)
}

/*
// 기존 init 샘플코드
func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	var A, B string    // Entities
	var Aval, Bval int // Asset holdings
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// Initialize the chaincode
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
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state to the ledger
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
*/



/////////////////////////////////////////////////////////////////////////////////////////////////////


// ChaincodeStubInterface로부터 체인코드 실행 시 넘겨 받은 인수, 즉 실행할 함수이름을 추출한다. 
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("################## b2b Invoke 호출 됨 ##################")

	// 체인코드 api인 shim에서 제공하는 GetFunctionAndParameters 함수를 통해서 넘어온 파라미터(함수명, 매개변수 등)를 받는다.
	// 사용자가 호출한 함수와 매개변수를 받아 function, args에 각각 넣어줌
	function, args := stub.GetFunctionAndParameters()

	fmt.Println("function 저장 ")
	fmt.Println("function : " + function)

	// 호출된 함수명에 맞게 분기처리 하는 부분
	if function == "registerOrder" {
		fmt.Println("######## registerOrder 함수 분기 if문 ################## ")
		return t.registerOrder(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// 운송장 등록
func (t *SmartContract) registerOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//var orderID string
	//var trnsprtPrdlstCode string
	fmt.Println("######## registerOrder 함수 진입 ################## ")


	// 운송장번호, 해시값 필수로 매개변수 2개 아닐 시 오류
	if len(args) != 2 {
		return shim.Error("registerOrder 함수 호출 : Incorrect number of arguments. Expecting 2")
	}

	var order = OrderLedger{OrderID: args[0], TrnsprtPrdlstCode: args[1]}
	//var order = OrderLedger{OrderID: args[1]}

	// json 데이터를 []byte형태로 변경 (마샬링)
	orderAsBytes, _ := json.Marshal(order)

	// 데이터 접근 확인
	//fmt.Println("######## 데이터 테스트 orderAsBytes : " + orderAsBytes)
	fmt.Println("######## 데이터 테스트 OrderID : " + order.OrderID)
	fmt.Println("######## 데이터 테스트 TrnsprtPrdlstCode : " + order.TrnsprtPrdlstCode)


	//stub.PutState(args[0], orderAsBytes)


	// state db에 저장
	err := stub.PutState(order.OrderID, orderAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}




// query callback representing the query of a chaincode
func (t *SmartContract) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var orderID string // Entities
	var err error
	orderID = args[0]


	fmt.Println("######## 쿼리 호출 (orderID :"+ orderID + ")########")


	orderAsBytes, _ := stub.GetState(orderID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + orderID + "\"}"
		return shim.Error(jsonResp)
	}


	jsonResp := "{\"orderID\":\"" + orderID + "\",\"data\":\"" + string(orderAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(orderAsBytes)


/*
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)

*/
}



/////////////////////////////////////////////////////////////////////////////////////////////////////



// Deletes an entity from state
func (t *SmartContract) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}


// 메인
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
