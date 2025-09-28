local bucket_key = KEYS[1]
local capacity = ARGV[1]
local refillRate = ARGV[2]
local refillInterval = ARGV[3]
local now = ARGV[4]

local data = redis.call("HMGET", bucket_key, "tokens", "lastRefill")
local tokens = tonumber(data[1]) or capacity
local lastRefill = tonumber(data[2]) or now

if not tokens then 
   tokens = capacity
   lastRefill = now
end

local elapsed = now - lastRefill
if elapsed > 0 then
   local intervals = math.floor(elapsed / refillInterval)
   local refill = intervals * refillRate
   if refill > 0 then
      tokens = math.min(capacity, tokens + refill)
      lastRefill = lastRefill + (intervals * refillInterval)
   end
end

local allowed = 0
if tokens > 0 then
   tokens = tokens - 1
   allowed = 1
end

redis.call("HSET", bucket_key, "tokens", tokens, "lastRefill", lastRefill)

return allowed