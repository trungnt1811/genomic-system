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
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	uploadDataEventABI = `[{
		"anonymous":false,
		"inputs":[
			{"indexed":false,"internalType": "string","name":"docId","type":"string"},
			{"indexed":false,"internalType": "uint256","name":"sessionId","type":"uint256"}
		],
		"name":"UploadData",
		"type":"event"
	}]`
)

type ControllerEventListener struct {
	client   *ethclient.Client
	auth     *bind.TransactOpts
	eventABI abi.ABI
}

// NewControllerEventListener initializes the ControllerEventListener by loading the ABI from the specified file.
func NewControllerEventListener(client *ethclient.Client, auth *bind.TransactOpts) (*ControllerEventListener, error) {
	/// Parse the ABI for the UploadData event
	eventABI, err := abi.JSON(strings.NewReader(uploadDataEventABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
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
		FromBlock: fromBlock,
	}

	// Poll for logs
	logs, err := s.client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	// Loop through the logs
	for _, vLog := range logs {
		// Unpack the non-indexed parameters
		events, err := s.eventABI.Unpack("UploadData", vLog.Data)
		if err != nil {
			fmt.Printf("Error unpacking event data: %v\n", err)
			continue
		}
		if len(events) == 2 {
			docID, ok1 := events[0].(string)
			sessionID, ok2 := events[1].(*big.Int)

			if ok1 && ok2 {
				fmt.Printf("Document ID: %s, Session ID: %s\n", docID, sessionID.String())
			} else {
				fmt.Println("Failed to cast event data to expected types")
			}

			// Compare the docId with the provided fileID
			if docID == fileID {
				fmt.Printf("Matching UploadData event found:\nDocument ID: %s\nSession ID: %s\n", docID, sessionID.String())
				return sessionID, nil
			}
		}
	}

	return nil, fmt.Errorf("no matching event found")
}
