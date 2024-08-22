package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ControllerEventListener struct {
	client   *ethclient.Client
	auth     *bind.TransactOpts
	eventABI abi.ABI
}

// NewControllerEventListener initializes the ControllerEventListener by loading the ABI from the specified file.
func NewControllerEventListener(client *ethclient.Client, auth *bind.TransactOpts) (*ControllerEventListener, error) {
	// Read the ABI file
	abiBytes, err := os.ReadFile("./services/blockchain/Controller.abi.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %w", err)
	}

	// Parse the ABI
	var eventABI abi.ABI
	if err := json.Unmarshal(abiBytes, &eventABI); err != nil {
		return nil, fmt.Errorf("failed to parse ABI JSON: %w", err)
	}

	return &ControllerEventListener{
		client:   client,
		auth:     auth,
		eventABI: eventABI,
	}, nil
}

// ListenForUploadDataEvents listens for the UploadData event and returns the sessionId if the fileID matches.
func (s *ControllerEventListener) ListenForUploadDataEvents(fileID string) (*big.Int, error) {
	// Get the latest block number
	header, err := s.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest block: %w", err)
	}
	latestBlock := header.Number

	// Ensure the block range is valid
	fromBlock := new(big.Int).Sub(latestBlock, big.NewInt(10))
	if fromBlock.Sign() < 0 {
		fromBlock = big.NewInt(0)
	}

	// Create a filter query for the contract address and UploadData event
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(controllerAddress)},
		Topics:    [][]common.Hash{{s.eventABI.Events["UploadData"].ID}},
		FromBlock: fromBlock,
	}

	// Poll for logs
	logs, err := s.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	// Loop through the logs
	for _, vLog := range logs {
		// Unpack the log data into the event struct
		event := struct {
			DocID     string
			SessionID *big.Int
		}{}

		// Unpack the non-indexed parameters
		err := s.eventABI.UnpackIntoInterface(&event, "UploadData", vLog.Data)
		if err != nil {
			continue
		}

		// Compare the docId with the provided fileID
		if event.DocID == fileID {
			fmt.Printf("Matching UploadData event found:\nDocument ID: %s\nSession ID: %s\n", event.DocID, event.SessionID.String())
			return event.SessionID, nil
		}
	}

	return nil, fmt.Errorf("no matching event found")
}
