package auth

import (
	"encoding/hex"
	"errors"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// User represents a user with a unique ID and their public key as a byte slice.
type User struct {
	UserID    uint64
	PublicKey []byte
}

// AuthService is a mock authentication service that uses an in-memory map to simulate a NoSQL database.
type AuthService struct {
	usersDB map[uint64]User
	mu      sync.Mutex
}

// NewAuthService creates and returns a new instance of AuthService.
func NewAuthService() *AuthService {
	return &AuthService{
		usersDB: make(map[uint64]User),
	}
}

// RegisterUser generates a new Ethereum key pair, stores the user's public key, and returns the private key and user ID.
func (s *AuthService) RegisterUser() (string, uint64, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", 0, err
	}

	userID := rand.Uint64()

	// Convert the public key to a byte slice
	publicKeyBytes := crypto.FromECDSAPub(&privateKey.PublicKey)

	s.mu.Lock()
	s.usersDB[userID] = User{
		UserID:    userID,
		PublicKey: publicKeyBytes,
	}
	s.mu.Unlock()

	// Return the private key in hex format
	privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

	return privateKeyHex, userID, nil
}

// AddExistingUser adds an existing user to the service with a provided public key in hex format.
func (s *AuthService) AddExistingUser(userID uint64, publicKeyBytes []byte) error {
	s.mu.Lock()
	s.usersDB[userID] = User{
		UserID:    userID,
		PublicKey: publicKeyBytes,
	}
	s.mu.Unlock()

	return nil
}

// Authenticate checks if the provided user ID and Ethereum address match the stored information.
func (s *AuthService) Authenticate(userID uint64, ethAddress string) bool {
	if !common.IsHexAddress(ethAddress) {
		return false
	}

	s.mu.Lock()
	user, exists := s.usersDB[userID]
	s.mu.Unlock()

	if !exists {
		return false
	}

	// Convert the stored public key bytes back to an ecdsa.PublicKey
	publicKey, err := crypto.UnmarshalPubkey(user.PublicKey)
	if err != nil {
		return false
	}

	// Convert the public key to an Ethereum address
	publicKeyAddress := crypto.PubkeyToAddress(*publicKey).Hex()

	return publicKeyAddress == ethAddress
}

// GetUserInfo returns the user ID and the Ethereum address for the given user ID.
func (s *AuthService) GetUserInfo(userID uint64) (uint64, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.usersDB[userID]
	if !exists {
		return 0, "", errors.New("user not found")
	}

	// Convert the stored public key bytes back to an ecdsa.PublicKey
	publicKey, err := crypto.UnmarshalPubkey(user.PublicKey)
	if err != nil {
		return 0, "", err
	}

	// Convert the public key to an Ethereum address
	ethAddress := crypto.PubkeyToAddress(*publicKey).Hex()

	return user.UserID, ethAddress, nil
}

func (s *AuthService) QueryUserByUserID(userID uint64) User {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.usersDB[userID]
}
