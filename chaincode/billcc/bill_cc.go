
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


type BillLedger struct {

	OrderID string `json:"OrderID"`				 // 화물ID		
	OMemID string `json:"OMemID"`				 // 화주 ID
	CMemID  string `json:"CMemID"`				 // 차주 ID
	Load string `json:"Load"`					 // 상차지
	LoadMN string `json:"LoadMN"`				 // 상차 담당자
	Off string `json:"Off"`						 // 하차지
	OffMN string `json:"OffMN"`					 // 하차 담당자
	ReCompanyNo string `json:"ReCompanyNo"`		 // 수취인 사업자 등록번호
	ReCompanyNM string `json:"ReCompanyNM"`		 // 수취인 상호명
	ReUserNM string `json:"ReUserNM"`			 // 수화인 성명
	ReCompanyAdd string `json:"ReCompanyAdd"`	 // 수화인 주소
	RegTime string `json:"RegTime"`				 // 최초 등록일
	UpdateTime string `json:"UpdateTime"`		 // 수정일

}


func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println(" ########### billcc chaincode Init ############")
	return shim.Success(nil)
}



func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "registerMKBill" {
		return t.registerMKBill(stub, args)       // 배송정보 등록
	}else if function == "readMKBill" {			  // 배송정보 조회
		return t.readMKBill(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} 

	return shim.Error("IF-BLC-301-004| Invalid Smart Contract function name.")
}



// 배송정보 등록
func (t *SmartContract) registerMKBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("######## registerMKBill 함수 진입 ################## ")

	// 운송장번호 필수값 없을 시 오류
	if args[0] == "" {
		return shim.Error("Error :[registerOrder] 'OrderID' does not exist")
	}

	//현재시간추출 
	time := time.Now()
	timeFmt := time.Format("20060102150405")

	var bill = BillLedger{
		OrderID: args[0], 
		OMemID: args[1], 
		CMemID: args[2], 	
		Load: args[3], 	
		LoadMN: args[4], 
		Off: args[5], 
		OffMN: args[6], 
		ReCompanyNo: args[7], 
		ReCompanyNM: args[8], 
		ReUserNM: args[9], 
		ReCompanyAdd: args[10], 
		RegTime: timeFmt, 
		UpdateTime: timeFmt, 	
	}

	billAsBytes, _ := json.Marshal(bill)

	fmt.Println("######## 데이터 테스트 OrderID : " + bill.OrderID)
	fmt.Println("######## 데이터 테스트 LoadMN : " + bill.LoadMN)


	// 이미 등록된 아이디가 있는 경우 에러
	checkOrderExists, err := stub.GetState(bill.OrderID)
	if err != nil {
		return shim.Error("Failed to getState")
	}
	if checkOrderExists != nil {
		return shim.Error("Error |" + bill.OrderID + " already exists.")
	}


	err = stub.PutState(bill.OrderID, billAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}

// 배송정보 단건 조회
func (t *SmartContract) readMKBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var orderID string 
	var err error
	orderID = args[0]


	fmt.Println("######## 쿼리 호출 (orderID :" + orderID + ")########")


	billAsBytes, _ := stub.GetState(orderID)
	if err != nil {
		return shim.Error("Error: " + err.Error())
	}



	// 데이터가 존재하지 않을  경우 메세지 처리
	if billAsBytes == nil {
		return shim.Error("IF-BLC-301-002| Data searched by orderID\"+\"" + orderID + " doesn't exists")
	}


	resultData := "{\"OrderID\":\"" + orderID + "\",\"data\":\"" + string(billAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", resultData)
	return shim.Success(billAsBytes)


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
