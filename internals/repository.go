package internals

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
)

type TxnBuilderRepository interface {
	GenerateRedeemScriptHex(preImage string) string
	GetAddressFromRedeemScriptHex(redeemScriptHex string) (string, error)
	ConstructTransaction(address string, privateKey *ecdsa.PrivateKey, amount int64) (string, error)
	SpendTransactionOriginalTxHex(txHex string, redeemScriptHex string, destAddress string,
		changeAddress string, privateKey *ecdsa.PrivateKey) (string, error)
}

func GenerateRedeemScriptHex(preImage string) string {
	lockingScript := sha256.Sum256([]byte(preImage))
	return fmt.Sprintf("OP_SHA256%xOP_EQUAL", lockingScript)
}

func GetAddressFromRedeemScriptHex(redeemScriptHex string) (string, error) {
	redeemScript, err := hex.DecodeString(fmt.Sprintf("%x", redeemScriptHex))
	if err != nil {
		return "", err
	}

	redeemScriptAddr, err := btcutil.NewAddressScriptHash(redeemScript, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}

	return redeemScriptAddr.EncodeAddress(), nil
}

func ConstructTransaction(address string, privateKey *btcec.PrivateKey, amount int64) (string, error) {

	keyDB := &SimpleKeyDB{privateKey: privateKey}

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add input to the transaction (unspent output from a previous transaction)
	prevOutHash, _ := chainhash.NewHashFromStr("67868b3b2b8e2c85893bd82e263c940bb18ee1371ffe4602313f41ddd675c590")
	prevOutPoint := wire.NewOutPoint(prevOutHash, 0)
	txIn := wire.NewTxIn(prevOutPoint, nil, nil)
	tx.AddTxIn(txIn)

	// Add output to the transaction (sending to the derived address)
	decodedAddress, err := btcutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	pkScript, _ := txscript.PayToAddrScript(decodedAddress)
	txOut := wire.NewTxOut(amount, pkScript)
	tx.AddTxOut(txOut)

	redeemScript := GenerateRedeemScriptHex("Btrust Builders")

	// sign the transaction
	signTx(tx, redeemScript, keyDB)

	// Serialize the transaction to hex
	buf := new(bytes.Buffer)
	err = tx.Serialize(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func SpendTransactionOriginalTxHex(txnHex string, redeemScriptHex string, destAddress string,
	changeAddress string, privateKey *btcec.PrivateKey) (string, error) {

	keyDB := &SimpleKeyDB{privateKey: privateKey}

	// Deserialize the original transaction
	originalTxBytes, _ := hex.DecodeString(txnHex)
	originalTx := wire.MsgTx{}
	err := originalTx.Deserialize(bytes.NewReader(originalTxBytes))
	if err != nil {
		return "", err
	}

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Get the transaction hash
	txHash := originalTx.TxHash()

	// Add input to the transaction (using the output of the previous transaction)
	outPoint := wire.NewOutPoint(&txHash, 0)
	txIn := wire.NewTxIn(outPoint, nil, nil)
	tx.AddTxIn(txIn)

	// Add output to the transaction (spending to the destination address)
	decodedDestAddress, err := btcutil.DecodeAddress(destAddress, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	destPkScript, _ := txscript.PayToAddrScript(decodedDestAddress)
	txOutDest := wire.NewTxOut(originalTx.TxOut[0].Value-1000, destPkScript)
	tx.AddTxOut(txOutDest)

	// Add another output for change
	decodedChangeAddress, err := btcutil.DecodeAddress(changeAddress, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}
	changePkScript, _ := txscript.PayToAddrScript(decodedChangeAddress)
	txOutChange := wire.NewTxOut(1000, changePkScript)
	tx.AddTxOut(txOutChange)

	//sign the transaction
	signTx(tx, redeemScriptHex, keyDB)

	// Serialize the spending transaction to hex
	buf := new(bytes.Buffer)
	err = tx.Serialize(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil

}

// SimpleKeyDB is a simple implementation of the KeyDB interface.
type SimpleKeyDB struct {
	privateKey *btcec.PrivateKey
}

// GetKey retrieves the private key associated with a given address.
func (k *SimpleKeyDB) GetKey(address btcutil.Address) (*btcec.PrivateKey, bool, error) {
	return k.privateKey, true, nil
}

func signTx(tx *wire.MsgTx, redeemScriptHex string, keyDB txscript.KeyDB) {
	redeemScript, err := hex.DecodeString(redeemScriptHex)
	if err != nil {
		log.Fatal("Error decoding redeem script:", err)
	}

	hashType := txscript.SigHashAll
	sigScript, err := txscript.SignTxOutput(&chaincfg.TestNet3Params, tx, 0, redeemScript, hashType, keyDB, nil, nil)
	if err != nil {
		log.Fatal("Error signing transaction:", err)
	}

	tx.TxIn[0].SignatureScript = sigScript
}
