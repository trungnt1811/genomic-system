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
)

const (
	uploadDataEventABI = `[{"anonymous":false,"inputs":[{"indexed":false,"name":"docId","type":"string"},{"indexed":false,"name":"sessionId","type":"uint256"}],"name":"UploadData","type":"event"}]`
)

type ControllerEventListener struct {
	client   *ethclient.Client
	auth     *bind.TransactOpts
	eventABI abi.ABI
}

func NewControllerEventListener(client *ethclient.Client, auth *bind.TransactOpts) (*ControllerEventListener, error) {
	// Parse the ABI for the UploadData event
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
