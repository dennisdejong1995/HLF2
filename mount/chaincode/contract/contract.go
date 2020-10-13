package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type State uint

const (
	PENDING State = iota + 1
	ACTIVE
	DONE
)

func (state State) String() string {
	names := []string{"PENDING", "ACTIVE", "DONE"}

	if state < PENDING || state > DONE {
		return "UNKNOWN"
	}

	return names[state-1]
}

type Asset struct {
	state            State  `metadata:"currentState"`
	Borrower         string `json:"issuer"`
	Lender           string `json:"owner"`
	Value        	 int    `json:"value"`
	Interest         int    `json:"interest"`
}

type QueryResult struct {
	Key    string `json:"Key"`
	Record *Asset
}


func (a *Asset) GetState() State {
	return a.state
}

func (a *Asset) SetPending() {
	a.state= PENDING
}

func (a *Asset) SetActive() {
	a.state= ACTIVE
}

func (a *Asset) SetDone() {
	a.state= DONE
}



func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	assets := []Asset{
		Asset{Borrower: "Borrower1", Lender: "Lender1", Value: 1000, Interest: 5},
		Asset{Borrower: "Borrower2", Lender: "Lender2", Value: 2000, Interest: 10},
		Asset{Borrower: "Borrower3", Lender: "Lender3", Value: 3000, Interest: 15},
	}


	// @TODO
	// Composite key toevoegen 
	for i, asset := range assets {
		asset.SetPending()
		assetAsBytes, _ := json.Marshal(asset)
		err := ctx.GetStub().PutState("ASSET"+strconv.Itoa(i), assetAsBytes)

		if err != nil {
			return fmt.Errorf("Failed.... %s", err.Error())
		}
	}

	return nil
}

// func (s *SmartContract) Issue(ctx contractapi.TransactionContextInterface, a *Asset) error {
// 	// Create offer
// }

// func (s *SmartContract) Accept(ctx contractapi.TransactionContextInterface, assetNumber string) error {
// 	// Accept offer
// }

// func (s *SmartContract) Amortize(ctx contractapi.TransactionContextInterface, assetNumber string, value int) error {
// 	// Pay of debt as a whole or partially
// }


func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create Asset chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
