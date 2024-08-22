package blockchain

import (
	"context"
	"fmt"
	"log"
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
func NewControllerService(client *ethclient.Client, auth *bind.TransactOpts, address common.Address) *ControllerService {
	controller, err := contracts.NewController(address, client)
	if err != nil {
		log.Fatalf("Failed to instantiate Controller contract: %v", err)
	}

	// Parse the ABI for the UploadData event
	eventABI, err := abi.JSON(strings.NewReader(uploadDataEventABI))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	return &ControllerService{
		client:     client,
		auth:       auth,
		controller: controller,
		eventABI:   eventABI,
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

// ListenForUploadDataEvents listens for UploadData events and processes them.
func (s *ControllerService) ListenForUploadDataEvents() {
	// Create a filter query
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(controllerAddress)},
		Topics:    [][]common.Hash{{s.eventABI.Events["UploadData"].ID}},
	}

	// Create a channel to receive logs
	logs := make(chan types.Log)
	sub, err := s.client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatalf("Failed to subscribe to logs: %v", err)
	}

	log.Println("Listening for UploadData events...")

	// Process logs
	for {
		select {
		case err := <-sub.Err():
			log.Fatalf("Subscription error: %v", err)
		case vLog := <-logs:
			// Parse the log data
			var docId string
			var sessionId *big.Int

			// Unpack the event data
			err := s.eventABI.UnpackIntoInterface(&docId, "docId", vLog.Data)
			if err != nil {
				log.Printf("Failed to unpack docId: %v", err)
				continue
			}
			err = s.eventABI.UnpackIntoInterface(&sessionId, "sessionId", vLog.Data)
			if err != nil {
				log.Printf("Failed to unpack sessionId: %v", err)
				continue
			}

			// Output the results
			log.Printf("Received UploadData event:\n")
			log.Printf("Document ID: %s\n", docId)
			log.Printf("Session ID: %s\n", sessionId.String())
		}
	}
}
