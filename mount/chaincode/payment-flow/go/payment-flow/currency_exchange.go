package payment_flow

import (
	"fmt"
	//token "github.com/dennisdejong1995/HLF2/mount/chaincode/payment-flow/go/contracts_erc20"
	//"github.com/ethereum/go-ethereum/accounts/abi/bind"
	//"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	//"math"
	//"math/big"
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

	fmt.Println("we have a connection2")
	_ = client // we'll use this in the upcoming sections

	//tokenAddress := common.HexToAddress("0xa2aa7BE85977168Ec15dAF221f1407b32d5036b9")
	//instance, err := token.NewToken(tokenAddress, client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//address := common.HexToAddress("0xada53a094bD017D1Fc0c80Eb19d8DA67dfd477A9")
	//bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//name, err := instance.Name(&bind.CallOpts{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//symbol, err := instance.Symbol(&bind.CallOpts{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//decimals, err := instance.Decimals(&bind.CallOpts{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Printf("name: %s\n", name)         // "name: Golem Network"
	//fmt.Printf("symbol: %s\n", symbol)     // "symbol: GNT"
	//fmt.Printf("decimals: %v\n", decimals) // "decimals: 18"
	//
	//fmt.Printf("wei: %s\n", bal) // "wei: 74605500647408739782407023"
	//
	//fbal := new(big.Float)
	//fbal.SetString(bal.String())
	//value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	//
	//fmt.Printf("balance: %f", value) // "balance: 74605500.647409"

	return "success", nil
}

func ExchangeEURS(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to EURS blockchain
	fmt.Printf("Exchanging %d in EURS to %s from %s\n", amount, receiver, sender)
	return "example_hash", nil
}
