package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/trungnt1811/blockchain-engineer-interview/backend/contracts"
)

const (
	controllerAddress  = "0x8A8937171197A78f47d8C2eE9A3C92FD33644B63"
	uploadDataEventABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"docId","type":"string"},{"indexed":false,"name":"sessionId","type":"uint256"}],"name":"UploadData","type":"event"}]`
)

type ControllerService struct {
	client     *ethclient.Client
	auth       *bind.TransactOpts
	controller *contracts.Controller
	eventABI   abi.ABI
}

// NewControllerService initializes a new ControllerService with the given client, authentication options, and contract address.
func NewControllerService(client *ethclient.Client, auth *bind.TransactOpts, address common.Address) (*ControllerService, error) {
	controller, err := contracts.NewController(address, client)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate Controller contract: %w", err)
	}

	// Parse the ABI for the UploadData event
	eventABI, err := abi.JSON(strings.NewReader(uploadDataEventABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	return &ControllerService{
		client:     client,
		auth:       auth,
		controller: controller,
		eventABI:   eventABI,
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

// ListenForUploadDataEvents listens for the UploadData event and returns the sessionId if the fileID matches.
func (s *ControllerService) ListenForUploadDataEvents(fileID string) (*big.Int, error) {
	// Create a filter query for the contract address and UploadData event
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(controllerAddress)},
		Topics:    [][]common.Hash{{s.eventABI.Events["UploadData"].ID}},
	}

	// Create a channel to receive logs
	logs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to logs: %w", err)
	}
	defer sub.Unsubscribe()

	// Loop through the logs
	for {
		select {
		case err := <-sub.Err():
			return nil, fmt.Errorf("subscription error: %w", err)
		case vLog := <-logs:
			// Unpack the log data into the event struct
			event := struct {
				DocID     string
				SessionID *big.Int
			}{}

			err := s.eventABI.UnpackIntoInterface(&event, "UploadData", vLog.Data)
			if err != nil {
				return nil, fmt.Errorf("failed to unpack event data: %w", err)
			}

			// Compare the docId with the provided fileID
			if event.DocID == fileID {
				fmt.Printf("Matching UploadData event found:\nDocument ID: %s\nSession ID: %s\n", event.DocID, event.SessionID.String())
				return event.SessionID, nil
			}
		}
	}
}
