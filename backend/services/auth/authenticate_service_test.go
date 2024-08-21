package auth_test

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	service "github.com/trungnt1811/blockchain-engineer-interview/backend/services/auth"
)

func TestRegisterUser(t *testing.T) {
	authService := service.NewAuthService()

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
	storedUser := authService.QueryUserByUserID(userID)
	require.Equal(t, userID, storedUser.UserID)
	require.Equal(t, publicKeyBytes, storedUser.PublicKey)
}

func TestAddExistingUser(t *testing.T) {
	authService := service.NewAuthService()

	// Generate a new key pair for the user
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)
	userID := uint64(1)

	// Add the user to the AuthService
	err = authService.AddExistingUser(userID, publicKeyBytes)
	require.NoError(t, err)

	// Check if the user is correctly stored in the AuthService
	storedUser := authService.QueryUserByUserID(userID)
	require.Equal(t, userID, storedUser.UserID)
	require.Equal(t, publicKeyBytes, storedUser.PublicKey)
}

func TestAuthenticate(t *testing.T) {
	authService := service.NewAuthService()

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

func TestGetUserInfo(t *testing.T) {
	authService := service.NewAuthService()

	// Register a new user
	privateKeyHex, userID, err := authService.RegisterUser()
	require.NoError(t, err)
	require.NotEmpty(t, privateKeyHex)
	require.NotZero(t, userID)

	// Retrieve user info
	retrievedUserID, ethAddress, err := authService.GetUserInfo(userID)
	require.NoError(t, err)
	require.Equal(t, userID, retrievedUserID)
	require.NotEmpty(t, ethAddress)

	// Verify that the Ethereum address matches the public key derived from the private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	require.NoError(t, err)

	expectedEthAddress := crypto.PubkeyToAddress(privateKey.PublicKey).Hex()
	require.Equal(t, expectedEthAddress, ethAddress)

	// Test case where the user does not exist
	_, _, err = authService.GetUserInfo(999999)
	require.Error(t, err)
	require.Equal(t, "user not found", err.Error())
}
