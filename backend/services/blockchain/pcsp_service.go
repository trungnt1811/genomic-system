package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/trungnt1811/blockchain-engineer-interview/backend/contracts"
)

type PCSPService struct {
	client    *ethclient.Client
	auth      *bind.TransactOpts
	pcspToken *contracts.PCSP
}

// NewPCSPService initializes a new PostCovidStrokePreventionService with the given client, authentication options, and contract address.
func NewPCSPService(client *ethclient.Client, auth *bind.TransactOpts, address common.Address) (*PCSPService, error) {
	pcspToken, err := contracts.NewPCSP(address, client)
	if err != nil {
		return nil, err
	}

	return &PCSPService{
		client:    client,
		auth:      auth,
		pcspToken: pcspToken,
	}, nil
}

// GetBalance retrieves the balance of a specific address.
func (s *PCSPService) GetBalance(address common.Address) (*big.Int, error) {
	// Call the balanceOf function from the ERC20 contract
	balance, err := s.pcspToken.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return nil, err
	}
	return balance, nil
}
