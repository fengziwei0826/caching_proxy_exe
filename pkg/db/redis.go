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

type CacheProxy interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	SAdd(key string, value string) error
	GetSet(key string) ([]string, error)
	Del(keys ...string) (int64, error)
	LPush(key, value string) error
	GetList(key string) ([]string, error)
	Close() error
}

type requestCacheProxy struct {
	ctx    context.Context
	client *redis.Client
}

func (r *requestCacheProxy) Set(key string, value string) error {
	err := r.client.Set(r.ctx, key, value, time.Hour*3).Err()
	if err != nil {
		return err
	}
	return nil
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

func (r *requestCacheProxy) SAdd(key string, value string) error {
	err := r.client.SAdd(r.ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *requestCacheProxy) GetSet(key string) ([]string, error) {
	val, err := r.client.SMembers(r.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *requestCacheProxy) Del(keys ...string) (int64, error) {
	cnt := int64(0)
	for _, key := range keys {
		n, err := r.client.Del(r.ctx, key).Result()
		if err != nil {
			return cnt, err
		}
		cnt += n
	}
	return cnt, nil
}

func (r *requestCacheProxy) LPush(key, value string) error {
	err := r.client.LPush(r.ctx, key, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *requestCacheProxy) GetList(key string) ([]string, error) {
	val, err := r.client.LRange(r.ctx, key, 0, -1).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return val, nil
}

func (r *requestCacheProxy) Close() error {
	if err := r.client.Close(); err != nil {
		log.Printf("Failed to close Redis client: %v \n", err)
		return err
	}
	return nil
}

func NewCacheProxy(ctx context.Context, config *conf.GlobalConfig) CacheProxy {
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
