package web

import (
	"net/http"

	"github.com/chdwlch/spark-pool/internal/miner"
	"github.com/chdwlch/spark-pool/internal/pool"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// API represents the REST API server
type API struct {
	poolManager    *pool.Manager
	minerManager   *miner.Manager
	upgrader       websocket.Upgrader
	clients        map[*websocket.Conn]bool
	broadcast      chan types.WebSocketMessage
	logger         *logrus.Logger
}

// NewAPI creates a new API server
func NewAPI(poolManager *pool.Manager, minerManager *miner.Manager) *API {
	return &API{
		poolManager:  poolManager,
		minerManager: minerManager,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for demo
			},
		},
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan types.WebSocketMessage, 100),
		logger:    logrus.New(),
	}
}

// SetupRoutes sets up the API routes
func (api *API) SetupRoutes(r *gin.Engine) {
	// API routes
	apiGroup := r.Group("/api/v1")
	{
		// Pool routes
		apiGroup.GET("/pool/stats", api.GetPoolStats)
		apiGroup.GET("/pool/miners", api.GetAllMiners)
		apiGroup.GET("/pool/channels", api.GetAllChannels)
		apiGroup.POST("/pool/block-reward", api.ProcessBlockReward)

		// Miner routes
		apiGroup.POST("/miners", api.AddMiner)
		apiGroup.GET("/miners/:id", api.GetMiner)
		apiGroup.PUT("/miners/:id/start", api.StartMiner)
		apiGroup.PUT("/miners/:id/stop", api.StopMiner)
		apiGroup.GET("/miners/:id/stats", api.GetMinerStats)

		// Channel routes
		apiGroup.GET("/channels/:id", api.GetChannel)
		apiGroup.POST("/channels/:id/close", api.CloseChannel)
	}

	// WebSocket route
	r.GET("/ws", api.HandleWebSocket)

	// Serve static files
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")
	
	// Web routes
	r.GET("/", api.ServeDashboard)
	r.GET("/miner/:id", api.ServeMinerDashboard)
}

// GetPoolStats returns pool statistics
func (api *API) GetPoolStats(c *gin.Context) {
	stats := api.poolManager.GetPoolStats()
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// GetAllMiners returns all miners
func (api *API) GetAllMiners(c *gin.Context) {
	miners := api.poolManager.GetAllMiners()
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    miners,
	})
}

// GetAllChannels returns all channels
func (api *API) GetAllChannels(c *gin.Context) {
	channels := api.poolManager.GetAllChannels()
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    channels,
	})
}

// ProcessBlockReward processes a block reward
func (api *API) ProcessBlockReward(c *gin.Context) {
	blockReward, err := api.poolManager.ProcessBlockReward(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Broadcast update to WebSocket clients
	api.broadcast <- types.WebSocketMessage{
		Type: "block_reward",
		Payload: blockReward,
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    blockReward,
	})
}

// AddMiner adds a new miner
func (api *API) AddMiner(c *gin.Context) {
	var req types.JoinPoolRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	simulator, err := api.minerManager.AddSimulator(req.MinerName, req.Address, req.HashRate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Broadcast new miner to WebSocket clients
	api.broadcast <- types.WebSocketMessage{
		Type: "miner_added",
		Payload: simulator,
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    simulator,
	})
}

// GetMiner returns a miner by ID
func (api *API) GetMiner(c *gin.Context) {
	minerID := c.Param("id")
	
	simulator, exists := api.minerManager.GetSimulator(minerID)
	if !exists {
		c.JSON(http.StatusNotFound, types.APIResponse{
			Success: false,
			Error:   "miner not found",
		})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    simulator,
	})
}

// StartMiner starts a miner
func (api *API) StartMiner(c *gin.Context) {
	minerID := c.Param("id")
	
	simulator, exists := api.minerManager.GetSimulator(minerID)
	if !exists {
		c.JSON(http.StatusNotFound, types.APIResponse{
			Success: false,
			Error:   "miner not found",
		})
		return
	}

	if err := simulator.StartMining(); err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Broadcast miner started
	api.broadcast <- types.WebSocketMessage{
		Type: "miner_started",
		Payload: map[string]string{"miner_id": minerID},
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    map[string]string{"status": "started"},
	})
}

// StopMiner stops a miner
func (api *API) StopMiner(c *gin.Context) {
	minerID := c.Param("id")
	
	simulator, exists := api.minerManager.GetSimulator(minerID)
	if !exists {
		c.JSON(http.StatusNotFound, types.APIResponse{
			Success: false,
			Error:   "miner not found",
		})
		return
	}

	if err := simulator.StopMining(); err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Broadcast miner stopped
	api.broadcast <- types.WebSocketMessage{
		Type: "miner_stopped",
		Payload: map[string]string{"miner_id": minerID},
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    map[string]string{"status": "stopped"},
	})
}

// GetMinerStats returns miner statistics
func (api *API) GetMinerStats(c *gin.Context) {
	minerID := c.Param("id")
	
	simulator, exists := api.minerManager.GetSimulator(minerID)
	if !exists {
		c.JSON(http.StatusNotFound, types.APIResponse{
			Success: false,
			Error:   "miner not found",
		})
		return
	}

	stats := simulator.GetStats()
	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    stats,
	})
}

// GetChannel returns a channel by ID
func (api *API) GetChannel(c *gin.Context) {
	channelID := c.Param("id")
	
	channel, exists := api.poolManager.GetChannel(channelID)
	if !exists {
		c.JSON(http.StatusNotFound, types.APIResponse{
			Success: false,
			Error:   "channel not found",
		})
		return
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    channel,
	})
}

// CloseChannel closes a channel
func (api *API) CloseChannel(c *gin.Context) {
	channelID := c.Param("id")
	
	// Find miner for this channel
	miners := api.poolManager.GetAllMiners()
	var minerID string
	for _, m := range miners {
		if m.ChannelID == channelID {
			minerID = m.ID
			break
		}
	}

	if minerID == "" {
		c.JSON(http.StatusNotFound, types.APIResponse{
			Success: false,
			Error:   "miner not found for channel",
		})
		return
	}

	if err := api.poolManager.CloseMinerChannel(minerID); err != nil {
		c.JSON(http.StatusInternalServerError, types.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Broadcast channel closed
	api.broadcast <- types.WebSocketMessage{
		Type: "channel_closed",
		Payload: map[string]string{"channel_id": channelID, "miner_id": minerID},
	}

	c.JSON(http.StatusOK, types.APIResponse{
		Success: true,
		Data:    map[string]string{"status": "closed"},
	})
}

// HandleWebSocket handles WebSocket connections
func (api *API) HandleWebSocket(c *gin.Context) {
	conn, err := api.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		api.logger.Errorf("Failed to upgrade connection: %v", err)
		return
	}

	// Register client
	api.clients[conn] = true

	// Send initial data
	initialData := map[string]interface{}{
		"pool_stats": api.poolManager.GetPoolStats(),
		"miners":     api.poolManager.GetAllMiners(),
		"channels":   api.poolManager.GetAllChannels(),
	}

	conn.WriteJSON(types.WebSocketMessage{
		Type:    "initial_data",
		Payload: initialData,
	})

	// Handle client disconnect
	defer func() {
		conn.Close()
		delete(api.clients, conn)
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// StartBroadcaster starts the WebSocket broadcaster
func (api *API) StartBroadcaster() {
	go func() {
		for message := range api.broadcast {
			for client := range api.clients {
				err := client.WriteJSON(message)
				if err != nil {
					api.logger.Errorf("Failed to send message to client: %v", err)
					client.Close()
					delete(api.clients, client)
				}
			}
		}
	}()
}

// ServeDashboard serves the main dashboard
func (api *API) ServeDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Mining Pool Demo",
	})
}

// ServeMinerDashboard serves a miner's dashboard
func (api *API) ServeMinerDashboard(c *gin.Context) {
	minerID := c.Param("id")
	c.HTML(http.StatusOK, "miner.html", gin.H{
		"title":   "Miner Dashboard",
		"miner_id": minerID,
	})
} 