package payment_flow

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/hyperledger/fabric-samples/commercial-paper/organization/digibank/contract-go/ledger-api"
)

// State enum for asset token state property
type State uint

const (
	// ISSUED state for when an AssetToken token is issued
	ISSUED State = iota + 1
	// TRADING state for when an AssetToken token is trading
	TRADING
	// REDEEMED state for when an AssetToken token has been redeemed
	REDEEMED
)

func (state State) String() string {
	names := []string{"ISSUED", "TRADING", "REDEEMED"}

	if state < ISSUED || state > REDEEMED {
		return "UNKNOWN"
	}

	return names[state-1]
}

// Currency enum for asset token state property
type Currency uint

const (
	// USDT currency for when an AssetToken token uses USDT as currency
	USDT Currency = iota + 1
	// EURS currency for when an AssetToken token uses EURS as currency
	EURS
)

func (cur Currency) String() string {
	names := []string{"USDT", "EURS"}

	if cur < USDT || cur > EURS {
		return "UNKNOWN"
	}

	return names[cur-1]
}

// CreateTokenKey creates a key for asset tokens
func CreateTokenKey(borrower string, tokenID string) string {
	return ledgerapi.MakeKey(borrower, tokenID)
}

// Used for managing the fact status is private but want it in world state
type assetTokenAlias AssetToken
type jsonAssetToken struct {
	*assetTokenAlias
	State State  `json:"currentState"`
	Class string `json:"class"`
	Key   string `json:"key"`
}

// CommercialPaper defines a AssetToken token
type AssetToken struct {
	TokenID          string   `json:"tokenID"`
	Borrower         string   `json:"borrower"`
	Investor         string   `json:"investor"`
	IssueDateTime    string   `json:"issueDateTime"`
	FaceValue        int      `json:"faceValue"`
	MaturityDateTime string   `json:"maturityDateTime"`
	RedeemDateTime   string   `json:"redeemDateTime"`
	Owner            string   `json:"owner"`
	Interest         int      `json:"interest"`
	BorrowerAddress  string   `json:"senderAddress"`
	InvestorAddress  string   `json:"investorAddress"`
	OwnerAddress     string   `json:"receiverAddress"`
	PaymentHashes    []string `json:"paymentHashes"`
	Currency         Currency `json:"currency"`
	state            State    `metadata:"currentState"`
	class            string   `metadata:"class"`
	key              string   `metadata:"key"`
}

// UnmarshalJSON special handler for managing JSON marshalling
func (cp *AssetToken) UnmarshalJSON(data []byte) error {
	jcp := jsonAssetToken{assetTokenAlias: (*assetTokenAlias)(cp)}

	err := json.Unmarshal(data, &jcp)

	if err != nil {
		return err
	}

	cp.state = jcp.State

	return nil
}

// MarshalJSON special handler for managing JSON marshalling
func (cp AssetToken) MarshalJSON() ([]byte, error) {
	jcp := jsonAssetToken{assetTokenAlias: (*assetTokenAlias)(&cp), State: cp.state, Class: "org.dealblock.erc721", Key: ledgerapi.MakeKey(cp.Borrower, cp.TokenID)}

	return json.Marshal(&jcp)
}

// GetState returns the state
func (cp *AssetToken) GetState() State {
	return cp.state
}

// SetIssued returns the state to issued
func (cp *AssetToken) SetIssued() {
	cp.state = ISSUED
}

// SetTrading sets the state to trading
func (cp *AssetToken) SetTrading() {
	cp.state = TRADING
}

// SetRedeemed sets the state to redeemed
func (cp *AssetToken) SetRedeemed() {
	cp.state = REDEEMED
}

// IsIssued returns true if state is issued
func (cp *AssetToken) IsIssued() bool {
	return cp.state == ISSUED
}

// IsTrading returns true if state is trading
func (cp *AssetToken) IsTrading() bool {
	return cp.state == TRADING
}

// IsRedeemed returns true if state is redeemed
func (cp *AssetToken) IsRedeemed() bool {
	return cp.state == REDEEMED
}

// GetCurrency returns the currency
func (cp *AssetToken) GetCurrency() Currency {
	return cp.Currency
}

// SetUSDT sets the currency to USDT
func (cp *AssetToken) SetUSDT() {
	cp.Currency = USDT
}

// SetEURS sets the currency to EURS
func (cp *AssetToken) SetEURS() {
	cp.Currency = EURS
}

// IsUSDT returns true if currency is USDT
func (cp *AssetToken) IsUSDT() bool {
	return cp.Currency == USDT
}

// IsEURS returns true if currency is EURS
func (cp *AssetToken) IsEURS() bool {
	return cp.Currency == EURS
}

// GetSplitKey returns values which should be used to form key
func (cp *AssetToken) GetSplitKey() []string {
	return []string{cp.Borrower, cp.TokenID}
}

// Serialize formats the AssetToken token as JSON bytes
func (cp *AssetToken) Serialize() ([]byte, error) {
	return json.Marshal(cp)
}

// Deserialize formats the AssetToken token from JSON bytes
func Deserialize(bytes []byte, cp *AssetToken) error {
	err := json.Unmarshal(bytes, cp)

	if err != nil {
		return fmt.Errorf("Error deserializing AssetToken token %s", err.Error())
	}

	return nil
}
