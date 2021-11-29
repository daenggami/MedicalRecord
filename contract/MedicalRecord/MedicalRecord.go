package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}

type MedicalRecord struct {
	PatNo      string    `json:"PatNo"`      //환자 ID
	RecordHash string    `json:"RecordHash"` //의무기록 사본 HASH
	PatName    string    `json:"PatName"`    //환자 이름
	TimeStamp  time.Time `json:"TimeStamp"`  //시간정보
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
	} else if function == "queryTX" {
		return s.queryTX(MeRe, args)
	} else if function == "createRecordCopy" {
		return s.createRecordCopy(MeRe, args)
	}
	return shim.Error("Invalid Smart Contract function name.")
}

// 초기 데이터 입력 # 타임스템프 값은 없음.
func (s *SmartContract) initLedger(MeRe shim.ChaincodeStubInterface) sc.Response {
	medicalrecord := []MedicalRecord{
		MedicalRecord{PatNo: "1024512", PatName: "홍길동", RecordHash: "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb"},
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

//환자 ID를 통해 해당 환자에게 증명서 티켓번호가 어떤것이 있는지 확인
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

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for PatNoTicket.HasNext() {
		responseRange, err := PatNoTicket.Next()
		_, compositeKeyParts, err := MeRe.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"PatNo\":")
		buffer.WriteString("\"")
		buffer.WriteString(string(compositeKeyParts[0]))
		buffer.WriteString("\"")

		buffer.WriteString(", \"TicketNumber\":")
		buffer.WriteString("\"")
		buffer.WriteString(string(compositeKeyParts[1]))
		buffer.WriteString("\"")
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- QueryPatNo:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

//증명서 티켓번호를 입력하여 정보호출
func (s *SmartContract) queryTicket(MeRe shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	} else if len(args[0]) != 16 {
		return shim.Error("Ticket Number is 16-digit.")
	}
	RecordAsBytes, _ := MeRe.GetState(args[0])
	return shim.Success(RecordAsBytes)
}

// 증명서 정보 입력
func (s *SmartContract) createRecordCopy(MeRe shim.ChaincodeStubInterface, args []string) sc.Response {
	//createRecordCopy PatNo, PatName, RecordHash
	//ex) createRecordCopy "1234567", "김수동", "d11b8fa4d028090bfe3fe174a1e769eb35c901a4983d9c4248cd7cd9f8386431"
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	} else if len(args[0]) != 7 {
		return shim.Error("Patient Number is 7-digit.")
	} else if len(args[2]) != 64 {
		return shim.Error("Record Hash is 64-digit.")
	}
	var Hash = args[2]
	t := time.Now()
	var TicketNumber = args[0] + Hash[0:9]
	var Record = MedicalRecord{PatNo: args[0], PatName: args[1], RecordHash: args[2], TimeStamp: t}
	RecordAsBytes, _ := json.Marshal(Record)
	MeRe.PutState(TicketNumber, RecordAsBytes)
	//Composite Key 등록
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

// 발급된 증명서 트랙젝션 확인
func (t *SmartContract) queryTX(MeRe shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	key := args[0]

	fmt.Printf("- start getHistoryForKey: %s\n", key)

	resultsIterator, err := MeRe.GetHistoryForKey(key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")

		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForKey returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
