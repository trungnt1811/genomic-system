package auth

import (
	"errors"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// User represents a user with a unique ID and their public key as a byte slice.
// The User struct simulates a user entity in the system.
type User struct {
	UserID    uint64
	PublicKey []byte
}

// AuthService is a mock authentication service that uses an in-memory map to simulate a NoSQL database.
// This service manages user registration, authentication, and retrieval of user data.
type AuthService struct {
	usersDB map[uint64]User // In-memory map simulating a NoSQL database
	mu      sync.Mutex      // Mutex to ensure thread-safe operations on the usersDB map
}

// NewAuthService creates and returns a new instance of AuthService.
// This function initializes the usersDB map for storing user data.
func NewAuthService() *AuthService {
	return &AuthService{
		usersDB: make(map[uint64]User),
	}
}

// RegisterUserWithPubkey registers a new user by storing their public key and generating a unique user ID.
// It stores the user information in the usersDB map and returns the generated user ID.
func (s *AuthService) RegisterUserWithPubkey(publicKeyBytes []byte) uint64 {
	// Generate a random unique user ID
	userID := rand.Uint64()

	// Lock the mutex to ensure thread-safe access to usersDB
	s.mu.Lock()
	// Store the user ID and public key in the usersDB map
	s.usersDB[userID] = User{
		UserID:    userID,
		PublicKey: publicKeyBytes,
	}
	// Unlock the mutex after updating usersDB
	s.mu.Unlock()

	// Return the generated user ID
	return userID
}

// Authenticate checks if the provided user ID and Ethereum address match the stored information.
// It validates the Ethereum address and compares it with the one derived from the stored public key.
func (s *AuthService) Authenticate(userID uint64, ethAddress string) bool {
	// Check if the provided Ethereum address is a valid hex address
	if !common.IsHexAddress(ethAddress) {
		return false
	}

	// Lock the mutex to ensure thread-safe access to usersDB
	s.mu.Lock()
	// Retrieve the user information from usersDB using the provided user ID
	user, exists := s.usersDB[userID]
	s.mu.Unlock()

	// If the user does not exist, return false
	if !exists {
		return false
	}

	// Convert the stored public key bytes back to an ecdsa.PublicKey
	publicKey, err := crypto.UnmarshalPubkey(user.PublicKey)
	if err != nil {
		return false
	}

	// Derive the Ethereum address from the public key
	publicKeyAddress := crypto.PubkeyToAddress(*publicKey).Hex()

	// Compare the derived Ethereum address with the provided one
	return publicKeyAddress == ethAddress
}

// GetUserPubkey returns the public key bytes for the given user ID.
// It retrieves the user's public key from the usersDB map.
func (s *AuthService) GetUserPubkey(userID uint64) ([]byte, error) {
	// Lock the mutex to ensure thread-safe access to usersDB
	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieve the user information from usersDB using the provided user ID
	user, exists := s.usersDB[userID]
	if !exists {
		return nil, errors.New("user not found")
	}

	// Return the user's public key bytes
	return user.PublicKey, nil
}

// QueryUserByUserID retrieves the user information associated with the provided user ID.
// It returns the User struct containing the user ID and public key.
func (s *AuthService) QueryUserByUserID(userID uint64) User {
	// Lock the mutex to ensure thread-safe access to usersDB
	s.mu.Lock()
	defer s.mu.Unlock()

	// Return the User struct for the given user ID
	return s.usersDB[userID]
}
