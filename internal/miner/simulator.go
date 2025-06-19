package miner

import (
	"context"
	cryptorand "crypto/rand"
	"encoding/hex"
	"fmt"
	mathrand "math/rand"
	"sync"
	"time"

	"github.com/chdwlch/spark-pool/pkg/types"
)

// Simulator simulates a Bitcoin miner
type Simulator struct {
	ID           string
	Name         string
	Address      string
	HashRate     float64
	IsMining     bool
	mu           sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	stats        *MiningStats
}

// MiningStats tracks mining statistics
type MiningStats struct {
	TotalShares    uint64
	AcceptedShares uint64
	RejectedShares uint64
	LastShareTime  time.Time
	Uptime         time.Duration
	StartTime      time.Time
}

// NewSimulator creates a new miner simulator
func NewSimulator(name, address string, hashRate float64) *Simulator {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Simulator{
		ID:       generateID(),
		Name:     name,
		Address:  address,
		HashRate: hashRate,
		IsMining: false,
		ctx:      ctx,
		cancel:   cancel,
		stats: &MiningStats{
			StartTime: time.Now(),
		},
	}
}

// StartMining starts the mining simulation
func (ms *Simulator) StartMining() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if ms.IsMining {
		return fmt.Errorf("miner is already mining")
	}

	ms.IsMining = true
	ms.stats.StartTime = time.Now()

	// Start mining goroutine
	go ms.miningLoop()

	return nil
}

// StopMining stops the mining simulation
func (ms *Simulator) StopMining() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if !ms.IsMining {
		return fmt.Errorf("miner is not mining")
	}

	ms.IsMining = false
	ms.cancel()

	return nil
}

// miningLoop simulates the mining process
func (ms *Simulator) miningLoop() {
	ticker := time.NewTicker(1 * time.Second) // Simulate shares every second
	defer ticker.Stop()

	for {
		select {
		case <-ms.ctx.Done():
			return
		case <-ticker.C:
			ms.submitShare()
		}
	}
}

// submitShare simulates submitting a mining share
func (ms *Simulator) submitShare() {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	if !ms.IsMining {
		return
	}

	// Simulate share submission
	ms.stats.TotalShares++
	ms.stats.LastShareTime = time.Now()

	// Simulate share acceptance/rejection (90% acceptance rate)
	if mathrand.Float64() < 0.9 {
		ms.stats.AcceptedShares++
	} else {
		ms.stats.RejectedShares++
	}

	// Update uptime
	ms.stats.Uptime = time.Since(ms.stats.StartTime)
}

// GetStats returns current mining statistics
func (ms *Simulator) GetStats() *MiningStats {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	stats := *ms.stats
	stats.Uptime = time.Since(ms.stats.StartTime)
	return &stats
}

// GetHashRate returns current hash rate
func (ms *Simulator) GetHashRate() float64 {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.HashRate
}

// SetHashRate updates the hash rate
func (ms *Simulator) SetHashRate(hashRate float64) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.HashRate = hashRate
}

// IsActive returns whether the miner is currently mining
func (ms *Simulator) IsActive() bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.IsMining
}

// GetEfficiency returns mining efficiency (accepted shares / total shares)
func (ms *Simulator) GetEfficiency() float64 {
	stats := ms.GetStats()
	if stats.TotalShares == 0 {
		return 0
	}
	return float64(stats.AcceptedShares) / float64(stats.TotalShares) * 100
}

// GetSharesPerSecond returns shares per second based on current stats
func (ms *Simulator) GetSharesPerSecond() float64 {
	stats := ms.GetStats()
	if stats.Uptime.Seconds() == 0 {
		return 0
	}
	return float64(stats.TotalShares) / stats.Uptime.Seconds()
}

// generateID generates a random ID
func generateID() string {
	bytes := make([]byte, 16)
	cryptorand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Pool represents a mining pool interface
type Pool interface {
	AddMiner(ctx context.Context, minerName, minerAddress string, hashRate float64) (*types.Miner, error)
	ProcessBlockReward(ctx context.Context) (*types.BlockReward, error)
	GetPoolStats() *types.MiningStats
}

// Manager manages multiple miner simulators
type Manager struct {
	simulators map[string]*Simulator
	pool       Pool
	mu         sync.RWMutex
}

// NewManager creates a new miner manager
func NewManager(pool Pool) *Manager {
	return &Manager{
		simulators: make(map[string]*Simulator),
		pool:       pool,
	}
}

// AddSimulator adds a new miner simulator
func (mm *Manager) AddSimulator(name, address string, hashRate float64) (*Simulator, error) {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	simulator := NewSimulator(name, address, hashRate)
	mm.simulators[simulator.ID] = simulator

	// Add to pool
	_, err := mm.pool.AddMiner(context.Background(), name, address, hashRate)
	if err != nil {
		delete(mm.simulators, simulator.ID)
		return nil, fmt.Errorf("failed to add miner to pool: %w", err)
	}

	return simulator, nil
}

// GetSimulator returns a simulator by ID
func (mm *Manager) GetSimulator(id string) (*Simulator, bool) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	simulator, exists := mm.simulators[id]
	return simulator, exists
}

// GetAllSimulators returns all simulators
func (mm *Manager) GetAllSimulators() map[string]*Simulator {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	result := make(map[string]*Simulator)
	for id, simulator := range mm.simulators {
		result[id] = simulator
	}
	return result
}

// StartAllSimulators starts all simulators
func (mm *Manager) StartAllSimulators() error {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	for _, simulator := range mm.simulators {
		if err := simulator.StartMining(); err != nil {
			return fmt.Errorf("failed to start simulator %s: %w", simulator.ID, err)
		}
	}

	return nil
}

// StopAllSimulators stops all simulators
func (mm *Manager) StopAllSimulators() error {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	for _, simulator := range mm.simulators {
		if err := simulator.StopMining(); err != nil {
			return fmt.Errorf("failed to stop simulator %s: %w", simulator.ID, err)
		}
	}

	return nil
} 