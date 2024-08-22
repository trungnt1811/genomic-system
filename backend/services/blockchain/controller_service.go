package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/contracts"
)

type ControllerService struct {
	client     *ethclient.Client
	auth       *bind.TransactOpts
	controller *contracts.Controller
}

// NewControllerService initializes a new ControllerService with the given client, authentication options, and contract address.
func NewControllerService(client *ethclient.Client, auth *bind.TransactOpts, address common.Address) *ControllerService {
	controller, err := contracts.NewController(address, client)
	if err != nil {
		log.Fatalf("Failed to instantiate Controller contract: %v", err)
	}

	return &ControllerService{
		client:     client,
		auth:       auth,
		controller: controller,
	}
}

// UploadData uploads data to the contract and returns the transaction hash for tracking.
func (s *ControllerService) UploadData(docId string) common.Hash {
	// Send the transaction to upload data
	tx, err := s.controller.UploadData(s.auth, docId)
	if err != nil {
		log.Fatalf("Failed to upload data: %v", err)
	}

	// Wait for the transaction receipt to ensure it's processed
	_, err = bind.WaitMined(context.Background(), s.client, tx)
	if err != nil {
		log.Fatalf("Failed to mine transaction: %v", err)
	}

	// Log and return the transaction hash
	log.Printf("Upload data transaction sent: %s\n", tx.Hash().Hex())
	return tx.Hash()
}

// Confirm confirms the uploaded data and logs the transaction hash.
// It returns an error if the transaction fails or if mining the transaction fails.
func (s *ControllerService) Confirm(docId, contentHash, proof string, sessionId *big.Int, riskScore uint8) error {
	// Send the transaction to confirm data
	tx, err := s.controller.Confirm(s.auth, docId, contentHash, proof, sessionId, big.NewInt(int64(riskScore)))
	if err != nil {
		return fmt.Errorf("failed to confirm session: %w", err)
	}

	// Wait for the transaction receipt to ensure it's processed
	_, err = bind.WaitMined(context.Background(), s.client, tx)
	if err != nil {
		return fmt.Errorf("failed to mine transaction: %w", err)
	}

	// Log the transaction hash
	log.Printf("Confirm transaction sent: %s\n", tx.Hash().Hex())
	return nil
}
