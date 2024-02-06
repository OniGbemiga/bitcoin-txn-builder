package main

import (
	"flag"
	"fmt"
	"github/OniGbemiga/bitcoin-txn-builder/internals"
	"log"
)

func main() {
	// Parse command line arguments
	taskPtr := flag.String("task", "", "Specify the task: redeem, derive, transaction, spending")
	preImagePtr := flag.String("preimage", "", "Specify the pre-image for redeeming")
	destAddressPtr := flag.String("dest", "", "Specify the destination address for transaction")
	amountPtr := flag.Int64("amount", 0, "Specify the amount for transaction")
	privateKeyPtr := flag.String("private-key", "", "Specify the private key for transaction")
	txHexPtr := flag.String("txhex", "", "Specify the transaction hex for transaction")
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
		destAddress := *destAddressPtr
		privateKey := *privateKeyPtr
		redeemScriptHex := internals.GenerateRedeemScriptHex(*preImagePtr)
		transactionHex, err := internals.ConstructTransaction(redeemScriptHex, destAddress, privateKey, *amountPtr)
		if err != nil {
			log.Fatal("Error constructing transaction:", err)
		}
		fmt.Println("Constructed Transaction Hex:", transactionHex)

	case "spending":
		privateKey := *privateKeyPtr
		originalTransactionHex := *txHexPtr
		redeemScriptHex := internals.GenerateRedeemScriptHex(*preImagePtr)
		changeAddress, err := internals.GetAddressFromRedeemScriptHex(redeemScriptHex)
		if err != nil {
			log.Fatal("Error deriving change address:", err)
		}
		spendingTransactionHex, err := internals.SpendTransactionOriginalTxHex(originalTransactionHex, redeemScriptHex, changeAddress, privateKey)
		if err != nil {
			log.Fatal("Error constructing spending transaction:", err)
		}
		fmt.Println("Constructed Spending Transaction Hex:", spendingTransactionHex)

	default:
		fmt.Println("Invalid task. Please specify a valid task: redeem, derive, transaction, spending")
	}
}

//# Example for redeem task
//go run main.go -task=redeem -preimage="Btrust Builders"
//
//# Example for derive task
//go run main.go -task=derive -preimage="Btrust Builders"
//
//# Example for transaction task
//go run main.go -task=transaction -preimage="Btrust Builders" -dest="" private-key="" -amount=100000
//
//# Example for spending task
//go run main.go -task=spending -preimage="Btrust Builders" private-key="", txHexPtr=""

//2N4dyn5ZzEuw61YJjSVxHM19EbGxjT8v5Ze
//2NEGYjZXxXxM8Vg4X1MpfRCxSdjrzM8sZgT

//private-key
