
package main


import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}


type HistLedger struct {

	OrderID string `json:"OrderID"`		// 화물 ID		
	Code string `json:"Code"`		// 화물 상태 코드
	HString string `json:"HString"`		// 검증 해시값
	RegID string `json:"RegID"`		// 등록 ID
	RegTime string `json:"RegTime"`		// 최초 등록일

}


func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println(" ########### histcc chaincode Init ############")
	return shim.Success(nil)
}



func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "regHist" {
		return t.regHist(stub, args)       	// 운송이력 등록
	}else if function == "readHist" {		// 운송이력 조회
		return t.readHist(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} 

	return shim.Error("IF-BLC-601-001| Invalid Smart Contract function name.")
}



// 운송이력 등록
func (t *SmartContract) regHist(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("######## regHist 함수 진입 ################## ")

	// 화물ID 필수값 없을 시 오류
	if args[0] == "" {
		return shim.Error("Error :[registerOrder] 'OrderID' does not exist")
	}

	//현재시간추출 
	time := time.Now()
	timeFmt := time.Format("20060102150405")

	var hist = HistLedger{
		OrderID: args[0], 
		Code: args[1], 
		HString: args[2], 	
		RegID: args[3], 	
		RegTime: timeFmt, 	
	}

	histAsBytes, _ := json.Marshal(hist)

	fmt.Println("######## 데이터 테스트 OrderID : " + hist.OrderID)


	// 이미 등록된 아이디가 있는 경우 에러
	checkOrderExists, err := stub.GetState(hist.OrderID)
	if err != nil {
		return shim.Error("Failed to getState")
	}
	if checkOrderExists != nil {
		return shim.Error("Error |" + hist.OrderID + " already exists.")
	}


	err = stub.PutState(hist.OrderID, histAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}


// 운송이력 단건 조회
func (t *SmartContract) readHist(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var orderID string 
	var err error
	orderID = args[0]


	fmt.Println("######## 쿼리 호출 (orderID :" + orderID + ")########")


	histAsBytes, _ := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Error: " + err.Error())
	}



	// 데이터가 존재하지 않을  경우 메세지 처리
	if histAsBytes == nil {
		return shim.Error("IF-BLC-601-002| Data searched by orderID\"+\"" + orderID + " doesn't exists")
	}


	resultData := "{\"OrderID\":\"" + orderID + "\",\"data\":\"" + string(histAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", resultData)
	return shim.Success(histAsBytes)


}



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

