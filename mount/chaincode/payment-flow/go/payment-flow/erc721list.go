package payment_flow

import ledgerapi "github.com/hyperledger/fabric-samples/commercial-paper/organization/digibank/contract-go/ledger-api"

// ListInterface defines functionality needed
// to interact with the world state on behalf
// of a commercial paper
type ListInterface interface {
	AddToken(*ERC721) error
	GetToken(string, string) (*ERC721, error)
	UpdateToken(*ERC721) error
}

type list struct {
	stateList ledgerapi.StateListInterface
}

func (cpl *list) AddToken(erc721 *ERC721) error {
	return cpl.stateList.AddState(erc721)
}

func (cpl *list) GetToken(borrower string, tokenID string) (*ERC721, error) {
	cp := new(ERC721)

	err := cpl.stateList.GetState(CreateTokenKey(borrower, tokenID), cp)

	if err != nil {
		return nil, err
	}

	return cp, nil
}

func (cpl *list) UpdateToken(token *ERC721) error {
	return cpl.stateList.UpdateState(token)
}

// NewList create a new list from context
func newList(ctx TransactionContextInterface) *list {
	stateList := new(ledgerapi.StateList)
	stateList.Ctx = ctx
	stateList.Name = "org.dealblock.erc721list"
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*ERC721))
	}

	list := new(list)
	list.stateList = stateList

	return list
}
