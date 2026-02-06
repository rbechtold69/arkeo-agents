package sentinel

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

// X402 Route Constants
const (
	// x402 payment requirements endpoint
	RouteX402Requirements = "/x402/requirements/{service}"
	
	// x402-enabled RPC endpoint (alternative to standard routes)
	RouteX402RPC = "/x402/{service}/{path:.*}"
)

// x402Handler is the handler instance for the proxy
var x402Handler *X402Handler

// InitX402 initializes the x402 payment handler
func (p *Proxy) InitX402(providerAddress string) {
	x402Handler = NewX402Handler(providerAddress)
	p.logger.Info("x402 payment handler initialized")
}

// RegisterX402Routes adds x402-specific routes to the router
func (p *Proxy) RegisterX402Routes(router *mux.Router) {
	// Payment requirements endpoint - agents query this to know how to pay
	router.HandleFunc(RouteX402Requirements, p.handleX402Requirements).Methods(http.MethodGet)
	
	// x402-enabled RPC proxy - checks payment before routing
	router.PathPrefix("/x402/").Handler(p.x402Middleware(http.HandlerFunc(p.handleX402Proxy)))
	
	p.logger.Info("x402 routes registered")
}

// handleX402Requirements returns the payment requirements for a service
func (p *Proxy) handleX402Requirements(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	
	if service == "" {
		http.Error(w, "service required", http.StatusBadRequest)
		return
	}
	
	// Check if service exists
	p.serviceMu.RLock()
	_, exists := p.serviceIDs[service]
	p.serviceMu.RUnlock()
	
	if !exists {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	
	// Return payment requirements
	if x402Handler == nil {
		http.Error(w, "x402 not initialized", http.StatusInternalServerError)
		return
	}
	
	requirements := x402Handler.BuildPaymentRequirements(service, r.URL.String())
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-X402-Version", "2")
	json.NewEncoder(w).Encode(requirements)
}

// x402Middleware checks for valid payment before allowing access
func (p *Proxy) x402Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if x402Handler == nil {
			http.Error(w, "x402 not initialized", http.StatusInternalServerError)
			return
		}
		
		// Extract service from path: /x402/{service}/...
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		service := ""
		if len(pathParts) >= 2 {
			service = pathParts[1]
		}
		
		// Check for payment header
		hasPayment, paymentPayload := x402Handler.CheckPaymentHeader(r)
		
		if !hasPayment {
			// Return 402 Payment Required
			x402Handler.WritePaymentRequired(w, service, r.URL.String())
			return
		}
		
		// Verify payment
		verified, settlementID, err := x402Handler.VerifyPayment(paymentPayload)
		if err != nil {
			p.logger.Error("x402 payment verification failed", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "payment verification failed",
				"details": err.Error(),
			})
			return
		}
		
		if !verified {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusPaymentRequired)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "payment not verified",
			})
			return
		}
		
		// Payment verified - add settlement ID to response headers
		w.Header().Set("X-Settlement-ID", settlementID)
		w.Header().Set("X-Payment-Status", "verified")
		
		// Log successful payment
		p.logger.Info("x402 payment verified", 
			"service", service,
			"settlement_id", settlementID,
		)
		
		// Proceed to actual handler
		next.ServeHTTP(w, r)
	})
}

// handleX402Proxy handles x402-enabled RPC requests
func (p *Proxy) handleX402Proxy(w http.ResponseWriter, r *http.Request) {
	// Extract service from path: /x402/{service}/...
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	
	if len(pathParts) < 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	
	service := pathParts[1]
	
	// Get the remaining path after /x402/{service}/
	remainingPath := ""
	if len(pathParts) > 2 {
		remainingPath = "/" + strings.Join(pathParts[2:], "/")
	}
	
	// Look up the proxy for this service
	p.proxyMu.RLock()
	targetURL, exists := p.proxies[service]
	p.proxyMu.RUnlock()
	
	if !exists || targetURL == nil {
		http.Error(w, "service not found", http.StatusNotFound)
		return
	}
	
	// Modify the request to proxy to the actual service
	r.URL.Path = remainingPath
	r.URL.Host = targetURL.Host
	r.URL.Scheme = targetURL.Scheme
	r.Host = targetURL.Host
	
	// Use the existing proxy logic
	p.handleHTTPWithURL(w, r, targetURL)
}

// handleHTTPWithURL proxies the request to the given URL
// This wraps the existing proxy functionality
func (p *Proxy) handleHTTPWithURL(w http.ResponseWriter, r *http.Request, target *url.URL) {
	// Create a reverse proxy
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = target.Scheme
			req.URL.Host = target.Host
			req.Host = target.Host
			
			// Preserve the path from the original request
			if r.URL.Path != "" {
				req.URL.Path = r.URL.Path
			}
			
			// Copy headers
			if r.Header.Get("Content-Type") != "" {
				req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
			}
		},
	}
	
	proxy.ServeHTTP(w, r)
}
