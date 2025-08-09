package server

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fengziwei0826/caching_proxy_exe/internal/cachemanager"
	"github.com/fengziwei0826/caching_proxy_exe/internal/conf"
)

const (
	cacheHeader = "X-Cache"
	cacheHit    = "HIT"
	cacheMiss   = "MISS"
)

type httpProxyServer struct {
	ctx context.Context
	srv http.Server
	mgr cachemanager.CacheManager
}

func (p *httpProxyServer) Start() error {
	go func() {
		if err := p.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("HTTP server failed to start: %v \n", err)
		}
	}()
	return nil
}

func (p *httpProxyServer) Stop() error {
	ctx, cancel := context.WithTimeout(p.ctx, 5*time.Second)
	defer cancel()
	if err := p.srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown failed: %v \n", err)
	}
	p.mgr.Close()
	return nil
}

func (p *httpProxyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s \n", r.Method, r.URL.Path)
	path := r.URL.Path
	res, err := p.mgr.GetResponse(path)
	if err == nil && res != "" {
		log.Printf("Cache hit for path: %s [%s]\n", path, res)
		w.WriteHeader(http.StatusOK)
		w.Header().Set(cacheHeader, cacheHit)
		w.Write([]byte(res))
		return
	}
	req, err := http.NewRequest(r.Method, conf.GetGlobalConfig().HTTPConfig.ProxyAddr+r.URL.Path, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
	if err != nil {
		log.Printf("Failed to create request: %v \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to forward request: %v \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	log.Printf("Forwarded request to %s, received status: %s resp: %v \n", req.URL, resp.Status, resp)
	respBytes, err := io.ReadAll(resp.Body) // Read the response body to avoid resource leak
	if err != nil {
		log.Printf("Failed to marshal response: %v \n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if err := p.mgr.CacheResponse(path, string(respBytes)); err != nil {
		log.Printf("Failed to cache response: %v \n", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set(cacheHeader, cacheMiss)
	w.Write(respBytes)
}

func NewHttpProxyServer(ctx context.Context, cache cachemanager.CacheManager, config *conf.GlobalConfig) HttpProxyServer {
	proxy := &httpProxyServer{
		ctx: ctx,
		mgr: cache,
	}
	proxy.srv = http.Server{
		Addr:    "localhost:" + strconv.Itoa(config.HTTPConfig.Port),
		Handler: proxy,
	}
	return proxy
}
