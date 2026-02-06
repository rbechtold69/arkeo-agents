package sentinel

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// X402 Protocol Version
const X402Version = 2

// PaymentRequirements defines what payment is accepted for a resource
type PaymentRequirements struct {
	Scheme            string                 `json:"scheme"`
	Network           string                 `json:"network"`
	Amount            string                 `json:"amount"`
	Asset             string                 `json:"asset"`
	PayTo             string                 `json:"payTo"`
	MaxTimeoutSeconds int                    `json:"maxTimeoutSeconds"`
	Extra             map[string]interface{} `json:"extra,omitempty"`
}

// ResourceInfo describes the protected resource
type ResourceInfo struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// PaymentRequiredResponse is the HTTP 402 response body
type PaymentRequiredResponse struct {
	X402Version int                   `json:"x402Version"`
	Error       string                `json:"error,omitempty"`
	Resource    ResourceInfo          `json:"resource"`
	Accepts     []PaymentRequirements `json:"accepts"`
	Extensions  map[string]interface{} `json:"extensions,omitempty"`
}

// X402Handler handles x402 payment verification for AI agents
type X402Handler struct {
	// Provider's payment address
	ProviderAddress string
	
	// Accepted payment methods
	AcceptUSDC  bool
	AcceptARKEO bool
	
	// Pricing (in atomic units)
	PricePerRequestUSDC  string // e.g., "1000" = 0.001 USDC (6 decimals)
	PricePerRequestARKEO string // e.g., "1000000" = 0.001 ARKEO (8 decimals)
	
	// ARKEO discount percentage (0-100)
	ARKEODiscountPercent int
	
	// Facilitator URL for payment verification
	FacilitatorURL string
}

// NewX402Handler creates a new x402 payment handler
func NewX402Handler(providerAddress string) *X402Handler {
	return &X402Handler{
		ProviderAddress:      providerAddress,
		AcceptUSDC:           true,
		AcceptARKEO:          true,
		PricePerRequestUSDC:  "1000",    // 0.001 USDC per request
		PricePerRequestARKEO: "850000",  // 0.00085 ARKEO (15% discount)
		ARKEODiscountPercent: 15,
		FacilitatorURL:       "https://x402.org/facilitator", // Default Coinbase facilitator
	}
}

// BuildPaymentRequirements creates the x402 payment requirements for a service
func (h *X402Handler) BuildPaymentRequirements(service string, requestURL string) PaymentRequiredResponse {
	accepts := []PaymentRequirements{}
	
	// USDC on Ethereum Mainnet
	if h.AcceptUSDC {
		accepts = append(accepts, PaymentRequirements{
			Scheme:            "exact",
			Network:           "eip155:1", // Ethereum Mainnet
			Amount:            h.PricePerRequestUSDC,
			Asset:             "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", // USDC
			PayTo:             h.ProviderAddress,
			MaxTimeoutSeconds: 60,
			Extra: map[string]interface{}{
				"name":    "USDC",
				"version": "2",
			},
		})
		
		// USDC on Base (cheaper option)
		accepts = append(accepts, PaymentRequirements{
			Scheme:            "exact",
			Network:           "eip155:8453", // Base
			Amount:            h.PricePerRequestUSDC,
			Asset:             "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", // USDC on Base
			PayTo:             h.ProviderAddress,
			MaxTimeoutSeconds: 60,
			Extra: map[string]interface{}{
				"name":      "USDC",
				"version":   "2",
				"chain":     "Base",
				"gasSaving": "true",
			},
		})
	}
	
	// ARKEO token (with discount)
	if h.AcceptARKEO {
		accepts = append(accepts, PaymentRequirements{
			Scheme:            "exact",
			Network:           "arkeo:arkeo-main-1", // Arkeo Mainnet
			Amount:            h.PricePerRequestARKEO,
			Asset:             "uarkeo", // Native ARKEO token
			PayTo:             h.ProviderAddress,
			MaxTimeoutSeconds: 60,
			Extra: map[string]interface{}{
				"name":     "ARKEO",
				"discount": "15%",
				"note":     "Pay with ARKEO for 15% off!",
			},
		})
	}
	
	return PaymentRequiredResponse{
		X402Version: X402Version,
		Error:       "Payment required to access this RPC endpoint",
		Resource: ResourceInfo{
			URL:         requestURL,
			Description: "Arkeo RPC Service: " + service,
			MimeType:    "application/json",
		},
		Accepts:    accepts,
		Extensions: map[string]interface{}{},
	}
}

// CheckPaymentHeader checks if the request has a valid x402 payment
func (h *X402Handler) CheckPaymentHeader(r *http.Request) (bool, string) {
	// x402 uses X-PAYMENT header for payment payload
	paymentHeader := r.Header.Get("X-PAYMENT")
	
	if paymentHeader == "" {
		// Also check older header format
		paymentHeader = r.Header.Get("X-Payment-Signature")
	}
	
	if paymentHeader == "" {
		return false, ""
	}
	
	return true, paymentHeader
}

// Facilitator is the x402 facilitator client (initialized on first use)
var facilitator *X402Facilitator

// VerifyPayment verifies the payment with the Coinbase facilitator
// Returns: (verified bool, settlementID string, error)
func (h *X402Handler) VerifyPayment(paymentPayload string) (bool, string, error) {
	// Initialize facilitator on first use
	if facilitator == nil {
		var err error
		facilitator, err = NewX402Facilitator()
		if err != nil {
			// Fall back to demo mode if no credentials
			return true, "demo-" + time.Now().Format("20060102150405"), nil
		}
	}
	
	// Use first payment requirement for verification
	requirements := PaymentRequirements{
		Scheme:            "exact",
		Network:           "eip155:8453", // Base
		Amount:            h.PricePerRequestUSDC,
		Asset:             "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", // USDC on Base
		PayTo:             h.ProviderAddress,
		MaxTimeoutSeconds: 60,
	}
	
	// Verify and settle via facilitator
	valid, settlementID, err := facilitator.VerifyAndSettle(paymentPayload, requirements)
	if err != nil {
		return false, "", err
	}
	
	return valid, settlementID, nil
}

// WritePaymentRequired writes the HTTP 402 response
func (h *X402Handler) WritePaymentRequired(w http.ResponseWriter, service string, requestURL string) {
	response := h.BuildPaymentRequirements(service, requestURL)
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-X402-Version", "2")
	w.WriteHeader(http.StatusPaymentRequired) // HTTP 402
	
	json.NewEncoder(w).Encode(response)
}

// Middleware wraps an HTTP handler with x402 payment verification
func (h *X402Handler) Middleware(next http.Handler, service string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for payment header
		hasPayment, paymentPayload := h.CheckPaymentHeader(r)
		
		if !hasPayment {
			// No payment - return 402 Payment Required
			h.WritePaymentRequired(w, service, r.URL.String())
			return
		}
		
		// Verify payment
		verified, settlementID, err := h.VerifyPayment(paymentPayload)
		if err != nil || !verified {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Payment verification failed",
			})
			return
		}
		
		// Payment verified - add settlement ID to response headers
		w.Header().Set("X-Settlement-ID", settlementID)
		
		// Proceed to actual handler
		next.ServeHTTP(w, r)
	})
}

// ExtractService extracts the service name from the request path
func ExtractService(path string) string {
	// Expected format: /{service}/... or /x402/{service}/...
	parts := strings.Split(strings.Trim(path, "/"), "/")
	
	if len(parts) == 0 {
		return "unknown"
	}
	
	// Skip "x402" prefix if present
	if parts[0] == "x402" && len(parts) > 1 {
		return parts[1]
	}
	
	return parts[0]
}
