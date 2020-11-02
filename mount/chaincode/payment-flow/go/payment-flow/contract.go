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

func (c *Contract) InitiatePayment(ctx TransactionContextInterface, borrower string, lender string, maturityDateTime string, faceValue int, interest float32) error {
	var tokenID = "000001"
	var issueDateTime = "02-11-2020"

	// Issue token under borrower
	fmt.Printf("Issuing ERC-721 token %s for borrower %s\n", tokenID, borrower)
	erc721, err := IssueToken(ctx, borrower, tokenID, issueDateTime, maturityDateTime, faceValue, interest)
	if err != nil {
		return err
	} else {
		fmt.Printf("Succesfully created ERC-721 token %s for borrower %s\n", erc721.TokenID, erc721.Owner)
	}

	// Exchange currency for token
	fmt.Printf("Exchanging ERC-721 token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Owner, lender)
	err = Exchange(ctx, borrower, lender, tokenID)
	if err != nil {
		return err
	} else {
		fmt.Printf("Succesfully exchanged ERC-721 token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Owner, lender)
	}

	return nil
}

//func (c *Contract) InitiateRepayment(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int, interest float32)  error {
//}

// Payment flow functions

func Exchange(ctx TransactionContextInterface, borrower string, lender string, tokenID string) error {
	// TODO: Add atomic swap functionality
	erc721, err := ExchangeToken(ctx, borrower, lender, tokenID)
	fmt.Println(erc721)
	if err != nil {
		return err
	}
	return nil
}

// ERC-721 functions

func IssueToken(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int, interest float32) (*ERC721, error) {
	token := ERC721{TokenID: tokenID, Borrower: borrower, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: borrower, Interest: interest}
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
