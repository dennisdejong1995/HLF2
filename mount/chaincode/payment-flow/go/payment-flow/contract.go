package payment_flow

import (
	"errors"
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

func (c *Contract) InitiatePayment(ctx TransactionContextInterface, assetID string, borrower string, lender string, maturityDateTime string, faceValue int, currencyID int, interest int) (*AssetToken, error) {
	currentTime := time.Now()
	var issueDateTime string = currentTime.Format("2006-01-02")

	// Issue token under borrower
	fmt.Printf("Issuing asset token %s for borrower %s\n", assetID, borrower)
	erc721, err := IssueToken(ctx, borrower, assetID, issueDateTime, maturityDateTime, faceValue, currencyID, interest)

	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Succesfully created asset token %s for borrower %s\n", erc721.TokenID, erc721.Owner)
	}

	fmt.Printf("Exchanging asset token %s from borrower %s to investor %s\n", erc721.TokenID, erc721.Owner, lender)
	erc721, err = Exchange(ctx, borrower, lender, erc721)
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Successfully exchanged asset token %s from borrower %s to lender %s\n", erc721.TokenID, erc721.Borrower, erc721.Owner)
	}

	return erc721, nil
}

func (c *Contract) InitiateRepayment(ctx TransactionContextInterface, borrower string, lender string, tokenID string) (*AssetToken, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)

	// Exchange token back to borrower for loaned currency
	token, err = Exchange(ctx, lender, borrower, token)
	if err != nil {
		return nil, err
	}

	// Redeem asset token after exchange back
	currentTime := time.Now()
	var redeemDateTime string = currentTime.Format("2006-01-02")

	token, err = RedeemToken(ctx, borrower, token, redeemDateTime)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Payment flow functions

func Exchange(ctx TransactionContextInterface, receiver string, sender string, token *AssetToken) (*AssetToken, error) {
	// TODO: Add atomic swap functionality
	token, err := ExchangeToken(ctx, receiver, sender, token)
	if err != nil {
		return nil, err
	}
	token, hash, err := ExchangeCurrency(receiver, sender, token)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Received hash %s\n", hash)
	return token, nil
}

// AssetToken functions

func IssueToken(ctx TransactionContextInterface, borrower string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int, currencyID int, interest int) (*AssetToken, error) {
	var currency string = ""
	switch currencyID {
	case 0:
		currency = "USDT"
	case 1:
		currency = "EURS"
	default:
		return nil, errors.New("No valid currency selected")
	}

	token := AssetToken{TokenID: tokenID, Borrower: borrower, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: borrower, Currency: currency, Interest: interest}
	token.SetIssued()
	err := ctx.GetTokenList().AddToken(&token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func RedeemToken(ctx TransactionContextInterface, borrower string, token *AssetToken, redeemDateTime string) (*AssetToken, error) {

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

func ExchangeToken(ctx TransactionContextInterface, currentOwner string, futureOwner string, token *AssetToken) (*AssetToken, error) {

	if token.Owner != currentOwner {
		return nil, fmt.Errorf("Asset token %s:%s is not owned by %s", token.Borrower, token.TokenID, currentOwner)
	}

	token.Owner = futureOwner

	err := ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExchangeCurrency(receiver string, sender string, token *AssetToken) (*AssetToken, string, error) {
	fmt.Printf("Exchanging currency from %s to %s\n", sender, receiver)

	// Determining amount
	var amount int = 0
	switch token.Borrower {
	case receiver:
		amount = token.FaceValue
	case sender:
		amount = token.FaceValue * (100 + token.Interest) / 100
	default:
		return nil, "", fmt.Errorf("Unknown transaction type")
	}

	// Select exchange function for currency type
	switch token.Currency {
	case "USDT":
		hash, err := ExchangeUSDT(amount, receiver, "receiver_address", sender, "sender_address")
		if err != nil {
			return nil, "", err
		}
		return token, hash, nil
	case "EURS":
		hash, err := ExchangeEURS(amount, receiver, "receiver_address", sender, "sender_address")
		if err != nil {
			return nil, "", err
		}
		return token, hash, nil
	default:
		return nil, "", fmt.Errorf("%s:%s No valid currency chosen for exchange", token.TokenID, token.Borrower)
	}

}
