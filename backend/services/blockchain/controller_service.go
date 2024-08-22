package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/contracts"
)

const (
	controllerAddress = "0x8A8937171197A78f47d8C2eE9A3C92FD33644B63"
)

type ControllerService struct {
	client     *ethclient.Client
	auth       *bind.TransactOpts
	controller *contracts.Controller
}

// NewControllerService initializes a new ControllerService with the given client, authentication options, and contract address.
func NewControllerService(client *ethclient.Client, auth *bind.TransactOpts, address common.Address) (*ControllerService, error) {
	controller, err := contracts.NewController(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Controller contract: %w", err)
	}

	return &ControllerService{
		client:     client,
		auth:       auth,
		controller: controller,
	}, nil
}

// UploadData uploads data to the contract and returns the transaction hash for tracking.
func (s *ControllerService) UploadData(docId string) (common.Hash, error) {
	// Send the transaction to upload data
	tx, err := s.controller.UploadData(s.auth, docId)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to upload data: %w", err)
	}

	// Wait for the transaction receipt to ensure it's processed
	_, err = bind.WaitMined(context.Background(), s.client, tx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to mine transaction: %w", err)
	}

	// Log and return the transaction hash
	return tx.Hash(), nil
}

// Confirm confirms the uploaded data and logs the transaction hash.
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
	return nil
}
