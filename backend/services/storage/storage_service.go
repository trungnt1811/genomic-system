package storage

import (
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

// GeneDataStorageService manages the storage of encrypted gene data using an in-memory map to simulate a NoSQL database.
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
func (s *GeneDataStorageService) StoreGeneData(
	userID uint64,
	encryptedData []byte,
	signatureBytes []byte,
	hashBytes []byte,
) (string, error) {
	// Check the signature length (should be 65 bytes, including the recovery ID)
	if len(signatureBytes) != crypto.SignatureLength {
		return "", errors.New("invalid signature length")
	}

	dataHash := hashBytes

	// Use part of the hash as a unique file identifier
	fileID := hex.EncodeToString(dataHash[:16])

	// Store the gene data in the in-memory map
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

	// Extract the signature without the recovery ID
	signatureNoRecoverID := data.Signature[:len(data.Signature)-1]

	// Verify the signature using the public key and the data hash
	// Note: crypto.VerifySignature expects the signature to be 64 bytes (without the recovery ID)
	isValid := crypto.VerifySignature(publicKeyBytes, data.DataHash, signatureNoRecoverID)
	return isValid, nil
}
