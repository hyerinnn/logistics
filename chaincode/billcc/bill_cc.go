
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


type BillLedger struct {

	OrderID string `json:"OrderID"`							
	OMemID string `json:"OMemID"`
	CMemID  string `json:"CMemID"`	
	Load string `json:"Load"`	
	LoadMN string `json:"LoadMN"`	
	Off string `json:"Off"`	
	OffMN string `json:"OffMN"`	
	ReCompanyNo string `json:"ReCompanyNo"`	
	ReCompanyNM string `json:"ReCompanyNM"`	
	ReUserNM string `json:"ReUserNM"`	
	ReCompanyAdd string `json:"ReCompanyAdd"`	
	RegTime string `json:"RegTime"`	
	UpdateTime string `json:"UpdateTime"`	

}


func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println(" ########### bill_cc chaincode Init ############")
	return shim.Success(nil)
}


func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if function == "registerMKBill" {
		return t.registerMKBill(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "readMKBill" {
		return t.readMKBill(stub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}



func (t *SmartContract) registerMKBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if args[0] == "" {
		return shim.Error("Error :[registerOrder] 'OrderID' does not exist")
	}

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
		RegTime: args[11], 
		UpdateTime: args[12], 	
	}

	billAsBytes, _ := json.Marshal(bill)


	fmt.Println("######## 데이터 테스트 OrderID : " + bill.OrderID)
	fmt.Println("######## 데이터 테스트 LoadMN : " + bill.LoadMN)

	err := stub.PutState(bill.OrderID, billAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}


	return shim.Success(nil)
}


func (t *SmartContract) readMKBill(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var orderID string 
	var err error
	orderID = args[0]


	fmt.Println("######## 쿼리 호출 (orderID :" + orderID + ")########")


	orderAsBytes, _ := stub.GetState(orderID)
	if err != nil {
		resultData := "{\"Error\":\"Failed to get state for " + orderID + "\"}"
		return shim.Error(resultData)
	}


	resultData := "{\"orderID\":\"" + orderID + "\",\"data\":\"" + string(orderAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", resultData)
	return shim.Success(orderAsBytes)


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
