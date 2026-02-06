# Arkeo Marketplace 2.0

**Upgraded marketplace with x402 payment support for AI agents and humans.**

## What's New

### For AI Agents
- x402 payment support (USDC via HTTP)
- Instant access - no token buying needed
- Discovery API for finding providers
- Integration guides for Eliza, AutoGPT, LangChain

### For Humans
- Pay-as-you-go with USDC (new)
- Traditional ARKEO contracts (existing)
- Better pricing comparison
- Reputation scores

### For Providers
- Enable/disable x402 payments
- Set USD pricing
- Earn from both AI agents AND humans
- Real-time earnings dashboard

---

## New Features Added

| Feature | Status | File |
|---------|--------|------|
| x402 payment toggle | 🔜 | admin/src/dashboard.jsx |
| USD pricing display | 🔜 | admin/src/dashboard.jsx |
| "Pay with USDC" button | 🔜 | admin/src/dashboard.jsx |
| AI Agent integration tab | 🔜 | admin/src/agents.jsx |
| Live activity feed | 🔜 | admin/src/dashboard.jsx |
| Reputation scores | 🔜 | admin/src/dashboard.jsx |
| AI Setup Wizard | 🔜 | admin/src/wizard.jsx |

---

## Development

```bash
# Install dependencies
npm install

# Build CSS and JS
npm run build

# Serve locally (requires Python backend or static server)
python3 -m http.server 8077 --directory admin
```

---

## Architecture

```
marketplace/
├── admin/
│   ├── src/
│   │   ├── dashboard.jsx    # Main marketplace UI
│   │   ├── agents.jsx       # AI agent integration (NEW)
│   │   ├── wizard.jsx       # Setup wizard (NEW)
│   │   └── styles.css
│   ├── index.html
│   └── images/
├── package.json
└── tailwind.config.cjs
```

---

## Backend Integration

This UI connects to:
- Arkeo chain (existing) - provider registry, contracts
- x402 Sentinel (NEW) - payment verification
- Discovery API (NEW) - agent-friendly provider lookup

---

*Part of the Arkeo Agents project - making Arkeo the #1 RPC marketplace for AI agents.*
