package sentinel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// X402Facilitator handles payment verification via Coinbase facilitator
type X402Facilitator struct {
	APIKey       string
	APISecret    string
	FacilitatorURL string
	HTTPClient   *http.Client
}

// VerifyRequest is sent to the facilitator
type VerifyRequest struct {
	PaymentPayload string `json:"paymentPayload"`
	Requirements   PaymentRequirements `json:"requirements"`
}

// VerifyResponse from the facilitator
type VerifyResponse struct {
	Valid        bool   `json:"valid"`
	SettlementID string `json:"settlementId,omitempty"`
	Error        string `json:"error,omitempty"`
	TxHash       string `json:"txHash,omitempty"`
}

// SettleRequest to settle a payment
type SettleRequest struct {
	PaymentPayload string `json:"paymentPayload"`
}

// SettleResponse from facilitator
type SettleResponse struct {
	Success      bool   `json:"success"`
	SettlementID string `json:"settlementId,omitempty"`
	TxHash       string `json:"txHash,omitempty"`
	Error        string `json:"error,omitempty"`
}

// NewX402Facilitator creates a facilitator client from environment/file
func NewX402Facilitator() (*X402Facilitator, error) {
	// Try to load from file first
	apiKey, apiSecret, err := loadCredentialsFromFile()
	if err != nil {
		// Fall back to environment variables
		apiKey = os.Getenv("X402_API_KEY")
		apiSecret = os.Getenv("X402_API_SECRET")
	}
	
	if apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("x402 credentials not found")
	}
	
	facilitatorURL := os.Getenv("X402_FACILITATOR_URL")
	if facilitatorURL == "" {
		facilitatorURL = "https://x402.coinbase.com"
	}
	
	return &X402Facilitator{
		APIKey:         apiKey,
		APISecret:      apiSecret,
		FacilitatorURL: facilitatorURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// loadCredentialsFromFile reads credentials from ~/.x402-credentials
func loadCredentialsFromFile() (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", "", err
	}
	
	data, err := os.ReadFile(home + "/.x402-credentials")
	if err != nil {
		return "", "", err
	}
	
	var apiKey, apiSecret string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "X402_API_KEY=") {
			apiKey = strings.TrimPrefix(line, "X402_API_KEY=")
		} else if strings.HasPrefix(line, "X402_API_SECRET=") {
			apiSecret = strings.TrimPrefix(line, "X402_API_SECRET=")
		}
	}
	
	if apiKey == "" || apiSecret == "" {
		return "", "", fmt.Errorf("credentials incomplete")
	}
	
	return apiKey, apiSecret, nil
}

// Verify checks if a payment payload is valid
func (f *X402Facilitator) Verify(paymentPayload string, requirements PaymentRequirements) (*VerifyResponse, error) {
	reqBody := VerifyRequest{
		PaymentPayload: paymentPayload,
		Requirements:   requirements,
	}
	
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", f.FacilitatorURL+"/verify", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", f.APIKey)
	req.Header.Set("Authorization", "Bearer "+f.APISecret)
	
	resp, err := f.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("facilitator request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("facilitator returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var verifyResp VerifyResponse
	if err := json.Unmarshal(body, &verifyResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &verifyResp, nil
}

// Settle settles a verified payment on-chain
func (f *X402Facilitator) Settle(paymentPayload string) (*SettleResponse, error) {
	reqBody := SettleRequest{
		PaymentPayload: paymentPayload,
	}
	
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", f.FacilitatorURL+"/settle", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", f.APIKey)
	req.Header.Set("Authorization", "Bearer "+f.APISecret)
	
	resp, err := f.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("facilitator request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("facilitator returned status %d: %s", resp.StatusCode, string(body))
	}
	
	var settleResp SettleResponse
	if err := json.Unmarshal(body, &settleResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	return &settleResp, nil
}

// VerifyAndSettle does both verification and settlement in one call
func (f *X402Facilitator) VerifyAndSettle(paymentPayload string, requirements PaymentRequirements) (bool, string, error) {
	// First verify
	verifyResp, err := f.Verify(paymentPayload, requirements)
	if err != nil {
		return false, "", fmt.Errorf("verification failed: %w", err)
	}
	
	if !verifyResp.Valid {
		return false, "", fmt.Errorf("payment not valid: %s", verifyResp.Error)
	}
	
	// Then settle
	settleResp, err := f.Settle(paymentPayload)
	if err != nil {
		return false, "", fmt.Errorf("settlement failed: %w", err)
	}
	
	if !settleResp.Success {
		return false, "", fmt.Errorf("settlement failed: %s", settleResp.Error)
	}
	
	return true, settleResp.SettlementID, nil
}
