//go:build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/fengziwei0826/caching_proxy_exe/internal/conf"
	"github.com/fengziwei0826/caching_proxy_exe/internal/server"
	"github.com/fengziwei0826/caching_proxy_exe/pkg/db"
)

func InitCacheProxyServer() *server.CacheProxyServer {
	wire.Build(context.Background, conf.GetGlobalConfig, db.NewRequestCacheProxy, server.NewHttpProxyServer, server.NewCacheProxyServer)
	return new(server.CacheProxyServer)
}
