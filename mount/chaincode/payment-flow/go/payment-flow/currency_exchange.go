package payment_flow

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func ExchangeUSDT(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to Ethereum for exchanging tether
	fmt.Printf("Exchanging %d in USDT to %s from %s\n", amount, receiver, sender)

	url := "https://mainnet.infura.io/v3/189920b69bd147cbbee96ca2c36e5ea3"
	fmt.Println("URL:>", url)

	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Printf("response Status: %s", resp.Status)
	fmt.Printf("response Headers: %s", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("response Body: %s", string(body))

	return "example_hash", nil
}

func ExchangeEURS(amount int, receiver string, receiverAddress string, sender string, senderAddress string) (string, error) {
	// TODO: Add API connection to EURS blockchain
	fmt.Printf("Exchanging %d in EURS to %s from %s\n", amount, receiver, sender)
	return "example_hash", nil
}
