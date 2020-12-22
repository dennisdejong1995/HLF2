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

func (c *Contract) InitiatePayment(ctx TransactionContextInterface, assetID string, borrower string, borrowerAddress string, investor string, investorAddress string, maturityDateTime string, faceValue int, currencyID int, interest int) (*AssetToken, error) {
	currentTime := time.Now()
	var issueDateTime string = currentTime.Format("2006-01-02")

	// Issue token under borrower
	fmt.Printf("Issuing asset token %s for borrower %s\n", assetID, borrower)
	token, err := IssueToken(ctx, borrower, borrowerAddress, investor, investorAddress, assetID, issueDateTime, maturityDateTime, faceValue, currencyID, interest)

	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Succesfully created asset token %s for borrower %s\n", token.TokenID, token.Owner)
	}

	fmt.Printf("Exchanging asset token %s from borrower %s to investor %s\n", token.TokenID, token.Owner, investor)
	token, err = Exchange(ctx, borrower, investor, token.OwnerAddress, token)
	if err != nil {
		return nil, err
	} else {
		fmt.Printf("Successfully exchanged asset token %s from borrower %s to lender %s\n", token.TokenID, token.Borrower, token.Owner)
	}

	return token, nil
}

func (c *Contract) InitiateRepayment(ctx TransactionContextInterface, borrower string, lender string, tokenID string) (*AssetToken, error) {
	token, err := ctx.GetTokenList().GetToken(borrower, tokenID)
	if err != nil {
		return nil, err
	}

	// Exchange token back to borrower for loaned currency
	token, err = Exchange(ctx, lender, borrower, token.BorrowerAddress, token)
	if err != nil {
		return nil, err
	}

	// Redeem asset token after exchange back
	token, err = RedeemToken(ctx, borrower, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Payment flow functions

func Exchange(ctx TransactionContextInterface, receiver string, sender string, senderAddress string, token *AssetToken) (*AssetToken, error) {

	hash, err := ExchangeCurrency(receiver, sender, token)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Received hash %s\n", hash)
	token.PaymentHashes = append(token.PaymentHashes, hash)
	token, err = ExchangeToken(ctx, receiver, sender, senderAddress, token)
	if err != nil {
		return nil, err
	}
	return token, nil
}

// AssetToken functions

func IssueToken(ctx TransactionContextInterface, borrower string, borrowerAddress string, investor string, investorAddress string, tokenID string, issueDateTime string, maturityDateTime string, faceValue int, currencyID int, interest int) (*AssetToken, error) {

	token := AssetToken{TokenID: tokenID, Borrower: borrower, BorrowerAddress: borrowerAddress, Investor: investor, InvestorAddress: investorAddress, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: borrower, Interest: interest}
	token.SetIssued()
	switch currencyID {
	case 0:
		token.SetUSDT()
	case 1:
		token.SetEURS()
	default:
		return nil, errors.New("no valid currency selected")
	}
	err := ctx.GetTokenList().AddToken(&token)

	if err != nil {
		return nil, err
	}

	return &token, nil
}

func RedeemToken(ctx TransactionContextInterface, borrower string, token *AssetToken) (*AssetToken, error) {
	currentTime := time.Now()
	var redeemDateTime string = currentTime.Format("2006-01-02")

	if token.Owner != token.Borrower {
		return nil, fmt.Errorf("asset token %s:%s is not owned by %s", token.Borrower, token.TokenID, borrower)
	}

	if token.IsRedeemed() {
		return nil, fmt.Errorf("Token %s:%s is already redeemed", token.Borrower, token.TokenID, borrower)
	}

	token.RedeemDateTime = redeemDateTime
	token.SetRedeemed()

	err := ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExchangeToken(ctx TransactionContextInterface, currentOwner string, futureOwner string, futureOwnerAddress string, token *AssetToken) (*AssetToken, error) {

	if token.Owner != currentOwner {
		return nil, fmt.Errorf("asset token %s:%s is not owned by %s", token.Borrower, token.TokenID, currentOwner)
	}

	token.Owner = futureOwner
	token.OwnerAddress = futureOwnerAddress

	err := ctx.GetTokenList().UpdateToken(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func ExchangeCurrency(receiver string, sender string, token *AssetToken) (string, error) {
	fmt.Printf("Exchanging currency from %s to %s\n", sender, receiver)

	// Determining amount and check for valid transaction
	var receiverAddress = ""
	var senderAddress = ""
	var amount = 0
	switch token.Borrower {
	case receiver:
		if sender == token.Investor {
			amount = token.FaceValue
			receiverAddress = token.BorrowerAddress
			senderAddress = token.InvestorAddress
		} else {
			return "", fmt.Errorf("invalid transaction: Sender is not Investor")
		}
	case sender:
		if token.Owner == token.Investor {
			amount = token.FaceValue * (100 + token.Interest) / 100
			receiverAddress = token.InvestorAddress
			senderAddress = token.BorrowerAddress
		} else if token.Owner == receiver {
			amount = token.FaceValue * (100 + token.Interest) / 100
			receiverAddress = token.OwnerAddress
			senderAddress = token.BorrowerAddress
		} else {
			return "", fmt.Errorf("invalid transaction: Receiver is neither Investor nor Owner")
		}
	default:
		return "", fmt.Errorf("invalid transaction: Borrower is neither sender nor receiver")
	}

	hash, err := CurrencyExchange{}.Exchange(amount, receiver, receiverAddress, sender, senderAddress, token.Currency)

	if err != nil {
		return "", err
	}
	return hash, nil
}
