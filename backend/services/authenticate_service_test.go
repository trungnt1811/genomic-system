package services_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/services"
)

func TestRegisterUser(t *testing.T) {
	authService := services.NewAuthService()

	// Register a new user
	privateKeyHex, userID, err := authService.RegisterUser()
	require.NoError(t, err)
	require.NotEmpty(t, privateKeyHex)
	require.NotZero(t, userID)

	// Decode the private key from hex
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	require.NoError(t, err)

	// Extract the public key from the private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	require.NoError(t, err)
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)

	// Check if the user is correctly stored in the AuthService
	storedUser := authService.GetUserInfo(userID)
	require.Equal(t, userID, storedUser.UserID)
	require.Equal(t, publicKeyBytes, storedUser.PublicKey)
}

func TestAddExistingUser(t *testing.T) {
	authService := services.NewAuthService()

	// Generate a new key pair for the user
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	userID := uint64(1)

	// Add the user to the AuthService
	err = authService.AddExistingUser(userID, publicKeyBytes)
	require.NoError(t, err)

	// Check if the user is correctly stored in the AuthService
	storedUser := authService.GetUserInfo(userID)
	require.Equal(t, userID, storedUser.UserID)
	require.Equal(t, publicKeyBytes, storedUser.PublicKey)
}

func TestAuthenticate(t *testing.T) {
	authService := services.NewAuthService()

	// Register a new user
	privateKeyHex, userID, err := authService.RegisterUser()
	require.NoError(t, err)

	// Decode the private key from hex
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	require.NoError(t, err)

	// Extract the public key and Ethereum address
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	require.NoError(t, err)
	publicKey := &privateKey.PublicKey
	ethAddress := crypto.PubkeyToAddress(*publicKey).Hex()

	// Test successful authentication
	isAuthenticated := authService.Authenticate(userID, ethAddress)
	require.True(t, isAuthenticated)

	// Test failed authentication with wrong address
	isAuthenticated = authService.Authenticate(userID, "0x0000000000000000000000000000000000000000")
	require.False(t, isAuthenticated)

	// Test failed authentication with non-existent user ID
	isAuthenticated = authService.Authenticate(999, ethAddress)
	require.False(t, isAuthenticated)
}
