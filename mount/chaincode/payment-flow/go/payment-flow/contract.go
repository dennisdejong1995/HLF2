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

func (c *Contract) Issue(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int) (*ERC721, error) {
	token := ERC721{TokenID: tokenID, Borrower: borrower, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: borrower}
	token.SetIssued()
	err := ctx.GetTokenList().AddToken(&token)
	fmt.Println("Issuing ERC-721 token")

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (c *Contract) GetOne(ctx TransactionContextInterface, borrower string, tokenID string) (*ERC721, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (c *Contract) InitiatePayment(ctx TransactionContextInterface, borrower string, lender string, maturityDateTime string, faceValue int) (*ERC721, error) {
	var tokenID string = "00007"
	var issueDateTime string = "2020-11-01"

	// Issue token under borrower
	fmt.Printf("Issuing ERC-721 token %s for borrower %s\n", tokenID, borrower)
	erc721, err := IssueToken(ctx, borrower, tokenID, issueDateTime, maturityDateTime, faceValue)

	//if err != nil {
	//	return nil, err
	//} else {
	//	fmt.Printf("Succesfully created ERC-721 token %s for borrower %s\n", erc721.TokenID, erc721.Owner)
	//}

	// Exchange currency for token
	fmt.Printf("Before")
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)
	fmt.Printf("After")

	if err != nil {
		fmt.Printf("Error")
		//return nil, err
	} else {
		fmt.Printf("Succesfully got ERC-721 token %s for borrower %s\n", token.TokenID, token.Owner)
	}

	//if token.Owner != borrower {
	//	return nil, fmt.Errorf("token %s:%s is not owned by %s", borrower, tokenID, borrower)
	//}

	//if token.IsIssued() {
	//	token.SetTrading()
	//}
	//
	//if !token.IsTrading() {
	//	return nil, fmt.Errorf("token %s:%s is not trading. Current state = %s", borrower, tokenID, token.GetState())
	//}
	//
	//token.Owner = lender
	//
	//err = ctx.GetTokenList().UpdateToken(token)
	//
	//if err != nil {
	//	return nil, err
	//}

	//fmt.Printf("Exchanging ERC-721 token %s from borrower %s to lender %s\n", token.TokenID, token.Owner, lender)
	//erc721, err = Exchange(ctx, borrower, lender, tokenID)
	//if err != nil {
	//	return nil, err
	//} else {
	//	fmt.Printf("Successfully exchanged ERC-721 token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Owner, lender)
	//}

	return token, nil
}

//func (c *Contract) InitiateRepayment(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int, interest float32)  error {
//}

// Payment flow functions

func Exchange(ctx TransactionContextInterface, borrower string, lender string, tokenID string) (*ERC721, error) {
	// TODO: Add atomic swap functionality
	erc721, err := ExchangeToken(ctx, borrower, lender, tokenID)
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

func ExchangeToken(ctx TransactionContextInterface, borrower string, lender string, tokenID string) (*ERC721, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)

	if err != nil {
		return nil, err
	}
	if token.Borrower != borrower {
		return nil, fmt.Errorf("ERC-721 token %s:%s is not owned by %s", token.Borrower, token.TokenID, borrower)
	}

	token.Owner = lender

	err = ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}
