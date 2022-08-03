
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


type ShippLedger struct {

	DlvId string `json:"dlvId"`			 		// 운송장 ID		
	SlipNo string `json:"slipNo"`				// 운송장 번호
	ComTcd string `json:"comTcd"`				// 택배사 코드
	State string `json:"state"`		 			// 배송상태
	ScanTime string `json:"scanTime"`			// 배송시각
	Place string `json:"place"`					// 배송위치
	Level string `json:"level"`					// 배송단계
	SalesNm string `json:"salesNm"`				// 배송기사
	SalesTelNo1 string `json:"salesTelNo1"`		// 배송기사 핸드폰번호1
	SalesTelNo2 string `json:"salesTelNo2"`		// 배송기사 핸드폰번호2
	Remark string `json:"remark"`				// 비고
	Id string `json:"id"`			 			// 운송장 세부 ID
	RegDate string `json:"regDate"`				// 등록일시

}



func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {

	fmt.Println(" ########### shippcc chaincode Init ############")
	return shim.Success(nil)
}


func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()


    switch function {
		case "registerSP":  // 배송정보 등록
			return t.registerSP(stub, args)
		case "readSP":  	// 배송정보 등록
			return t.readSP(stub, args)	
		case "delete":  	// 배송정보 삭제
			return t.delete(stub, args)	
	//	case "update":  	// 배송정보 수정
	//		return t.update(stub, args)	
	
	default:
		return shim.Error("IF-BLC-202| Invalid Smart Contract function name.")
	}
		
}



// 배송정보 등록
func (t *SmartContract) registerSP(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("######## registerSP 함수 진입 ################## ")

	// 운송장ID 필수값 없을 시 오류
	if args[0] == "" {
		return shim.Error("Error :[registerSP] 'DlvId' does not exist")
	}

	//현재시간추출 
	time := time.Now()
	timeFmt := time.Format("20060102150405")

	var shipp = ShippLedger{
		DlvId: args[0], 
		SlipNo: args[1], 
		ComTcd: args[2], 	
		State: args[3], 
		ScanTime: args[4], 
		Place: args[5], 	
		Level: args[6], 
		SalesNm: args[7], 
		SalesTelNo1: args[8], 
		SalesTelNo2: args[9], 
		Remark: args[10], 
		Id: args[11], 
		RegDate: timeFmt, 	
	}

	shippccAsBytes, _ := json.Marshal(shipp)


	// 등록된 아이디인지 확인
	existIdCheck, returnMessage := _searchData(stub, shipp.DlvId); 
	if returnMessage != "" {
		return shim.Error(returnMessage)
	}	

	fmt.Println("########existIdCheck : " + string(existIdCheck))

	// 이미 아이디가 있으면 에러
	if existIdCheck == "true" {
		return shim.Error("Error | \"+\"" + shipp.DlvId + " already exists.")
	}


	//블록에 저장
	err := stub.PutState(shipp.DlvId, shippccAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}


// 배송정보 단건 조회
func (t *SmartContract) readSP(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var dlvId string 
	dlvId = args[0]

	fmt.Println("######## 쿼리 호출 (dlvId :" + dlvId + ")########")

	shippccAsBytes, _ := stub.GetState(dlvId)
	if err != nil {
		return shim.Error("Error: " + err.Error())
	}
	if shippccAsBytes == nil {
		return shim.Error("IF-BLC-301-002| Data searched by orderID\"+\"" + dlvId + " doesn't exists")
	}

	resultData := "{\"DlvId\":\"" + dlvId + "\",\"data\":\"" + string(shippccAsBytes) + "\"}"
	fmt.Printf("Query Response:%s\n", resultData)
	return shim.Success(shippccAsBytes)

}


// 배송정보 삭제
func (t *SmartContract) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var dlvId string
	dlvId = args[0]

	// 등록된 아이디인지 확인
	//existIdCheck, _ := _searchData(stub, dlvId); 
/*
	// 아이디가 없으면 에러
	if existIdCheck == false {
		return shim.Error("IF-BLC-601-002| Data searched by dlvId(" + dlvId + ") doesn't exists")
	}

*/	// 스테이트 db에서 데이터 삭제
	err := stub.DelState(dlvId)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	fmt.Printf("delete success")
	return shim.Success(nil)
}


// 이미 등록된 데이터인지 체크하는 함수
func  _searchData(stub shim.ChaincodeStubInterface, dlvId string) (string, string) {



	// 해당 배송정보(운송장ID)가 존재하는지 조회
	searchData, err := stub.GetState(dlvId)
	if err != nil {
		//returnMessage = "GetState Error: " + err.Error()
				return "false", "Failed to getState"

	}
/*	
	//배송정보(운송장ID)가 존재하지 않는 경우 메세지 처리
	if searchData == nil {
		returnMessage = "IF-BLC-601-002| Data searched by dlvId(" + dlvId + ") doesn't exists"
	}
*/

	// 존재하면 true , 없으면 false
	if searchData != nil {
		return "true", ""
	}
	return "false", ""
}


// 배송정보 수정
//func (t *SmartContract) update(stub shim.ChaincodeStubInterface, args []string) pb.Response {}





// 메인
func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}

}

