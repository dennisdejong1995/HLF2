package main

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"gitlab.com/dealblock/payment-package/payment-flow"
)

func main() {

	contract := new(payment_flow.Contract)
	contract.TransactionContextHandler = new(payment_flow.TransactionContext)
	contract.Name = "org.dealblock.payment_flow"
	contract.Info.Version = "0.0.1"

	chaincode, err := contractapi.NewChaincode(contract)

	if err != nil {
		panic(fmt.Sprintf("Error creating chaincode. %s", err.Error()))
	}

	chaincode.Info.Title = "PaymentFlowChaincode"
	chaincode.Info.Version = "0.0.1"

	err = chaincode.Start()

	if err != nil {
		panic(fmt.Sprintf("Error starting chaincode. %s", err.Error()))
	}
}
