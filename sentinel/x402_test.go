package sentinel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestX402Handler_BuildPaymentRequirements(t *testing.T) {
	handler := NewX402Handler("0x1234567890abcdef1234567890abcdef12345678")
	
	requirements := handler.BuildPaymentRequirements("eth", "/eth/blockNumber")
	
	// Check version
	if requirements.X402Version != 2 {
		t.Errorf("Expected X402Version 2, got %d", requirements.X402Version)
	}
	
	// Check we have payment options
	if len(requirements.Accepts) < 2 {
		t.Errorf("Expected at least 2 payment options, got %d", len(requirements.Accepts))
	}
	
	// Check USDC option exists
	hasUSDC := false
	for _, opt := range requirements.Accepts {
		if opt.Extra != nil {
			if name, ok := opt.Extra["name"].(string); ok && name == "USDC" {
				hasUSDC = true
				break
			}
		}
	}
	if !hasUSDC {
		t.Error("Expected USDC payment option")
	}
	
	// Check ARKEO option exists with discount
	hasARKEO := false
	for _, opt := range requirements.Accepts {
		if opt.Extra != nil {
			if name, ok := opt.Extra["name"].(string); ok && name == "ARKEO" {
				hasARKEO = true
				if discount, ok := opt.Extra["discount"].(string); !ok || discount != "15%" {
					t.Error("Expected ARKEO 15% discount")
				}
				break
			}
		}
	}
	if !hasARKEO {
		t.Error("Expected ARKEO payment option")
	}
	
	t.Logf("Payment requirements built successfully with %d options", len(requirements.Accepts))
}

func TestX402Handler_WritePaymentRequired(t *testing.T) {
	handler := NewX402Handler("0x1234567890abcdef1234567890abcdef12345678")
	
	// Create a test request
	req := httptest.NewRequest(http.MethodGet, "/eth/blockNumber", nil)
	rec := httptest.NewRecorder()
	
	// Write payment required response
	handler.WritePaymentRequired(rec, "eth", req.URL.String())
	
	// Check status code
	if rec.Code != http.StatusPaymentRequired {
		t.Errorf("Expected status 402, got %d", rec.Code)
	}
	
	// Check content type
	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
	
	// Check X-X402-Version header
	x402Version := rec.Header().Get("X-X402-Version")
	if x402Version != "2" {
		t.Errorf("Expected X-X402-Version 2, got %s", x402Version)
	}
	
	// Parse response body
	var response PaymentRequiredResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	
	t.Logf("402 response generated successfully: %d payment options", len(response.Accepts))
}

func TestX402Handler_CheckPaymentHeader(t *testing.T) {
	handler := NewX402Handler("0x1234567890abcdef1234567890abcdef12345678")
	
	// Test without payment header
	req1 := httptest.NewRequest(http.MethodGet, "/eth/blockNumber", nil)
	hasPayment, _ := handler.CheckPaymentHeader(req1)
	if hasPayment {
		t.Error("Expected no payment header")
	}
	
	// Test with X-PAYMENT header
	req2 := httptest.NewRequest(http.MethodGet, "/eth/blockNumber", nil)
	req2.Header.Set("X-PAYMENT", "test-payment-payload")
	hasPayment, payload := handler.CheckPaymentHeader(req2)
	if !hasPayment {
		t.Error("Expected payment header to be found")
	}
	if payload != "test-payment-payload" {
		t.Errorf("Expected payload 'test-payment-payload', got '%s'", payload)
	}
	
	// Test with legacy X-Payment-Signature header
	req3 := httptest.NewRequest(http.MethodGet, "/eth/blockNumber", nil)
	req3.Header.Set("X-Payment-Signature", "legacy-payment")
	hasPayment, payload = handler.CheckPaymentHeader(req3)
	if !hasPayment {
		t.Error("Expected legacy payment header to be found")
	}
	if payload != "legacy-payment" {
		t.Errorf("Expected payload 'legacy-payment', got '%s'", payload)
	}
	
	t.Log("Payment header detection working correctly")
}

func TestExtractService(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"/x402/eth/blockNumber", "eth"},
		{"/x402/cosmos/status", "cosmos"},
		{"/eth/blockNumber", "eth"},
		{"/thorchain/pools", "thorchain"},
		{"", "unknown"},
	}
	
	for _, tt := range tests {
		result := ExtractService(tt.path)
		if result != tt.expected {
			t.Errorf("ExtractService(%s) = %s, expected %s", tt.path, result, tt.expected)
		}
	}
	
	t.Log("Service extraction working correctly")
}
