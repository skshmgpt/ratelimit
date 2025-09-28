package main

// import "time"

// type TokenBucket struct {
// 	Config
// 	lastRefill time.Time
// }

// func NewTokenBucket(cfg Config) *TokenBucket {
// 	return &TokenBucket{
// 		Config:     cfg,
// 		lastRefill: time.Now(),
// 	}
// }

// func (b *TokenBucket) Allow() bool {
// 	now := time.Now()
// 	elapsed := now.Sub(b.lastRefill)
// 	refillUnits := int8(elapsed / b.RefillInterval)
// 	if refillUnits > 0 {
// 		b.Tokens += int8(refillUnits) * b.RefillRate
// 		if b.Tokens > b.Capacity {
// 			b.Tokens = b.Capacity
// 		}
// 		b.lastRefill = b.lastRefill.Add(time.Duration(refillUnits) * b.RefillInterval)
// 	}
// 	if b.Tokens > 0 {
// 		b.Tokens--
// 		return true
// 	}
// 	return false
// }
