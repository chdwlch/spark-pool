package pool

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/chdwlch/spark-pool/internal/channel"
	"github.com/chdwlch/spark-pool/pkg/types"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// Manager handles mining pool operations
type Manager struct {
	pool           *types.MiningPool
	channelManager *channel.Manager
	mu             sync.RWMutex
	blockHeight    uint64
	lastBlockTime  time.Time
	blockInterval  time.Duration
}

// NewManager creates a new mining pool manager
func NewManager(poolName string, operatorAddress string, serverPubKey *secp256k1.PublicKey) *Manager {
	pool := &types.MiningPool{
		ID:              generateID(),
		Name:            poolName,
		OperatorAddress: operatorAddress,
		TotalHashRate:   0,
		BlockReward:     625000000, // 6.25 BTC in satoshis
		Miners:          make(map[string]*types.Miner),
		ActiveChannels:  make(map[string]*types.Channel),
		CreatedAt:       time.Now(),
	}

	return &Manager{
		pool:           pool,
		channelManager: channel.NewManager(serverPubKey),
		blockHeight:    100000,
		lastBlockTime:  time.Now(),
		blockInterval:  10 * time.Minute, // 10 minutes per block
	}
}

// AddMiner adds a new miner to the pool
func (pm *Manager) AddMiner(ctx context.Context, minerName, minerAddress string, hashRate float64) (*types.Miner, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Generate miner key (in real implementation, miner would provide this)
	minerPrivKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate miner key: %w", err)
	}

	// Get pool operator key (simulated)
	poolOperatorKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate pool operator key: %w", err)
	}

	// Create miner
	miner := &types.Miner{
		ID:            generateID(),
		Address:       minerAddress,
		Name:          minerName,
		HashRate:      hashRate,
		TotalEarned:   0,
		CurrentBalance: 0,
		JoinedAt:      time.Now(),
		LastActivity:  time.Now(),
		IsActive:      true,
	}

	// Create Virtual Channel for the miner
	initialFunding := uint64(1000000) // 0.01 BTC initial funding
	channel, err := pm.channelManager.CreateMiningPoolChannel(
		poolOperatorKey.PubKey(),
		minerPrivKey.PubKey(),
		initialFunding,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	channel.MinerID = miner.ID
	channel.MinerAddress = minerAddress

	// Add to pool
	pm.pool.Miners[miner.ID] = miner
	pm.pool.ActiveChannels[channel.ID] = channel
	pm.pool.TotalHashRate += hashRate

	miner.ChannelID = channel.ID

	return miner, nil
}

// ProcessBlockReward processes a block reward and distributes it to miners
func (pm *Manager) ProcessBlockReward(ctx context.Context) (*types.BlockReward, error) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.blockHeight++
	pm.lastBlockTime = time.Now()

	// Calculate total hash rate
	totalHashRate := pm.pool.TotalHashRate
	if totalHashRate == 0 {
		return nil, fmt.Errorf("no active miners")
	}

	// Create block reward
	blockReward := &types.BlockReward{
		ID:           generateID(),
		BlockHeight:  pm.blockHeight,
		TotalReward:  pm.pool.BlockReward,
		Distributions: make(map[string]uint64),
		CreatedAt:    time.Now(),
	}

	// Distribute rewards to each miner based on hash rate
	for minerID, miner := range pm.pool.Miners {
		if !miner.IsActive {
			continue
		}

		// Calculate miner's share
		share := float64(pm.pool.BlockReward) * (miner.HashRate / totalHashRate)
		minerReward := uint64(math.Floor(share))

		if minerReward > 0 {
			blockReward.Distributions[minerID] = minerReward
			
			// Update miner stats
			miner.TotalEarned += minerReward
			miner.LastActivity = time.Now()

			// Process payment through Virtual Channel
			err := pm.processMinerPayment(ctx, miner, minerReward)
			if err != nil {
				return nil, fmt.Errorf("failed to process payment for miner %s: %w", minerID, err)
			}
		}
	}

	return blockReward, nil
}

// processMinerPayment processes a payment to a miner through their Virtual Channel
func (pm *Manager) processMinerPayment(ctx context.Context, miner *types.Miner, amount uint64) error {
	channel, exists := pm.pool.ActiveChannels[miner.ChannelID]
	if !exists {
		return fmt.Errorf("channel not found for miner %s", miner.ID)
	}

	// Create payment update
	_, err := pm.channelManager.CreatePaymentUpdate(
		channel,
		amount,
		"pool_operator",
	)
	if err != nil {
		return fmt.Errorf("failed to create payment update: %w", err)
	}

	// Update miner's current balance
	miner.CurrentBalance += amount

	// In a real implementation, you would:
	// 1. Sign the PSBT as pool operator
	// 2. Send to miner for storage
	// 3. Miner would store the latest payment update

	return nil
}

// GetPoolStats returns pool statistics
func (pm *Manager) GetPoolStats() *types.MiningStats {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	activeMiners := 0
	for _, miner := range pm.pool.Miners {
		if miner.IsActive {
			activeMiners++
		}
	}

	return &types.MiningStats{
		TotalMiners:     len(pm.pool.Miners),
		ActiveMiners:    activeMiners,
		TotalHashRate:   pm.pool.TotalHashRate,
		TotalEarned:     pm.calculateTotalEarned(),
		ActiveChannels:  len(pm.pool.ActiveChannels),
		LastBlockReward: pm.pool.BlockReward,
	}
}

// GetMiner returns a miner by ID
func (pm *Manager) GetMiner(minerID string) (*types.Miner, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	miner, exists := pm.pool.Miners[minerID]
	return miner, exists
}

// GetChannel returns a channel by ID
func (pm *Manager) GetChannel(channelID string) (*types.Channel, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	channel, exists := pm.pool.ActiveChannels[channelID]
	return channel, exists
}

// GetAllMiners returns all miners
func (pm *Manager) GetAllMiners() map[string]*types.Miner {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*types.Miner)
	for id, miner := range pm.pool.Miners {
		result[id] = miner
	}
	return result
}

// GetAllChannels returns all channels
func (pm *Manager) GetAllChannels() map[string]*types.Channel {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	result := make(map[string]*types.Channel)
	for id, channel := range pm.pool.ActiveChannels {
		result[id] = channel
	}
	return result
}

// CloseMinerChannel closes a miner's channel
func (pm *Manager) CloseMinerChannel(minerID string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	miner, exists := pm.pool.Miners[minerID]
	if !exists {
		return fmt.Errorf("miner not found")
	}

	channel, exists := pm.pool.ActiveChannels[miner.ChannelID]
	if !exists {
		return fmt.Errorf("channel not found")
	}

	// Close channel
	err := pm.channelManager.CloseChannel(channel)
	if err != nil {
		return fmt.Errorf("failed to close channel: %w", err)
	}

	// Remove from active channels
	delete(pm.pool.ActiveChannels, miner.ChannelID)
	pm.pool.TotalHashRate -= miner.HashRate
	miner.IsActive = false

	return nil
}

// calculateTotalEarned calculates total earned by all miners
func (pm *Manager) calculateTotalEarned() uint64 {
	total := uint64(0)
	for _, miner := range pm.pool.Miners {
		total += miner.TotalEarned
	}
	return total
}

// generateID generates a random ID
func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
} 