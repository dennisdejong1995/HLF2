package erc721

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the ERC-721
// token
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetTokenList() ListInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// ERC-721 token
type TransactionContext struct {
	contractapi.TransactionContext
	paperList *list
}

// GetTokenList return ERC-721 token list
func (tc *TransactionContext) GetTokenList() ListInterface {
	if tc.paperList == nil {
		tc.paperList = newList(tc)
	}

	return tc.paperList
}
