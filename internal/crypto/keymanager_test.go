package crypto_test

import (
	"testing"

	"github.com/abaxxtech/abaxx-id-go/internal/crypto"
	"github.com/abaxxtech/abaxx-id-go/internal/crypto/dsa"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	keyManager := crypto.NewLocalKeyManager()

	keyID, err := keyManager.GeneratePrivateKey(dsa.AlgorithmIDSECP256K1)
	assert.NoError(t, err)

	assert.True(t, keyID != "", "keyID is empty")
}

func TestGetPublicKey(t *testing.T) {
	keyManager := crypto.NewLocalKeyManager()

	keyID, err := keyManager.GeneratePrivateKey(dsa.AlgorithmIDSECP256K1)
	assert.NoError(t, err)

	publicKey, err := keyManager.GetPublicKey(keyID)
	assert.NoError(t, err)

	thumbprint, err := publicKey.ComputeThumbprint()
	assert.NoError(t, err)

	assert.Equal(t, keyID, thumbprint, "unexpected keyID")
}

func TestSign(t *testing.T) {
	keyManager := crypto.NewLocalKeyManager()

	keyID, err := keyManager.GeneratePrivateKey(dsa.AlgorithmIDSECP256K1)
	assert.NoError(t, err)

	payload := []byte("hello world")
	signature, err := keyManager.Sign(keyID, payload)
	assert.NoError(t, err)

	if signature == nil {
		t.Errorf("signature is nil")
	}

	assert.True(t, signature != nil, "signature is nil")
}
