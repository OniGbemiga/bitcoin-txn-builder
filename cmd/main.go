package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github/OniGbemiga/bitcoin-txn-builder/internals"
	"log"
)

func main() {
	// Parse command line arguments
	taskPtr := flag.String("task", "", "Specify the task: redeem, derive, transaction, spending")
	preImagePtr := flag.String("preimage", "", "Specify the pre-image for redeeming")
	destAddressPtr := flag.String("dest", "", "Specify the destination address for spending")
	changeAddressPtr := flag.String("change", "", "Specify the change address for spending")
	amountPtr := flag.Int64("amount", 0, "Specify the amount for transaction")
	privateKey := flag.String("private-key", "", "Specify the private key for transaction")
	flag.Parse()

	switch *taskPtr {
	case "redeem":
		redeemScriptHex := internals.GenerateRedeemScriptHex(*preImagePtr)
		fmt.Println("Redeem Script Hex:", redeemScriptHex)

	case "derive":
		redeemScriptHex := internals.GenerateRedeemScriptHex(*preImagePtr)
		address, err := internals.GetAddressFromRedeemScriptHex(redeemScriptHex)
		if err != nil {
			log.Fatal("Error deriving address:", err)
		}
		fmt.Println("Derived Address:", address)

	case "transaction":
		privateKey := parsePrivateKey(*privateKey)
		address, _ := internals.GetAddressFromRedeemScriptHex(internals.GenerateRedeemScriptHex(*preImagePtr))
		transactionHex, err := internals.ConstructTransaction(address, privateKey, *amountPtr)
		if err != nil {
			log.Fatal("Error constructing transaction:", err)
		}
		fmt.Println("Constructed Transaction Hex:", transactionHex)

	case "spending":
		privateKey := parsePrivateKey(*privateKey)
		originalTransactionHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0" +
			"401e80300000000001600140f0754d8e59a359ba426e103ea4cf4a084c190c4f53e1f02000000001600140f0754d8e59a359ba426" +
			"e103ea4cf4a084c190c4ffffffff012c6c0000000000001976a914a72ecf0ca3a2a9d8545a0e49ef7898e64c4c76cf88ac00000000"
		redeemScriptHex := internals.GenerateRedeemScriptHex(*preImagePtr)
		destAddress := *destAddressPtr
		changeAddress := *changeAddressPtr
		spendingTransactionHex, err := internals.SpendTransactionOriginalTxHex(originalTransactionHex, redeemScriptHex, destAddress, changeAddress, privateKey)
		if err != nil {
			log.Fatal("Error constructing spending transaction:", err)
		}
		fmt.Println("Constructed Spending Transaction Hex:", spendingTransactionHex)

	default:
		fmt.Println("Invalid task. Please specify a valid task: redeem, derive, transaction, spending")
	}
}

func parsePrivateKey(privateKeyHx string) *btcec.PrivateKey {
	// Decode the private key from hex
	privateKeyBytes, err := hex.DecodeString(privateKeyHx)
	if err != nil {
		return nil
	}

	// Parse the private key
	privateKey, _ := btcec.PrivKeyFromBytes(privateKeyBytes)

	return privateKey

}

//# Example for redeem task
//go run main.go -task=redeem -preimage="Btrust Builders"
//
//# Example for derive task
//go run main.go -task=derive -preimage="Btrust Builders"
//
//# Example for transaction task
//go run main.go -task=transaction -preimage="Btrust Builders" -amount=100000
//
//# Example for spending task
//go run main.go -task=spending -preimage="Btrust Builders" -dest="your_destination_address" -change="your_change_address"
