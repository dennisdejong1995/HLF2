package main

import payment_flow "github.com/dennisdejong1995/HLF2/mount/chaincode/payment-flow/go/payment-flow"

func main() {
	_, _ = payment_flow.ExchangeUSDT(1, "bla", "blabla", "blablabla", "blablalbla")
}
