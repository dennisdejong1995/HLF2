package payment_flow

import (
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"time"
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
	currentTime := time.Now()
	var issueDateTime string = currentTime.Format("2006-01-02")

	// Issue token under borrower
	fmt.Printf("Issuing ERC-721 token %s for borrower %s\n", tokenID, borrower)
	erc721, err := IssueToken(ctx, borrower, tokenID, issueDateTime, maturityDateTime, faceValue)

	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Succesfully created ERC-721 token %s for borrower %s\n", erc721.TokenID, erc721.Owner)
	}

	fmt.Printf("Exchanging ERC-721 token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Owner, lender)
	erc721, err = Exchange(ctx, borrower, lender, erc721)
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Successfully exchanged ERC-721 token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Owner, lender)
	}

	return erc721, nil
}

func (c *Contract) InitiateRepayment(ctx TransactionContextInterface, borrower string, lender string, tokenID string) (*ERC721, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)

	// Exchange token back to borrower for loaned currency
	token, err = Exchange(ctx, lender, borrower, token)
	if err != nil {
		return nil, err
	}

	// Redeem ERC-721 token after exchange back
	currentTime := time.Now()
	var redeemDateTime string = currentTime.Format("2006-01-02")

	token, err = RedeemToken(ctx, borrower, token, redeemDateTime)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Payment flow functions

func Exchange(ctx TransactionContextInterface, receiver string, sender string, erc721 *ERC721) (*ERC721, error) {
	// TODO: Add atomic swap functionality
	erc721, err := ExchangeToken(ctx, receiver, sender, erc721)
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

func RedeemToken(ctx TransactionContextInterface, borrower string, token *ERC721, redeemDateTime string) (*ERC721, error) {

	if token.Owner != token.Borrower {
		return nil, fmt.Errorf("Token %s:%s is not owned by %s", token.Borrower, token.TokenID, borrower)
	}

	if token.IsRedeemed() {
		return nil, fmt.Errorf("Token %s:%s is already redeemed", token.Borrower, token.TokenID, borrower)
	}

	token.Owner = token.Borrower
	token.RedeemDateTime = redeemDateTime
	token.SetRedeemed()

	err := ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExchangeToken(ctx TransactionContextInterface, currentOwner string, futureOwner string, token *ERC721) (*ERC721, error) {

	if token.Borrower != currentOwner {
		return nil, fmt.Errorf("ERC-721 token %s:%s is not owned by %s", token.Borrower, token.TokenID, currentOwner)
	}

	token.Owner = futureOwner

	err := ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}
