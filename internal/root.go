package internal

import (
	"catyousha/caching-proxy/internal/cache"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	httpServer *http.Server
	client = &http.Client{
		Timeout: 30 * time.Second,
	}
)

func SetupProxy(port int, origin string) error {
	if origin == "" {
		return fmt.Errorf("origin is required")
	}
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	addr := fmt.Sprintf("localhost:%d", port)

	httpServer = &http.Server{
		Addr:    addr,
		Handler: createProxyHandler(origin),
		ReadTimeout: 15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second,
	}
	fmt.Printf("Listening on %s ...\n", addr)
	return httpServer.ListenAndServe()
}

func createProxyHandler(origin string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, fmt.Errorf("method not allowed: %s", r.Method))
			return
		}

		// Find Cached
		cacheKey := r.URL.String()
		cached, err := cache.Get(cacheKey);
		if err != nil {
			writeError(w, fmt.Errorf("error getting cache: %v", err))
			return
		}
		if cached != nil {
			w.Header().Set("X-Cache", "HIT")
			w.Write([]byte(*cached))
			return
		}

		// Create Request
		proxyReq, err := http.NewRequest(r.Method, origin+r.URL.String(), r.Body)
		if err != nil {
			writeError(w, fmt.Errorf("error creating proxy request: %v", err))
			return
		}

		copyHeaders(proxyReq.Header, r.Header)
		w.Header().Set("X-Cache", "MISS")
		
		resp, err := client.Do(proxyReq)
		if err != nil {
			writeError(w, fmt.Errorf("error forwarding request: %v", err))
			return
		}
		defer resp.Body.Close()
		
		// Handle Response
		copyHeaders(w.Header(), resp.Header)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			writeError(w, fmt.Errorf("error reading response body: %v", err))
			return
		}
		w.Write(body)
		cache.Set(cacheKey, string(body))
	})
}

func ClearCache() {
	cache.ClearAll()
}
