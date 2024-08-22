package auth_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/services/auth"
)

func TestRegisterUserWithPubkey(t *testing.T) {
	// Initialize AuthService
	authService := auth.NewAuthService()

	// Generate a new ECDSA key pair
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err, "Failed to generate ECDSA key")

	// Extract the public key bytes
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)

	// Register the user with the public key
	userID := authService.RegisterUserWithPubkey(publicKeyBytes)

	// Query the registered user and verify the public key is stored correctly
	registeredUser := authService.QueryUserByUserID(userID)
	assert.Equal(t, publicKeyBytes, registeredUser.PublicKey, "Public key mismatch")
}

func TestAuthenticate(t *testing.T) {
	// Initialize AuthService
	authService := auth.NewAuthService()

	// Generate a new ECDSA key pair
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err, "Failed to generate ECDSA key")

	// Extract the public key bytes and Ethereum address
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	ethAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()

	// Register the user with the public key
	userID := authService.RegisterUserWithPubkey(publicKeyBytes)

	// Test successful authentication
	isAuthenticated := authService.Authenticate(userID, ethAddress)
	assert.True(t, isAuthenticated, "User should be authenticated successfully")

	// Test unsuccessful authentication with a wrong Ethereum address
	wrongEthAddress := "0x0000000000000000000000000000000000000000"
	isAuthenticated = authService.Authenticate(userID, wrongEthAddress)
	assert.False(t, isAuthenticated, "Authentication should fail with a wrong Ethereum address")
}

func TestGetUserPubkey(t *testing.T) {
	// Initialize AuthService
	authService := auth.NewAuthService()

	// Generate a new ECDSA key pair
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err, "Failed to generate ECDSA key")

	// Extract the public key bytes
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)

	// Register the user with the public key
	userID := authService.RegisterUserWithPubkey(publicKeyBytes)

	// Retrieve the public key for the registered user
	retrievedPubkey, err := authService.GetUserPubkey(userID)
	assert.NoError(t, err, "Failed to get user public key")
	assert.Equal(t, publicKeyBytes, retrievedPubkey, "Public key mismatch")

	// Attempt to retrieve the public key for a non-existent user
	_, err = authService.GetUserPubkey(9999)
	assert.Error(t, err, "Expected error for non-existent user")
}
