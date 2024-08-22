package blockchain

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/contracts"
)

type PostCovidStrokePreventionService struct {
	client    *ethclient.Client
	auth      *bind.TransactOpts
	pcspToken *contracts.PCSP
}

// NewPostCovidStrokePreventionService initializes a new PostCovidStrokePreventionService with the given client, authentication options, and contract address.
func NewPostCovidStrokePreventionService(client *ethclient.Client, auth *bind.TransactOpts, address common.Address) *PostCovidStrokePreventionService {
	pcspToken, err := contracts.NewPCSP(address, client)
	if err != nil {
		log.Fatalf("Failed to instantiate PostCovidStrokePrevention contract: %v", err)
	}

	return &PostCovidStrokePreventionService{
		client:    client,
		auth:      auth,
		pcspToken: pcspToken,
	}
}

// GetBalance retrieves the balance of a specific address.
func (s *PostCovidStrokePreventionService) GetBalance(address common.Address) *big.Int {
	// Call the balanceOf function from the ERC20 contract
	balance, err := s.pcspToken.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}
	return balance
}
