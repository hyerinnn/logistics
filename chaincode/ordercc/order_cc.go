
package main


import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)


type SmartContract struct {
}


// 운송장 구조체 정의
type OrderLedger struct {

	OrderID string `json:"OrderID"`							 // 화물(운송장) ID
	OrderHash string `json:"OrderHash"`						 // 화물정보 hash
	FrghtInfo  string `json:"FrghtInfo"`					 // 화물정보
	TrnsprtPrdlstCode  string `json:"TrnsprtPrdlstCode"`	 // 운송품목
	OnAdd  string `json:OnAdd"`								 // 상차지 주소
	OffAdd  string `json:"OffAdd"`							 // 하차지 주소
	LoadDe  string `json:"LoadDe"`						     // 상차일시
	UnloadDe  string `json:"UnloadDe"`						 // 하차일시
	MemID  string `json:"MemID"`							 // 화주 ID
	//Tel  string `json:"Tel"`								 // 화주 연락처
	Pay  string `json:"Pay"`								 // 화물 운송료
	RegTime  string `json:"RegTime"`						 // 최초 등록일

}



// 체인코드 인스턴스화(초기화) 하는 부분(Instantiate/upgrade)으로 모든 데이터를 초기화 한다.
// 필수적으로 초기화를 해야하는 것들이 있다면 init내에 구현하고 없다면 빈 값을 리턴하면 된다.

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	// stub : rpc 프로토콜 버퍼의 메시지 타입을 호출할 때 파라미터 값을 가리키는 객체 변수
	// shim : 하이퍼렛저 패키지로 트랜잭션을 발생시킬 수 있고, 상태값을 조회할 수 있는 기능을 제공
	fmt.Println(" ############ orderercc chaincode Init ############")
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

	// 체인코드 api인 shim에서 제공하는 GetFunctionAndParameters 함수를 통해서 넘어온 파라미터(함수명, 매개변수 등)를 받는다.
	// 사용자가 호출한 함수와 매개변수를 받아 function, args에 각각 넣어줌
	function, args := stub.GetFunctionAndParameters()

	/*
	* 데이터 테스트용 로그
	*/
	fmt.Println("function 저장 ")
	fmt.Println("function : " + function)


	// 호출된 함수명에 맞게 분기처리 하는 부분
	if function == "registerOrder" {
		return t.registerOrder(stub, args)		// 운송장 등록
	} else if function == "readOrder" {
		return t.readOrder(stub, args)			// 운송장 조회
	} else if function == "delete" {
		return t.delete(stub, args)
	}
	return shim.Error("IF-BLC-301-004| Invalid Smart Contract function name.")
}

// 운송장 등록
func (t *SmartContract) registerOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("######## registerOrder 함수 진입 ################## ")


	// 운송장번호 필수값 없을 시 오류
	if args[0] == "" {
		return shim.Error("Error :[registerOrder] 'OrderID' does not exist")
	}

	//현재시간추출 
	time := time.Now()
	timeFmt := time.Format("20060102150405")

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
		RegTime: timeFmt,
	}


	// 마샬링: 구조체 데이터(또는 map)를 json 포맷으로 인코딩
	// 밑줄문자 : 사용하지 않을 값은 밑줄문자로 대신하여 무시할 수 있음
	// json으로 인코딩된 바이트배열과 에러객체를 리턴하는데, 에러객체는 밀줄문자로 무시함 
	orderAsBytes, _ := json.Marshal(order)

	// 데이터 접근 확인
	fmt.Println("######## 데이터 테스트 OrderID : " + order.OrderID)
	fmt.Println("######## 데이터 테스트 TrnsprtPrdlstCode : " + order.TrnsprtPrdlstCode)


	
	// 이미 등록된 아이디가 있는지 조회
	// ChaincodeStubInterface.GetState 함수 사용
	checkOrderExists, err := stub.GetState(order.OrderID)
	if err != nil {
		return shim.Error("Failed to getState")
	}
	// 이미 등록된 아이디가 있는 경우 에러
	if checkOrderExists != nil {
		return shim.Error("Error |" + order.OrderID + " already exists.")
	}


	// state db에 저장
	// := 선언과 동시에 초기화 해줌( 변수의 자료형은 입력된 값에 따라 자동으로 자료형이 결정됨),  : 없을땐 err변수를 따로 선언하지 않아서 못찾는다고 오류남
	err = stub.PutState(order.OrderID, orderAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}



// 운송장 단건 조회
func (t *SmartContract) readOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	orderID := args[0]


	fmt.Println("######## 쿼리 호출 (OrderID :" + orderID + ")########")

	// 파라미터 체크 (필수값 없을 시, 에러)
	/*
	if len(args) <= 1 {
		return shim.Error("Error|Incorrect number of arguments")
	}
	*/

	orderAsBytes, _ := stub.GetState(orderID)
	if err != nil {
		//resultData := "{\"Error\":\"Failed to get state for " + orderID + "\"}"
		//return shim.Error(resultData)
		return shim.Error("Error: " + err.Error())
	}

	// 데이터가 존재하지 않을  경우 메세지 처리
	if orderAsBytes == nil {
		return shim.Error("IF-BLC-301-002| Data searched by orderID\"+\"" + orderID + " doesn't exists")
	}


	resultData := "{\"OrderID\":\"" + orderID + "\",\"data\":\"" + string(orderAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", resultData)
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
