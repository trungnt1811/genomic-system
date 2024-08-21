package tee_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	service "github.com/trungnt1811/blockchain-engineer-interview/backend/services/tee"
)

func TestEncryptGeneData(t *testing.T) {
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

func TestSignEncryptedGeneData(t *testing.T) {
	// Generate a new private key for testing.
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	require.NotNil(t, privateKey)

	// Initialize the TEEService.
	service := service.NewTEEService()

	// Define the test data.
	encryptedData := []byte("sample encrypted gene data")

	// Sign the encrypted data.
	hash, signature, err := service.SignEncryptedGeneData(privateKey, encryptedData)
	require.NoError(t, err)
	require.NotNil(t, hash)
	require.NotNil(t, signature)

	// Verify that the hash matches the expected Keccak256 hash.
	expectedHash := crypto.Keccak256Hash(encryptedData).Bytes()
	require.Equal(t, expectedHash, hash)

	// Verify the signature using the public key.
	publicKey := &privateKey.PublicKey
	valid := crypto.VerifySignature(crypto.CompressPubkey(publicKey), hash, signature[:len(signature)-1])
	require.True(t, valid)
}
