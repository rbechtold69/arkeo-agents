# Arkeo Onboarding Analysis & Recommendations

## Executive Summary

**The Problem:** Arkeo's onboarding is 10-20x more complex than competitors. While Infura/Alchemy take 2 minutes to get an API key, Arkeo requires running infrastructure, CLI commands, and blockchain transactions.

**The Solution:** Create a "Hosted Sentinel" tier that lets providers onboard in minutes, not hours. Keep the self-hosted option for power users.

---

## Part 1: Competitor Analysis

### Infura (Consumer Onboarding)

**Time to first API call: ~2 minutes**

| Step | Action | Time |
|------|--------|------|
| 1 | Sign up with email/Google | 30 sec |
| 2 | Verify email | 30 sec |
| 3 | Create project (1 click) | 10 sec |
| 4 | Copy API key | 5 sec |
| 5 | Use endpoint | Immediate |

**Endpoint format:**
```
https://mainnet.infura.io/v3/YOUR_API_KEY
```

**What makes it easy:**
- No blockchain transactions
- No wallet required for basic tier
- No infrastructure to run
- Instant activation

---

### Alchemy (Consumer Onboarding)

**Time to first API call: ~2 minutes**

| Step | Action | Time |
|------|--------|------|
| 1 | Sign up with email/Google | 30 sec |
| 2 | Create app (name + chain) | 30 sec |
| 3 | Copy API key | 5 sec |
| 4 | Use endpoint | Immediate |

**Endpoint format:**
```
https://eth-mainnet.g.alchemy.com/v2/YOUR_API_KEY
```

**What makes it easy:**
- Dashboard shows all chains
- Usage metrics built-in
- One-click upgrade paths
- Free tier = 300M compute units/month

---

### QuickNode (Consumer Onboarding)

**Time to first API call: ~3 minutes**

| Step | Action | Time |
|------|--------|------|
| 1 | Sign up | 30 sec |
| 2 | Pick chain + network | 30 sec |
| 3 | Choose region | 10 sec |
| 4 | Get endpoint | Immediate |

**Endpoint format:**
```
https://your-endpoint-name.quiknode.pro/YOUR_KEY/
```

---

## Part 2: Current Arkeo Process

### Provider Onboarding (Current)

**Time to first listing: 2-4 hours (best case)**

| Step | Action | Complexity | Time |
|------|--------|------------|------|
| 1 | Set up server | High | 30 min |
| 2 | Install arkeod | Medium | 15 min |
| 3 | Sync blockchain node | High | 1-4 hours |
| 4 | Create wallet | Medium | 5 min |
| 5 | Get ARKEO tokens | High | Variable |
| 6 | Bond provider (TX) | High | 10 min |
| 7 | Configure Sentinel | High | 30 min |
| 8 | Set metadata URI | Medium | 10 min |
| 9 | Mod-provider TX | High | 10 min |
| 10 | Verify listing | Medium | 10 min |

**CLI commands required:**
```bash
# Create wallet
arkeod keys add provider --keyring-backend file

# Get pubkey (complicated)
RAWPUBKEY=$(arkeod keys show provider -p | jq -r .key)
PUBKEY=$(arkeod debug pubkey-raw $RAWPUBKEY | grep 'Bech32 Acc:' | sed "s|Bech32 Acc: ||g")

# Bond provider
arkeod tx arkeo bond-provider -y --from provider --fees 50uarkeo \
  --keyring-backend file -- "$PUBKEY" "$SERVICE" "$BOND"

# Modify provider
arkeod tx arkeo mod-provider -y --from provider --keyring-backend file \
  -- "$PUBKEY" "$SERVICE" "$METADATAURI" $STATUS $MIN_DURATION \
  $MAX_DURATION $SUB_RATE $PAYG_RATE $SETTLEMENT_DURATION
```

**Pain points:**
1. Must run own infrastructure (server, node, sentinel)
2. Multiple complex CLI commands
3. Need to acquire ARKEO tokens first
4. Must understand nonces, pubkeys, signing
5. No web UI for provider registration
6. Manual metadata.json hosting
7. No validation or error checking
8. No guided wizard

---

### Consumer/Subscriber Onboarding (Current)

**Time to first API call: 30-60 minutes**

| Step | Action | Complexity | Time |
|------|--------|------------|------|
| 1 | Install arkeod CLI | Medium | 10 min |
| 2 | Create wallet | Medium | 5 min |
| 3 | Get ARKEO tokens | High | Variable |
| 4 | Find provider (CLI) | High | 10 min |
| 5 | Evaluate provider | Medium | 5 min |
| 6 | Open contract (TX) | High | 10 min |
| 7 | Get claim signature | High | 5 min |
| 8 | Make API call | Medium | 5 min |

**Pain points:**
1. Need CLI installed locally
2. Need ARKEO tokens before anything
3. No simple "just give me an endpoint" option
4. Contract mechanics are confusing
5. Claim/nonce system is opaque

---

## Part 3: What Phil Has Built

### arkeo-data-engine/provider-core

Phil created a Docker-based solution that bundles:
- arkeod (chain daemon)
- Sentinel (proxy)
- Admin UI (web interface)
- Claim automation scripts

**Improvements:**
- One Docker container instead of manual setup
- Web UI for provider management
- Automated claim settlements
- Configuration via JSON file

**Still missing:**
- Still requires own server
- Still requires ARKEO tokens
- Still requires manual chain sync
- No hosted option

---

## Part 4: Gap Analysis

| Feature | Competitors | Arkeo Current | Gap |
|---------|-------------|---------------|-----|
| Time to start | 2 min | 2-4 hours | 60-120x slower |
| Wallet required | No (basic) | Yes | Friction |
| Tokens required | No (basic) | Yes | Huge barrier |
| Infrastructure | None | Required | Massive barrier |
| Web signup | Yes | No | Missing |
| One-click deploy | Yes | No | Missing |
| Free tier | Yes | No | Missing |
| API key model | Yes | No | Missing |
| Usage dashboard | Yes | Partial | Needs work |

---

## Part 5: Recommendations

### Tier 1: "Instant Access" (Like Infura)

**For consumers who just want an endpoint:**

1. **Web signup** → Email + password
2. **Pick chain** → Dropdown (Ethereum, Bitcoin, etc.)
3. **Get API key** → Instant
4. **Use endpoint** → `https://rpc.arkeo.network/v1/{api_key}/eth-mainnet`

**How it works behind the scenes:**
- Arkeo runs hosted Sentinels
- API keys map to internal contracts
- Usage tracked, billed in USDC (x402) or fiat
- No ARKEO tokens needed for consumer

**Implementation:**
- Build API gateway service
- Handle auth/rate limiting
- Auto-provision contracts against hosted providers
- Accept credit card or USDC

---

### Tier 2: "Self-Custody" (Current model, simplified)

**For users who want decentralized + control:**

Keep the current model but simplify:

1. **Web UI for provider registration**
   - Connect Keplr wallet
   - Fill form (service, rates, metadata)
   - One-click bond + register transaction

2. **One-line provider setup**
   ```bash
   curl -sSL https://get.arkeo.network/provider | bash
   ```
   - Downloads Docker image
   - Prompts for mnemonic
   - Auto-configures everything

3. **Web UI for subscriber**
   - Browse providers in marketplace
   - Click "Subscribe"
   - Signs transaction via Keplr
   - Get endpoint URL instantly

---

### Tier 3: "Hosted Provider" (New)

**For node operators who don't want to run Sentinel:**

1. Provider runs their RPC node
2. Registers on Arkeo web UI
3. Points to their node URL
4. Arkeo runs hosted Sentinel in front of it
5. Provider still earns, less infrastructure

**Benefits:**
- Lowers barrier for existing node operators
- They just need to expose an RPC endpoint
- Arkeo handles billing/claims/contracts

---

## Part 6: Priority Implementation Order

### Phase 1: Consumer Quick-Start (Week 1-2)
- [ ] API Gateway service (hosted sentinel cluster)
- [ ] Web signup flow (email/password)
- [ ] API key generation
- [ ] Simple endpoint format
- [ ] x402 USDC payments (already built!)

### Phase 2: Provider Web UI (Week 2-3)
- [ ] Keplr wallet connection
- [ ] Provider registration form
- [ ] One-click bond transaction
- [ ] Status dashboard

### Phase 3: One-Line Install (Week 3-4)
- [ ] Install script
- [ ] Docker Compose auto-config
- [ ] Wizard for mnemonic/settings
- [ ] Auto-registration option

### Phase 4: Hosted Provider Tier (Week 4-5)
- [ ] Hosted Sentinel cluster
- [ ] Provider points to external RPC
- [ ] Arkeo proxies and handles contracts

---

## Part 7: Proposed New Flows

### New Consumer Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    ARKEO MARKETPLACE                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   [Sign Up Free]     or     [Connect Wallet]               │
│                                                             │
│   ┌─────────────────────────────────────────────────────┐  │
│   │ Pick your chain:  [Ethereum ▼]                      │  │
│   │                                                     │  │
│   │ Your endpoint:                                      │  │
│   │ https://rpc.arkeo.network/v1/abc123/eth-mainnet    │  │
│   │                                          [Copy]     │  │
│   │                                                     │  │
│   │ Free tier: 100,000 requests/month                   │  │
│   │ Need more? [Upgrade to Pro]                         │  │
│   └─────────────────────────────────────────────────────┘  │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### New Provider Flow

```
┌─────────────────────────────────────────────────────────────┐
│                 BECOME A PROVIDER                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│   Step 1: Connect Wallet                                    │
│   [Connect Keplr]  ✓ Connected: arkeo1abc...               │
│                                                             │
│   Step 2: Choose Your Setup                                 │
│   ○ I'll run my own Sentinel (advanced)                    │
│   ● Host my RPC through Arkeo (easy)                       │
│                                                             │
│   Step 3: Configure                                         │
│   Service: [Ethereum Mainnet ▼]                            │
│   Your RPC URL: [https://my-node.com:8545_______]          │
│   Rate (per 1000 req): [0.10 USDC ▼]                       │
│                                                             │
│   Step 4: Bond & Register                                   │
│   Bond amount: 1000 ARKEO (min: 1)                         │
│   [Register Provider] ← Signs 1 transaction                │
│                                                             │
│   ✓ You're live! Earning starts immediately.               │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Part 8: Technical Requirements

### API Gateway Service

```go
// Core components needed:
type APIGateway struct {
    SentinelPool    *SentinelCluster  // Pool of hosted sentinels
    AuthService     *AuthService       // API key validation
    RateLimiter     *RateLimiter       // Per-key limits
    UsageTracker    *UsageTracker      // Request counting
    X402Handler     *X402Handler       // Payment processing
    ContractManager *ContractManager   // Auto-contract provisioning
}

// Endpoint routing:
// GET /v1/{api_key}/{chain}/{method}
// POST /v1/{api_key}/{chain}  (JSON-RPC body)
```

### Web Registration Service

```typescript
// Provider registration flow:
interface ProviderRegistration {
  wallet: string;           // From Keplr
  service: ServiceType;     // Dropdown selection
  rpcEndpoint?: string;     // For hosted tier
  subscriptionRate: Coin;   // From form
  payAsYouGoRate: Coin;     // From form
  bondAmount: Coin;         // From form
}

// Auto-generates:
// 1. bond-provider TX
// 2. mod-provider TX
// 3. metadata.json (hosted on Arkeo CDN)
```

---

## Summary

**Do this first:** API Gateway + Web signup for consumers. Get people using Arkeo in 2 minutes, not 2 hours.

**Do this second:** Web UI for provider registration. Get those 8-10 providers listed without CLI pain.

**Do this third:** Hosted provider tier. Let node operators join without running Sentinel.

The x402 integration we built today enables the payment layer. Now you need the onboarding layer to match competitor simplicity.
