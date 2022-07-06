
package main


import (
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}



// 운송장 구조체 정의
type OrderLedger struct {

	OrderID string `json:"OrderID"`							
	OrderHash string `json:"OrderHash"`
	FrghtInfo  string `json:"FrghtInfo"`
	TrnsprtPrdlstCode  string `json:"TrnsprtPrdlstCode"`
	OnAdd  string `json:OnAdd"`
	OffAdd  string `json:"OffAdd"`
	LoadDe  string `json:"LoadDe"`
	UnloadDe  string `json:"UnloadDe"`
	MemID  string `json:"MemID"`
	//Tel  string `json:"Tel"`
	Pay  string `json:"Pay"`
	RegTime  string `json:"RegTime"`	

	
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



/////////////////////////////////////////////////////////////////////////////////////////////////////


// ChaincodeStubInterface로부터 체인코드 실행 시 넘겨 받은 인수, 즉 실행할 함수이름을 추출한다. 
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

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
	} else if function == "readOrder" {
		return t.readOrder(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// 운송장 등록
func (t *SmartContract) registerOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//var orderID string
	//var trnsprtPrdlstCode string
	fmt.Println("######## registerOrder 함수 진입 ################## ")


	// 운송장번호 필수값 없을 시 오류
	if args[0] == "" {
		return shim.Error("Error :[registerOrder] 'OrderID' does not exist")
	}


	var order = OrderLedger{
		OrderID: args[0], 
		OrderHash: args[1], 
		FrghtInfo: args[2], 
		TrnsprtPrdlstCode: args[3],  
		OnAdd: args[4],  
		OffAdd: args[5],  
		LoadDe: args[6],  
		UnloadDe: args[7], 
		MemID: args[8], 
		Pay: args[9],
		RegTime: args[10],
	}



	// json 데이터를 []byte형태로 변경 (마샬링)
	// 밑줄문자 : 사용하지 않을 값은 밑줄문자로 대신하여 무시할 수 있음
	orderAsBytes, _ := json.Marshal(order)

	// 데이터 접근 확인
	//fmt.Println("######## 데이터 테스트 orderAsBytes : " + orderAsBytes)
	fmt.Println("######## 데이터 테스트 OrderID : " + order.OrderID)
	fmt.Println("######## 데이터 테스트 TrnsprtPrdlstCode : " + order.TrnsprtPrdlstCode)


	//stub.PutState(args[0], orderAsBytes)


	// state db에 저장
	// := 선언과 동시에 초기화 해줌( 변수의 자료형은 입력된 값에 따라 자동으로 자료형이 결정됨),  : 없을땐 err변수를 따로 선언하지 않아서 못찾는다고 오류남
	err := stub.PutState(order.OrderID, orderAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}




func (t *SmartContract) readOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var orderID string // Entities
	var err error
	orderID = args[0]


	fmt.Println("######## 쿼리 호출 (OrderID :" + orderID + ")########")

	orderAsBytes, _ := stub.GetState(orderID)
	if err != nil {
		resultData := "{\"Error\":\"Failed to get state for " + orderID + "\"}"
		return shim.Error(resultData)
	}


	resultData := "{\"OrderID\":\"" + orderID + "\",\"data\":\"" + string(orderAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", resultData)
	return shim.Success(orderAsBytes)



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
