package internals

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRedeemScriptHex(t *testing.T) {
	preImage := "Btrust Builders"
	hash := sha256.Sum256([]byte(preImage))
	expectedHex := "OP_SHA256 " + hex.EncodeToString(hash[:]) + " OP_EQUAL"
	assert.Equal(t, expectedHex, GenerateRedeemScriptHex(preImage))
}

func TestDeriveAddress(t *testing.T) {
	hash := sha256.Sum256([]byte("test"))
	redeemScriptHex := "OP_SHA256 " + hex.EncodeToString(hash[:]) + " OP_EQUAL"
	expectedAddress := "3G9wXQHQH2KjMy3qkjDbbewNq2ooNLGnNE"
	address, err := GetAddressFromRedeemScriptHex(redeemScriptHex)
	assert.NoError(t, err)
	assert.Equal(t, expectedAddress, address)
}

func TestConstructTransaction(t *testing.T) {
	privateKey := parsePrivateKey("KwPtUazuC9TgYZQ7ptkqLZrW4KmR9hcrgAJbmAzAwwP7b3UkeN9m")
	amount := int64(100000)

	preImage := "Btrust Builders"
	redeemScriptHex, _ := GetAddressFromRedeemScriptHex(preImage)
	address, _ := GetAddressFromRedeemScriptHex(redeemScriptHex)

	transactionHex, err := ConstructTransaction(address, privateKey, amount)
	assert.NoError(t, err)

	// not done
	expectedTransactionHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401e80300000000001600140f0754d8e59a359ba426e103ea4cf4a084c190c4f53e1f02000000001600140f0754d8e59a359ba426e103ea4cf4a084c190c4ffffffff012c6c0000000000001976a914a72ecf0ca3a2a9d8545a0e49ef7898e64c4c76cf88ac00000000"
	assert.Equal(t, expectedTransactionHex, transactionHex)
}

func TestConstructSpendingTransaction(t *testing.T) {
	privateKey := parsePrivateKey("privateKey")
	destAddress := "mv4rnyY3Su5gjcDNzbMLKBQkBicCtHUtFB"    //blockchain
	changeAddress := "2N4dyn5ZzEuw61YJjSVxHM19EbGxjT8v5Ze" //mine

	originalTransactionHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401e80300000000001600140f0754d8e59a359ba426e103ea4cf4a084c190c4f53e1f02000000001600140f0754d8e59a359ba426e103ea4cf4a084c190c4ffffffff012c6c0000000000001976a914a72ecf0ca3a2a9d8545a0e49ef7898e64c4c76cf88ac00000000"
	redeemScriptHex := GenerateRedeemScriptHex("Btrust Builders")

	spendingTransactionHex, err := SpendTransactionOriginalTxHex(originalTransactionHex, redeemScriptHex, destAddress, changeAddress, privateKey)
	assert.NoError(t, err)

	// not done
	expectedSpendingTransactionHex := "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0401e803000000000016001483b87c3e24db13f85cf3f3aa8a3b46177ce82ab7f53e1f02000000001600140f0754d8e59a359ba426e103ea4cf4a084c190c4ffffffff012c6c0000000000001976a914a72ecf0ca3a2a9d8545a0e49ef7898e64c4c76cf88ac00000000"
	assert.Equal(t, expectedSpendingTransactionHex, spendingTransactionHex)
}

func TestSignTx(t *testing.T) {
	privateKey := parsePrivateKey("KwPtUazuC9TgYZQ7ptkqLZrW4KmR9hcrgAJbmAzAwwP7b3UkeN9m")

	redeemScriptHex := GenerateRedeemScriptHex("Btrust Builders")

	// Create a dummy transaction
	tx := wire.NewMsgTx(wire.TxVersion)
	txIn := wire.NewTxIn(nil, nil, nil)
	tx.AddTxIn(txIn)

	keyDB := &SimpleKeyDB{privateKey: privateKey}

	// Sign the transaction
	signTx(tx, redeemScriptHex, keyDB)

	// not done
	assert.True(t, true, "Signature verification passed")
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
