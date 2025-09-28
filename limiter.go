package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	capacity       int // also the burst rate
	refillRate     int
	refillInterval time.Duration
	redisAddr      string
}

type RateLimiter struct {
	Config
	redisClient *redis.Client
	scriptSHA   string
}

func NewRateLimiter(cfg Config) *RateLimiter {
	if cfg.redisAddr == "" {
		log.Fatal("redis addr expected")
		return nil
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.redisAddr,
	})
	scriptRaw, err := os.ReadFile("limit.lua")
	if err != nil {
		panic(err)
	}

	script := string(scriptRaw)

	sha, err := rdb.ScriptLoad(ctx, script).Result()
	if err != nil {
		log.Fatal(err)
	}

	return &RateLimiter{
		Config:      cfg,
		redisClient: rdb,
		scriptSHA:   sha,
	}
}

var ctx = context.Background()

func (r *RateLimiter) Allow(id string) bool {
	key := fmt.Sprintf("bucket:user:%s", id)
	now := time.Now().UnixMilli()

	res, err := r.redisClient.EvalSha(ctx, r.scriptSHA, []string{key},
		r.capacity,
		r.refillRate,
		r.refillInterval.Milliseconds(),
		now,
	).Result()

	if err != nil {
		log.Println("evalsha error")
		log.Print(err)
		return false
	}
	allowed, ok := res.(int64)
	return ok && allowed == 1
}

func (r *RateLimiter) GetQuota(id string) string {
	key := fmt.Sprintf("bucket:user:%s", id)
	res, err := r.redisClient.HGet(ctx, key, "tokens").Result()
	if err != nil {
		log.Println("error getting quota")
		return ""
	}
	return res
}
