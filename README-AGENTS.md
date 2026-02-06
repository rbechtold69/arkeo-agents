# Arkeo Agents 🤖

**Permissionless RPC for AI Agents — Pay with USDC, no signup required.**

A fork of [Arkeo Protocol](https://github.com/arkeonetwork/arkeo) with native x402 payment support for AI agents.

---

## Why Arkeo Agents?

| Problem | Solution |
|---------|----------|
| Agents can't sign up for Infura/Alchemy | No signup needed — just pay |
| Agents can't do KYC | No KYC — just a wallet |
| POKT requires buying their token | Pay with USDC (stablecoins) |
| Complex API key management | x402 — payment per request |

---

## How It Works

```
1. Agent requests RPC data
2. Arkeo returns HTTP 402 + payment requirements
3. Agent signs x402 payment (USDC or ARKEO)
4. Arkeo verifies payment, returns data
```

**One HTTP request. No accounts. No API keys.**

---

## Supported Payment Methods

| Token | Network | Discount |
|-------|---------|----------|
| USDC | Ethereum Mainnet | — |
| USDC | Base L2 | Lower gas |
| ARKEO | Arkeo Chain | **15% off** |

---

## Quick Start (For Agent Developers)

### 1. Make a Request

```bash
curl https://rpc.arkeo.network/eth/v1/blockNumber
```

Response (HTTP 402):
```json
{
  "x402Version": 2,
  "error": "Payment required to access this RPC endpoint",
  "resource": {
    "url": "/eth/v1/blockNumber",
    "description": "Arkeo RPC Service: eth"
  },
  "accepts": [
    {
      "scheme": "exact",
      "network": "eip155:8453",
      "amount": "1000",
      "asset": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
      "payTo": "0x...",
      "maxTimeoutSeconds": 60,
      "extra": { "name": "USDC", "chain": "Base" }
    }
  ]
}
```

### 2. Sign Payment

Using x402 SDK:
```typescript
import { createPayment } from '@x402/core';

const payment = await createPayment({
  requirements: response.accepts[0],
  wallet: agentWallet,
});
```

### 3. Retry with Payment

```bash
curl https://rpc.arkeo.network/eth/v1/blockNumber \
  -H "X-PAYMENT: ${payment}"
```

Response (HTTP 200):
```json
{
  "jsonrpc": "2.0",
  "result": "0x134e82a",
  "id": 1
}
```

---

## Agent Framework Integration

### Eliza Framework
```typescript
// Coming soon
```

### LangChain
```typescript
// Coming soon
```

### AutoGPT
```python
# Coming soon
```

---

## Pricing

| Tier | Price per Request | Best For |
|------|-------------------|----------|
| Standard | $0.001 USDC | Most agents |
| Bulk | $0.0008 USDC | High volume |
| ARKEO | 15% discount | Token holders |

---

## Supported Chains

| Chain | RPC Endpoint | Status |
|-------|--------------|--------|
| Ethereum | /eth/ | 🟢 Live |
| Cosmos Hub | /cosmos/ | 🟢 Live |
| THORChain | /thorchain/ | 🟢 Live |
| Osmosis | /osmosis/ | 🟢 Live |
| More coming... | — | 🟡 Soon |

---

## Why Pay with ARKEO?

- **15% discount** on all requests
- **Support the network** — your payments go to providers
- **Governance rights** — ARKEO holders vote on protocol changes

---

## For Providers

Want to earn by running RPC nodes?

1. Run your node
2. Register on Arkeo
3. Set your prices
4. Receive USDC/ARKEO payments automatically

[Provider Documentation →](./docs/providers.md)

---

## Technical Specs

- **Protocol:** x402 v2
- **Settlement:** Ethereum, Base, Arkeo Chain
- **Payment:** Stablecoins (USDC) or ARKEO token
- **Latency:** ~50ms overhead for payment verification
- **Trust:** Non-custodial — facilitator never holds funds

---

## Links

- [x402 Protocol](https://x402.org)
- [ERC-8004: Trustless Agents](https://eips.ethereum.org/EIPS/eip-8004)
- [Arkeo Protocol](https://arkeo.network)
- [Discord](https://discord.gg/arkeo)

---

## Status

🚧 **Under Development** — Not production ready yet.

Built by the Arkeo community.
