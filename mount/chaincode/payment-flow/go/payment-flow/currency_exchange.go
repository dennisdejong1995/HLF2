package payment_flow

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/sha3"
	"hash"
	"io/ioutil"
	"net/http"
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

type KeccakState interface {
	hash.Hash
	Read([]byte) (int, error)
}

func Keccak256(data ...[]byte) []byte {
	b := make([]byte, 32)
	d := sha3.NewLegacyKeccak256().(KeccakState)
	for _, b := range data {
		d.Write(b)
	}
	d.Read(b)
	return b
}

func ExchangeUSDT(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to Ethereum for exchanging tether
	fmt.Printf("Exchanging %d in USDT to %s from %s\n", amount, receiver, sender)

	url := "https://rinkeby.infura.io/v3/189920b69bd147cbbee96ca2c36e5ea3"
	fmt.Printf("URL:>%s\n", url)

	callHash := Keccak256([]byte("balanceOf(address)"))
	fmt.Printf("callHash: %s\n", string(callHash))

	//jsonCall = `{"jsonrpc":"2.0","method":"eth_call","params": [{"to":"0xa2aa7BE85977168Ec15dAF221f1407b32d5036b9"}],"id":1}`
	//jsonCall := `{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`
	jsonCall := `{"jsonrpc":"2.0","method":"eth_call","params": [{"from":"0xada53a094bD017D1Fc0c80Eb19d8DA67dfd477A9", to":"0xa2aa7BE85977168Ec15dAF221f1407b32d5036b9", "data","` + string(callHash) + receiverAddress + `"}],"id":1}`
	var jsonStr = []byte(jsonCall)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("response Status: %s\n", resp.Status)
	fmt.Printf("response Headers: %s\n", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response Body: %s\n", string(body))

	return resp.Status, nil
}

func ExchangeEURS(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to EURS blockchain
	fmt.Printf("Exchanging %d in EURS to %s from %s\n", amount, receiver, sender)
	return "example_hash", nil
}
