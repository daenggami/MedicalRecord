package main

import (
	//"bytes"
	"encoding/json"
	"fmt"
	"unicode/utf8"

	//"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type MedicalRecord struct {
	PatNo string `json:"PatNo"` //환자 ID
	//TicketNumber string `json:"TicketNumber"` //티켓번호
	RecordHash string `json:"RecordHash"` //의무기록 사본 HASH
	PatName    string `json:"PatName"`    //환자 이름
	//Deprivacy    string `json:"Deprivacy"`    //비식별화
	//Hash         string `json:"Hash"`         //환자 이름
	//	medical string `json:"medical"` 			//
	//	test string `json:"test"`	 				//Composite key test
	//	MainAilments string `json:"MainAilments"`	//주상병
	//	SubAilments string `json:"SubAilments"`		//부상병
	//	OccurDate string `json:"OccurDate"`			//발병일
	//	DiagnosisDate string `json:"DiagnosisDate"`	//진단일
	//	Opinion string `json:"Opinion"` 			//소견
}

func (s *SmartContract) Init(MeRe shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(MeRe shim.ChaincodeStubInterface) sc.Response {
	function, args := MeRe.GetFunctionAndParameters()
	if function == "queryPatNo" {
		return s.queryPatNo(MeRe, args)
	} else if function == "initLedger" {
		return s.initLedger(MeRe)
	} else if function == "queryTicket" {
		return s.queryTicket(MeRe, args)
	} else if function == "createRecordCopy" {
		return s.createRecordCopy(MeRe, args)
	} else if function == "VerificationTicket" {
		return s.VerificationTicket(MeRe, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) initLedger(MeRe shim.ChaincodeStubInterface) sc.Response {
	medicalrecord := []MedicalRecord{
		MedicalRecord{PatNo: "1024512", PatName: "hong", RecordHash: "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
		MedicalRecord{PatNo: "1124501", PatName: "임꺽정", RecordHash: "3e23e8160039594a33894f6564e1b1348bbd7a0088d42c4acb73eeaed59c009d"},
		MedicalRecord{PatNo: "1310414", PatName: "성기훈", RecordHash: "2e7d2c03a9507ae265ecf5b5356885a53393a2029d241394997265a1a25aefc6"},
		MedicalRecord{PatNo: "1489582", PatName: "강새벽", RecordHash: "18ac3e7343f016890c510e93f935261169d9e3f565436429830faf0934f4f8e4"},
		MedicalRecord{PatNo: "1294503", PatName: "조상우", RecordHash: "3f79bb7b435b05321651daefd374cdc681dc06faa65e374e38337b88ca046dea"},
		MedicalRecord{PatNo: "1630498", PatName: "오일남", RecordHash: "252f10c83610ebca1a059c0bae8255eba2f95be4d1d7bcfa89d7248a82d9f111"},
		MedicalRecord{PatNo: "1472390", PatName: "한미녀", RecordHash: "cd0aa9856147b6c5b4ff2b7dfee5da20aa38253099ef1b4a64aced233c9afe29"},
		MedicalRecord{PatNo: "2038451", PatName: "황준호", RecordHash: "aaa9402664f1a41f40ebbc52c9993eb66aeb366602958fdfaa283b71e64db123"},
		MedicalRecord{PatNo: "1492831", PatName: "지영", RecordHash: "de7d1b721a1e0632b7cf04edf5032c8ecffa9f9a08492152b926f1a5a7e765d7"},
		MedicalRecord{PatNo: "2093829", PatName: "압둘알리", RecordHash: "189f40034be7a199f1fa9891668ee3ab6049f82d38c68be70f596eab2e1857b7"},
	}
	//createRecordCopy PatNo, PatName, RecordHash
	i := 0
	for i < len(medicalrecord) {
		fmt.Println("i is ", i)
		medical := medicalrecord[i]
		var Hash = medical.RecordHash
		var TicketNumber = medical.PatNo + Hash[0:9]
		RecordAsBytes, _ := json.Marshal(medicalrecord[i])
		//MeRe.PutState("CAR"+strconv.Itoa(i), RecordAsBytes)
		MeRe.PutState(TicketNumber, RecordAsBytes)
		fmt.Println("Added", medicalrecord[i])
		i = i + 1
		indexName := "PatNo~TicketNumber"
		PatNoTicketNumberIndexKey, err := MeRe.CreateCompositeKey(indexName, []string{medical.PatNo, TicketNumber})
		fmt.Println(PatNoTicketNumberIndexKey)
		if err != nil {
			return shim.Error(err.Error())
		}
		value := []byte{0x00}
		MeRe.PutState(PatNoTicketNumberIndexKey, value)
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryPatNo(MeRe shim.ChaincodeStubInterface, args []string) sc.Response {
	//ex) queryPatNo "1234567"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	} else if len(args[0]) != 7 {
		return shim.Error("Patient Number is 7-digit.")
	}
	//RecordAsBytes, _ := MeRe.GetState(args[0])

	patno := args[0] //환자 아이디 넣기
	PatNoTicket, err := MeRe.GetStateByPartialCompositeKey("PatNo~TicketNumber", []string{patno})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer PatNoTicket.Close()
	var i int
	//var returnedPatNo string
	//var returnedTicketNumber string
	var responsePayload string
	for i = 1; PatNoTicket.HasNext(); i++ {
		responseRange, err := PatNoTicket.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		_, compositeKeyParts, err := MeRe.SplitCompositeKey(responseRange.Key)
		returnedPatNo := compositeKeyParts[0]
		returnedTicketNumber := compositeKeyParts[1]

		//const aa = '{"PatNo":compositeKeyParts[0], "TicketNumber":compositeKeyParts[1] }';
		//fmt.Println("11111 : %s ", compositeKeyParts)
		if err != nil {
			return shim.Error(err.Error())
		}
		//var res := []MedicalRecord{PatNo: compositeKeyParts[0], PatName: compositeKeyParts[1]}
		if returnedTicketNumber != "" {
			//responseargs := []MedicalRecord{returninfo{returnedPatNo:compositeKeyParts[0], returnedTicketNumber:compositeKeyParts[1]}}
			//responseargs, _ := args{returnedPatNo, returnedTicketNumber}

			//var responseargs += MedicalRecord{PatNo: compositeKeyParts[0], TicketNumber: compositeKeyParts[1]}
			//MedicalRecord += []string{i, returnedPatNo, returnedTicketNumber}
			responsePayload += fmt.Sprintf("%d - PatNo : %s / TicketNumber : %s \n", i, returnedPatNo, returnedTicketNumber)
			//responsePayload += args{returnedPatNo, returnedTicketNumber}
			//RecordAsBytes, _ := json.Marshal(responseargs)
			//responsePayload := args[returnedPatNo, returnedTicketNumber]
		}
	}

	fmt.Println(responsePayload)
	return shim.Success([]byte(responsePayload))
}

func (s *SmartContract) queryTicket(MeRe shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	} else if len(args[0]) != 16 {
		return shim.Error("Ticket Number is 16-digit.")
	}
	RecordAsBytes, _ := MeRe.GetState(args[0])
	return shim.Success(RecordAsBytes)
}

func (s *SmartContract) createRecordCopy(MeRe shim.ChaincodeStubInterface, args []string) sc.Response {
	//createRecordCopy PatNo, PatName, RecordHash
	//ex) createRecordCopy "1234567", "김수동", "d11b8fa4d028090bfe3fe174a1e769eb35c901a4983d9c4248cd7cd9f8386431"
	if len(args) != 3 { // PatNo, PatName, RecordHash
		return shim.Error("Incorrect number of arguments. Expecting 3")
	} else if len(args[0]) != 7 {
		return shim.Error("Patient Number is 7-digit.")
	} else if len(args[2]) != 64 {
		return shim.Error("Record Hash is 64-digit.")
	}
	//objectType := "Record"
	var Hash = args[2]

	var TicketNumber = args[0] + Hash[0:9]
	var Record = MedicalRecord{PatNo: args[0], PatName: args[1], RecordHash: args[2]}
	RecordAsBytes, _ := json.Marshal(Record)
	MeRe.PutState(TicketNumber, RecordAsBytes)
	indexName := "PatNo~TicketNumber"
	PatNoTicketNumberIndexKey, err := MeRe.CreateCompositeKey(indexName, []string{Record.PatNo, TicketNumber})
	fmt.Println(PatNoTicketNumberIndexKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	value := []byte{0x00}
	MeRe.PutState(PatNoTicketNumberIndexKey, value)
	return shim.Success(nil)
}

func (s *SmartContract) VerificationTicket(MeRe shim.ChaincodeStubInterface, args []string) sc.Response {
	// VerificationTicket TicketNumber
	//ex) VerificationTicket "1234567d11b8fa4d"
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	} else if len(args[0]) != 16 {
		return shim.Error("Ticket Number is 16-digit.")
	}

	RecordAsBytes, _ := MeRe.GetState(args[0])
	//TicketNumber, PatNo, PatName, RecordHash
	record := MedicalRecord{}
	json.Unmarshal(RecordAsBytes, &record)
	DePatNo := record.PatNo[:1] + "*****" + record.PatNo[6:]
	//record.PatNo = args[1]
	//record.PatName = args[2]
	name := record.PatName
	fmt.Println(name)

	//r, size := utf8.DecodeRune(name[:1])

	var lastname = []byte(record.PatName)
	r, _ := utf8.DecodeRune(lastname)
	fmt.Println("lastname: %c", r)
	r, _ = utf8.DecodeLastRune(lastname)
	fmt.Println("lastname: %c", r)
	//fmt.Println("lastname: %s / string : %s /rune : %s / string(rune): %s / rune(string)", lastname, string(lastname), []rune(lastname), string([]rune(lastname)), []rune(string(lastname)))
	fmt.Println("string : %c", string(lastname))
	//fmt.Println(" string(rune): %c", []rune(lastname))
	fmt.Println("rune(string) %c", []rune(string(lastname)))
	//fmt.Println("string(rune) %c", string([]rune(lastname)))

	var DePatName string
	if len(record.PatName) == 2 {
		DePatName = fmt.Sprintf("%s*", lastname)
	} else if len(record.PatName) >= 3 {
		var lastchrNum int
		var starNum int
		var star string
		lastchrNum = len(record.PatName) - 1
		starNum = len(record.PatName) - 2
		//aaa = star * starNum
		for i := 0; i < starNum; i++ {
			star += "*"
		}
		DePatName = fmt.Sprintf("%s,%s,%s", string(lastname), star, string(record.PatName[lastchrNum:]))
	}
	//DePatName2 := []rune(DePatName)
	//fmt.Println(DePatName2)
	fmt.Println("DePatName: %s_/name[:1]: %s_/record.PatName: %s / record.PatName[:1] : %s / rune : %s", DePatName, name[:1], record.PatName, record.PatName[:1], string([]rune(name)))
	responsePayload := fmt.Sprintf("- PatNo : %s / PatName : %s", DePatNo, DePatName)

	//DePatName := record.PatName[:1] + "***" + record.PatName[6:]
	//fmt.Printf("- objectType : %s / PatNo : %s / TicketNumber : %s \n", objectType, returnedPatNo, returnedTicketNumber)
	// fmt.Printf("PatNo : %s / PatName : %s", DePatNo, DePatName)
	// 대충 비식별화 하는 기능 추가
	// TicketNumber(16자리), DeName(홍*동), DePatNo(12****5),TimeStamp
	return shim.Success([]byte(responsePayload))
}
func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
