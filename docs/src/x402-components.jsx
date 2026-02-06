// x402 Payment Components for Arkeo Marketplace 2.0
// These components add USDC payment support for AI agents and humans

import React, { useState } from 'react';

// =============================================================================
// ICONS
// =============================================================================

const X402Icons = {
  USDC: () => (
    <svg className="w-5 h-5" viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle cx="16" cy="16" r="16" fill="#2775CA"/>
      <path d="M20.5 18.5c0-2-1.5-2.75-4.5-3.25-2.25-.375-2.75-.875-2.75-1.875 0-1 .875-1.625 2.25-1.625 1.25 0 2.125.5 2.5 1.5.125.25.375.375.625.375h.75c.375 0 .625-.25.625-.625v-.125c-.375-1.375-1.5-2.375-3-2.625V8.5c0-.375-.25-.625-.625-.625h-.75c-.375 0-.625.25-.625.625v1.625c-2 .25-3.25 1.5-3.25 3.125 0 2 1.5 2.75 4.5 3.25 2 .375 2.75.875 2.75 1.875s-1 1.75-2.375 1.75c-1.75 0-2.375-.75-2.625-1.75-.125-.25-.375-.375-.625-.375h-.75c-.375 0-.625.25-.625.625v.125c.375 1.625 1.5 2.625 3.5 2.875v1.625c0 .375.25.625.625.625h.75c.375 0 .625-.25.625-.625V21.5c2-.25 3.25-1.5 3.25-3z" fill="#fff"/>
    </svg>
  ),
  Robot: () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
    </svg>
  ),
  Lightning: () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
    </svg>
  ),
  Code: () => (
    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
    </svg>
  ),
  Copy: () => (
    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
    </svg>
  ),
};

// =============================================================================
// X402 PRICING DISPLAY
// =============================================================================

export const X402PricingBadge = ({ priceUSD = 0.001, arkeoPrice, arkeoDiscount = 15 }) => {
  return (
    <div className="flex items-center gap-2 bg-blue-500/10 border border-blue-500/20 rounded-lg px-3 py-2">
      <X402Icons.USDC />
      <div>
        <p className="text-blue-400 text-sm font-bold">${priceUSD.toFixed(4)}</p>
        <p className="text-blue-300 text-xs">per request (USDC)</p>
      </div>
      {arkeoDiscount > 0 && (
        <div className="ml-2 bg-emerald-500/20 text-emerald-400 text-xs px-2 py-1 rounded">
          {arkeoDiscount}% off with ARKEO
        </div>
      )}
    </div>
  );
};

// =============================================================================
// X402 TOGGLE (FOR PROVIDER SETTINGS)
// =============================================================================

export const X402Toggle = ({ enabled, onToggle, priceUSD, onPriceChange }) => {
  return (
    <div className="card-surface rounded-xl p-4 border border-[var(--border)]">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-3">
          <X402Icons.USDC />
          <div>
            <p className="text-white font-medium">x402 Payments</p>
            <p className="text-secondaryText text-sm">Accept USDC from AI agents</p>
          </div>
        </div>
        <button
          onClick={() => onToggle?.(!enabled)}
          className={`relative w-12 h-6 rounded-full transition-colors ${enabled ? 'bg-blue-500' : 'bg-gray-600'}`}
        >
          <span className={`absolute top-1 w-4 h-4 bg-white rounded-full transition-transform ${enabled ? 'left-7' : 'left-1'}`} />
        </button>
      </div>
      
      {enabled && (
        <div className="mt-4 pt-4 border-t border-[var(--border)]">
          <label className="text-secondaryText text-sm mb-2 block">Price per request (USD)</label>
          <div className="flex items-center gap-2">
            <span className="text-white">$</span>
            <input
              type="number"
              step="0.0001"
              min="0.0001"
              value={priceUSD}
              onChange={(e) => onPriceChange?.(parseFloat(e.target.value))}
              className="bg-[#1E222C] border border-[var(--border)] rounded-lg px-3 py-2 text-white w-32 focus:border-blue-500 focus:ring-1 focus:ring-blue-500"
            />
            <span className="text-secondaryText text-sm">per request</span>
          </div>
        </div>
      )}
    </div>
  );
};

// =============================================================================
// PAY WITH USDC BUTTON
// =============================================================================

export const PayWithUSDCButton = ({ provider, onPay }) => {
  const [isLoading, setIsLoading] = useState(false);
  
  const handlePay = async () => {
    setIsLoading(true);
    try {
      await onPay?.(provider);
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <button
      onClick={handlePay}
      disabled={isLoading}
      className="w-full bg-gradient-to-r from-blue-500 to-blue-600 text-white font-semibold py-3 rounded-xl flex items-center justify-center gap-2 hover:from-blue-600 hover:to-blue-700 transition-all shadow-lg shadow-blue-500/25 disabled:opacity-50"
    >
      {isLoading ? (
        <div className="w-5 h-5 border-2 border-white border-t-transparent rounded-full animate-spin" />
      ) : (
        <>
          <X402Icons.USDC />
          <span>Pay with USDC</span>
          <X402Icons.Lightning />
        </>
      )}
    </button>
  );
};

// =============================================================================
// AI AGENT INTEGRATION PANEL
// =============================================================================

export const AIAgentIntegration = ({ provider, endpoint }) => {
  const [copied, setCopied] = useState(null);
  
  const copyToClipboard = (text, key) => {
    navigator.clipboard.writeText(text);
    setCopied(key);
    setTimeout(() => setCopied(null), 2000);
  };
  
  const curlExample = `curl -X POST ${endpoint || 'https://eth.arkeo.network'} \\
  -H "Content-Type: application/json" \\
  -H "X-PAYMENT: <your-x402-signature>" \\
  -d '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}'`;

  const pythonExample = `from x402 import Client

client = Client(
    wallet_key="YOUR_PRIVATE_KEY",
    network="base"  # or "ethereum"
)

response = client.request(
    url="${endpoint || 'https://eth.arkeo.network'}",
    method="POST",
    json={"jsonrpc": "2.0", "method": "eth_blockNumber", "params": [], "id": 1}
)

print(response.json())`;

  const elizaExample = `// eliza.config.ts
import { arkeoPlugin } from '@arkeo/eliza-plugin';

export default {
  plugins: [
    arkeoPlugin({
      walletKey: process.env.WALLET_KEY,
      defaultProvider: "${provider?.moniker || 'auto'}",
      network: "base"
    })
  ]
};`;

  const CodeBlock = ({ code, language, copyKey }) => (
    <div className="relative">
      <pre className="bg-[#0a0c10] rounded-lg p-4 overflow-x-auto text-sm text-gray-300 border border-[var(--border)]">
        <code>{code}</code>
      </pre>
      <button
        onClick={() => copyToClipboard(code, copyKey)}
        className="absolute top-2 right-2 p-2 bg-[#1E222C] rounded-lg text-secondaryText hover:text-white transition-colors"
      >
        {copied === copyKey ? '✓' : <X402Icons.Copy />}
      </button>
    </div>
  );

  return (
    <div className="card-surface rounded-2xl p-6 border border-[var(--border)]">
      <div className="flex items-center gap-3 mb-6">
        <div className="p-2.5 rounded-xl bg-purple-500/15 text-purple-400">
          <X402Icons.Robot />
        </div>
        <div>
          <h3 className="text-xl font-bold text-white">AI Agent Integration</h3>
          <p className="text-secondaryText text-sm">Connect your agent in minutes</p>
        </div>
      </div>

      <div className="space-y-6">
        {/* Quick Start */}
        <div>
          <h4 className="text-white font-medium mb-3 flex items-center gap-2">
            <X402Icons.Lightning />
            Quick Start (cURL)
          </h4>
          <CodeBlock code={curlExample} language="bash" copyKey="curl" />
        </div>

        {/* Python SDK */}
        <div>
          <h4 className="text-white font-medium mb-3 flex items-center gap-2">
            <X402Icons.Code />
            Python SDK
          </h4>
          <CodeBlock code={pythonExample} language="python" copyKey="python" />
        </div>

        {/* Eliza Plugin */}
        <div>
          <h4 className="text-white font-medium mb-3 flex items-center gap-2">
            <X402Icons.Robot />
            Eliza Framework
          </h4>
          <CodeBlock code={elizaExample} language="typescript" copyKey="eliza" />
        </div>

        {/* More Resources */}
        <div className="flex flex-wrap gap-3 pt-4 border-t border-[var(--border)]">
          <a href="https://docs.arkeo.network/agents" className="text-arkeo text-sm font-medium hover:underline flex items-center gap-1">
            📖 Full Documentation
          </a>
          <a href="https://github.com/arkeonetwork/arkeo-agents" className="text-arkeo text-sm font-medium hover:underline flex items-center gap-1">
            💻 GitHub
          </a>
          <a href="https://discord.gg/arkeo" className="text-arkeo text-sm font-medium hover:underline flex items-center gap-1">
            💬 Discord Support
          </a>
        </div>
      </div>
    </div>
  );
};

// =============================================================================
// LIVE ACTIVITY FEED
// =============================================================================

export const LiveActivityFeed = ({ activities = [] }) => {
  // Demo activities if none provided
  const demoActivities = activities.length > 0 ? activities : [
    { id: 1, type: 'request', agent: 'ElizaBot', provider: 'Red_5', chain: 'ethereum', amount: 0.001, time: '2s ago' },
    { id: 2, type: 'request', agent: 'AutoGPT-7', provider: 'Node_42', chain: 'bitcoin', amount: 0.0005, time: '5s ago' },
    { id: 3, type: 'payment', agent: 'LangChain-Agent', provider: 'Arkeo_1', chain: 'cosmos', amount: 0.001, time: '8s ago' },
    { id: 4, type: 'request', agent: 'CrewAI', provider: 'Red_5', chain: 'ethereum', amount: 0.001, time: '12s ago' },
  ];

  return (
    <div className="card-surface rounded-2xl p-5 border border-[var(--border)]">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <div className="w-2 h-2 bg-emerald-400 rounded-full animate-pulse" />
          <h3 className="text-white font-bold">Live Activity</h3>
        </div>
        <span className="text-secondaryText text-xs">Real-time x402 payments</span>
      </div>

      <div className="space-y-3 max-h-64 overflow-y-auto">
        {demoActivities.map((activity) => (
          <div key={activity.id} className="flex items-center justify-between p-3 bg-[#1E222C] rounded-lg">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-blue-500/15 rounded-lg text-blue-400">
                <X402Icons.Robot />
              </div>
              <div>
                <p className="text-white text-sm font-medium">{activity.agent}</p>
                <p className="text-secondaryText text-xs">{activity.chain} → {activity.provider}</p>
              </div>
            </div>
            <div className="text-right">
              <p className="text-emerald-400 text-sm font-medium">${activity.amount.toFixed(4)}</p>
              <p className="text-secondaryText text-xs">{activity.time}</p>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

// =============================================================================
// HERO BANNER FOR AI AGENTS
// =============================================================================

export const AIAgentHeroBanner = ({ onLearnMore }) => {
  return (
    <div className="relative overflow-hidden rounded-2xl bg-gradient-to-r from-purple-900/50 via-blue-900/50 to-cyan-900/50 border border-purple-500/20 p-8 mb-8">
      {/* Animated background */}
      <div className="absolute inset-0 bg-[url('data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iNjAiIGhlaWdodD0iNjAiIHZpZXdCb3g9IjAgMCA2MCA2MCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48ZyBmaWxsPSJub25lIiBmaWxsLXJ1bGU9ImV2ZW5vZGQiPjxnIGZpbGw9IiNmZmYiIGZpbGwtb3BhY2l0eT0iMC4wNSI+PGNpcmNsZSBjeD0iMzAiIGN5PSIzMCIgcj0iMiIvPjwvZz48L2c+PC9zdmc+')] opacity-50" />
      
      <div className="relative z-10 flex flex-col md:flex-row items-center justify-between gap-6">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-3">
            <span className="px-3 py-1 bg-purple-500/20 text-purple-300 text-sm font-medium rounded-full border border-purple-500/30">
              🆕 NEW
            </span>
          </div>
          <h2 className="text-3xl font-bold text-white mb-3">
            AI Agents Welcome
          </h2>
          <p className="text-gray-300 text-lg mb-4">
            Pay for blockchain data with <span className="text-blue-400 font-semibold">USDC</span> via x402. 
            No token buying. No accounts. Just HTTP requests.
          </p>
          <div className="flex flex-wrap gap-4">
            <div className="flex items-center gap-2 text-gray-300">
              <X402Icons.Lightning />
              <span>Instant access</span>
            </div>
            <div className="flex items-center gap-2 text-gray-300">
              <X402Icons.USDC />
              <span>Pay with stablecoins</span>
            </div>
            <div className="flex items-center gap-2 text-gray-300">
              <X402Icons.Robot />
              <span>Agent-native</span>
            </div>
          </div>
        </div>
        
        <div className="flex flex-col gap-3">
          <button
            onClick={onLearnMore}
            className="px-6 py-3 bg-gradient-to-r from-purple-500 to-blue-500 text-white font-semibold rounded-xl hover:from-purple-600 hover:to-blue-600 transition-all shadow-lg shadow-purple-500/25"
          >
            Get Started →
          </button>
          <a
            href="https://docs.arkeo.network/agents"
            className="px-6 py-3 bg-white/10 text-white font-medium rounded-xl hover:bg-white/20 transition-all text-center"
          >
            View Docs
          </a>
        </div>
      </div>
    </div>
  );
};

// =============================================================================
// EXPORT ALL
// =============================================================================

export default {
  X402PricingBadge,
  X402Toggle,
  PayWithUSDCButton,
  AIAgentIntegration,
  LiveActivityFeed,
  AIAgentHeroBanner,
};
