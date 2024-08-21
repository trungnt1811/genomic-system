package tee_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	service "github.com/trungnt1811/blockchain-engineer-interview/backend/services/tee"
)

func TestEncryptAndDecryptGeneData(t *testing.T) {
	// Initialize the TEE service
	teeService := service.NewTEEService()

	// Generate a new ECDSA private key
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Extract the public key from the private key
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)

	// Sample gene data to encrypt and decrypt
	originalGeneData := "This is a test gene data."

	// Encrypt the gene data using the public key
	encryptedData, err := teeService.EncryptGeneData(publicKeyBytes, originalGeneData)
	require.NoError(t, err)
	require.NotEmpty(t, encryptedData)

	// Decrypt the gene data using the private key
	decryptedData, err := teeService.DecryptGeneData(privateKey, encryptedData)
	require.NoError(t, err)
	require.Equal(t, originalGeneData, decryptedData)
}

func TestEncryptGeneData_InvalidPublicKey(t *testing.T) {
	// Initialize the TEE service
	teeService := service.NewTEEService()

	// Provide an invalid public key
	invalidPublicKey := []byte("invalid public key")

	// Sample gene data to encrypt
	geneData := "This is a test gene data."

	// Attempt to encrypt with the invalid public key
	_, err := teeService.EncryptGeneData(invalidPublicKey, geneData)
	require.Error(t, err)
}

func TestDecryptGeneData_InvalidEncryptedData(t *testing.T) {
	// Initialize the TEE service
	teeService := service.NewTEEService()

	// Generate a new ECDSA private key
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Provide invalid encrypted data (too short to contain a valid nonce)
	invalidEncryptedData := []byte("short data")

	// Attempt to decrypt the invalid encrypted data
	_, err = teeService.DecryptGeneData(privateKey, invalidEncryptedData)
	require.Error(t, err)
}
