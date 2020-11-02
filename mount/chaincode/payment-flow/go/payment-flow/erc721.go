package payment_flow

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/hyperledger/fabric-samples/commercial-paper/organization/digibank/contract-go/ledger-api"
)

// State enum for ERC-721 state property
type State uint

const (
	// ISSUED state for when an ERC-721 token is issued
	ISSUED State = iota + 1
	// TRADING state for when an ERC-721 token is trading
	TRADING
	// REDEEMED state for when an ERC-721 token has been redeemed
	REDEEMED
)

func (state State) String() string {
	names := []string{"ISSUED", "TRADING", "REDEEMED"}

	if state < ISSUED || state > REDEEMED {
		return "UNKNOWN"
	}

	return names[state-1]
}

// CreateTokenKey creates a key for ERC-721 tokens
func CreateTokenKey(issuer string, paperNumber string) string {
	return ledgerapi.MakeKey(issuer, paperNumber)
}

// Used for managing the fact status is private but want it in world state
type erc721Alias ERC721
type jsonERC721 struct {
	*erc721Alias
	State State  `json:"currentState"`
	Class string `json:"class"`
	Key   string `json:"key"`
}

// CommercialPaper defines a ERC-721 token
type ERC721 struct {
	TokenID          string  `json:"tokenID"`
	Borrower         string  `json:"borrower"`
	IssueDateTime    string  `json:"issueDateTime"`
	FaceValue        int     `json:"faceValue"`
	MaturityDateTime string  `json:"maturityDateTime"`
	Owner            string  `json:"owner"`
	Interest         float32 `json:"interest"`
	state            State   `metadata:"currentState"`
	class            string  `metadata:"class"`
	key              string  `metadata:"key"`
}

// UnmarshalJSON special handler for managing JSON marshalling
func (cp *ERC721) UnmarshalJSON(data []byte) error {
	jcp := jsonERC721{erc721Alias: (*erc721Alias)(cp)}

	err := json.Unmarshal(data, &jcp)

	if err != nil {
		return err
	}

	cp.state = jcp.State

	return nil
}

// MarshalJSON special handler for managing JSON marshalling
func (cp ERC721) MarshalJSON() ([]byte, error) {
	jcp := jsonERC721{erc721Alias: (*erc721Alias)(&cp), State: cp.state, Class: "org.dealblock.erc721", Key: ledgerapi.MakeKey(cp.Borrower, cp.TokenID)}

	return json.Marshal(&jcp)
}

// GetState returns the state
func (cp *ERC721) GetState() State {
	return cp.state
}

// SetIssued returns the state to issued
func (cp *ERC721) SetIssued() {
	cp.state = ISSUED
}

// SetTrading sets the state to trading
func (cp *ERC721) SetTrading() {
	cp.state = TRADING
}

// SetRedeemed sets the state to redeemed
func (cp *ERC721) SetRedeemed() {
	cp.state = REDEEMED
}

// IsIssued returns true if state is issued
func (cp *ERC721) IsIssued() bool {
	return cp.state == ISSUED
}

// IsTrading returns true if state is trading
func (cp *ERC721) IsTrading() bool {
	return cp.state == TRADING
}

// IsRedeemed returns true if state is redeemed
func (cp *ERC721) IsRedeemed() bool {
	return cp.state == REDEEMED
}

// GetSplitKey returns values which should be used to form key
func (cp *ERC721) GetSplitKey() []string {
	return []string{cp.Borrower, cp.TokenID}
}

// Serialize formats the ERC-721 token as JSON bytes
func (cp *ERC721) Serialize() ([]byte, error) {
	return json.Marshal(cp)
}

// Deserialize formats the ERC-721 token from JSON bytes
func Deserialize(bytes []byte, cp *ERC721) error {
	err := json.Unmarshal(bytes, cp)

	if err != nil {
		return fmt.Errorf("Error deserializing ERC-721 token %s", err.Error())
	}

	return nil
}
