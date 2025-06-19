package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/chdwlch/spark-pool/internal/miner"
	"github.com/chdwlch/spark-pool/internal/pool"
	"github.com/chdwlch/spark-pool/web"
)

func main() {
	// Parse command line flags
	var (
		port         = flag.String("port", "8080", "Server port")
		poolName     = flag.String("pool-name", "Ark Virtual Channels Demo Pool", "Mining pool name")
		operatorAddr = flag.String("operator-addr", "bc1qdemooperatoraddress", "Pool operator address")
		blockInterval = flag.Duration("block-interval", 30*time.Second, "Block reward interval for demo")
	)
	flag.Parse()

	// Setup logging
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Generate server key (in real implementation, this would come from Ark server)
	serverPrivKey, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		logger.Fatalf("Failed to generate server key: %v", err)
	}
	serverPubKey := serverPrivKey.PubKey()

	// Create pool manager
	poolManager := pool.NewManager(*poolName, *operatorAddr, serverPubKey)

	// Create miner manager
	minerManager := miner.NewManager(poolManager)

	// Create API server
	api := web.NewAPI(poolManager, minerManager)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Setup routes
	api.SetupRoutes(router)

	// Start WebSocket broadcaster
	api.StartBroadcaster()

	// Start block reward simulator
	go startBlockRewardSimulator(context.Background(), poolManager, *blockInterval, logger)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + *port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		logger.Infof("Starting mining pool demo server on port %s", *port)
		logger.Infof("Pool name: %s", *poolName)
		logger.Infof("Operator address: %s", *operatorAddr)
		logger.Infof("Block interval: %v", *blockInterval)
		logger.Infof("Dashboard available at: http://localhost:%s", *port)
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// startBlockRewardSimulator simulates block rewards at regular intervals
func startBlockRewardSimulator(ctx context.Context, poolManager *pool.Manager, interval time.Duration, logger *logrus.Logger) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	blockHeight := uint64(100000)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			blockHeight++
			
			logger.Infof("Processing block reward for height %d", blockHeight)
			
			blockReward, err := poolManager.ProcessBlockReward(ctx)
			if err != nil {
				logger.Errorf("Failed to process block reward: %v", err)
				continue
			}

			// Log block reward details
			logger.Infof("Block %d reward distributed:", blockHeight)
			for minerID, amount := range blockReward.Distributions {
				logger.Infof("  Miner %s: %d sats", minerID, amount)
			}

			// Log pool stats
			stats := poolManager.GetPoolStats()
			logger.Infof("Pool stats - Total miners: %d, Active: %d, Hash rate: %.2f TH/s, Total earned: %d sats",
				stats.TotalMiners, stats.ActiveMiners, stats.TotalHashRate/1e12, stats.TotalEarned)
		}
	}
} 