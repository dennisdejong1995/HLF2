package erc721

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing commercial
// paper
type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("ERC-721 module instantiated")

}

// Issue creates a new commercial paper and stores it in the world state
func (c *Contract) Issue(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int, interest float32) (*ERC721, error) {
	token := ERC721{TokenID: tokenID, Borrower: borrower, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: borrower, Interest: interest}
	token.SetIssued()
	err := ctx.GetTokenList().AddToken(&token)
	fmt.Println("Issuing ERC-721 token")

	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Buy updates a commercial paper to be in trading status and sets the new owner
func (c *Contract) Buy(ctx TransactionContextInterface, borrower string, tokenID string, currentOwner string, newOwner string, price int, purchaseDateTime string) (*ERC721, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)

	if err != nil {
		return nil, err
	}

	if token.Owner != currentOwner {
		return nil, fmt.Errorf("Paper %s:%s is not owned by %s", borrower, tokenID, currentOwner)
	}

	if token.IsIssued() {
		token.SetTrading()
	}

	if !token.IsTrading() {
		return nil, fmt.Errorf("Token %s:%s is not trading. Current state = %s", borrower, tokenID, token.GetState())
	}

	token.Owner = newOwner

	err = ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

// Redeem updates a commercial paper status to be redeemed
func (c *Contract) Redeem(ctx TransactionContextInterface, borrower string, paperNumber string, redeemingOwner string, redeemDateTime string) (*ERC721, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, paperNumber)

	if err != nil {
		return nil, err
	}

	if token.Owner != redeemingOwner {
		return nil, fmt.Errorf("Token %s:%s is not owned by %s", borrower, paperNumber, redeemingOwner)
	}

	if token.IsRedeemed() {
		return nil, fmt.Errorf("Token %s:%s is already redeemed", borrower, paperNumber)
	}

	token.Owner = token.Borrower
	token.SetRedeemed()

	err = ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (c *Contract) GetOne(ctx TransactionContextInterface, borrower string, tokenID string) (*ERC721, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)

	if err != nil {
		return nil, err
	}

	return token, nil
}
