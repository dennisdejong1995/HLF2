package payment_flow

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dennisdejong1995/HLF2/contracts_erc20/SUSD"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

const infuraKEY = "189920b69bd147cbbee96ca2c36e5ea3"
const network = "rinkeby"
const CONTRACT_ADDRESS = "0x566160489E6548120385Fd5397BEf0E2b608C602"

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

	accounts := genAccounts()

	transfer(accounts[2], accounts[0].address, 100)

	org1_ball := queryBalance(accounts[0])
	org2_ball := queryBalance(accounts[1])
	deal_ball := queryBalance(accounts[2])

	fmt.Println(org1_ball)
	fmt.Println(org2_ball)
	fmt.Println(deal_ball)

	org1_allow1 := queryAllowance(accounts[0], accounts[2].address)
	org2_allow1 := queryAllowance(accounts[1], accounts[2].address)

	fmt.Println(org1_allow1)
	fmt.Println(org2_allow1)

	aproveAllowance(accounts[0], accounts[2].address, 300)

	transferFrom(accounts[2], accounts[0].address, accounts[1].address, 100)
	return "success", nil
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

func deployContract() {

	accounts := genAccounts()

	dealblock := accounts[2]

	client, err := ethclient.Dial("https://" + network + ".infura.io/v3/" + infuraKEY)
	if err != nil {
		log.Fatal(err)
	}

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

func queryBalance(owner Account) *big.Int {
	client, err := ethclient.Dial("https://" + network + ".infura.io/v3/" + infuraKEY)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contract, err := SUSD.NewSUSD(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	value, _ := contract.BalanceOf(&bind.CallOpts{}, owner.address)

	return value

}

func queryAllowance(owner Account, spender common.Address) *big.Int {
	client, err := ethclient.Dial("https://" + network + ".infura.io/v3/" + infuraKEY)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contract, err := SUSD.NewSUSD(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	value, _ := contract.Allowance(&bind.CallOpts{}, owner.address, spender)

	return value

}

func aproveAllowance(owner Account, spender common.Address, amount int64) {
	client, err := ethclient.Dial("https://" + network + ".infura.io/v3/" + infuraKEY)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contract, err := SUSD.NewSUSD(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), owner.address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(owner.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	contract.Approve(auth, spender, big.NewInt(amount))

}

func transfer(owner Account, to common.Address, amount int64) {
	client, err := ethclient.Dial("https://" + network + ".infura.io/v3/" + infuraKEY)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contract, err := SUSD.NewSUSD(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), owner.address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(owner.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	contract.Transfer(auth, to, big.NewInt(amount))

}

func transferFrom(owner Account, from common.Address, to common.Address, amount int64) {
	client, err := ethclient.Dial("https://" + network + ".infura.io/v3/" + infuraKEY)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
	}

	contract, err := SUSD.NewSUSD(common.HexToAddress(CONTRACT_ADDRESS), client)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), owner.address)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(owner.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice

	contract.TransferFrom(auth, from, to, big.NewInt(amount))
}
