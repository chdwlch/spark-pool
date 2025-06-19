package types

import (
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

// MiningPool represents a Bitcoin mining pool
type MiningPool struct {
	ID              string                    `json:"id"`
	Name            string                    `json:"name"`
	OperatorAddress string                    `json:"operator_address"`
	TotalHashRate   float64                   `json:"total_hash_rate"`
	BlockReward     uint64                    `json:"block_reward"`
	Miners          map[string]*Miner         `json:"miners"`
	ActiveChannels  map[string]*Channel       `json:"active_channels"`
	CreatedAt       time.Time                 `json:"created_at"`
}

// Miner represents a miner in the pool
type Miner struct {
	ID             string    `json:"id"`
	Address        string    `json:"address"`
	Name           string    `json:"name"`
	HashRate       float64   `json:"hash_rate"`
	TotalEarned    uint64    `json:"total_earned"`
	CurrentBalance uint64    `json:"current_balance"`
	JoinedAt       time.Time `json:"joined_at"`
	LastActivity   time.Time `json:"last_activity"`
	IsActive       bool      `json:"is_active"`
	ChannelID      string    `json:"channel_id"`
}

// Channel represents a Virtual Channel between pool operator and miner
type Channel struct {
	ID              string                    `json:"id"`
	PoolOperatorKey *secp256k1.PublicKey     `json:"pool_operator_key"`
	MinerKey        *secp256k1.PublicKey     `json:"miner_key"`
	InitialFunding  uint64                   `json:"initial_funding"`
	CurrentBalance  uint64                   `json:"current_balance"`
	Status          string                   `json:"status"`
	CreatedAt       time.Time                `json:"created_at"`
	LastUpdated     time.Time                `json:"last_updated"`
	PaymentHistory  []*PaymentUpdate         `json:"payment_history"`
	MinerID         string                   `json:"miner_id"`
	MinerAddress    string                   `json:"miner_address"`
}

// PaymentUpdate represents a payment update in a Virtual Channel
type PaymentUpdate struct {
	ID          string    `json:"id"`
	ChannelID   string    `json:"channel_id"`
	Amount      uint64    `json:"amount"`
	FromParty   string    `json:"from_party"`
	ToParty     string    `json:"to_party"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	SequenceNum uint64    `json:"sequence_num"`
}

// BlockReward represents a block reward distribution
type BlockReward struct {
	ID            string            `json:"id"`
	BlockHeight   uint64            `json:"block_height"`
	TotalReward   uint64            `json:"total_reward"`
	Distributions map[string]uint64 `json:"distributions"`
	CreatedAt     time.Time         `json:"created_at"`
}

// MiningStats represents pool statistics
type MiningStats struct {
	TotalMiners     int     `json:"total_miners"`
	ActiveMiners    int     `json:"active_miners"`
	TotalHashRate   float64 `json:"total_hash_rate"`
	TotalEarned     uint64  `json:"total_earned"`
	ActiveChannels  int     `json:"active_channels"`
	LastBlockReward uint64  `json:"last_block_reward"`
}

// MinerStats represents individual miner statistics
type MinerStats struct {
	MinerID         string  `json:"miner_id"`
	Name            string  `json:"name"`
	HashRate        float64 `json:"hash_rate"`
	TotalEarned     uint64  `json:"total_earned"`
	CurrentBalance  uint64  `json:"current_balance"`
	ChannelBalance  uint64  `json:"channel_balance"`
	IsActive        bool    `json:"is_active"`
	LastActivity    string  `json:"last_activity"`
}

// ChannelStats represents channel statistics
type ChannelStats struct {
	ChannelID       string `json:"channel_id"`
	MinerID         string `json:"miner_id"`
	MinerName       string `json:"miner_name"`
	InitialFunding  uint64 `json:"initial_funding"`
	CurrentBalance  uint64 `json:"current_balance"`
	Status          string `json:"status"`
	PaymentCount    int    `json:"payment_count"`
	CreatedAt       string `json:"created_at"`
	LastUpdated     string `json:"last_updated"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// AddMinerRequest represents a request to add a miner
type AddMinerRequest struct {
	Name     string  `json:"name"`
	Address  string  `json:"address"`
	HashRate float64 `json:"hash_rate"`
}

// JoinPoolRequest represents a request to join the mining pool
type JoinPoolRequest struct {
	MinerName string  `json:"miner_name"`
	Address   string  `json:"address"`
	HashRate  float64 `json:"hash_rate"`
}

// ProcessBlockRequest represents a request to process a block reward
type ProcessBlockRequest struct {
	BlockHeight uint64 `json:"block_height,omitempty"`
	Reward      uint64 `json:"reward,omitempty"`
}

// WebSocketMessage represents a WebSocket message
type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
} 