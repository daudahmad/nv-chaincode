package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args []string) {
	_, err := stub.MockInit("1", "init", args)
	if err != nil {
		fmt.Println("Init failed", err)
		t.FailNow()
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	var rfid FinancialInst
	err := json.Unmarshal(bytes, &rfid)
	if err != nil {
		fmt.Println("State value", name, "was not as expected")
		t.FailNow()
	}
}

func checkOwner(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	var rfid FinancialInst
	err := json.Unmarshal(bytes, &rfid)
	if err != nil {
		fmt.Println("State value", name, "failed Unmarshal")
		t.FailNow()
	}
	if rfid.Owner != value {
		fmt.Println("State ", name, "owner was not BANKA as expected")
		t.FailNow()
	}
}

func TestNostroVostro_Init(t *testing.T) {
	scc := new(SimpleChaincode)
	stub := shim.NewMockStub("nostrovostro_cc", scc)

	checkInit(t, stub, []string{"BANKA"})

	checkState(t, stub, "BANKA")
	checkOwner(t, stub, "BANKA", "BANKA")
	checkOwner(t, stub, "BANKB", "BANKB")
	checkOwner(t, stub, "BANKC", "BANKC")
}
