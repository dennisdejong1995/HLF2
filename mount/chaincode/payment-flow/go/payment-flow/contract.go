package payment_flow

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

// Public functions

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Payment module instantiated")

}

func (c *Contract) InitiatePayment(ctx TransactionContextInterface, borrower string, lender string, maturityDateTime string, faceValue int) (*ERC721, error) {
	var tokenID string = "00007"
	var issueDateTime string = "2020-11-01"

	// Issue token under borrower
	fmt.Printf("Issuing ERC-721 token %s for borrower %s\n", tokenID, borrower)
	erc721, err := IssueToken(ctx, borrower, tokenID, issueDateTime, maturityDateTime, faceValue)

	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Succesfully created ERC-721 token %s for borrower %s\n", erc721.TokenID, erc721.Owner)
	}
	//
	//// Exchange currency for token
	//if erc721.Owner != borrower {
	//	return nil, fmt.Errorf("token %s:%s is not owned by %s", borrower, tokenID, borrower)
	//}
	//
	//if erc721.IsIssued() {
	//	erc721.SetTrading()
	//}
	//
	//if !erc721.IsTrading() {
	//	return nil, fmt.Errorf("token %s:%s is not trading. Current state = %s", borrower, tokenID, erc721.GetState())
	//}
	//
	//erc721.Owner = lender
	//
	//err = ctx.GetTokenList().UpdateToken(erc721)
	//
	//if err != nil {
	//	return nil, err
	//}

	fmt.Printf("Exchanging ERC-721 token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Owner, lender)
	erc721, err = Exchange(ctx, borrower, lender, erc721)
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Successfully exchanged ERC-721 token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Owner, lender)
	}

	return erc721, nil
}

//func (c *Contract) InitiateRepayment(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int, interest float32)  error {
//}

// Payment flow functions

func Exchange(ctx TransactionContextInterface, borrower string, lender string, erc721 *ERC721) (*ERC721, error) {
	// TODO: Add atomic swap functionality
	erc721, err := ExchangeToken(ctx, borrower, lender, erc721)
	if err != nil {
		return nil, err
	}
	return erc721, nil
}

// ERC-721 functions

func IssueToken(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int) (*ERC721, error) {
	token := ERC721{TokenID: tokenID, Borrower: borrower, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: borrower}
	token.SetIssued()
	err := ctx.GetTokenList().AddToken(&token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func ExchangeToken(ctx TransactionContextInterface, borrower string, lender string, token *ERC721) (*ERC721, error) {

	if token.Borrower != borrower {
		return nil, fmt.Errorf("ERC-721 token %s:%s is not owned by %s", token.Borrower, token.TokenID, borrower)
	}

	token.Owner = lender

	err := ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}
