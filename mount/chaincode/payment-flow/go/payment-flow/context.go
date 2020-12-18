package payment_flow

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the AssetToken
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetTokenList() ListInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// AssetToken
type TransactionContext struct {
	contractapi.TransactionContext
	tokenList *list
}

// GetTokenList return AssetToken list
func (tc *TransactionContext) GetTokenList() ListInterface {
	if tc.tokenList == nil {
		tc.tokenList = newList(tc)
	}

	return tc.tokenList
}
