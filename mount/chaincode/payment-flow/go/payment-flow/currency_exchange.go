package payment_flow

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type CurrencyExchange struct {
}

func (ce CurrencyExchange) Exchange(amount int, receiver string, receiverAddress string, sender string, senderAddress string, currency Currency) (string, error) {
	// Select exchange function for currency type
	switch currency {
	case USDT:
		transHash, err := ExchangeUSDT(amount, receiver, receiverAddress, sender, senderAddress)
		if err != nil {
			return "", err
		}
		return transHash, nil
	case EURS:
		transHash, err := ExchangeEURS(amount, receiver, receiverAddress, sender, senderAddress)
		if err != nil {
			return "", err
		}
		return transHash, nil
	default:
		return "", fmt.Errorf("no valid currency chosen for exchange")
	}
}

func ExchangeUSDT(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to Ethereum for exchanging tether
	fmt.Printf("Exchanging %d in USDT to %s from %s\n", amount, receiver, sender)

	url := "https://rinkeby.infura.io/v3/189920b69bd147cbbee96ca2c36e5ea3"

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections

	return "success", nil
}

func ExchangeEURS(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to EURS blockchain
	fmt.Printf("Exchanging %d in EURS to %s from %s\n", amount, receiver, sender)
	return "example_hash", nil
}
