package internals

import (
	"bytes"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

type TxnBuilderRepository interface {
	GenerateRedeemScriptHex(preImage string) []byte
	GetAddressFromRedeemScriptHex(redeemScriptHex []byte) (string, error)
	ConstructTransaction(redeemScriptHex []byte, address string, privateKey string, amount int64) (string, error)
	SpendTransactionOriginalTxHex(txnHex string, redeemScriptHex []byte, changeAddress, privateKey string) (string, error)
}

func GenerateRedeemScriptHex(preImage string) []byte {
	script, _ := txscript.NewScriptBuilder().
		AddOp(txscript.OP_SHA256).
		AddData([]byte(preImage)).
		AddOp(txscript.OP_EQUAL).
		Script()

	return script
}

func GetAddressFromRedeemScriptHex(redeemScriptHex []byte) (string, error) {
	// Create a P2SH address
	redeemScriptAddr, err := btcutil.NewAddressScriptHash(redeemScriptHex, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	return redeemScriptAddr.EncodeAddress(), nil
}

func ConstructTransaction(redeemScriptHex []byte, address string, privateKey string, amount int64) (string, error) {
	//derive the private Key
	privateKeyEn := parsePrivateKey(privateKey)

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add an output with the P2SH redeem script
	tx.AddTxOut(&wire.TxOut{
		Value:    amount,
		PkScript: redeemScriptHex,
	})

	// Get the transaction hash
	txHash := tx.TxHash()

	// Create the funding transaction
	fundingTx := wire.NewMsgTx(wire.TxVersion)
	fundingTx.AddTxIn(&wire.TxIn{
		PreviousOutPoint: *wire.NewOutPoint(&txHash, 0),
		SignatureScript:  redeemScriptHex,
	})
	fundingTx.AddTxOut(&wire.TxOut{
		Value:    90000000, // 0.9 BTC change
		PkScript: []byte(address),
	})

	// Sign the funding transaction
	signature, err := txscript.RawTxInSignature(fundingTx, 0, redeemScriptHex, txscript.SigHashAll, privateKeyEn)
	if err != nil {
		log.Fatal("Error signing transaction:", err)
	}
	fundingTx.TxIn[0].SignatureScript, _ = txscript.NewScriptBuilder().
		AddData(signature).
		Script()

	// Serialize the transaction to hex
	buf := new(bytes.Buffer)
	err = fundingTx.Serialize(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func SpendTransactionOriginalTxHex(txnHex string, redeemScriptHex []byte, changeAddress, privateKey string) (string, error) {

	//derive the private Key
	privateKeyEn := parsePrivateKey(privateKey)

	//derive the output
	prevTxHash, err := chainhash.NewHashFromStr(txnHex)
	if err != nil {
		log.Fatal(err)
	}
	// Create a new wire.OutPoint instance
	outpoint := wire.NewOutPoint(prevTxHash, 0)

	// Create a spending transaction
	spendingTx := wire.NewMsgTx(wire.TxVersion)

	// Add an input referencing the output of the funding transaction
	spendingTx.AddTxIn(&wire.TxIn{
		PreviousOutPoint: *outpoint,
		SignatureScript:  redeemScriptHex,
	})

	// Add an output
	spendingTx.AddTxOut(&wire.TxOut{
		Value:    90000000, // 0.9 BTC
		PkScript: []byte(changeAddress),
	})

	// Sign the spending transaction
	signature, err := txscript.RawTxInSignature(spendingTx, 0, redeemScriptHex, txscript.SigHashAll, privateKeyEn)
	if err != nil {
		log.Fatal(err)
	}
	spendingTx.TxIn[0].SignatureScript, _ = txscript.NewScriptBuilder().
		AddData(signature).
		Script()

	// Serialize the spending transaction to hex
	buf := new(bytes.Buffer)
	err = spendingTx.Serialize(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil

}

func signTx(tx *wire.MsgTx, privKey *btcec.PrivateKey, redeemScript []byte) (*wire.MsgTx, *wire.OutPoint) {
	txid := tx.TxHash()
	signature, err := txscript.RawTxInSignature(tx, 0, redeemScript, txscript.SigHashAll, privKey)
	if err != nil {
		log.Fatal(err)
	}
	tx.TxIn[0].SignatureScript, _ = txscript.NewScriptBuilder().
		AddData(signature).
		Script()

	return tx, wire.NewOutPoint(&txid, 0)
}

func parsePrivateKey(privateKeyHx string) *btcec.PrivateKey {
	// Parse the private key
	privateKey, _ := btcec.PrivKeyFromBytes([]byte(privateKeyHx))
	return privateKey
}
