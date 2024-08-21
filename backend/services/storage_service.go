package services

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

// GeneData represents the structure to hold gene data and associated information.
type GeneData struct {
	FileID        string
	UserID        uint64
	DataHash      []byte // Hash of the encrypted gene data.
	Signature     []byte // Digital signature of the gene data.
	EncryptedData []byte // The encrypted gene data.
}

// GeneDataStorageService manages the storage of encrypted gene data that uses an in-memory map to simulate a NoSQL database.
type GeneDataStorageService struct {
	dataStore map[string]GeneData
	mu        sync.Mutex
}

// NewGeneDataStorageService creates a new instance of GeneDataStorageService.
func NewGeneDataStorageService() *GeneDataStorageService {
	return &GeneDataStorageService{
		dataStore: make(map[string]GeneData),
	}
}

// StoreGeneData stores the encrypted gene data, its hash, and the associated signature.
func (s *GeneDataStorageService) StoreGeneData(userID uint64, encryptedData []byte, signatureBytes []byte) (string, error) {
	// Calculate the SHA-256 hash of the encrypted data.
	dataHash := sha256.Sum256(encryptedData)

	// Use part of the hash as a unique file identifier.
	fileID := hex.EncodeToString(dataHash[:6])

	// Store the gene data in the in-memory map.
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.dataStore[fileID]; exists {
		return "", errors.New("gene data with the same hash already exists")
	}

	s.dataStore[fileID] = GeneData{
		FileID:        fileID,
		UserID:        userID,
		DataHash:      dataHash[:],
		Signature:     signatureBytes,
		EncryptedData: encryptedData,
	}

	return fileID, nil
}

// RetrieveGeneData retrieves the original encrypted gene data based on the file ID.
func (s *GeneDataStorageService) RetrieveGeneData(fileID string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.dataStore[fileID]
	if !exists {
		return nil, errors.New("gene data not found")
	}

	return data.EncryptedData, nil
}

// VerifyGeneDataSignature verifies the digital signature of the stored gene data.
func (s *GeneDataStorageService) VerifyGeneDataSignature(fileID string, publicKeyBytes []byte) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.dataStore[fileID]
	if !exists {
		return false, errors.New("gene data not found")
	}

	publicKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return false, errors.New("invalid public key")
	}

	// Verify the signature using the hash of the encrypted data and the public key.
	isValid := crypto.VerifySignature(crypto.CompressPubkey(publicKey), data.DataHash, data.Signature[:len(data.Signature)-1])

	return isValid, nil
}
