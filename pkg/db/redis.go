package db

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/fengziwei0826/caching_proxy_exe/internal/conf"
)

const (
	poolSize = 200
	minIdle  = 20
	maxIdle  = 50
)

type RequestCacheProxy interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type requestCacheProxy struct {
	ctx    context.Context
	client *redis.Client
}

func (r *requestCacheProxy) Get(key string) (string, error) {
	log.Printf("Fetching key from cache: %s \n", key)
	val, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *requestCacheProxy) Set(key string, value string) error {
	err := r.client.Set(r.ctx, key, value, time.Hour*3).Err()
	if err != nil {
		return err
	}
	return nil
}

func NewRequestCacheProxy(ctx context.Context, config *conf.GlobalConfig) RequestCacheProxy {
	r := &requestCacheProxy{
		ctx: ctx,
	}
	r.client = redis.NewClient(&redis.Options{
		Addr:         config.RedisConfig.Addr,
		Password:     config.RedisConfig.Password,
		DB:           config.RedisConfig.DB,
		PoolSize:     poolSize,
		MinIdleConns: minIdle,
		MaxIdleConns: maxIdle,
	})
	return r

}
