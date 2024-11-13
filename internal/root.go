package internal

import (
	"catyousha/caching-proxy/internal/cache"
	"fmt"
	"io"
	"net/http"
)

var httpServer *http.Server

func SetupProxy(port int, origin string) (*string, error) {
	if origin == "" {
		return nil, fmt.Errorf("origin is required")
	}
	if port < 1 || port > 65535 {
		return nil, fmt.Errorf("port must be between 1 and 65535")
	}

	addr := fmt.Sprintf("localhost:%d", port)

	httpServer = &http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "GET" {
				writeError(w, fmt.Errorf("method not allowed: %s", r.Method))
				return
			}
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

			proxyReq, err := http.NewRequest(r.Method, origin+r.URL.String(), r.Body)
			if err != nil {
				writeError(w, fmt.Errorf("error creating proxy request: %v", err))
				return
			}

			client := http.Client{}
			w.Header().Set("X-Cache", "MISS")
			resp, err := client.Do(proxyReq)
			if err != nil {
				writeError(w, fmt.Errorf("error forwarding request: %v", err))
				return
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				writeError(w, fmt.Errorf("error reading response body: %v", err))
				return
			}
			w.Write(body)
			cache.Set(cacheKey, string(body))
		}),
	}
	fmt.Printf("Listening on %s ...\n", addr)
	httpServer.ListenAndServe()
	return &origin, nil
}

func ClearCache() {
	cache.ClearAll()
}
