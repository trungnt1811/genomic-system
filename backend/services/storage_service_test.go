package services_test

import (
	"crypto/ecdsa"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/services"
)

func TestStoreGeneData(t *testing.T) {
	service := services.NewGeneDataStorageService()

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Create a test user ID and gene data
	userID := uint64(1)
	encryptedData := []byte("encrypted_gene_data")
	hashData := crypto.Keccak256Hash(encryptedData).Bytes()
	signature, err := crypto.Sign(crypto.Keccak256Hash(encryptedData).Bytes(), privateKey)
	require.NoError(t, err)

	// Test storing the gene data
	fileID, err := service.StoreGeneData(userID, encryptedData, signature, hashData)
	require.NoError(t, err)
	require.NotEmpty(t, fileID)

	// Attempt to store the same data again, should return an error
	_, err = service.StoreGeneData(userID, encryptedData, signature, hashData)
	require.Error(t, err)
	require.Equal(t, "gene data with the same hash already exists", err.Error())
}

func TestRetrieveGeneData(t *testing.T) {
	service := services.NewGeneDataStorageService()

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Create and store test gene data
	userID := uint64(1)
	encryptedData := []byte("encrypted_gene_data")
	hashData := crypto.Keccak256Hash(encryptedData).Bytes()
	signature, err := crypto.Sign(crypto.Keccak256Hash(encryptedData).Bytes(), privateKey)
	require.NoError(t, err)

	fileID, err := service.StoreGeneData(userID, encryptedData, signature, hashData)
	require.NoError(t, err)

	// Retrieve the gene data by file ID
	retrievedData, err := service.RetrieveGeneData(fileID)
	require.NoError(t, err)
	require.Equal(t, encryptedData, retrievedData)

	// Attempt to retrieve data with an invalid file ID
	_, err = service.RetrieveGeneData("invalid_file_id")
	require.Error(t, err)
	require.Equal(t, "gene data not found", err.Error())
}

func TestVerifyGeneDataSignature(t *testing.T) {
	service := services.NewGeneDataStorageService()

	// Generate new key pair
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	// Create and store test gene data
	userID := uint64(1)
	encryptedData := []byte("encrypted_gene_data")
	hashData := crypto.Keccak256Hash(encryptedData).Bytes()
	signature, err := crypto.Sign(hashData, privateKey)
	require.NoError(t, err)

	fileID, err := service.StoreGeneData(userID, encryptedData, signature, hashData)
	require.NoError(t, err)

	// Verify the signature
	isValid, err := service.VerifyGeneDataSignature(fileID, publicKeyBytes)
	require.NoError(t, err)
	require.True(t, isValid)

	// Verify the signature with an incorrect public key
	anotherPrivateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	anotherPublicKey := anotherPrivateKey.Public().(*ecdsa.PublicKey)

	isValid, err = service.VerifyGeneDataSignature(fileID, crypto.CompressPubkey(anotherPublicKey))
	require.NoError(t, err)
	require.False(t, isValid)

	// Attempt to verify signature with an invalid file ID
	isValid, err = service.VerifyGeneDataSignature("invalid_file_id", publicKeyBytes)
	require.Error(t, err)
	require.Equal(t, "gene data not found", err.Error())
	require.False(t, isValid)
}
