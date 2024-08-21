package tee

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"
)

// TEEService simulates a Trusted Execution Environment (TEE) service.
type TEEService struct{}

// NewTEEService creates a new instance of TEEService.
func NewTEEService() *TEEService {
	return &TEEService{}
}

// CalculateRiskScore calculates a risk score based on the provided gene data.
// The score is a simple mock and will return 1, 2, 3, or 4 based on the hash of the gene data.
func (s *TEEService) CalculateRiskScore(geneData string) uint {
	hash := sha256.Sum256([]byte(geneData))

	hashInt := new(big.Int).SetBytes(hash[:])

	riskScore := new(big.Int).Mod(hashInt, big.NewInt(4)).Uint64() + 1

	return uint(riskScore)
}

// EncryptGeneData encrypts the gene data using the user's public key.
func (s *TEEService) EncryptGeneData(publicKeyBytes []byte, geneData string) ([]byte, error) {
	publicKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	sharedSecret := generateSharedSecret(publicKey)

	// Encrypt the gene data using AES-256-GCM with the shared secret
	encryptedData, err := encryptAES(sharedSecret, []byte(geneData))
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

// generateSharedSecret generates a shared secret using the user's public key.
func generateSharedSecret(publicKey *ecdsa.PublicKey) []byte {
	// Hash the public key to generate the shared secret
	xBytes := publicKey.X.Bytes()
	sharedSecret := sha256.Sum256(xBytes)
	return sharedSecret[:]
}

// encryptAES encrypts data using AES-256-GCM.
func encryptAES(key, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

// DecryptGeneData decrypts the encrypted gene data using the user's private key.
func (s *TEEService) DecryptGeneData(privateKey *ecdsa.PrivateKey, encryptedData []byte) (string, error) {
	// Extract the nonce and ciphertext
	nonceSize := 12 // GCM standard nonce size
	if len(encryptedData) < nonceSize {
		return "", errors.New("invalid encrypted data")
	}
	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// Generate the shared secret
	xBytes := privateKey.PublicKey.X.Bytes()
	sharedSecret := sha256.Sum256(xBytes)

	// Decrypt using AES-256-GCM
	block, err := aes.NewCipher(sharedSecret[:])
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// SignEncryptedGeneData generates a hash of the encrypted data and then signs it using the provided private key.
// It returns the hash and the corresponding digital signature.
func (s *TEEService) SignEncryptedGeneData(privateKey *ecdsa.PrivateKey, encryptedData []byte) ([]byte, []byte, error) {
	// Generate the hash of the encrypted data
	hashData := crypto.Keccak256Hash(encryptedData).Bytes()

	// Sign the hash using the provided private key
	signature, err := crypto.Sign(hashData, privateKey)
	if err != nil {
		return nil, nil, err
	}

	return hashData, signature, nil
}
