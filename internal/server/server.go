package server

import (
	"context"

	"github.com/fengziwei0826/caching_proxy_exe/internal/conf"
)

type CacheProxyServer struct {
	ctx context.Context
	srv HttpProxyServer
}

func (c *CacheProxyServer) Start() error {
	return c.srv.Start()
}

func (c *CacheProxyServer) Stop() error {
	return c.srv.Stop()
}

func NewCacheProxyServer(ctx context.Context, srv HttpProxyServer, config *conf.GlobalConfig) *CacheProxyServer {
	return &CacheProxyServer{
		ctx: ctx,
		srv: srv,
	}
}
