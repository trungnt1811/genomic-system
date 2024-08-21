package blockchain

type Transaction struct {
	SessionID   uint64
	DocID       string
	User        string
	ContentHash string
	Proof       string
	TokenID     uint64
	RiskScore   uint64
}

type BlockchainMockService struct {
	transactions []Transaction
	userNFTs     map[string]uint64 // Maps user address to their NFT token ID
	userPCSP     map[string]uint64 // Maps user address to their PCSP token balance
}

func NewBlockchainMockService() *BlockchainMockService {
	return &BlockchainMockService{
		transactions: []Transaction{},
		userNFTs:     make(map[string]uint64),
		userPCSP:     make(map[string]uint64),
	}
}

// Mock function to simulate a data upload transaction
func (s *BlockchainMockService) UploadData(docId string, user string) uint64 {
	sessionID := uint64(len(s.transactions) + 1)
	s.transactions = append(s.transactions, Transaction{
		SessionID: sessionID,
		DocID:     docId,
		User:      user,
	})
	return sessionID
}

// Mock function to confirm upload data
func (s *BlockchainMockService) Confirm(sessionID uint64, contentHash string, proof string, riskScore uint64) {
	for i, tx := range s.transactions {
		if tx.SessionID == sessionID {
			s.transactions[i].ContentHash = contentHash
			s.transactions[i].Proof = proof
			s.transactions[i].RiskScore = riskScore

			// Minting a mock token if the user doesn't already have an NFT
			if _, exists := s.userNFTs[tx.User]; !exists {
				tokenID := s.mintNFT(tx.User)
				s.transactions[i].TokenID = tokenID
				s.userNFTs[tx.User] = tokenID
			} else {
				s.transactions[i].TokenID = s.userNFTs[tx.User]
			}

			// Reward tokens based on the risk score
			s.rewardPCSP(tx.User, riskScore)
		}
	}
}

func (s *BlockchainMockService) mintNFT(user string) uint64 {
	// Mock minting of a new NFT token
	tokenID := uint64(len(s.userNFTs) + 1)
	// Store the minted NFT for the user
	s.userNFTs[user] = tokenID
	return tokenID
}

func (s *BlockchainMockService) rewardPCSP(user string, riskScore uint64) {
	// Mock rewarding PCSP tokens based on the risk score
	// Increase the user's PCSP balance based on the risk score
	reward := s.calculatePCSPReward(riskScore)
	s.userPCSP[user] += reward
}

func (s *BlockchainMockService) calculatePCSPReward(riskScore uint64) uint64 {
	// Define the rewards based on risk score
	switch riskScore {
	case 1:
		return 15000
	case 2:
		return 3000
	case 3:
		return 225
	case 4:
		return 30
	default:
		return 0
	}
}

func (s *BlockchainMockService) GetUserNFT(user string) uint64 {
	return s.userNFTs[user]
}

func (s *BlockchainMockService) GetUserPCSPBalance(user string) uint64 {
	return s.userPCSP[user]
}
