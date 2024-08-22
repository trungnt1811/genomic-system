package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/services/auth"
	"github.com/trungnt1811/blockchain-engineer-interview/backend/services/blockchain"
	"github.com/trungnt1811/blockchain-engineer-interview/backend/services/storage"
	"github.com/trungnt1811/blockchain-engineer-interview/backend/services/tee"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Get user private key from .env
	userPrivateKeyHex := os.Getenv("PRIVATE_KEY")
	if userPrivateKeyHex == "" {
		fmt.Println("PRIVATE_KEY not found in .env file")
		return
	}

	// Convert the private key hex string to an ECDSA private key
	ecdsaPrivateKey, err := hexToECDSAPrivateKey(userPrivateKeyHex)
	if err != nil {
		fmt.Println("Error converting private key hex to ECDSA:", err)
		return
	}

	// Generate user pubkey from private key
	userPubkeyBytes := crypto.FromECDSAPub(&ecdsaPrivateKey.PublicKey)

	// Convert the public key bytes to an Ethereum address for identification
	userETHAddress, err := pubkeyToETHAddress(userPubkeyBytes)
	if err != nil {
		fmt.Println("Error converting public key to Ethereum address:", err)
		return
	}

	// Initialize necessary services for the application
	authService := auth.NewAuthService()
	geneDataStorageService := storage.NewGeneDataStorageService()
	teeService := tee.NewTEEService()
	blockchainService := blockchain.NewBlockchainMockService()

	// Step 1: Register a new user with public key
	fmt.Println("Registering a new user...")
	userID := authService.RegisterUserWithPubkey(userPubkeyBytes)
	fmt.Printf("User registered with UserID: %d and Ethereum address: %s\n", userID, userETHAddress)
	// Note: The private key is typically stored securely by the user; the service only retains the public key.

	// Step 2: Authenticate the user using their Ethereum address derived from the public key
	isAuthenticated := authService.Authenticate(userID, userETHAddress)
	if !isAuthenticated {
		fmt.Println("User authentication failed!")
		return
	}
	fmt.Println("User authenticated successfully with Ethereum address:", userETHAddress)

	// Step 3: Encrypt gene data using the user's public key via the TEE service
	geneData := "AGTCAGTCAGTC..." // Example gene data to be encrypted
	fmt.Printf("Original gene data: %s\n", geneData)
	fmt.Println("Encrypting gene data...")
	encryptedData, err := teeService.EncryptGeneData(userPubkeyBytes, geneData)
	if err != nil {
		fmt.Println("Error encrypting gene data:", err)
		return
	}
	fmt.Println("Gene data encrypted successfully.")

	// Step 4: Sign the encrypted gene data using the user's private key via the TEE service
	fmt.Println("Signing encrypted gene data...")
	hash, signature, err := teeService.SignEncryptedGeneData(ecdsaPrivateKey, encryptedData)
	if err != nil {
		fmt.Println("Error signing gene data:", err)
		return
	}
	fmt.Println("Gene data signed successfully.")

	// Step 5: Store the encrypted gene data, signature, and hash in the storage service
	fmt.Println("Storing encrypted gene data...")
	fileID, err := geneDataStorageService.StoreGeneData(userID, encryptedData, signature, hash)
	if err != nil {
		fmt.Println("Error storing gene data:", err)
		return
	}
	fmt.Printf("Gene data stored successfully with FileID: %s\n", fileID)

	// Step 6: Verify the gene data signature to ensure its integrity
	fmt.Println("Verifying gene data signature...")
	isSignatureValid, err := geneDataStorageService.VerifyGeneDataSignature(fileID, userPubkeyBytes)
	if err != nil {
		fmt.Println("Error verifying signature:", err)
		return
	}
	if isSignatureValid {
		fmt.Println("Gene data signature is valid.")
	} else {
		fmt.Println("Gene data signature is invalid.")
	}

	// Step 7: Calculate the risk score based on the gene data using the TEE service
	fmt.Println("Calculating risk score...")
	riskScore := teeService.CalculateRiskScore(geneData)
	fmt.Printf("Risk score calculated: %d\n", riskScore)

	// Step 8: Upload the gene data to the blockchain for secure storage
	fmt.Println("Uploading gene data to blockchain...")
	sessionID := blockchainService.UploadData(fileID, userETHAddress)
	fmt.Printf("Gene data uploaded with SessionID: %d\n", sessionID)

	// Step 9: Confirm the blockchain transaction, mint an NFT, and reward PCSP tokens
	fmt.Println("Confirming transaction on blockchain...")
	blockchainService.Confirm(sessionID, fmt.Sprintf("%x", hash), fmt.Sprintf("%x", signature), uint64(riskScore))
	fmt.Println("Transaction confirmed, NFT minted, and PCSP tokens rewarded.")

	// Step 10: Retrieve the user's NFT and PCSP balance from the blockchain
	userNFT := blockchainService.GetUserNFT(userETHAddress)
	userPCSPBalance := blockchainService.GetUserPCSPBalance(userETHAddress)
	fmt.Printf("User's NFT TokenID: %d, PCSP Balance: %d\n", userNFT, userPCSPBalance)

	// Step 11: Retrieve and decrypt the original gene data using the user's private key
	fmt.Println("Retrieving and decrypting original gene data...")

	// Step 12: Verify the gene data signature again before decryption
	fmt.Println("Verifying gene data signature...")
	isSignatureValid, err = geneDataStorageService.VerifyGeneDataSignature(fileID, userPubkeyBytes)
	if err != nil {
		fmt.Println("Error verifying signature:", err)
		return
	}
	if isSignatureValid {
		fmt.Println("Gene data signature is valid.")
	} else {
		fmt.Println("Gene data signature is invalid.")
		return
	}

	// Step 13: Retrieve the encrypted gene data from storage using the fileID
	retrievedEncryptedData, err := geneDataStorageService.RetrieveGeneData(fileID)
	if err != nil {
		fmt.Println("Error retrieving gene data:", err)
		return
	}
	fmt.Println("Encrypted gene data retrieved successfully.")

	// Step 14: Decrypt the gene data using the user's private key
	decryptedGeneData, err := decryptGeneData(ecdsaPrivateKey, retrievedEncryptedData)
	if err != nil {
		fmt.Println("Error decrypting gene data:", err)
		return
	}
	fmt.Printf("Original gene data retrieved and decrypted successfully: %s\n", decryptedGeneData)
}

// pubkeyToETHAddress converts a public key byte slice to an Ethereum address.
func pubkeyToETHAddress(publicKeyBytes []byte) (string, error) {
	// Convert the public key bytes back to an ECDSA public key
	publicKey, err := crypto.UnmarshalPubkey(publicKeyBytes)
	if err != nil {
		return "", err
	}

	// Derive the Ethereum address from the public key
	ethAddress := crypto.PubkeyToAddress(*publicKey).Hex()
	return ethAddress, nil
}

// hexToECDSAPrivateKey converts a hexadecimal string to an ECDSA private key.
func hexToECDSAPrivateKey(privateKeyHex string) (*ecdsa.PrivateKey, error) {
	// Decode the hex string to a byte slice
	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, errors.New("failed to decode hex string: " + err.Error())
	}

	// Convert the byte slice to an ECDSA private key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, errors.New("failed to convert to ECDSA private key: " + err.Error())
	}

	return privateKey, nil
}

// decryptGeneData decrypts the encrypted gene data using the user's private key.
func decryptGeneData(privateKey *ecdsa.PrivateKey, encryptedData []byte) (string, error) {
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
