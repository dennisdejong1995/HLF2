package payment_flow

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dennisdejong1995/HLF2/mount/chaincode/payment-flow/go/contracts_erc20/SUSD"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
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

	accounts := genAccounts()

	dealblock := accounts[2]

	nonce, err := client.PendingNonceAt(context.Background(), dealblock.address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(dealblock.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	address, tx, instance, err := SUSD.DeploySUSD(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())

	_ = instance
}

func ExchangeEURS(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to EURS blockchain
	fmt.Printf("Exchanging %d in EURS to %s from %s\n", amount, receiver, sender)
	return "example_hash", nil
}

type Account struct {
	holder     string
	privateKey *ecdsa.PrivateKey
	address    common.Address
}

func genAccounts() []Account {

	var accounts []Account

	pk1, err := crypto.HexToECDSA("f238a37e42b7062bdbc062a1833a6361f9a6d0e324a95ca2f7c4c3034e67ee5c")
	if err != nil {
		fmt.Println("error")
	}

	a1 := Account{holder: "Org1"}
	a1.privateKey = pk1
	a1.address = common.HexToAddress("0x559BC07434C89c5496d790DFD2885dC966F9113a")

	accounts = append(accounts, a1)

	pk2, err := crypto.HexToECDSA("6c0081a5b9511910a6cec018a99d3031197f079cde51c1a78124750a990cdd08")
	if err != nil {
		fmt.Println("error")
	}

	a2 := Account{holder: "Org2"}
	a2.privateKey = pk2
	a2.address = common.HexToAddress("0x54806DD512b21814aa560D627432a75720ed6bB3")

	accounts = append(accounts, a2)

	pk3, err := crypto.HexToECDSA("a3968111221303d38214eb7b2db9b04cefb2300b72771d65d258e08322dc573d")
	if err != nil {
		fmt.Println("error")
	}

	a3 := Account{holder: "Dealblock"}
	a3.privateKey = pk3
	a3.address = common.HexToAddress("0x6dc89393fa30b64c56deff31daacf10cedcd852d")

	accounts = append(accounts, a3)

	return accounts
}
