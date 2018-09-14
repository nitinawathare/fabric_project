/*
docker stop $(docker ps -aq)
docker rm -f $(docker ps -aq)
docker rmi
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

//Init and Invoke
type SmartContract struct {
}

/*
type Citizen struct {
	FirstName  string `json:"fname"`
	MiddleName string `json:"mname"`
	LastName   string `json:"lname"`
	Address    string `json:"addr"`
}
*/

type StampPaper struct {
	StampID      string `json:"stampid"`
	AdhaarID     string `json:"uid"`
	StampHolder  string `json:"stamp_holder"`
	Location     string `json:"location"`
	DocumentType string `json:"doc_type"`
	Content      string `json:"doc_content"`
}

/*
 * The Init method *
 called when the Smart Contract is instantiated by the network
 -- see initLedger()
*/
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract
 The app also specifies the specific smart contract function to call with args
*/
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	//Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryStampPaper" {
		return s.queryStampPaper(APIstub, args)
	} else if function == "initLedger" {
		fmt.Println("I was called")
		return s.initLedger(APIstub)
	} else if function == "queryAllStampPaper" {
		return s.queryAllStampPaper(APIstub)
	} else if function == "recordStamp" {
		return s.recordStamp(APIstub, args)
	} else if function == "changeOwner" {
		return s.changeOwner(APIstub, args)
	} else {
		return shim.Error("Invalid Smart Contract function name.")
	}
	//return shim.Success(nil)
}

func (s *SmartContract) recordStamp(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments")
	}

	var sptemp = StampPaper{StampID: args[0], AdhaarID: args[1], StampHolder: args[2], Location: args[3], DocumentType: args[4], Content: args[5]}
	spAsBytes, _ := json.Marshal(sptemp)
	err := APIstub.PutState(args[0], spAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to add Stamp Paper: %s", args[0]))
	}
	return shim.Success(nil)
}

func (s *SmartContract) queryStampPaper(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Please enter the Stamp ID")
	}

	spAsBytes, _ := APIstub.GetState(args[0])
	if spAsBytes == nil {
		return shim.Error("Stamp Paper not found")
	}
	return shim.Success(spAsBytes)
}

func (s *SmartContract) queryAllStampPaper(APIstub shim.ChaincodeStubInterface) sc.Response {
	startKey := "0"
	endKey := "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)

	if err != nil {
		return shim.Error(err.Error())
	}

	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer

	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false

	for resultsIterator.HasNext() {

		queryResponse, err := resultsIterator.Next()
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
	fmt.Printf("- queryAllStampPaper:\n%s\n", buffer.String())
	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	sp := []StampPaper{
		StampPaper{StampID: "1000", AdhaarID: "123456789012", StampHolder: "Mahima Manik", Location: "Varanasi", DocumentType: "Company Registration", Content: "I want to register my Block chain company"},
		StampPaper{StampID: "1001", AdhaarID: "456712906777", StampHolder: "Khushboo Goel", Location: "Haridwar", DocumentType: "Affidavit", Content: "Need certificate for Anti- Ragging"},
		StampPaper{StampID: "1002", AdhaarID: "931072345612", StampHolder: "Jyoti Dhakla", Location: "Nalanda", DocumentType: "Land Registry", Content: "Transfer the land from me to my brother"},
		StampPaper{StampID: "1003", AdhaarID: "188905671234", StampHolder: "Vyom Manik", Location: "New Delhi", DocumentType: "Driving Licence", Content: "I want to my Driving License"},
		StampPaper{StampID: "1004", AdhaarID: "786548193639", StampHolder: "Shresth Manik", Location: "Mumbai", DocumentType: "Vehicle Registry", Content: "Registration of bike"},
	}

	i := 0
	for i < len(sp) {
		fmt.Println("i is ", i)
		spAsBtyes, _ := json.Marshal(sp[i])
		APIstub.PutState(string(sp[i].StampID), spAsBtyes)
		fmt.Println("Added ", sp[i])
		i = i + 1
	}

	return shim.Success(nil)
}
func (s *SmartContract) changeOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments")
	}

	spAsBytes, _ := APIstub.GetState(args[0])
	if spAsBytes == nil {
		return shim.Error("Stamp Paper not found")
	}
	sp := StampPaper{}

	json.Unmarshal(spAsBytes, &sp)
	sp.StampHolder = args[1]
	spAsBytes, _ = json.Marshal(sp)
	err := APIstub.PutState(args[0], spAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change Stamp Paper: %s", args[0]))
	}
	return shim.Success(nil)
}

func main() {
	//first_stamp := StampPaper{"1001", "123456789012", "Mahima Manik", "New Delhi", "Company Registration", "I want to register my Block chain company"}
	//theJson, _ := json.Marshal(first_stamp) //returns JSON encoding of first_stamp

	//fmt.Printf("%+v\n", string(theJson))
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	} else {
		fmt.Println("Success")
	}
}
