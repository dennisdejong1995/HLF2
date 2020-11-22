package payment_flow

import "fmt"

func ExchangeUSDT(erc721 *ERC721, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to Ethereum for exchanging tether
	fmt.Printf("Exchanging %d in USDT to %s from %s\n", erc721.FaceValue, receiver, sender)
	return "example_hash", nil
}

func ExchangeEURS(erc721 *ERC721, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to EURS blockchain
	fmt.Printf("Exchanging %d in EURS to %s from %s\n", erc721.FaceValue, receiver, sender)
	return "example_hash", nil
}
