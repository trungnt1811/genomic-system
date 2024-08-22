package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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

	// Initialize Optimism rpc client
	client, err := ethclient.Dial("https://sepolia.optimism.io")
	if err != nil {
		fmt.Println("Error connecting to Ethereum rpc client:", err)
		return
	}

	// Get the chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("Error getting chain ID:", err)
		return
	}

	// Create a new transactor with chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(ecdsaPrivateKey, chainID)
	if err != nil {
		fmt.Println("Error creating keyed transactor:", err)
		return
	}

	// Initialize Controller and PCSP services
	controllerAddress := common.HexToAddress("0x8A8937171197A78f47d8C2eE9A3C92FD33644B63")
	pcspAddress := common.HexToAddress("0x7bc91a89bb437fBB199fB3D1d0dc3a9D913d4f9F")
	controllerService, err := blockchain.NewControllerService(client, auth, controllerAddress)
	if err != nil {
		fmt.Println("Error initializing Controller service:", err)
		return
	}
	pcspService, err := blockchain.NewPCSPService(client, auth, pcspAddress)
	if err != nil {
		fmt.Println("Error initializing PCSP service:", err)
		return
	}

	// Initialize Controller event listener
	controllerEventListener, err := blockchain.NewControllerEventListener(client, auth)
	if err != nil {
		fmt.Println("Error initializing controller event listener:", err)
		return
	}

	// Step 1: Register a new user with public key
	fmt.Println("\nStep 1")
	fmt.Println("Registering a new user...")
	userID := authService.RegisterUserWithPubkey(userPubkeyBytes)
	fmt.Printf("User registered with UserID: %d and Ethereum address: %s\n", userID, userETHAddress)

	// Step 2: Authenticate the user using their Ethereum address derived from the public key
	fmt.Println("\nStep 2")
	isAuthenticated := authService.Authenticate(userID, userETHAddress)
	if !isAuthenticated {
		fmt.Println("User authentication failed!")
		return
	}
	fmt.Println("User authenticated successfully with Ethereum address:", userETHAddress)

	// Step 3: Encrypt gene data using the user's public key via the TEE service
	fmt.Println("\nStep 3")
	geneData, err := randomStringWithRandomLength(10, 50) // Example gene data to be encrypted
	if err != nil {
		fmt.Println("Error creating random gene data:", err)
		return
	}
	fmt.Printf("Original gene data: %s\n", geneData)
	fmt.Println("Encrypting gene data...")
	encryptedData, err := teeService.EncryptGeneData(userPubkeyBytes, geneData)
	if err != nil {
		fmt.Println("Error encrypting gene data:", err)
		return
	}
	fmt.Println("Gene data encrypted successfully.")

	// Step 4: Sign the encrypted gene data using the user's private key via the TEE service
	fmt.Println("\nStep 4")
	fmt.Println("Signing encrypted gene data...")
	hash, signature, err := teeService.SignEncryptedGeneData(ecdsaPrivateKey, encryptedData)
	if err != nil {
		fmt.Println("Error signing gene data:", err)
		return
	}
	fmt.Println("Gene data signed successfully.")

	// Step 5: Store the encrypted gene data, signature, and hash in the storage service
	fmt.Println("\nStep 5")
	fmt.Println("Storing encrypted gene data...")
	fileID, err := geneDataStorageService.StoreGeneData(userID, encryptedData, signature, hash)
	if err != nil {
		fmt.Println("Error storing gene data:", err)
		return
	}
	fmt.Printf("Gene data stored successfully with FileID: %s\n", fileID)

	// Step 6: Verify the gene data signature to ensure its integrity
	fmt.Println("\nStep 6")
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
	fmt.Println("\nStep 7")
	fmt.Println("Calculating risk score...")
	riskScore := teeService.CalculateRiskScore(geneData)
	fmt.Printf("Risk score calculated: %d\n", riskScore)

	// Step 8: Upload the gene data to the blockchain for secure storage
	fmt.Println("\nStep 8")
	fmt.Println("Uploading gene data to blockchain...")
	txHash, err := controllerService.UploadData(fileID)
	if err != nil {
		fmt.Println("Error uploading data to blockchain:", err)
		return
	}
	fmt.Printf("Gene data uploaded at txHash: %s\n", txHash.Hex())

	// Step 9: Listen for the UploadData event to get the sessionID
	fmt.Println("\nStep 9")
	fmt.Println("Listening for UploadData event to get sessionID...")
	sessionID, err := controllerEventListener.ListenForUploadDataEvents(fileID)
	if err != nil {
		fmt.Println("Error listening for UploadData event:", err)
		return
	}
	fmt.Printf("Received sessionID: %s\n", sessionID)

	// Step 9.1: Confirm the blockchain transaction, mint an NFT, and reward PCSP tokens
	fmt.Println("Confirming transaction on blockchain...")
	err = controllerService.Confirm(fileID, fmt.Sprintf("%x", hash), fmt.Sprintf("%x", signature), sessionID, uint8(riskScore))
	if err != nil {
		fmt.Println("Error confirming transaction on blockchain:", err)
		return
	}
	fmt.Println("Transaction confirmed, NFT minted, and PCSP tokens rewarded.")

	// Step 9.2: Retrieve the user's PCSP balance from the blockchain
	userPCSPBalance, err := pcspService.GetBalance(common.HexToAddress(userETHAddress))
	if err != nil {
		fmt.Println("Error retrieving PCSP balance:", err)
		return
	}
	fmt.Printf("User's PCSP Balance: %d\n", userPCSPBalance)

	// Step 10: Retrieve and decrypt the original gene data using the user's private key
	fmt.Println("\nStep 10")
	fmt.Println("Retrieving and decrypting original gene data...")

	// Step 10.1: Verify the gene data signature again before decryption
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

	// Step 10.2: Retrieve the encrypted gene data from storage using the fileID
	retrievedEncryptedData, err := geneDataStorageService.RetrieveGeneData(fileID)
	if err != nil {
		fmt.Println("Error retrieving gene data:", err)
		return
	}
	fmt.Println("Encrypted gene data retrieved successfully.")

	// Step 10.3: Decrypt the gene data using the user's private key
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
func hexToECDSAPrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.HexToECDSA(hexKey)
	if err != nil {
		return nil, err
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

func randomStringWithRandomLength(minLength, maxLength int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	// Generate a random length within the given range
	length, err := rand.Int(rand.Reader, big.NewInt(int64(maxLength-minLength+1)))
	if err != nil {
		return "", err
	}
	length = length.Add(length, big.NewInt(int64(minLength)))

	// Build the random string
	var sb strings.Builder
	for i := 0; i < int(length.Int64()); i++ {
		charIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		sb.WriteByte(charset[charIndex.Int64()])
	}

	return sb.String(), nil
}
