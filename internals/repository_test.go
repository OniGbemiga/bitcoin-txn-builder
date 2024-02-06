package internals

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateRedeemScriptHex(t *testing.T) {
	preImage := "Btrust Builders"
	lockingScript := GenerateRedeemScriptHex(preImage)
	expectedScriptHex := hex.EncodeToString(lockingScript[:])
	assert.Equal(t, expectedScriptHex, hex.EncodeToString(lockingScript))
}

func TestDeriveAddress(t *testing.T) {
	preImage := "Btrust Builders"
	redeemScript := GenerateRedeemScriptHex(preImage)
	decodeHex, err := hex.DecodeString(string(redeemScript))
	assert.NoError(t, err)
	address, err := GetAddressFromRedeemScriptHex(decodeHex)
	assert.NoError(t, err)
	expectedAddress := "2N8hwP1WmJrFF5QWABn38y63uYLhnJYJYTF"
	assert.Equal(t, expectedAddress, address)
}

func TestConstructTransaction(t *testing.T) {
	privateKey := "KwPtUazuC9TgYZQ7ptkqLZrW4KmR9hcrgAJbmAzAwwP7b3UkeN9m"
	amount := int64(100000)
	redeemScriptHex := GenerateRedeemScriptHex("Btrust Builders")
	address := "2MytaPKkM6FYRt7PgUSSfwvMwYsHrQLbH9W"

	transactionHex, err := ConstructTransaction(redeemScriptHex, address, privateKey, amount)
	assert.NoError(t, err)

	expectedTransactionHex := "01000000013c9c9c1d8478f0230d715209de85f7b91e2623996203b8cb5e69c5e7c391499e00" +
		"0000004847304402203cb10cfe93201f53996bc8f9c107c807b8ab57f06da1da904e9192f2e6684a9c02203c63d76717f1" +
		"86cdbde3679987fb58c8074cb7a6c802c8b855507968434eb906010000000001804a5d050000000023324d797461504b6b" +
		"4d36465952743750675553536677764d7759734872514c6248395700000000"
	assert.Equal(t, expectedTransactionHex, transactionHex)
}

func TestConstructSpendingTransaction(t *testing.T) {
	privateKey := "KwPtUazuC9TgYZQ7ptkqLZrW4KmR9hcrgAJbmAzAwwP7b3UkeN9m"
	redeemScriptHex := GenerateRedeemScriptHex("Btrust Builders")
	changeAddress, err := GetAddressFromRedeemScriptHex(redeemScriptHex)
	assert.NoError(t, err)

	originalTransactionHex := "01000000013c9c9c1d8478f0230d715209de85f7b91e2623996203b8cb5e69c5e7c391499e00" +
		"0000004847304402203cb10cfe93201f53996bc8f9c107c807b8ab57f06da1da904e9192f2e6684a9c02203c63d76717f1" +
		"86cdbde3679987fb58c8074cb7a6c802c8b855507968434eb906010000000001804a5d050000000023324d797461504b6b" +
		"4d36465952743750675553536677764d7759734872514c6248395700000000"

	spendingTransactionHex, err := SpendTransactionOriginalTxHex(originalTransactionHex, redeemScriptHex, changeAddress, privateKey)
	assert.NoError(t, err)

	expectedSpendingTransactionHex := "010000000100000000000000000000000000000000000000000000000000000000000" +
		"00000000000004847304402202ee75860ead3fc57a838e6434c4f617d01c94d9bd431c79064d049ec0502bffa022030d838" +
		"bcf52fc648c73d0950e45ff0f7617cdedbea8d3030475aa562e5bbeed0010000000001804a5d050000000023324d7974615" +
		"04b6b4d36465952743750675553536677764d7759734872514c6248395700000000"
	assert.Equal(t, expectedSpendingTransactionHex, spendingTransactionHex)
}
