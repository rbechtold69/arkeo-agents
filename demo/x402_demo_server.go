package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// PaymentRequirements defines what payment is accepted
type PaymentRequirements struct {
	Scheme            string                 `json:"scheme"`
	Network           string                 `json:"network"`
	Amount            string                 `json:"amount"`
	Asset             string                 `json:"asset"`
	PayTo             string                 `json:"payTo"`
	MaxTimeoutSeconds int                    `json:"maxTimeoutSeconds"`
	Extra             map[string]interface{} `json:"extra,omitempty"`
}

// PaymentRequiredResponse is the HTTP 402 response
type PaymentRequiredResponse struct {
	X402Version int                   `json:"x402Version"`
	Error       string                `json:"error,omitempty"`
	Resource    ResourceInfo          `json:"resource"`
	Accepts     []PaymentRequirements `json:"accepts"`
}

type ResourceInfo struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// Demo RPC response
type RPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
}

func main() {
	fmt.Println("🚀 Arkeo x402 Demo Server")
	fmt.Println("========================")
	fmt.Println("")
	fmt.Println("Starting server on http://localhost:8402")
	fmt.Println("")
	fmt.Println("Try these endpoints:")
	fmt.Println("  1. GET http://localhost:8402/            - Welcome page")
	fmt.Println("  2. GET http://localhost:8402/x402/eth    - Without payment (gets 402)")
	fmt.Println("  3. GET http://localhost:8402/x402/eth with X-PAYMENT header - With payment (gets data)")
	fmt.Println("")

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/x402/", handleX402)

	log.Fatal(http.ListenAndServe(":8402", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Arkeo x402 Demo</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        h1 { color: #333; }
        .box { background: #f5f5f5; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .success { background: #d4edda; border: 1px solid #28a745; }
        .error { background: #f8d7da; border: 1px solid #dc3545; }
        .payment { background: #fff3cd; border: 1px solid #ffc107; }
        pre { background: #2d2d2d; color: #f8f8f2; padding: 15px; border-radius: 5px; overflow-x: auto; }
        button { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; margin: 5px; }
        button:hover { background: #0056b3; }
        #response { margin-top: 20px; }
    </style>
</head>
<body>
    <h1>🤖 Arkeo x402 Demo</h1>
    <p>This demo shows how AI agents can pay for RPC data using the x402 protocol.</p>
    
    <div class="box">
        <h3>Test the x402 Flow</h3>
        <button onclick="testWithoutPayment()">1. Request WITHOUT Payment</button>
        <button onclick="testWithPayment()">2. Request WITH Payment</button>
    </div>
    
    <div id="response"></div>
    
    <script>
        async function testWithoutPayment() {
            const response = await fetch('/x402/eth');
            const data = await response.json();
            
            document.getElementById('response').innerHTML = ` + "`" + `
                <div class="box payment">
                    <h3>⚠️ HTTP 402 - Payment Required</h3>
                    <p>The server requires payment. Here are the accepted payment methods:</p>
                    <pre>${JSON.stringify(data, null, 2)}</pre>
                </div>
            ` + "`" + `;
        }
        
        async function testWithPayment() {
            const response = await fetch('/x402/eth', {
                headers: {
                    'X-PAYMENT': 'demo-payment-token-12345'
                }
            });
            const data = await response.json();
            
            document.getElementById('response').innerHTML = ` + "`" + `
                <div class="box success">
                    <h3>✅ HTTP 200 - Success!</h3>
                    <p>Payment verified. Here's your RPC data:</p>
                    <pre>${JSON.stringify(data, null, 2)}</pre>
                </div>
            ` + "`" + `;
        }
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func handleX402(w http.ResponseWriter, r *http.Request) {
	// Check for payment header
	paymentHeader := r.Header.Get("X-PAYMENT")

	if paymentHeader == "" {
		// No payment - return 402 Payment Required
		response := PaymentRequiredResponse{
			X402Version: 2,
			Error:       "Payment required to access this RPC endpoint",
			Resource: ResourceInfo{
				URL:         r.URL.String(),
				Description: "Arkeo RPC Service - Ethereum",
			},
			Accepts: []PaymentRequirements{
				{
					Scheme:            "exact",
					Network:           "eip155:8453",
					Amount:            "1000",
					Asset:             "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
					PayTo:             "0xYourProviderAddress",
					MaxTimeoutSeconds: 60,
					Extra: map[string]interface{}{
						"name":  "USDC",
						"chain": "Base",
					},
				},
				{
					Scheme:            "exact",
					Network:           "eip155:1",
					Amount:            "1000",
					Asset:             "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
					PayTo:             "0xYourProviderAddress",
					MaxTimeoutSeconds: 60,
					Extra: map[string]interface{}{
						"name":  "USDC",
						"chain": "Ethereum",
					},
				},
				{
					Scheme:            "exact",
					Network:           "arkeo:arkeo-main-1",
					Amount:            "850000",
					Asset:             "uarkeo",
					PayTo:             "arkeo1yourprovideraddress",
					MaxTimeoutSeconds: 60,
					Extra: map[string]interface{}{
						"name":     "ARKEO",
						"discount": "15%",
						"note":     "Pay with ARKEO for 15% off!",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-X402-Version", "2")
		w.WriteHeader(http.StatusPaymentRequired)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Payment provided - return mock RPC response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Settlement-ID", fmt.Sprintf("settlement-%d", time.Now().Unix()))
	w.Header().Set("X-Payment-Status", "verified")

	response := RPCResponse{
		JSONRPC: "2.0",
		ID:      1,
		Result:  "0x134e82a", // Mock block number
	}

	json.NewEncoder(w).Encode(response)
}
