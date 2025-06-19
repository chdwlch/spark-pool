# Ark Virtual Channels Mining Pool Demo

A demonstration of Bitcoin mining pool payouts using Ark's Virtual Channels technology. This project simulates a mining pool where miners join, contribute hash power, and receive payments through trustless Virtual Channels.

## 🎯 Demo Overview

This demo showcases:

- **Virtual Channels**: Trustless payment channels between pool operator and miners
- **Real-time Mining Simulation**: Simulated miners with configurable hash rates
- **Automatic Payouts**: Block rewards distributed through Virtual Channels
- **Live Dashboard**: Real-time monitoring of pool, miners, and channels
- **WebSocket Updates**: Live updates across all connected clients

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Pool Operator │    │   Virtual       │    │   Individual    │
│   Dashboard     │◄──►│   Channels      │◄──►│   Miner         │
│                 │    │   (VTXOs)       │    │   Dashboards    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   REST API      │    │   WebSocket     │    │   Mining        │
│   Endpoints     │    │   Real-time     │    │   Simulator     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or later
- Access to Ark Network codebase (for dependencies)

### Installation

1. **Clone and setup dependencies**:

```bash
# Make sure you're in the arkd directory
cd mining-pool-demo

# Initialize Go modules
go mod tidy
```

2. **Run the pool operator**:

```bash
go run cmd/pool-operator/main.go
```

3. **Access the dashboard**:
   - Pool Operator Dashboard: http://localhost:8080
   - Individual Miner Dashboard: http://localhost:8080/miner/{miner-id}

## 🎮 Demo Scenarios

### Scenario 1: Basic Mining Pool Operation

1. **Start the pool operator** (runs on port 8080)
2. **Add miners** through the web interface:
   - Miner Name: "Alice"
   - Bitcoin Address: "bc1qalice..."
   - Hash Rate: 50 TH/s
3. **Start mining** for each miner
4. **Watch automatic payouts** every 30 seconds
5. **Monitor Virtual Channels** being created and updated

### Scenario 2: Multiple Miners with Different Hash Rates

1. **Add multiple miners** with varying hash rates:
   - Alice: 100 TH/s
   - Bob: 75 TH/s
   - Charlie: 25 TH/s
2. **Start all miners** simultaneously
3. **Observe proportional payouts** based on hash rate contribution
4. **Monitor channel balances** updating in real-time

### Scenario 3: Channel Management

1. **Create channels** for multiple miners
2. **Process several block rewards** to build up balances
3. **Close a channel** and observe the withdrawal process
4. **Verify final settlement** on the blockchain

## 📊 Key Features

### Virtual Channels (Spillman Channels)

- **Trustless**: No need to trust the pool operator
- **Instant**: Off-chain payments with on-chain settlement
- **Efficient**: No on-chain fees for each payment
- **Secure**: Time-locked contracts prevent fraud

### Real-time Dashboard

- **Pool Statistics**: Total miners, hash rate, earnings
- **Miner Management**: Add, start, stop miners
- **Channel Monitoring**: View all active Virtual Channels
- **Live Updates**: WebSocket-powered real-time updates

### Mining Simulation

- **Configurable Hash Rates**: Adjust miner performance
- **Realistic Statistics**: Shares, efficiency, uptime
- **Automatic Rewards**: Block rewards every 30 seconds
- **Individual Dashboards**: Per-miner monitoring

## 🔧 Configuration

### Command Line Options

```bash
go run cmd/pool-operator/main.go \
  --port 8080 \
  --pool-name "My Ark Mining Pool" \
  --operator-addr "bc1qpooloperator..." \
  --block-interval 30s
```

### Environment Variables

- `PORT`: Server port (default: 8080)
- `POOL_NAME`: Mining pool name
- `OPERATOR_ADDR`: Pool operator Bitcoin address
- `BLOCK_INTERVAL`: Block reward interval for demo

## 📡 API Endpoints

### Pool Management

- `GET /api/v1/pool/stats` - Get pool statistics
- `GET /api/v1/pool/miners` - List all miners
- `GET /api/v1/pool/channels` - List all channels
- `POST /api/v1/pool/block-reward` - Process block reward

### Miner Management

- `POST /api/v1/miners` - Add new miner
- `GET /api/v1/miners/:id` - Get miner details
- `PUT /api/v1/miners/:id/start` - Start miner
- `PUT /api/v1/miners/:id/stop` - Stop miner
- `GET /api/v1/miners/:id/stats` - Get miner statistics

### Channel Management

- `GET /api/v1/channels/:id` - Get channel details
- `POST /api/v1/channels/:id/close` - Close channel

### WebSocket

- `GET /ws` - WebSocket connection for real-time updates

## 🎨 Web Interface

### Pool Operator Dashboard

- **Overview**: Pool statistics and controls
- **Miner Management**: Add and control miners
- **Channel Monitoring**: View all Virtual Channels
- **Block Rewards**: Manual and automatic processing

### Individual Miner Dashboard

- **Mining Controls**: Start/stop mining, adjust hash rate
- **Statistics**: Real-time mining performance
- **Channel Information**: Virtual Channel details
- **Payment History**: Track all received payments

## 🔍 Understanding Virtual Channels

### How It Works

1. **Channel Creation**: Pool operator creates a VTXO with a multi-signature script
2. **Initial Funding**: Pool operator funds the channel with initial capital
3. **Off-chain Updates**: Payments happen through partially signed transactions
4. **Final Settlement**: Channel closes with on-chain transaction

### Channel Script Structure

```go
// Collaborative close (pool operator + miner + server)
MultisigClosure{
    PubKeys: [poolOperator, miner, server],
    Type: MultisigTypeChecksig,
}

// Pool operator unilateral close after 24 hours
CSVMultisigClosure{
    Locktime: 144 blocks,
    MultisigClosure: {
        PubKeys: [poolOperator, server],
        Type: MultisigTypeChecksig,
    },
}

// Miner unilateral close after 48 hours
CSVMultisigClosure{
    Locktime: 288 blocks,
    MultisigClosure: {
        PubKeys: [miner, server],
        Type: MultisigTypeChecksig,
    },
}
```

## 🧪 Testing

### Manual Testing

1. **Start the server** and add a few miners
2. **Open multiple browser tabs** for different miners
3. **Start mining** and observe real-time updates
4. **Process block rewards** and watch channel updates
5. **Close channels** and verify final settlement

### API Testing

```bash
# Add a miner
curl -X POST http://localhost:8080/api/v1/miners \
  -H "Content-Type: application/json" \
  -d '{"miner_name":"TestMiner","address":"bc1qtest","hash_rate":50}'

# Get pool stats
curl http://localhost:8080/api/v1/pool/stats

# Process block reward
curl -X POST http://localhost:8080/api/v1/pool/block-reward
```

## 🚨 Important Notes

### Demo Limitations

- **Simulated Environment**: This is a demo, not production code
- **Fake Keys**: All cryptographic keys are generated for demo purposes
- **Simulated Mining**: Hash rates and shares are simulated
- **No Real Bitcoin**: All amounts are in satoshis but not real transactions

### Production Considerations

- **Real Ark Integration**: Connect to actual Ark Network
- **Secure Key Management**: Proper key storage and management
- **Database Persistence**: Store state in a real database
- **Error Handling**: Robust error handling and recovery
- **Monitoring**: Production monitoring and alerting

## 🤝 Contributing

This is a demo project for hackathon purposes. For production use:

1. **Security Review**: Thorough security audit required
2. **Ark Integration**: Proper integration with Ark Network
3. **Testing**: Comprehensive test suite
4. **Documentation**: Production deployment guide

## 📚 Resources

- [Ark Network Documentation](https://docs.ark.network/)
- [Virtual Channels Blog Post](https://blog.arklabs.xyz/bitcoin-virtual-channels/)
- [Bitcoin Lightning Network](https://lightning.network/)
- [Spillman Channels Paper](https://arxiv.org/abs/1702.05812)

## 📄 License

This project is for educational and demo purposes. Please refer to the Ark Network license for production use.

---

**Happy Mining! 🚀⛏️**
